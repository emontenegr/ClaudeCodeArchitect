package skill

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

// GetEmbeddedVersion extracts version from embedded skill
func GetEmbeddedVersion() string {
	return parseVersion(GetEmbeddedContent())
}

// GetInstalledVersion reads version from installed skill at path
func GetInstalledVersion(skillDir string) (string, error) {
	path := filepath.Join(skillDir, skillDirName, skillFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return parseVersion(string(data)), nil
}

// parseVersion extracts version from YAML frontmatter
func parseVersion(content string) string {
	re := regexp.MustCompile(`(?m)^version:\s*(.+)$`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return "unknown"
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

// CheckUpdate compares installed vs embedded versions
// Returns: needsUpdate, installedVersion, embeddedVersion
func CheckUpdate(skillDir string) (bool, string, string) {
	embeddedVer := GetEmbeddedVersion()
	installedVer, err := GetInstalledVersion(skillDir)
	if err != nil {
		return false, "", embeddedVer
	}
	return installedVer != embeddedVer, installedVer, embeddedVer
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

// PromptYesNo asks user for y/n confirmation
func PromptYesNo(prompt string) bool {
	fmt.Print(prompt + " [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
