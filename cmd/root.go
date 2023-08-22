package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "find-dupe-music",
	Short: "A simple tool to find duplicate music files",
	Long: `A tool for finding duplicate music files based on ID3 tags:

The application will scan thru multiple folders looking
for duplicate arist/album's based on their ID3 tags.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
