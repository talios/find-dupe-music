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
	Short: "Scan paths for duplicate music files",
	Long: `Scan paths for duplicate music files based on
configured paths, with configurable album skipping.`,
	Run: func(_ *cobra.Command, args []string) {
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
