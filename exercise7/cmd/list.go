package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	listTasksCmd = &cobra.Command{
		Use:  "list",
		Long: "Lists all tasks in your to-do list.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("====Task List====")
			for i, v := range taskManager.ListTasks() {
				fmt.Printf("  %d. Title: %s, Description: %s, Complete: %t\n", i+1, v.Name, v.Description, v.Complete)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listTasksCmd)
}
