package compiler

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// Compile compiles the full spec to Markdown
func Compile(specPath string) (string, error) {
	html, err := CompileToHTML(specPath)
	if err != nil {
		return "", err
	}

	return HTMLToMarkdown(html)
}

// CompileToHTML compiles the spec to HTML using asciidoctor CLI
func CompileToHTML(specPath string) (string, error) {
	if !IsAsciidoctorAvailable() {
		return "", fmt.Errorf("asciidoctor not found in PATH\n\nInstall with: gem install asciidoctor\nOr: brew install asciidoctor")
	}

	absPath, err := filepath.Abs(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %v", err)
	}

	// asciidoctor -b html5 -o - file.adoc
	cmd := exec.Command("asciidoctor", "-b", "html5", "-o", "-", absPath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to compile spec: %v\n%s", err, stderr.String())
	}

	return stdout.String(), nil
}

// HTMLToMarkdown converts HTML to Markdown
func HTMLToMarkdown(html string) (string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	if err != nil {
		return "", fmt.Errorf("failed to convert to markdown: %v", err)
	}

	return markdown, nil
}

// CompileContent compiles AsciiDoc content string to Markdown
// This is useful for compiling sections or fragments
func CompileContent(content string, baseDir string) (string, error) {
	if !IsAsciidoctorAvailable() {
		return "", fmt.Errorf("asciidoctor not found in PATH\n\nInstall with: gem install asciidoctor\nOr: brew install asciidoctor")
	}

	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base dir: %v", err)
	}

	// asciidoctor -b html5 -B basedir -o - -
	// -B sets base directory for includes
	// - at end means read from stdin
	cmd := exec.Command("asciidoctor", "-b", "html5", "-B", absBaseDir, "-o", "-", "-")
	cmd.Stdin = strings.NewReader(content)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to compile content: %v\n%s", err, stderr.String())
	}

	return HTMLToMarkdown(stdout.String())
}

// IsAsciidoctorAvailable checks if asciidoctor CLI is installed
func IsAsciidoctorAvailable() bool {
	_, err := exec.LookPath("asciidoctor")
	return err == nil
}
