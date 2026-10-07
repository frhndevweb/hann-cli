package cmd

import (
	"bytes"
	"testing"
)

func TestRootCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("rootCmd --help failed: %v", err)
	}

	out := buf.String()
	if out == "" {
		t.Error("Expected non-empty output for rootCmd --help")
	}
}

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}
}

func TestDoctorCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"doctor"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("doctor command failed: %v", err)
	}
}

func TestSystemCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"system"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("system command failed: %v", err)
	}
}

func TestPortsCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"ports"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("ports command failed: %v", err)
	}
}

func TestProjectsCommandMissingDir(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"projects", "--path", "/non/existent/path/xyz_test_123"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("projects command failed: %v", err)
	}
}
