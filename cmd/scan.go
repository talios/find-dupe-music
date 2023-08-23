package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/talios/find-dupe-music/find"
)

var IgnoreSkips bool

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		skipped := viper.GetStringSlice("skip")
		paths := viper.GetStringSlice("path")
		allPaths := []string{}
		allPaths = append(allPaths, paths...)
		allPaths = append(allPaths, args...)

		for _, skip := range skipped {
			slog.Info("Skipping album", "album", skip)
		}

		find.ScanFiles(IgnoreSkips, allPaths)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().BoolVarP(&IgnoreSkips, "ignore-skips", "s", false, "Ignore defined skips")
}
