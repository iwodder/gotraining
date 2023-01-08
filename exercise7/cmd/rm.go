package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	rmTaskCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a task from your list.",
		Run: func(cmd *cobra.Command, args []string) {
			taskManager.Remove(name)
			fmt.Printf("You removed task %s.\n", name)
		},
	}
)

func init() {
	rmTaskCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the task being removed.")
	rmTaskCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(rmTaskCmd)
}
