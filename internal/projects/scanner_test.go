package projects

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectType(t *testing.T) {
	tempDir := t.TempDir()

	goDir := filepath.Join(tempDir, "go-app")
	if err := os.Mkdir(goDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(goDir, "go.mod"), []byte("module myapp"), 0644); err != nil {
		t.Fatal(err)
	}
	if typ := DetectProjectType(goDir); typ != "Go" {
		t.Errorf("Expected Go, got %s", typ)
	}

	nodeDir := filepath.Join(tempDir, "node-app")
	if err := os.Mkdir(nodeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(nodeDir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if typ := DetectProjectType(nodeDir); typ != "Node.js" {
		t.Errorf("Expected Node.js, got %s", typ)
	}

	rustDir := filepath.Join(tempDir, "rust-app")
	if err := os.Mkdir(rustDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rustDir, "Cargo.toml"), []byte("[package]"), 0644); err != nil {
		t.Fatal(err)
	}
	if typ := DetectProjectType(rustDir); typ != "Rust" {
		t.Errorf("Expected Rust, got %s", typ)
	}
}

func TestScannerScan(t *testing.T) {
	root := t.TempDir()

	// Create 2 mock projects
	proj1 := filepath.Join(root, "api-sentinel")
	os.Mkdir(proj1, 0755)
	os.WriteFile(filepath.Join(proj1, "go.mod"), []byte("module api-sentinel"), 0644)

	proj2 := filepath.Join(root, "website")
	os.Mkdir(proj2, 0755)
	os.WriteFile(filepath.Join(proj2, "package.json"), []byte("{}"), 0644)

	// Create ignored directory with mock file
	ignored := filepath.Join(root, "node_modules")
	os.Mkdir(ignored, 0755)
	os.WriteFile(filepath.Join(ignored, "package.json"), []byte("{}"), 0644)

	scanner := NewScanner()
	results, err := scanner.Scan(root)
	if err != nil {
		t.Fatalf("Unexpected scan error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 projects found, got %d", len(results))
	}

	if results[0].Name != "api-sentinel" || results[0].Type != "Go" {
		t.Errorf("Unexpected project 0: %+v", results[0])
	}
	if results[1].Name != "website" || results[1].Type != "Node.js" {
		t.Errorf("Unexpected project 1: %+v", results[1])
	}
}

func TestScanNonExistentPath(t *testing.T) {
	scanner := NewScanner()
	_, err := scanner.Scan("/non/existent/path/xyz_998877")
	if err == nil {
		t.Fatal("Expected error for non-existent path")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected os.ErrNotExist, got %v", err)
	}
}

func TestExpandPath(t *testing.T) {
	t.Setenv("HOME", "/root")
	expanded, err := ExpandPath("~/projects")
	if err != nil {
		t.Fatalf("Unexpected error expanding path: %v", err)
	}
	if expanded != "/root/projects" {
		t.Errorf("Expected /root/projects, got %s", expanded)
	}
}
