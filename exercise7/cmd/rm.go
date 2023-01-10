package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rmTaskCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove a task from your list.",
		Run: func(cmd *cobra.Command, args []string) {
			if err := taskManager.Remove(name); err != nil {
				fmt.Println("Unable to remove task", err)
				os.Exit(1)
			}
			fmt.Printf("You removed task %s.\n", name)
		},
	}
)

func init() {
	rmTaskCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the task being removed.")
	rmTaskCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(rmTaskCmd)
}
