package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var (
		typeName   string
		outputFile string
		pkgName    string
		dir        string
	)

	flag.StringVar(&typeName, "type", "", "Type name to generate matcher for (required)")
	flag.StringVar(&outputFile, "output", "", "Output file name (default: {type}_matcher.gen.go)")
	flag.StringVar(&pkgName, "pkg", "", "Package name for the generated code (default: current directory name)")
	flag.StringVar(&dir, "dir", ".", "Directory to search for the type definition")
	flag.Parse()

	if typeName == "" {
		fmt.Fprintln(os.Stderr, "Error: -type flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if outputFile == "" {
		outputFile = fmt.Sprintf("%s_matcher.gen.go", typeName)
	}

	// Resolve absolute path for directory
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving directory path: %v\n", err)
		os.Exit(1)
	}

	if pkgName == "" {
		pkgName = filepath.Base(absDir)
	}

	fmt.Printf("Searching for type %s in directory %s\n", typeName, absDir)
	structInfo, err := findStructInDir(absDir, typeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding struct: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found struct %s with %d fields\n", structInfo.Name, len(structInfo.Fields))
	fmt.Printf("Generating matcher in package %s\n", pkgName)
	fmt.Printf("Output will be written to %s\n", outputFile)

	if err := generateMatcher(structInfo, outputFile, pkgName); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating matcher: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated matcher for %s\n", typeName)
}