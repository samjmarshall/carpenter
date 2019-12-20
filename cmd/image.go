package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var driver string
var layers string

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build and manage virtual machine images",
	Long:  `Build and manage virtual machine images for AWS, VirtualBox or Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand")
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "Image driver [vagrant|packer|docker]")
}

func getImageName(args []string) string {
	if len(args) == 0 {
		if viper.IsSet("image.name") {
			return viper.GetString("image.name")
		} else {
			fmt.Println("Image name was not set")
			os.Exit(1)
		}
	}

	return args[0]
}

func getLayers(imageName string) []string {
	if layers == "" {
		if viper.IsSet("image.layers") {
			return viper.GetStringSlice("image.layers")
		}

		return []string{"base", imageName}
	}

	return strings.Split(layers, ",")
}

func inspecLocations(imageName string) string {
	var locations strings.Builder

	for _, layer := range getLayers(imageName) {
		locations.WriteString(fmt.Sprintf("/tmp/inspec/layer/%s ", layer))
	}

	return locations.String()
}
