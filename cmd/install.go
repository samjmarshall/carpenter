package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the image command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install carpenter dependencies",
	Long:  `Install carpenter dependencies defined in carpenter.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Missing subcommand")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	installCmd.PersistentFlags().StringVarP(&driver, "driver", "d", "", "Image driver [vagrant|packer|docker]")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
