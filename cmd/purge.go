// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	filequeue "github.com/sugy/filequeue/internal"
)

// purgeCmd ...
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Purge queue (remove all queues)",
	Long:  `Purge queue (remove all queues)`,
	Run: func(cmd *cobra.Command, args []string) {
		d, _ := cmd.Flags().GetString("queuedir")

		if len(d) == 0 {
			d = getDefaultQueueDirPath()
		}

		f := filequeue.NewFileQueue(d)
		_ = f.Purge()
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	purgeCmd.Flags().StringP("queuedir", "d", "", "Queue directory.")
}
