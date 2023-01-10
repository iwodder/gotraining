package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	description string

	addTaskCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a task",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := taskManager.CreateTask(name, description)
			if err != nil {
				fmt.Printf("Unable to create task, \"%s\". Please try again later.", name)
				os.Exit(1)
			}
			fmt.Printf("Another task to crush! \"%s\" was added to your list\n", name)
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
