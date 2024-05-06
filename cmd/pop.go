// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "Pop queue (dequeue)",
	Long:  `Retrieve and execute stored queues.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pop called")
	},
}

func init() {
	rootCmd.AddCommand(popCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// popCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// popCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
