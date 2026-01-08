package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elijahmont3x/ClaudeCodeArchitect/cli/internal/compiler"
	"github.com/elijahmont3x/ClaudeCodeArchitect/cli/internal/config"
	"github.com/elijahmont3x/ClaudeCodeArchitect/cli/internal/differ"
	"github.com/elijahmont3x/ClaudeCodeArchitect/cli/internal/impact"
	"github.com/elijahmont3x/ClaudeCodeArchitect/cli/internal/skill"
	"github.com/elijahmont3x/ClaudeCodeArchitect/cli/internal/validator"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	var err error
	switch command {
	case "compile":
		err = runCompile()
	case "validate":
		err = runValidate()
	case "diff":
		err = runDiff()
	case "impact":
		err = runImpact()
	case "list":
		err = runList()
	case "skill":
		err = runSkill()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`spec-cli - Compile architecture specifications for AI consumption

Usage:
  spec-cli compile                      Compile entire spec to Markdown (stdout)
  spec-cli compile --section <name>     Compile specific section only
  spec-cli validate                     Full validation (structural + Claude semantic)
  spec-cli validate --quick             Structural checks only (no Claude)
  spec-cli validate --yes               Skip confirmation for large specs
  spec-cli diff [commit]                Diff compiled output vs commit (default: HEAD~1)
  spec-cli impact <attribute>           Show sections using attribute
  spec-cli list                         List all sections in spec
  spec-cli skill                        Install/update Claude Code skill (project)
  spec-cli skill --global               Install/update skill for all projects
  spec-cli help                         Show this help

Flags:
  --quick, -q     Structural checks only, skip Claude semantic validation
  --yes, -y       Skip interactive confirmation

Configuration:
  Create .spec.yaml in your project root:
    spec: ./MANIFEST.adoc

  Or use convention - spec-cli looks for:
    - MANIFEST.adoc
    - spec/MANIFEST.adoc
    - plan/MANIFEST.adoc

Examples:
  spec-cli compile                           # Full spec to stdout
  spec-cli compile --section "API Spec"      # Single section with attrs resolved
  spec-cli validate                          # Full validation with Claude
  spec-cli validate --quick                  # Fast structural checks only
  spec-cli validate --yes                    # Skip size confirmation (CI/scripts)
  spec-cli diff HEAD~1                       # Compare with previous commit
  spec-cli impact api-p99-latency            # Find attribute usages
`)
}

func runCompile() error {
	checkSkillUpdate()

	specPath, err := config.FindSpec()
	if err != nil {
		return err
	}

	// Check for --section flag
	sectionQuery := ""
	for i, arg := range os.Args {
		if arg == "--section" && i+1 < len(os.Args) {
			sectionQuery = os.Args[i+1]
			break
		}
		if strings.HasPrefix(arg, "--section=") {
			sectionQuery = strings.TrimPrefix(arg, "--section=")
			break
		}
	}

	var output string
	if sectionQuery != "" {
		output, err = compiler.CompileSection(specPath, sectionQuery)
	} else {
		output, err = compiler.Compile(specPath)
	}

	if err != nil {
		return err
	}

	fmt.Print(output)
	return nil
}

func runValidate() error {
	checkSkillUpdate()

	specPath, err := config.FindSpec()
	if err != nil {
		return err
	}

	// Parse flags
	quick := false
	opts := validator.ValidationOptions{}

	for _, arg := range os.Args {
		switch arg {
		case "--quick", "-q":
			quick = true
		case "--yes", "-y":
			opts.SkipConfirm = true
		}
	}

	if quick {
		result, err := validator.ValidateQuick(specPath)
		if err != nil {
			return err
		}
		fmt.Print(validator.FormatStructuralChecks(result.StructuralChecks))
		if !result.StructuralPassed {
			os.Exit(1)
		}
		return nil
	}

	// Full validation: structural + Claude
	result, err := validator.Validate(specPath, os.Stdout, opts)
	if err != nil {
		return err
	}

	if !result.StructuralPassed || result.Cancelled {
		os.Exit(1)
	}

	return nil
}

func runDiff() error {
	specPath, err := config.FindSpec()
	if err != nil {
		return err
	}

	// Get target commit (default: HEAD~1)
	targetCommit := "HEAD~1"
	if len(os.Args) > 2 {
		targetCommit = os.Args[2]
	}

	result, err := differ.DiffCompiled(specPath, targetCommit)
	if err != nil {
		return err
	}

	fmt.Println(differ.FormatDiffResult(result))
	return nil
}

func runImpact() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: spec-cli impact <attribute-name>")
	}

	attrName := os.Args[2]

	specPath, err := config.FindSpec()
	if err != nil {
		return err
	}

	result, err := impact.AnalyzeAttribute(specPath, attrName)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(specPath)
	fmt.Println(impact.FormatImpact(result, baseDir))
	return nil
}

func runList() error {
	specPath, err := config.FindSpec()
	if err != nil {
		return err
	}

	sections, err := compiler.ListSections(specPath)
	if err != nil {
		return err
	}

	fmt.Println("Sections in specification:\n")
	fmt.Print(compiler.FormatSectionList(sections))
	return nil
}

func runSkill() error {
	// Parse flags
	global := false
	yes := false
	for _, arg := range os.Args[2:] {
		switch arg {
		case "--global", "-g":
			global = true
		case "--yes", "-y":
			yes = true
		}
	}

	// Determine target directory
	var skillDir string
	var location string
	if global {
		var err error
		skillDir, err = skill.GetGlobalSkillDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		location = "global (~/.claude/skills)"
	} else {
		skillDir = skill.GetProjectSkillDir()
		location = "project (.claude/skills)"
	}

	embeddedVer := skill.GetEmbeddedVersion()

	// Check if already installed
	if skill.IsInstalled(skillDir) {
		installedVer, _ := skill.GetInstalledVersion(skillDir)

		if installedVer == embeddedVer {
			fmt.Printf("Skill already up to date (%s) at %s\n", embeddedVer, location)
			return nil
		}

		// Version differs - prompt for update
		if !yes {
			prompt := fmt.Sprintf("Update skill %s → %s at %s?", installedVer, embeddedVer, location)
			if !skill.PromptYesNo(prompt) {
				fmt.Println("Cancelled")
				return nil
			}
		}

		if err := skill.Install(skillDir); err != nil {
			return err
		}
		fmt.Printf("Updated skill %s → %s at %s\n", installedVer, embeddedVer, location)
		return nil
	}

	// Fresh install
	if err := skill.Install(skillDir); err != nil {
		return err
	}
	fmt.Printf("Installed skill %s at %s\n", embeddedVer, location)
	return nil
}

// checkSkillUpdate prints a notice if skill update is available
func checkSkillUpdate() {
	skillDir := skill.GetProjectSkillDir()
	if !skill.IsInstalled(skillDir) {
		return
	}

	needsUpdate, installed, embedded := skill.CheckUpdate(skillDir)
	if needsUpdate {
		fmt.Fprintf(os.Stderr, "Note: Skill update available (%s → %s). Run `spec-cli skill` to update.\n\n", installed, embedded)
	}
}
