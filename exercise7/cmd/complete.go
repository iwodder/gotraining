package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	compTaskCmd = &cobra.Command{
		Use:   "complete",
		Short: "List tasks completed within the last 6, 12, or 24 hours",
		Run: func(cmd *cobra.Command, args []string) {
			time, _ := cmd.Flags().GetInt("time")
			tasks, err := taskManager.ListCompletedTasks(time)
			if err != nil {
				fmt.Println("Unable to list all tasks", err)
				os.Exit(1)
			}
			fmt.Println("===Completed Tasks===")
			for i, v := range tasks {
				fmt.Printf("%d. %s", i+1, v.Name)
			}
		},
	}
)

func init() {
	compTaskCmd.Flags().Int("time", 24, "The window to check for completed tasks.")

	rootCmd.AddCommand(compTaskCmd)
}
