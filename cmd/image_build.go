package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageBuildCmd represents the build command
var imageBuildCmd = &cobra.Command{
	Use:   "build [image name]",
	Short: "Build image",
	Long:  `Build and configure a virtual machine image.`,
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
			build.Run()
		case "packer":
			build := new(Packer)
			build.Configure()
			build.Run()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageBuildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageBuildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
