package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var driver string

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	imageCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "Image driver [vagrant|packer|docker]")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
