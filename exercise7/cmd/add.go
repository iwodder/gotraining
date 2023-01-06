package cmd

import (
	"github.com/spf13/cobra"
)

var (
	description string

	addTaskCmd = &cobra.Command{
		Use:   "add",
		Short: "add a task",
		Long:  "Some long message here about what the version is called and why.",
		Run: func(cmd *cobra.Command, args []string) {
			taskManager.CreateTask(name, description)
		},
	}
)

func init() {
	addTaskCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the task being added.")
	addTaskCmd.Flags().StringVarP(&description, "description", "d", "", "The description of the task being added.")
	addTaskCmd.MarkFlagRequired("name")
	addTaskCmd.MarkFlagRequired("description")

	rootCmd.AddCommand(addTaskCmd)
}
