package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Field represents a struct field with its type and tag information
type Field struct {
	Name     string
	Type     string
	Tag      string
	DumbQLTag string
	IsExported bool
}

// StructInfo represents the information about a struct
type StructInfo struct {
	Name      string
	PkgName   string
	Fields    []Field
	ImportPaths []string
}

// parseFile parses a Go source file and extracts struct declarations
func parseFile(filename string, targetType string) (*StructInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing file %s: %w", filename, err)
	}

	var structInfo *StructInfo
	importPaths := make([]string, 0)

	// Extract import paths
	for _, imp := range node.Imports {
		if imp.Path != nil {
			path := strings.Trim(imp.Path.Value, `"`)
			importPaths = append(importPaths, path)
		}
	}

	// Find the target struct
	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != targetType {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		structInfo = &StructInfo{
			Name:      targetType,
			PkgName:   node.Name.Name,
			Fields:    make([]Field, 0),
			ImportPaths: importPaths,
		}

		// Extract fields
		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				// Skip embedded types for now
				continue
			}

			fieldName := field.Names[0].Name
			var tagValue string
			if field.Tag != nil {
				tagValue = field.Tag.Value
			}

			// Parse dumbql tag
			dumbqlTag := extractDumbQLTag(tagValue)

			// Get type as string
			typeStr := typeToString(field.Type)

			structInfo.Fields = append(structInfo.Fields, Field{
				Name:     fieldName,
				Type:     typeStr,
				Tag:      tagValue,
				DumbQLTag: dumbqlTag,
				IsExported: ast.IsExported(fieldName),
			})
		}

		return false
	})

	if structInfo == nil {
		return nil, fmt.Errorf("struct %s not found in file %s", targetType, filename)
	}

	return structInfo, nil
}

// extractDumbQLTag extracts the dumbql tag value from a struct tag
func extractDumbQLTag(tagValue string) string {
	if tagValue == "" {
		return ""
	}

	// Remove backticks
	tagValue = strings.Trim(tagValue, "`")

	// Find dumbql tag
	parts := strings.Split(tagValue, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "dumbql:") {
			return strings.Trim(strings.TrimPrefix(part, "dumbql:"), "\"")
		}
	}

	return ""
}

// typeToString converts an AST type to a string representation
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	default:
		return fmt.Sprintf("unsupported type: %T", expr)
	}
}

// findGoFilesInDir finds all Go source files in a directory
func findGoFilesInDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var goFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") && !strings.HasSuffix(file.Name(), "_test.go") {
			goFiles = append(goFiles, filepath.Join(dir, file.Name()))
		}
	}

	return goFiles, nil
}

// findStructInDir looks for a struct with the given name in all Go files in a directory
func findStructInDir(dir, typeName string) (*StructInfo, error) {
	files, err := findGoFilesInDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		structInfo, err := parseFile(file, typeName)
		if err == nil && structInfo != nil {
			return structInfo, nil
		}
	}

	return nil, fmt.Errorf("struct %s not found in directory %s", typeName, dir)
}

// generateMatcher generates a matcher for the given struct
func generateMatcher(structInfo *StructInfo, outputFile, packageName string) error {
	tmpl, err := template.New("matcher").Parse(matcherTemplate)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	data := struct {
		StructInfo  *StructInfo
		PackageName string
	}{
		StructInfo:  structInfo,
		PackageName: packageName,
	}

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	return nil
}

// See template.go for the matcher template