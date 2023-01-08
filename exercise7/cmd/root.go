package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gotraining/exercise7/internal/tasks"
	"log"
	"os"
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to locate home directory.")
	}
	r, err := tasks.NewRepository(home + "/.tasks")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to open file at %s/.tasks", home))
	}
	taskManager = tasks.NewTaskManager(r)
}

var (
	name        string
	taskManager *tasks.TaskManager
	rootCmd     = &cobra.Command{
		Use:   "task",
		Short: "Manage your daily to-do's",
		Long:  "Ensure those important items are never forgotten again and stay on top of all your work.",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Unable to execute, ", err)
	}
}
