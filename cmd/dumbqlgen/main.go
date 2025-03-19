package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS

type Config struct {
	StructType   string
	PackagePath  string
	Output       string
	PrintVersion bool
	PrintHelp    bool
}

type FieldInfo struct {
	Name    string
	TagName string
	Type    string
	Skip    bool
}

type StructInfo struct {
	Name       string
	Package    string
	Fields     []FieldInfo
	NestedInfo []StructInfo
}

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: dumbql [flags]")
		fmt.Fprintln(os.Stderr, "More documentation at https://pkg.go.dev/go.tomakado.io/dumbql/cmd/dumbqlgen")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		os.Exit(2) //nolint:mnd
	}

	var config Config
	flag.StringVar(&config.StructType, "type", "", "Struct type to generate router for")
	flag.StringVar(&config.PackagePath, "package", ".", "Package path containing the struct")
	flag.StringVar(&config.Output, "output", "", "Output file path")
	flag.BoolVar(&config.PrintVersion, "version", false, "Print version")
	flag.BoolVar(&config.PrintHelp, "help", false, "Print help")
	flag.Parse()

	if config.PrintHelp {
		flag.Usage()
	}

	if config.PrintVersion {
		printVersion()
	}

	if config.StructType == "" {
		handleError(errors.New("missing required flag: -type"))
	}

	// If output is not specified, use the current directory with structtype_matcher.gen.go
	if config.Output == "" {
		config.Output = strings.ToLower(config.StructType) + "_matcher.gen.go"
	}

	code, err := generate(config)
	if err != nil {
		handleError(fmt.Errorf("generate router: %w", err))
	}

	err = os.WriteFile(config.Output, []byte(code), 0o600) //nolint:mnd
	if err != nil {
		handleError(fmt.Errorf("write output file: %w", err))
	}
}

func printVersion() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Println("dumbql unknown version")
		return
	}

	fmt.Printf("dumbql %s\n", info.Main.Version)
	fmt.Printf("commit: %s at %s\n", getCommitHash(info), getCommitTime(info))
	fmt.Printf("go version: %s\n", info.GoVersion)
	os.Exit(0)
}

func getCommitHash(info *debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value
		}
	}

	return "(commit unknown)"
}

func getCommitTime(info *debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.time" {
			return setting.Value
		}
	}

	return "(build time unknown)"
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintln(os.Stderr, "Run dumbqlgen -help for more information")
		os.Exit(1)
	}
}

func generate(config Config) (string, error) {
	pkgs, err := loadPackages(config.PackagePath)
	if err != nil {
		return "", fmt.Errorf("failed to load package: %w", err)
	}

	if len(pkgs) == 0 {
		return "", fmt.Errorf("no packages found at %s", config.PackagePath)
	}

	pkg := pkgs[0]
	packageName := pkg.Name

	structInfo, err := findStruct(pkg, config.StructType)
	if err != nil {
		return "", err
	}

	tmpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("failed to read build info")
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "router.tmpl", map[string]any{
		"Package":    packageName,
		"StructInfo": structInfo,
		"Version":    buildInfo.Main.Version,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func loadPackages(path string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedSyntax | packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, err
	}

	return pkgs, nil
}

func findStruct(pkg *packages.Package, structName string) (StructInfo, error) {
	var structInfo StructInfo
	structInfo.Name = structName
	structInfo.Package = pkg.Name

	// Find the struct declaration
	var structType *types.Struct
	obj := pkg.Types.Scope().Lookup(structName)
	if obj == nil {
		return structInfo, fmt.Errorf("struct %s not found in package %s", structName, pkg.Name)
	}

	// Check if it's a named type and it's a struct
	named, ok := obj.Type().(*types.Named)
	if !ok {
		return structInfo, fmt.Errorf("%s is not a named type", structName)
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok {
		return structInfo, fmt.Errorf("%s is not a struct", structName)
	}
	structType = st

	// Extract field information
	for i := range structType.NumFields() {
		field := structType.Field(i)
		if !field.Exported() {
			continue
		}

		fieldName := field.Name()
		fieldType := field.Type().String()

		// Get the AST node for this field to extract struct tags
		var fieldInfo FieldInfo
		fieldInfo.Name = fieldName
		fieldInfo.TagName = fieldName // Default to field name
		fieldInfo.Type = fieldType

		findAndEnrichNode(pkg.Syntax, structName, fieldName, &fieldInfo)
		structInfo.Fields = append(structInfo.Fields, fieldInfo)
	}

	return structInfo, nil
}

func findAndEnrichNode(files []*ast.File, structName, fieldName string, fieldInfo *FieldInfo) {
	for _, file := range files {
		ast.Inspect(file, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				return true
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return false
			}

			if err := enrichFieldInfo(structType, fieldName, fieldInfo); err != nil {
				panic(err)
			}

			return false
		})
	}
}

func enrichFieldInfo(structType *ast.StructType, fieldName string, fieldInfo *FieldInfo) error {
	for _, field := range structType.Fields.List {
		for _, ident := range field.Names {
			if ident.Name != fieldName {
				continue
			}

			if field.Tag == nil {
				continue
			}

			tag := strings.Trim(field.Tag.Value, "`")
			dumbqlTag, err := extractTag(tag, "dumbql")
			if err != nil {
				return err
			}

			if dumbqlTag == "-" {
				fieldInfo.Skip = true
				break
			}

			fieldInfo.TagName = dumbqlTag

			break
		}
	}

	return nil
}

func extractTag(tag, key string) (string, error) {
	if tag == "" {
		return "", nil
	}

	for part := range strings.SplitSeq(tag, " ") {
		if tagKey, tagValue, found := strings.Cut(part, ":"); found && tagKey == key {
			return strconv.Unquote(tagValue)
		}
	}

	return "", nil
}
