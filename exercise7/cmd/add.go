package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	description string

	addTaskCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a task",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := taskManager.CreateTask(name, description)
			if err != nil {
				log.Fatal(fmt.Sprintf("Unable to create task, \"%s\". Please try again later.", name))
			}
			fmt.Printf("Another task to crush! \"%s\" was added to your list", name)
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
