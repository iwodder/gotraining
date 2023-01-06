package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print version and exit",
	Long:  "Some long message here about what the version is called and why.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Task Manager")
	},
}
