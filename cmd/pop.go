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

		d, _ := cmd.Flags().GetString("queuedir")

		if len(d) == 0 {
			d = getDefaultQueueDirPath()
		}

		f := filequeue.NewQueue(d)
		_ = f.Dequeue()
	},
}

func init() {
	rootCmd.AddCommand(popCmd)

	popCmd.Flags().StringP("queuedir", "d", "", "Queue directory.")
}
