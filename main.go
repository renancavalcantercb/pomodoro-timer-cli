package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "pomodoro"}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start a new pomodoro",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting a new pomodoro")
		},
	}

	rootCmd.AddCommand(startCmd)

	rootCmd.Execute()
}
