package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	listTasksCmd = &cobra.Command{
		Use:  "list",
		Long: "Lists all tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, v := range taskManager.ListTasks() {
				fmt.Printf("Title: %s, Description: %s\n", v.Name, v.Description)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listTasksCmd)
}
