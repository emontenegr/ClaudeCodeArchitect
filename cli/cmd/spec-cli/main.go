package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "compile":
		if err := compileSpec(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "validate":
		if err := validateSpec(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`spec-cli - Compile architecture specifications for AI consumption

Usage:
  spec-cli compile            Compile entire spec to Markdown (stdout)
  spec-cli validate           Validate spec compiles without errors
  spec-cli help               Show this help

Configuration:
  Create .spec.yaml in your project root:
    spec: ./MANIFEST.adoc

  Or use convention - spec-cli looks for:
    - MANIFEST.adoc
    - spec/MANIFEST.adoc
    - plan/MANIFEST.adoc

Examples:
  spec-cli compile > compiled-spec.md      # Save to file
  spec-cli compile                         # Output to stdout for AI
  spec-cli validate                        # Check spec is valid
`)
}

func compileSpec() error {
	spec, err := findSpec()
	if err != nil {
		return err
	}

	output, err := compileAsciiDoc(spec)
	if err != nil {
		return err
	}

	fmt.Print(output)
	return nil
}

func validateSpec() error {
	spec, err := findSpec()
	if err != nil {
		return err
	}

	errors := runValidation(spec)
	if len(errors) == 0 {
		fmt.Println("✓ Specification is valid")
		return nil
	}

	fmt.Println("✗ Specification validation failed:")
	for _, e := range errors {
		fmt.Printf("  - %s\n", e)
	}
	os.Exit(1)
	return nil
}
