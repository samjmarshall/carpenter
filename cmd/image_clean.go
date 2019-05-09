package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// imageCleanCmd represents the build command
var imageCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean up image build resources",
	Run: func(cmd *cobra.Command, args []string) {

		if driver == "" {
			driver = viper.GetStringMapString("driver")["name"]
		}

		switch driver {
		case "vagrant":
			build := new(Vagrant)
			build.Clean()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageCleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
