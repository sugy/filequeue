// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	filequeue "github.com/sugy/filequeue/internal"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push queue (enqueue)",
	Long: `Save queue to one file in directory.
One file is created per queue.`,
	Run: func(cmd *cobra.Command, args []string) {
		d, _ := cmd.Flags().GetString("queuedir")
		k, _ := cmd.Flags().GetString("kind")
		m, _ := cmd.Flags().GetString("message")

		if len(d) == 0 {
			d = getDefaultQueueDirPath()
		}

		if len(m) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				m = scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading stdin:", err)
			}
		}

		f, err := filequeue.NewFileQueue(d)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		if err := f.Enqueue(k, m); err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringP("queuedir", "d", "", "Queue directory.")
	pushCmd.Flags().StringP("kind", "k", "exec", "Queue kind.")
	pushCmd.Flags().StringP("message", "m", "", "Queue message. Usually received from stdin.")
}
