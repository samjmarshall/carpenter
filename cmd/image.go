package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var driver string
var imageName string
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

func getLayers() []string {
	if layers == "" {
		return []string{"base", imageName}
	}

	return strings.Split(layers, ",")
}

func inspecLocations() string {
	var locations strings.Builder

	for _, layer := range getLayers() {
		locations.WriteString(fmt.Sprintf("/tmp/inspec/layer/%s ", layer))
	}

	return locations.String()
}
