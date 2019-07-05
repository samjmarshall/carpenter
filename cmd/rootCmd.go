package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// cmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "carpenter",
	Short: "CLI to build machine images and infrastructure.",
	Long: `CLI abstraction for common image and infrastructure configuration
management and automation tool sets.

Carpenter allows you to build and configure you images and
infrastructure either locally or remote (AWS) via common commands.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the cmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $CWD/carpenter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find working directory.
		dir, err := os.Getwd()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in working directory with name ".carpenter" (without extension).
		viper.AddConfigPath(dir)
		viper.SetConfigName(".carpenter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig()
}
