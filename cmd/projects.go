package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/frhndevweb/hann-cli/internal/projects"
	"github.com/frhndevweb/hann-cli/internal/terminal"
	"github.com/spf13/cobra"
)

var targetPathFlag string

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Discover development projects in your workspace",
	Long:  "Scan target workspace directory for git repositories and language manifests (Go, Node, Rust, Python, etc.).",
	Run: func(cmd *cobra.Command, args []string) {
		ui := terminal.DefaultUI

		expandedPath, err := projects.ExpandPath(targetPathFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid target path: %v\n", err)
			return
		}

		displayPath := projects.DisplayPath(expandedPath)

		scanner := projects.NewScanner()
		found, err := scanner.Scan(expandedPath)

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("Projects directory not found:")
				fmt.Println(displayPath)
				fmt.Println()
				fmt.Println("Create it with:")
				fmt.Printf("mkdir -p %s\n", displayPath)
				return
			}
			fmt.Fprintf(os.Stderr, "Error: unable to scan directory: %v\n", err)
			return
		}

		ui.PrintHeader("Projects")
		fmt.Println(displayPath)
		fmt.Println()

		if len(found) == 0 {
			fmt.Printf("No projects found in %s.\n", displayPath)
			return
		}

		for i, proj := range found {
			numStr := fmt.Sprintf("%d.", i+1)
			fmt.Printf("%-3s %-20s %s\n", numStr, proj.Name, ui.Cyan(proj.Type))
		}
	},
}

func init() {
	projectsCmd.Flags().StringVarP(&targetPathFlag, "path", "p", "~/projects", "Target workspace directory to scan")
	rootCmd.AddCommand(projectsCmd)
}
