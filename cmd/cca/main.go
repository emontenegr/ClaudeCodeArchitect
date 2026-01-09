package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/elijahmont3x/ClaudeCodeArchitect/internal/compiler"
	"github.com/elijahmont3x/ClaudeCodeArchitect/internal/config"
	"github.com/elijahmont3x/ClaudeCodeArchitect/internal/differ"
	"github.com/elijahmont3x/ClaudeCodeArchitect/internal/impact"
	"github.com/elijahmont3x/ClaudeCodeArchitect/internal/skill"
	"github.com/elijahmont3x/ClaudeCodeArchitect/internal/validator"
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
	fmt.Println(`cca - Claude Code Architect CLI

Usage:
  cca compile                      Compile entire spec to Markdown (stdout)
  cca compile --section <name>     Compile specific section only
  cca validate                     Full validation (structural + Claude semantic)
  cca validate --quick             Structural checks only (no Claude)
  cca validate --ultra             Enhanced validation (3x + synthesis)
  cca validate --yes               Skip confirmation for large specs
  cca diff [commit]                Diff compiled output vs commit (default: HEAD~1)
  cca impact <attribute>           Show sections using attribute
  cca list                         List all sections in spec
  cca skill                        Install/update Claude Code skill
  cca skill --global               Install to ~/.claude/skills (all projects)
  cca help                         Show this help

Flags:
  --quick, -q     Structural checks only, skip Claude semantic validation
  --yes, -y       Skip interactive confirmation

Configuration:
  Create .spec.yaml in your project root:
    spec: ./MANIFEST.adoc

  Or use convention - cca looks for:
    - MANIFEST.adoc
    - spec/MANIFEST.adoc
    - plan/MANIFEST.adoc

Examples:
  cca compile                           # Full spec to stdout
  cca compile --section "API Spec"      # Single section with attrs resolved
  cca validate                          # Full validation with Claude
  cca validate --quick                  # Fast structural checks only
  cca validate --yes                    # Skip size confirmation (CI/scripts)
  cca diff HEAD~1                       # Compare with previous commit
  cca impact api-p99-latency            # Find attribute usages
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
		case "--ultra", "-u":
			opts.Ultra = true
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
		return fmt.Errorf("usage: cca impact <attribute-name>")
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
	for _, arg := range os.Args[2:] {
		if arg == "--global" || arg == "-g" {
			global = true
		}
	}

	// Determine target directory
	var skillDir string
	var suffix string
	if global {
		var err error
		skillDir, err = skill.GetGlobalSkillDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		suffix = " (global)"
	} else {
		skillDir = skill.GetProjectSkillDir()
		suffix = ""
	}

	// Check if already installed
	if skill.IsInstalled(skillDir) {
		installed, _ := skill.GetInstalledContent(skillDir)
		if installed == skill.GetEmbeddedContent() {
			fmt.Printf("Skill up to date%s\n", suffix)
			return nil
		}
		if err := skill.Install(skillDir); err != nil {
			return err
		}
		fmt.Printf("Skill updated%s\n", suffix)
		return nil
	}

	// Fresh install
	if err := skill.Install(skillDir); err != nil {
		return err
	}
	fmt.Printf("Skill installed%s\n", suffix)
	return nil
}

// checkSkillUpdate prints a notice if skill update is available
func checkSkillUpdate() {
	if skill.NeedsUpdate(skill.GetProjectSkillDir()) {
		fmt.Fprintf(os.Stderr, "Skill update available — run `cca skill`\n\n")
	}
}
