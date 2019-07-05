package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageTestCmd represents the build command
var imageTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test image configuration",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if driver == "" {
			driver = viper.GetString("driver.default")
		}

		imageName = args[0]

		switch driver {
		case "vagrant":
			build := new(Vagrant)
			build.Configure()
			build.Test()
		case "packer":
			build := new(Packer)
			build.Configure()
			build.Test()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageTestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageTestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageTestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
