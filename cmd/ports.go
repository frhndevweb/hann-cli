package cmd

import (
	"fmt"
	"strconv"

	"github.com/frhndevweb/hann-cli/internal/ports"
	"github.com/frhndevweb/hann-cli/internal/terminal"
	"github.com/spf13/cobra"
)

var listeningOnlyFlag bool

var portsCmd = &cobra.Command{
	Use:   "ports",
	Short: "List active listening ports and associated processes",
	Long:  "Inspect open network sockets by prioritizing direct Linux /proc/net decoding with fallback mechanisms.",
	Run: func(cmd *cobra.Command, args []string) {
		ui := terminal.DefaultUI

		entries, err := ports.ScanPorts(ports.ScannerOptions{
			ListeningOnly: listeningOnlyFlag,
		})

		if err != nil {
			ui.PrintHeader("Hann Ports")
			ui.PrintWarning("Notice", fmt.Sprintf("%v", err))
			fmt.Println()
			fmt.Println("Tip: In unrooted PRoot/Termux environments, direct /proc/net access may be restricted by Android SELinux.")
			fmt.Println("Try installing iproute2 (`apt install iproute2`) or running with elevated permissions.")
			return
		}

		if len(entries) == 0 {
			ui.PrintHeader("Hann Ports")
			fmt.Println("No listening ports detected on this system.")
			return
		}

		headers := []string{"PORT", "PROTOCOL", "ADDRESS", "PROCESS"}
		var rows [][]string
		for _, e := range entries {
			rows = append(rows, []string{
				strconv.Itoa(e.Port),
				e.Protocol,
				e.Address,
				e.Process,
			})
		}

		fmt.Print(ui.FormatTable(headers, rows))
	},
}

func init() {
	portsCmd.Flags().BoolVarP(&listeningOnlyFlag, "listening", "l", true, "Filter only listening sockets")
	rootCmd.AddCommand(portsCmd)
}
