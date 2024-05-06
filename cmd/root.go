// Package cmd implements CLI applications.
/*
Copyright © 2024 sugy <sugy.kz@gmail.com>
*/
package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "filequeue",
	Version: getVersion(),
	Short:   "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Debug mode flag
var Debug bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "Debug mode")
	//rootCmd.SetVersionTemplate(`version {{.Version}}`)
}

func initConfig() {
	customFmt := new(log.TextFormatter)
	customFmt.TimestampFormat = "2006-01-02 15:04:05"
	customFmt.FullTimestamp = true
	log.SetFormatter(customFmt)
	log.SetOutput(os.Stdout)

	if Debug {
		log.SetLevel(log.DebugLevel)
	}
}
