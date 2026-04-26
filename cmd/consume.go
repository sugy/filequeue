// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	filequeue "github.com/sugy/filequeue/internal"
)

// consumeCmd represents the consume command
var consumeCmd = &cobra.Command{
	Use:   "consume",
	Short: "Continuously process queue items",
	Long: `Consume queue items continuously. Watches the queue directory and
automatically processes items as they arrive. Runs until interrupted via
SIGTERM or SIGINT.`,
	Run: func(cmd *cobra.Command, args []string) {
		d, _ := cmd.Flags().GetString("queuedir")
		p, _ := cmd.Flags().GetString("pidfile")
		interval, _ := cmd.Flags().GetDuration("interval")

		if len(d) == 0 {
			d = getDefaultQueueDirPath()
		}
		if len(p) == 0 {
			p = getDefaultPidFilePath()
		}

		// Check for duplicate consumer
		if existingPid, err := filequeue.ReadPid(p); err == nil {
			if filequeue.IsRunning(existingPid) {
				fmt.Fprintf(os.Stderr, "consumer already running with PID %d\n", existingPid)
				os.Exit(1)
			}
			// Stale PID file, remove it
			if err := filequeue.RemovePid(p); err != nil {
				log.Warn(fmt.Sprintf("failed to remove stale pidfile: %v", err))
			}
		}

		f, err := filequeue.NewFileQueue(d)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		consumer, err := filequeue.NewConsumer(f, d, interval)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// Write PID file
		if err := filequeue.WritePid(p); err != nil {
			log.Fatal(err)
			os.Exit(2)
		}
		log.Info(fmt.Sprintf("consumer started with PID %d, pidfile: %s", os.Getpid(), p))

		// Setup signal handling
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
		defer cancel()

		// Run consumer
		if err := consumer.Start(ctx); err != nil {
			log.Error(fmt.Sprintf("consumer error: %v", err))
			os.Exit(3)
		}

		// Clean up PID file
		if err := filequeue.RemovePid(p); err != nil {
			log.Warn(fmt.Sprintf("failed to remove pidfile: %v", err))
		}

		log.Info("consumer exited")
	},
}

func init() {
	rootCmd.AddCommand(consumeCmd)

	consumeCmd.Flags().StringP("queuedir", "d", "", "Queue directory.")
	consumeCmd.Flags().StringP("pidfile", "p", "", "PID file path.")
	consumeCmd.Flags().DurationP("interval", "i", 5*time.Second, "Polling interval (used as fallback).")
}
