package validator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/emontenegro/ClaudeCodeArchitect/internal/compiler"
	"github.com/emontenegro/ClaudeCodeArchitect/internal/parser"
)

// StructuralCheck represents a fast pre-flight check
type StructuralCheck struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

// RunStructuralChecks performs fast pre-flight validation
// These checks don't require Claude - they're instant Go checks
func RunStructuralChecks(manifestPath string) ([]StructuralCheck, error) {
	var checks []StructuralCheck

	// Check 1: Spec compiles
	compileCheck := StructuralCheck{
		ID:   "compiles",
		Name: "Specification compiles",
	}
	_, err := compiler.Compile(manifestPath)
	if err != nil {
		compileCheck.Passed = false
		compileCheck.Message = fmt.Sprintf("Failed to compile: %v", err)
	} else {
		compileCheck.Passed = true
		compileCheck.Message = "OK"
	}
	checks = append(checks, compileCheck)

	// If compile failed, skip other checks
	if !compileCheck.Passed {
		return checks, nil
	}

	// Check 2: Can parse structure
	structureCheck := StructuralCheck{
		ID:   "parseable",
		Name: "Structure parseable",
	}
	structure, err := parser.BuildStructure(manifestPath)
	if err != nil {
		structureCheck.Passed = false
		structureCheck.Message = fmt.Sprintf("Failed to parse structure: %v", err)
		checks = append(checks, structureCheck)
		return checks, nil
	}
	structureCheck.Passed = true
	structureCheck.Message = "OK"
	checks = append(checks, structureCheck)

	// Check 3: Has sections
	sectionsCheck := StructuralCheck{
		ID:   "has-sections",
		Name: "Has defined sections",
	}
	if len(structure.Sections) == 0 {
		sectionsCheck.Passed = false
		sectionsCheck.Message = "No sections found - spec appears empty"
	} else {
		sectionsCheck.Passed = true
		sectionsCheck.Message = fmt.Sprintf("Found %d sections", len(structure.Sections))
	}
	checks = append(checks, sectionsCheck)

	// Check 4: Has attributes (optional but good indicator)
	attrsCheck := StructuralCheck{
		ID:   "has-attributes",
		Name: "Has reusable attributes",
	}
	if len(structure.Attributes) == 0 {
		attrsCheck.Passed = true // Not a failure, just a note
		attrsCheck.Message = "No attributes defined (consider using :attr: for reusable values)"
	} else {
		attrsCheck.Passed = true
		attrsCheck.Message = fmt.Sprintf("Found %d attributes", len(structure.Attributes))
	}
	checks = append(checks, attrsCheck)

	return checks, nil
}

// AllStructuralChecksPassed returns true if all checks passed
func AllStructuralChecksPassed(checks []StructuralCheck) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

// ANSI color codes
const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
	colorRed   = "\033[31m"
)

// FormatStructuralChecks formats checks for display
func FormatStructuralChecks(checks []StructuralCheck) string {
	return formatStructuralChecks(checks, true)
}

// FormatStructuralChecksPlain formats checks without color
func FormatStructuralChecksPlain(checks []StructuralCheck) string {
	return formatStructuralChecks(checks, false)
}

func formatStructuralChecks(checks []StructuralCheck, color bool) string {
	var sb strings.Builder

	sb.WriteString("Structural Checks:\n")
	for _, check := range checks {
		if check.Passed {
			if color {
				sb.WriteString(fmt.Sprintf("  %s✓%s %s: %s\n", colorGreen, colorReset, check.Name, check.Message))
			} else {
				sb.WriteString(fmt.Sprintf("  ✓ %s: %s\n", check.Name, check.Message))
			}
		} else {
			if color {
				sb.WriteString(fmt.Sprintf("  %s✗%s %s: %s\n", colorRed, colorReset, check.Name, check.Message))
			} else {
				sb.WriteString(fmt.Sprintf("  ✗ %s: %s\n", check.Name, check.Message))
			}
		}
	}

	return sb.String()
}

// FormatStructuralChecksJSON formats checks as JSON
func FormatStructuralChecksJSON(checks []StructuralCheck) string {
	data, _ := json.MarshalIndent(checks, "", "  ")
	return string(data)
}

