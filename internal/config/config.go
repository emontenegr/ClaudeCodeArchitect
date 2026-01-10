package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// SpecConfig represents the .spec.yaml configuration file
type SpecConfig struct {
	Spec string `yaml:"spec"`
}

// FindSpec discovers the specification file location in the current directory
func FindSpec() (string, error) {
	return FindSpecInDir(".")
}

// FindSpecInDir discovers the specification file location in the given directory
// It checks .spec.yaml first, then falls back to conventions
func FindSpecInDir(dir string) (string, error) {
	// 1. Check for .spec.yaml config
	configPath := filepath.Join(dir, ".spec.yaml")
	if data, err := os.ReadFile(configPath); err == nil {
		var config SpecConfig
		if err := yaml.Unmarshal(data, &config); err == nil {
			path := config.Spec
			if !filepath.IsAbs(path) {
				path = filepath.Join(dir, path)
			}
			if _, err := os.Stat(path); err == nil {
				absPath, _ := filepath.Abs(path)
				return absPath, nil
			}
			return "", fmt.Errorf("spec file not found at configured path: %s", path)
		}
	}

	// 2. Check conventions
	conventions := []string{
		"MANIFEST.adoc",
		"spec/MANIFEST.adoc",
		"plan/MANIFEST.adoc",
	}

	for _, name := range conventions {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			return absPath, nil
		}
	}

	return "", fmt.Errorf("spec not found in %s - checked: %v\nCreate .spec.yaml or use MANIFEST.adoc", dir, conventions)
}

// LoadSpecConfig loads .spec.yaml configuration from the current directory
func LoadSpecConfig() (*SpecConfig, error) {
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

// GetSpecRoot returns the directory containing the MANIFEST.adoc file
func GetSpecRoot(manifestPath string) string {
	return filepath.Dir(manifestPath)
}
