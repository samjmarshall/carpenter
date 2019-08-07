package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageDestroyCmd represents the build command
var imageDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy up image build resources",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if driver == "" {
			driver = viper.GetString("image.driver.default")
		}

		imageName = args[0]

		switch driver {
		case "vagrant":
			build := new(Vagrant)
			build.Configure()
			build.Destroy()
		case "packer":
			build := new(Packer)
			build.Configure()
			build.Destroy()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageDestroyCmd)
}
