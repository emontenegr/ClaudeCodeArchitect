package compiler

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/bytesparadise/libasciidoc"
	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/sirupsen/logrus"
)

func init() {
	// Suppress verbose logging from libasciidoc
	logrus.SetLevel(logrus.WarnLevel)
}

// Compile compiles the full spec to Markdown
func Compile(specPath string) (string, error) {
	html, err := CompileToHTML(specPath)
	if err != nil {
		return "", err
	}

	return HTMLToMarkdown(html)
}

// CompileToHTML compiles the spec to HTML
func CompileToHTML(specPath string) (string, error) {
	input, err := os.Open(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to open spec: %v", err)
	}
	defer input.Close()

	htmlBuf := &bytes.Buffer{}

	// Set the base directory for include resolution
	baseDir := filepath.Dir(specPath)
	config := configuration.NewConfiguration(
		configuration.WithBackEnd("html5"),
		configuration.WithFilename(specPath),
		configuration.WithAttribute("docdir", baseDir),
	)

	_, err = libasciidoc.Convert(input, htmlBuf, config)
	if err != nil {
		return "", fmt.Errorf("failed to compile spec: %v", err)
	}

	return htmlBuf.String(), nil
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
	htmlBuf := &bytes.Buffer{}

	config := configuration.NewConfiguration(
		configuration.WithBackEnd("html5"),
		configuration.WithAttribute("docdir", baseDir),
	)

	// libasciidoc.Convert requires an io.Reader, so wrap the string
	reader := strings.NewReader(content)
	_, err := libasciidoc.Convert(reader, htmlBuf, config)
	if err != nil {
		return "", fmt.Errorf("failed to compile content: %v", err)
	}

	return HTMLToMarkdown(htmlBuf.String())
}
