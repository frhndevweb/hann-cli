package cmd

import (
	"fmt"

	"github.com/frhndevweb/hann-cli/internal/terminal"
	"github.com/frhndevweb/hann-cli/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version, commit, and build environment info",
	Long:  "Display current hann-cli version, git commit, build date, and compiler details.",
	Run: func(cmd *cobra.Command, args []string) {
		ui := terminal.DefaultUI
		info := version.Get()

		ui.PrintHeader("hann-cli")
		fmt.Printf("Version:    %s\n", ui.Bold(info.Version))
		fmt.Printf("Commit:     %s\n", info.Commit)
		fmt.Printf("Built:      %s\n", info.BuildDate)
		fmt.Printf("Platform:   %s\n", info.Platform)
		fmt.Printf("Go runtime: %s\n", info.GoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
