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
			driver = viper.GetString("driver.default")
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageDestroyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageDestroyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
