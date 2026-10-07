package cmd

import (
	"fmt"

	"github.com/frhndevweb/hann-cli/internal/environment"
	"github.com/frhndevweb/hann-cli/internal/terminal"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Inspect Linux environment and verify developer toolchain",
	Long:  "Perform diagnostics on current Linux system properties, Termux/PRoot environment, and development CLI tools.",
	Run: func(cmd *cobra.Command, args []string) {
		ui := terminal.DefaultUI
		ui.PrintHeader("Hann Doctor")

		env := environment.DetectEnvironment()
		tools := environment.CheckDefaultTools()

		passed := 0
		failed := 0

		// Environment checks
		ui.PrintSection("Environment")

		if env.OS != "" {
			ui.PrintCheck(true, "OS", env.OS)
			passed++
		} else {
			ui.PrintCheck(false, "OS", "undetected")
			failed++
		}

		if env.Architecture != "" {
			ui.PrintCheck(true, "Architecture", env.Architecture)
			passed++
		} else {
			ui.PrintCheck(false, "Architecture", "undetected")
			failed++
		}

		if env.Kernel != "" {
			ui.PrintCheck(true, "Kernel", env.Kernel)
			passed++
		} else {
			ui.PrintCheck(false, "Kernel", "undetected")
			failed++
		}

		if env.Shell != "" {
			ui.PrintCheck(true, "Shell", env.Shell)
			passed++
		} else {
			ui.PrintCheck(false, "Shell", "undetected")
			failed++
		}

		if env.IsPRoot {
			ui.PrintCheck(true, "PRoot", "detected")
			passed++
		} else {
			ui.PrintCheck(true, "PRoot", "not detected")
			passed++
		}

		if env.IsTermux {
			ui.PrintCheck(true, "Termux", "detected")
			passed++
		} else {
			ui.PrintCheck(true, "Termux", "not detected")
			passed++
		}

		fmt.Println()

		// Tools checks
		ui.PrintSection("Tools")
		for _, tool := range tools {
			if tool.Available {
				ui.PrintCheck(true, tool.Name, "")
				passed++
			} else {
				ui.PrintCheck(false, tool.Name, "not installed")
				failed++
			}
		}

		fmt.Println()

		// Summary
		ui.PrintSection("Summary")
		passSuffix := "checks passed"
		if passed == 1 {
			passSuffix = "check passed"
		}
		failSuffix := "checks failed"
		if failed == 1 {
			failSuffix = "check failed"
		}

		fmt.Println(ui.Green(fmt.Sprintf("%d %s", passed, passSuffix)))
		if failed > 0 {
			fmt.Println(ui.Red(fmt.Sprintf("%d %s", failed, failSuffix)))
		} else {
			fmt.Println(ui.Green("0 checks failed"))
		}
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
