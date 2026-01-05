package validator

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"
)

//go:embed prompts/*.tmpl
var promptTemplates embed.FS

// Size thresholds for warnings
const (
	SizeWarningThreshold = 20 * 1024  // 20KB - warn user
	SizeLargeThreshold   = 50 * 1024  // 50KB - strongly warn
)

// ValidationOptions controls validation behavior
type ValidationOptions struct {
	SkipConfirm bool // --yes flag: skip size confirmation
}

// TemplateData holds data passed to prompt templates
type TemplateData struct {
	CompiledSpec string
}

// LoadPromptTemplate loads and parses a prompt template
func LoadPromptTemplate(name string) (*template.Template, error) {
	return template.ParseFS(promptTemplates, "prompts/"+name+".tmpl")
}

// RenderPrompt renders a prompt template with data
func RenderPrompt(templateName string, data TemplateData) (string, error) {
	tmpl, err := LoadPromptTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, templateName, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

// RunClaudeValidation shells out to claude CLI for semantic validation
// It streams output directly to the provided writer
func RunClaudeValidation(compiledSpec string, output io.Writer) error {
	// Render the prompt
	prompt, err := RenderPrompt("validate", TemplateData{
		CompiledSpec: compiledSpec,
	})
	if err != nil {
		return fmt.Errorf("failed to render prompt: %w", err)
	}

	// Check if claude CLI is available
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude CLI not found in PATH - install from https://claude.ai/code")
	}

	// Run claude with prompt via stdin (avoids command line length limits)
	// Using --print for non-interactive mode
	fmt.Fprint(output, "Running Claude validation ")

	cmd := exec.Command("claude", "--print")
	cmd.Stdin = strings.NewReader(prompt)

	// Capture stdout to buffer while showing spinner
	var resultBuf bytes.Buffer
	cmd.Stdout = &resultBuf
	cmd.Stderr = os.Stderr

	// Start spinner in goroutine
	done := make(chan bool)
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Fprintf(output, "\rRunning Claude validation %s", spinner[i%len(spinner)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	err = cmd.Run()
	done <- true

	// Clear spinner line and show result
	fmt.Fprint(output, "\r                                    \r")

	if err != nil {
		return fmt.Errorf("claude CLI failed: %w", err)
	}

	// Write the captured output
	fmt.Fprint(output, resultBuf.String())

	return nil
}

// RunClaudeValidationToString runs validation and returns result as string
func RunClaudeValidationToString(compiledSpec string) (string, error) {
	var buf bytes.Buffer
	if err := RunClaudeValidation(compiledSpec, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// IsClaudeAvailable checks if claude CLI is installed
func IsClaudeAvailable() bool {
	_, err := exec.LookPath("claude")
	return err == nil
}

// CheckSpecSize checks spec size and prompts for confirmation if large
// Returns true if should proceed, false if user cancelled
func CheckSpecSize(compiledSpec string, opts ValidationOptions, output io.Writer) (bool, error) {
	size := len(compiledSpec)

	if size < SizeWarningThreshold {
		return true, nil // Small spec, no warning needed
	}

	// Calculate approximate tokens (rough: 1 token ≈ 4 chars)
	approxTokens := size / 4

	if size >= SizeLargeThreshold {
		fmt.Fprintf(output, "⚠ Large spec detected: %dKB (~%d tokens)\n", size/1024, approxTokens)
		fmt.Fprintf(output, "  This will use significant Claude API capacity.\n\n")
	} else {
		fmt.Fprintf(output, "Note: Spec size %dKB (~%d tokens)\n\n", size/1024, approxTokens)
	}

	if opts.SkipConfirm {
		return true, nil
	}

	// Interactive confirmation
	fmt.Fprint(output, "Proceed with Claude validation? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

