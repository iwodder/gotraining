package cmd

import (
	"github.com/spf13/cobra"
)

var (
	doTaskCmd = &cobra.Command{
		Use:   "done",
		Short: "Mark task as done.",
		Run: func(cmd *cobra.Command, args []string) {
			taskManager.GetTask(name)
		},
	}
)

func init() {

	doTaskCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the task being added.")
	doTaskCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(doTaskCmd)
}
