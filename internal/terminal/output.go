package terminal

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ANSI color codes
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiWhite   = "\033[37m"
	ansiGray    = "\033[90m"
)

// UI manages clean, adaptive terminal styling.
type UI struct {
	colorEnabled   bool
	unicodeEnabled bool
	writer         io.Writer
}

// DefaultUI is the default terminal UI instance.
var DefaultUI = NewUI(os.Stdout)

// NewUI initializes a UI with terminal capability detection.
func NewUI(w io.Writer) *UI {
	return &UI{
		colorEnabled:   detectColorSupport(),
		unicodeEnabled: detectUnicodeSupport(),
		writer:         w,
	}
}

// SetColorExplicit allows overriding color detection for testing or flags.
func (u *UI) SetColorExplicit(enabled bool) {
	u.colorEnabled = enabled
}

// SetUnicodeExplicit allows overriding unicode detection for testing or flags.
func (u *UI) SetUnicodeExplicit(enabled bool) {
	u.unicodeEnabled = enabled
}

// detectColorSupport determines if terminal colors should be used.
func detectColorSupport() bool {
	// Respect NO_COLOR standard (https://no-color.org)
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}
	// Check if stdout is a character device
	fi, err := os.Stdout.Stat()
	if err == nil {
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			return false
		}
	}
	return true
}

// detectUnicodeSupport checks if the environment safely handles UTF-8 glyphs.
func detectUnicodeSupport() bool {
	lang := strings.ToUpper(os.Getenv("LANG") + os.Getenv("LC_ALL") + os.Getenv("LC_CTYPE"))
	if strings.Contains(lang, "UTF-8") || strings.Contains(lang, "UTF8") {
		return true
	}
	term := os.Getenv("TERM")
	if term == "linux" {
		return false
	}
	// On modern systems, default to true unless standard indicates ASCII/C
	if lang == "C" || lang == "POSIX" {
		return false
	}
	return true
}

// Colorize wraps text in ANSI sequence if colors are enabled.
func (u *UI) Colorize(code, text string) string {
	if !u.colorEnabled {
		return text
	}
	return code + text + ansiReset
}

// Bold returns bolded text.
func (u *UI) Bold(text string) string {
	return u.Colorize(ansiBold, text)
}

// Dim returns dimmed text.
func (u *UI) Dim(text string) string {
	return u.Colorize(ansiDim, text)
}

// Green returns green text.
func (u *UI) Green(text string) string {
	return u.Colorize(ansiGreen, text)
}

// Red returns red text.
func (u *UI) Red(text string) string {
	return u.Colorize(ansiRed, text)
}

// Yellow returns yellow text.
func (u *UI) Yellow(text string) string {
	return u.Colorize(ansiYellow, text)
}

// Cyan returns cyan text.
func (u *UI) Cyan(text string) string {
	return u.Colorize(ansiCyan, text)
}

// Gray returns gray text.
func (u *UI) Gray(text string) string {
	return u.Colorize(ansiGray, text)
}

// SymbolSuccess returns checkmark symbol with green color.
func (u *UI) SymbolSuccess() string {
	sym := "✓"
	if !u.unicodeEnabled {
		sym = "[+]"
	}
	return u.Green(sym)
}

// SymbolFailure returns cross symbol with red color.
func (u *UI) SymbolFailure() string {
	sym := "✗"
	if !u.unicodeEnabled {
		sym = "[-]"
	}
	return u.Red(sym)
}

// SymbolWarning returns warning symbol with yellow color.
func (u *UI) SymbolWarning() string {
	sym := "⚠"
	if !u.unicodeEnabled {
		sym = "[!]"
	}
	return u.Yellow(sym)
}

// SectionHeader formats a clean section title.
func (u *UI) SectionHeader(title string) string {
	return u.Bold(title)
}

// PrintHeader prints a top-level command header.
func (u *UI) PrintHeader(title string) {
	fmt.Fprintln(u.writer, u.Bold(title))
	fmt.Fprintln(u.writer)
}

// PrintSection prints a section header with spacing.
func (u *UI) PrintSection(title string) {
	fmt.Fprintln(u.writer, u.Bold(title))
}

// PrintCheck prints a single check item with status icon.
func (u *UI) PrintCheck(passed bool, label, value string) {
	icon := u.SymbolSuccess()
	if !passed {
		icon = u.SymbolFailure()
	}
	if value != "" {
		fmt.Fprintf(u.writer, "%s %-14s %s\n", icon, label, value)
	} else {
		fmt.Fprintf(u.writer, "%s %s\n", icon, label)
	}
}

// PrintWarning prints a warning item.
func (u *UI) PrintWarning(label, detail string) {
	icon := u.SymbolWarning()
	if detail != "" {
		fmt.Fprintf(u.writer, "%s %-14s %s\n", icon, label, u.Yellow(detail))
	} else {
		fmt.Fprintf(u.writer, "%s %s\n", icon, label)
	}
}

// FormatTable formats a tabular dataset with clean columnar padding.
func (u *UI) FormatTable(headers []string, rows [][]string) string {
	if len(headers) == 0 {
		return ""
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder

	// Header row
	for i, h := range headers {
		padding := widths[i] - len(h) + 4
		if i == len(headers)-1 {
			padding = 0
		}
		sb.WriteString(u.Bold(h) + strings.Repeat(" ", padding))
	}
	sb.WriteString("\n")

	// Data rows
	for _, row := range rows {
		for i, cell := range row {
			padding := widths[i] - len(cell) + 4
			if i == len(row)-1 {
				padding = 0
			}
			sb.WriteString(cell + strings.Repeat(" ", padding))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
