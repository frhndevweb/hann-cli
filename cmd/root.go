package cmd

import (
	"fmt"
	"os"

	"github.com/frhndevweb/hann-cli/internal/terminal"
	"github.com/frhndevweb/hann-cli/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hann",
	Short: "Developer and system utility CLI for Linux (Termux & PRoot optimized)",
	Long: `hann-cli is a developer and system utility CLI designed by Farhan ("hann")
for Linux environments, with deep optimizations for Ubuntu running on
Termux + proot-distro. Local-first, fast, and dependency-light.`,
	Version: version.Version,
	Run: func(cmd *cobra.Command, args []string) {
		ui := terminal.DefaultUI
		info := version.Get()

		ui.PrintHeader("hann-cli")
		fmt.Printf("Version:     %s\n", ui.Bold(info.Version))
		fmt.Printf("Description: Developer & system utility for Linux (Termux & PRoot)\n")
		fmt.Printf("Platform:    %s\n", info.Platform)
		fmt.Println()

		ui.PrintSection("Available Commands")
		commands := [][]string{
			{"doctor", "Inspect Linux environment and verify developer toolchain"},
			{"system", "Display real-time system, CPU, memory, and uptime metrics"},
			{"ports", "List active listening ports and associated processes"},
			{"projects", "Discover development projects in your workspace"},
			{"version", "Show version, commit, and build environment info"},
		}

		for _, c := range commands {
			fmt.Printf("  %-12s %s\n", ui.Cyan(c[0]), c[1])
		}

		fmt.Println()
		ui.PrintSection("Flags")
		fmt.Printf("  %-12s %s\n", ui.Dim("-h, --help"), "Display help for hann or any command")
		fmt.Printf("  %-12s %s\n", ui.Dim("-v, --version"), "Display version information")
		fmt.Println()
		fmt.Printf("Run %s for detailed usage.\n", ui.Cyan("hann [command] --help"))
	},
}

// Execute runs the root CLI command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("hann %s (%s, %s)\n", version.Version, version.Commit, version.BuildDate))
}
