package cmd

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type Facts struct {
	AwsRegion   string   `yaml:"aws_region"`
	ImageFormat string   `yaml:"image_format"`
	ImageLayers []string `yaml:"image_layers"`
}

// imageBuildCmd represents the build command
var imageBuildCmd = &cobra.Command{
	Use:   "build [image name]",
	Short: "Build image",
	Long:  `Build and configure a virtual machine image.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		imageName = args[0]

		if driver == "" {
			driver = viper.GetString("image.driver.default")
		}

		err := generatePuppetFacts()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Generate image/puppet/facts.yaml")
		}

		switch driver {
		case "vagrant":
			build := new(Vagrant)
			build.Configure()
			build.Run()
			build.Test()
		case "packer":
			build := new(Packer)
			build.Configure()
			build.Run()
		}
	},
}

func init() {
	imageCmd.AddCommand(imageBuildCmd)
	imageBuildCmd.Flags().StringVarP(&layers, "layers", "l", "", "Image layers e.g. --layers=base,php. Default `base,[image name]`.")
}

func generatePuppetFacts() error {
	f := Facts{
		AwsRegion:   "ap-southeast-2",
		ImageFormat: "ami",
		ImageLayers: getLayers(),
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
