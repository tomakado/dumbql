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
		os.Exit(2)
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
		fmt.Println("dumbqlgen version 1.0.0")
		return
	}

	if config.StructType == "" {
		handleError(errors.New("missing required flag: -type"))
	}

	// If output is not specified, use the current directory with structtype_matcher.gen.go
	if config.Output == "" {
		config.Output = fmt.Sprintf("%s_matcher.gen.go", strings.ToLower(config.StructType))
	}

	code, err := generate(config)
	if err != nil {
		handleError(fmt.Errorf("generate router: %w", err))
	}

	err = os.WriteFile(config.Output, []byte(code), 0644)
	if err != nil {
		handleError(fmt.Errorf("write output file: %w", err))
	}
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

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "router.tmpl", map[string]any{
		"Package":    packageName,
		"StructInfo": structInfo,
		"Version":    "1.0.0",
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
	for i := 0; i < structType.NumFields(); i++ {
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

		// Find the AST node for the struct to get the tags
		for _, file := range pkg.Syntax {
			ast.Inspect(file, func(n ast.Node) bool {
				typeSpec, ok := n.(*ast.TypeSpec)
				if !ok || typeSpec.Name.Name != structName {
					return true
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					return false
				}

				// Loop through fields to find the current field
				for _, field := range structType.Fields.List {
					for _, ident := range field.Names {
						if ident.Name == fieldName {
							if field.Tag != nil {
								tag := strings.Trim(field.Tag.Value, "`")
								dumbqlTag := extractTag(tag, "dumbql")
								if dumbqlTag == "-" {
									fieldInfo.Skip = true
								} else if dumbqlTag != "" {
									fieldInfo.TagName = dumbqlTag
								}
							}
							break
						}
					}
				}

				return false
			})
		}

		structInfo.Fields = append(structInfo.Fields, fieldInfo)
	}

	return structInfo, nil
}

func extractTag(tag, key string) string {
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		i = 0
		for i < len(tag) && tag[i] != ' ' && tag[i] != ':' && tag[i] > 0x20 {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		if key == name {
			value, _ := strconv.Unquote(qvalue)
			return value
		}
	}
	return ""
}
