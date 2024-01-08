package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "goscii",
		Short: "A simple CLI tool to convert images to ASCII art",
		Long:  `GoScii is a simple CLI tool to convert images to ASCII art. It is written in Go. You just need to pass the path to the image file and the magic happens!.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
