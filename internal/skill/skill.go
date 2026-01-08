package skill

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed SKILL.md
var embeddedSkill embed.FS

const skillFileName = "SKILL.md"
const skillDirName = "adoc"

// GetEmbeddedContent returns the embedded skill file content
func GetEmbeddedContent() string {
	data, _ := embeddedSkill.ReadFile(skillFileName)
	return string(data)
}

// GetInstalledContent reads installed skill content
func GetInstalledContent(skillDir string) (string, error) {
	path := filepath.Join(skillDir, skillDirName, skillFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsInstalled checks if skill exists at the given skills directory
func IsInstalled(skillDir string) bool {
	path := filepath.Join(skillDir, skillDirName, skillFileName)
	_, err := os.Stat(path)
	return err == nil
}

// Install writes the skill to the target directory
func Install(skillDir string) error {
	targetDir := filepath.Join(skillDir, skillDirName)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	targetPath := filepath.Join(targetDir, skillFileName)
	content := GetEmbeddedContent()

	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	return nil
}

// NeedsUpdate checks if installed skill differs from embedded
func NeedsUpdate(skillDir string) bool {
	if !IsInstalled(skillDir) {
		return false
	}
	installed, err := GetInstalledContent(skillDir)
	if err != nil {
		return false
	}
	return installed != GetEmbeddedContent()
}

// GetProjectSkillDir returns .claude/skills in current directory
func GetProjectSkillDir() string {
	return filepath.Join(".claude", "skills")
}

// GetGlobalSkillDir returns ~/.claude/skills
func GetGlobalSkillDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude", "skills"), nil
}
