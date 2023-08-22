package cmd

import (
	"github.com/spf13/cobra"
	"github.com/talios/find-dupe-music/find"
)

var SkipEditions bool

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		find.ScanFiles(SkipEditions, args)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().BoolVarP(&SkipEditions, "skip-editions", "s", false, "Skip different album editions from duplicate checks")
}
