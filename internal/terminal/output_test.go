package terminal

import (
	"bytes"
	"strings"
	"testing"
)

func TestUIColorsDisabled(t *testing.T) {
	var buf bytes.Buffer
	ui := NewUI(&buf)
	ui.SetColorExplicit(false)

	text := ui.Green("Success")
	if text != "Success" {
		t.Errorf("Expected raw text without ANSI codes, got %q", text)
	}

	ui.SetColorExplicit(true)
	colored := ui.Green("Success")
	if !strings.Contains(colored, "\033[32m") {
		t.Errorf("Expected green ANSI code, got %q", colored)
	}
}

func TestUISymbolsFallback(t *testing.T) {
	var buf bytes.Buffer
	ui := NewUI(&buf)
	ui.SetColorExplicit(false)

	ui.SetUnicodeExplicit(true)
	if ui.SymbolSuccess() != "✓" {
		t.Errorf("Expected unicode checkmark, got %s", ui.SymbolSuccess())
	}
	if ui.SymbolFailure() != "✗" {
		t.Errorf("Expected unicode cross, got %s", ui.SymbolFailure())
	}

	ui.SetUnicodeExplicit(false)
	if ui.SymbolSuccess() != "[+]" {
		t.Errorf("Expected ascii fallback [+], got %s", ui.SymbolSuccess())
	}
	if ui.SymbolFailure() != "[-]" {
		t.Errorf("Expected ascii fallback [-], got %s", ui.SymbolFailure())
	}
}

func TestFormatTable(t *testing.T) {
	var buf bytes.Buffer
	ui := NewUI(&buf)
	ui.SetColorExplicit(false)

	headers := []string{"PORT", "PROTOCOL", "ADDRESS"}
	rows := [][]string{
		{"3000", "TCP", "0.0.0.0"},
		{"8080", "TCP", "127.0.0.1"},
	}

	table := ui.FormatTable(headers, rows)
	if !strings.Contains(table, "PORT") || !strings.Contains(table, "3000") || !strings.Contains(table, "8080") {
		t.Errorf("Table output missing expected content:\n%s", table)
	}
}
