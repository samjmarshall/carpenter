package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var driver string
var imageName string

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Build and manage virtual machine images.",
	Long:  `Build and manage virtual machine images for AWS, VirtualBox or Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand")
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "Image driver [vagrant|packer|docker]")
}
