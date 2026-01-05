package validator

import (
	"fmt"
	"strings"

	"github.com/ClaudeCodeArchitect/spec-cli/internal/compiler"
	"github.com/ClaudeCodeArchitect/spec-cli/internal/parser"
)

// StructuralCheck represents a fast pre-flight check
type StructuralCheck struct {
	ID      string
	Name    string
	Passed  bool
	Message string
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

	// Check 5: Has key sections (heuristic)
	keySectionsCheck := checkKeySections(structure)
	checks = append(checks, keySectionsCheck)

	return checks, nil
}

// checkKeySections looks for commonly expected sections
func checkKeySections(structure *parser.SpecStructure) StructuralCheck {
	check := StructuralCheck{
		ID:   "key-sections",
		Name: "Key sections present",
	}

	// Common section keywords to look for
	keyPatterns := []string{
		"type", "api", "endpoint", "route",
		"deploy", "test", "performance",
	}

	foundPatterns := []string{}
	for _, section := range structure.Sections {
		lower := strings.ToLower(section.Title)
		for _, pattern := range keyPatterns {
			if strings.Contains(lower, pattern) {
				foundPatterns = append(foundPatterns, pattern)
				break
			}
		}
	}

	if len(foundPatterns) >= 3 {
		check.Passed = true
		check.Message = fmt.Sprintf("Found key sections: %v", unique(foundPatterns))
	} else if len(foundPatterns) > 0 {
		check.Passed = true
		check.Message = fmt.Sprintf("Found some key sections: %v (consider adding more)", unique(foundPatterns))
	} else {
		check.Passed = false
		check.Message = "No recognizable key sections (types, API, deployment, etc.)"
	}

	return check
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

// FormatStructuralChecks formats checks for display
func FormatStructuralChecks(checks []StructuralCheck) string {
	var sb strings.Builder

	sb.WriteString("Structural Checks:\n")
	for _, check := range checks {
		if check.Passed {
			sb.WriteString(fmt.Sprintf("  ✓ %s: %s\n", check.Name, check.Message))
		} else {
			sb.WriteString(fmt.Sprintf("  ✗ %s: %s\n", check.Name, check.Message))
		}
	}

	return sb.String()
}

func unique(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
