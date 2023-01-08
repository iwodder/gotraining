package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:  "version",
	Long: "Print version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Embed build number later...")
	},
}
