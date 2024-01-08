package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GoScii",
	Long:  `All software has versions. This is GoScii's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoScii v0.1 -- HEAD")
	},
}
