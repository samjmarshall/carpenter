package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageTestCmd represents the build command
var imageTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test image configuration",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if driver == "" {
			driver = viper.GetString("image.driver.default")
		}

		imageName := getImageName(args)

		switch driver {
		case "vagrant":
			build := new(Vagrant)
			build.Configure(imageName)
			build.Test()
		case "packer":
			build := new(Packer)
			build.Configure(imageName)
			build.Test()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageTestCmd)
	imageTestCmd.Flags().StringVarP(&layers, "layers", "l", "", "Image layers e.g. --layers=base,php. Default `base,[image name]`.")
}
