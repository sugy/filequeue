// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	filequeue "github.com/sugy/filequeue/lib"
)

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop queue (dequeue)",
	Long:  `Retrieve and execute stored queues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pop called")
		f := filequeue.NewQueue()
		_ = f.Dequeue()
	},
}

func init() {
	rootCmd.AddCommand(popCmd)
}
