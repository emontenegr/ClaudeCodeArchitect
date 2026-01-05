package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/bytesparadise/libasciidoc"
	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"gopkg.in/yaml.v3"
)

type SpecConfig struct {
	Spec string `yaml:"spec"`
}

// findSpec discovers the specification file location
func findSpec() (string, error) {
	// 1. Check for .spec.yaml config
	if config, err := loadSpecConfig(); err == nil {
		path := config.Spec
		if !filepath.IsAbs(path) {
			path, _ = filepath.Abs(path)
		}
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("spec file not found at configured path: %s", path)
	}

	// 2. Check conventions
	conventions := []string{
		"MANIFEST.adoc",
		"spec/MANIFEST.adoc",
		"plan/MANIFEST.adoc",
	}

	for _, path := range conventions {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath, nil
		}
	}

	return "", fmt.Errorf("spec not found - checked: %v\nCreate .spec.yaml or use MANIFEST.adoc", conventions)
}

// loadSpecConfig loads .spec.yaml configuration
func loadSpecConfig() (*SpecConfig, error) {
	data, err := os.ReadFile(".spec.yaml")
	if err != nil {
		return nil, err
	}

	var config SpecConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// compileAsciiDoc compiles the spec to Markdown
func compileAsciiDoc(specPath string) (string, error) {
	input, err := os.Open(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to open spec: %v", err)
	}
	defer input.Close()

	// Compile to HTML first
	htmlBuf := &bytes.Buffer{}
	config := configuration.NewConfiguration(
		configuration.WithBackEnd("html5"),
	)

	_, err = libasciidoc.Convert(input, htmlBuf, config)
	if err != nil {
		return "", fmt.Errorf("failed to compile spec: %v", err)
	}

	// Convert HTML to Markdown
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(htmlBuf.String())
	if err != nil {
		return "", fmt.Errorf("failed to convert to markdown: %v", err)
	}

	return markdown, nil
}

// runValidation runs validation checks on the spec
func runValidation(specPath string) []string {
	var errors []string

	// Check if spec compiles
	if _, err := compileAsciiDoc(specPath); err != nil {
		errors = append(errors, fmt.Sprintf("Spec does not compile: %v", err))
		return errors
	}

	// TODO: Add content validation checks here
	// - Check for required sections
	// - Check for conditionals (if/should/maybe)
	// - Check for unversioned dependencies
	// - Check for incomplete type definitions

	return errors
}
