package cmd

import (
	"fmt"

	"github.com/frhndevweb/hann-cli/internal/system"
	"github.com/frhndevweb/hann-cli/internal/terminal"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Display real-time system, CPU, memory, and uptime metrics",
	Long:  "Inspect live Linux hardware statistics, kernel attributes, memory usage from /proc, and runtime container indicators.",
	Run: func(cmd *cobra.Command, args []string) {
		ui := terminal.DefaultUI
		ui.PrintHeader("Hann System")

		info, err := system.CollectSystemInfo()
		if err != nil {
			fmt.Printf("Error collecting system information: %s\n", err)
			return
		}

		// System
		ui.PrintSection("System")
		fmt.Printf("%-13s %s\n", "OS", info.OS)
		fmt.Printf("%-13s %s\n", "Kernel", info.Kernel)
		fmt.Printf("%-13s %s\n", "Architecture", info.Architecture)
		fmt.Printf("%-13s %s\n", "Hostname", info.Hostname)
		fmt.Printf("%-13s %s\n", "User", info.User)
		fmt.Printf("%-13s %s\n", "Shell", info.Shell)
		fmt.Println()

		// Resources
		ui.PrintSection("Resources")
		cpuCoreText := "cores"
		if info.CPU.Cores == 1 {
			cpuCoreText = "core"
		}
		fmt.Printf("%-13s %d %s\n", "CPU", info.CPU.Cores, cpuCoreText)

		if info.Memory.TotalBytes > 0 {
			memStr := fmt.Sprintf("%s / %s",
				system.FormatBytes(info.Memory.UsedBytes),
				system.FormatBytes(info.Memory.TotalBytes),
			)
			fmt.Printf("%-13s %s\n", "Memory", memStr)
		} else {
			fmt.Printf("%-13s %s\n", "Memory", "unavailable")
		}

		fmt.Printf("%-13s %s\n", "Uptime", system.FormatUptime(info.Uptime))
		fmt.Println()

		// Environment
		ui.PrintSection("Environment")
		prootStatus := "no"
		if info.IsPRoot {
			prootStatus = "yes"
		}
		termuxStatus := "no"
		if info.IsTermux {
			termuxStatus = "yes"
		}
		fmt.Printf("%-13s %s\n", "PRoot", prootStatus)
		fmt.Printf("%-13s %s\n", "Termux", termuxStatus)
		fmt.Println()

		// Directory
		ui.PrintSection("Directory")
		fmt.Println(info.WorkingDir)
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
}
