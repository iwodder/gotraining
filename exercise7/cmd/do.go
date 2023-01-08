package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var (
	doTaskCmd = &cobra.Command{
		Use:   "done",
		Short: "Mark a task as complete.",
		Run: func(cmd *cobra.Command, args []string) {
			task, err := taskManager.GetTask(name)
			if err != nil {
				log.Fatal("Unable to complete task, try again later.")
			}
			taskManager.Complete(*task)
			fmt.Printf("Congrats! Task %s is off your list! Keep killin' it!\n", name)
		},
	}
)

func init() {

	doTaskCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the task you finished.")
	doTaskCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(doTaskCmd)
}
