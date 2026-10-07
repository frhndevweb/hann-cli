package projects

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

// Project represents a discovered development project.
type Project struct {
	Name string
	Path string
	Type string
}

// Scanner handles discovery of development projects within target paths.
type Scanner struct {
	MaxDepth int
}

// NewScanner creates a Scanner configured with default depth limit.
func NewScanner() *Scanner {
	return &Scanner{
		MaxDepth: 2,
	}
}

// ProjectTypeDetector inspects directory contents to deduce the project ecosystem.
func DetectProjectType(dir string) string {
	checks := []struct {
		filename string
		projType string
	}{
		{"go.mod", "Go"},
		{"package.json", "Node.js"},
		{"Cargo.toml", "Rust"},
		{"pyproject.toml", "Python"},
		{"requirements.txt", "Python"},
		{"Pipfile", "Python"},
		{"composer.json", "PHP"},
		{"pubspec.yaml", "Flutter/Dart"},
		{"pom.xml", "Java (Maven)"},
		{"build.gradle", "Gradle"},
		{"build.gradle.kts", "Gradle"},
		{"CMakeLists.txt", "C/C++"},
		{"Makefile", "C/Make"},
		{".git", "Git Repository"},
	}

	for _, c := range checks {
		target := filepath.Join(dir, c.filename)
		if _, err := os.Stat(target); err == nil {
			return c.projType
		}
	}

	return ""
}

// ExpandPath expands leading ~ to user's home directory.
func ExpandPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path cannot be empty")
	}

	if path == "~" || strings.HasPrefix(path, "~/") {
		home := os.Getenv("HOME")
		if home == "" {
			if u, err := user.Current(); err == nil {
				home = u.HomeDir
			}
		}
		if home == "" {
			return "", errors.New("unable to determine user home directory")
		}

		if path == "~" {
			return home, nil
		}
		return filepath.Join(home, path[2:]), nil
	}

	return filepath.Abs(path)
}

// DisplayPath formats an absolute path into ~/relative if inside HOME.
func DisplayPath(absPath string) string {
	home := os.Getenv("HOME")
	if home == "" {
		if u, err := user.Current(); err == nil {
			home = u.HomeDir
		}
	}

	if home != "" && strings.HasPrefix(absPath, home) {
		rel := strings.TrimPrefix(absPath, home)
		if rel == "" {
			return "~"
		}
		return "~" + rel
	}
	return absPath
}

var ignoredDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
	".cache":       true,
	"target":       true,
	"dist":         true,
	"build":        true,
	"__pycache__":  true,
	".idea":        true,
	".vscode":      true,
}

// Scan searches the target directory up to MaxDepth for project roots.
func (s *Scanner) Scan(targetPath string) ([]Project, error) {
	fi, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("unable to access path %s: %w", targetPath, err)
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("target path is not a directory: %s", targetPath)
	}

	var projects []Project

	// If targetPath itself is a project, add it
	if selfType := DetectProjectType(targetPath); selfType != "" {
		projects = append(projects, Project{
			Name: filepath.Base(targetPath),
			Path: targetPath,
			Type: selfType,
		})
		return projects, nil
	}

	var walk func(currentDir string, depth int)
	walk = func(currentDir string, depth int) {
		if depth > s.MaxDepth {
			return
		}

		entries, err := os.ReadDir(currentDir)
		if err != nil {
			return
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			name := entry.Name()
			if strings.HasPrefix(name, ".") && name != ".git" {
				continue
			}
			if ignoredDirectories[name] {
				continue
			}

			fullPath := filepath.Join(currentDir, name)
			projType := DetectProjectType(fullPath)
			if projType != "" {
				projects = append(projects, Project{
					Name: name,
					Path: fullPath,
					Type: projType,
				})
				// Do not recurse into already matched project
				continue
			}

			// If not a project, recurse deeper if within depth limit
			if depth < s.MaxDepth {
				walk(fullPath, depth+1)
			}
		}
	}

	walk(targetPath, 1)

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	return projects, nil
}
