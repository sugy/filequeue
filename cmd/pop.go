// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	filequeue "github.com/sugy/filequeue/internal"
)

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop queue (dequeue)",
	Long:  `Retrieve and execute stored queues.`,
	Run: func(cmd *cobra.Command, args []string) {
		d, _ := cmd.Flags().GetString("queuedir")

		if len(d) == 0 {
			d = getDefaultQueueDirPath()
		}

		f, err := filequeue.NewFileQueue(d)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if err := f.Dequeue(); err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(popCmd)

	popCmd.Flags().StringP("queuedir", "d", "", "Queue directory.")
}
