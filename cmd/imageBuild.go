package cmd

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Facts struct {
	ImageLayers []string `yaml:"image_layers"`
}

// imageBuildCmd represents the build command
var imageBuildCmd = &cobra.Command{
	Use:   "build [image name]",
	Short: "Build image",
	Long:  `Build and configure a virtual machine image.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if driver == "" {
			driver = viper.GetString("image.driver.default")
		}

		imageName := getImageName(args)

		if viper.GetString("image.provisioner") == "puppet" {
			err := generatePuppetFacts(imageName)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Generate image/puppet/facts.yaml")
			}
		}

		switch driver {
		case "vagrant":
			build := new(Vagrant)
			build.Configure(imageName)
			build.Run()
			build.Test()
		case "packer":
			build := new(Packer)
			build.Configure(imageName)
			build.Run()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
	imageBuildCmd.Flags().StringVarP(&layers, "layers", "l", "", "Image layers e.g. --layers=base,php. Default `base,[image name]`.")
}

func generatePuppetFacts(imageName string) error {
	f := Facts{
		ImageLayers: getLayers(imageName),
	}

	y, err := yaml.Marshal(&f)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("image/puppet/facts.yaml", y, 0644)
	if err != nil {
		return err
	}

	return nil
}
