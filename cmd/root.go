package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var SkipKeys []string
var ScanPaths []string

var rootCmd = &cobra.Command{
	Use:   "find-dupe-music",
	Short: "A simple tool to find duplicate music files",
	Long: `A tool for finding duplicate music files based on ID3 tags:

The application will scan thru multiple folders looking
for duplicate arist/album's based on their ID3 tags.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.find-dupes.yaml)")

	rootCmd.PersistentFlags().StringArrayVar(&SkipKeys, "skip", []string{},
		"list of artist:albums keys to skip/ignore for dupe checks")
	rootCmd.PersistentFlags().StringArrayVar(&ScanPaths, "path", []string{},
		"list of paths to scan for duplicate music")

	viper.BindPFlag("skip", rootCmd.PersistentFlags().Lookup("skip"))
	viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))

	viper.SetDefault("author", "Mark Derricutt <mark@talios.com>")
	viper.SetDefault("license", "apache")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".find-dupes")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}
