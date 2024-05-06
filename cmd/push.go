// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>

*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	filequeue "github.com/sugy/filequeue/lib"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push queue (enqueue)",
	Long: `Save queue to one file in directory.
One file is created per queue.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("push called")

		q, _ := cmd.Flags().GetString("queue")
		t, _ := cmd.Flags().GetString("type")

		if len(q) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				q = scanner.Text()
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading stdin:", err)
			}
		}
		f := filequeue.NewQueue()
		_ = f.Enqueue(t, q)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().StringP("type", "t", "exec", "Queue type.")
	pushCmd.Flags().StringP("queue", "q", "", "Queue strings. Usually received from stdin.")
}
