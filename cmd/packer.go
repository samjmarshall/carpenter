package cmd

import (
	"os"

	"github.com/spf13/viper"
)

type Packer struct {
	runArg string
}

func (p *Packer) Run() {
	// TODO
}

func (p *Packer) Clean() {
	// TODO
}

func (p *Packer) Init(imageName string) {
	os.Setenv("PACKER_IMAGE_NAME", imageName)

	if os.Getenv("PACKER_AWS_REGION") == "" {
		if os.Getenv("AWS_REGION") != "" {
			os.Setenv("PACKER_AWS_REGION", os.Getenv("AWS_REGION"))
		} else {
			os.Setenv("PACKER_AWS_REGION", viper.GetStringMapString("packer")["aws_region"])
		}
	}

	if os.Getenv("PACKER_AWS_ACCESS_KEY_ID") == "" {
		os.Setenv("PACKER_AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
		os.Setenv("PACKER_AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	}

	if os.Getenv("PACKER_INSTANCE_TYPE") == "" {
		os.Setenv("PACKER_INSTANCE_TYPE", viper.GetStringMapString("packer")["instance_type"])
	}

	if os.Getenv("PACKER_SECURITY_GROUP_ID") == "" {
		os.Setenv("PACKER_SECURITY_GROUP_ID", viper.GetStringMapString("packer")["security_group_id"])
	}

	if os.Getenv("PACKER_SOURCE_AMI") == "" {
		os.Setenv("PACKER_SOURCE_AMI", viper.GetStringMapString("packer")["source_ami"])
	}

	if os.Getenv("PACKER_SPOT_PRICE") == "" {
		os.Setenv("PACKER_SPOT_PRICE", viper.GetStringMapString("packer")["spot_price"])
	}

	if os.Getenv("PACKER_SUBNET_ID") == "" {
		os.Setenv("PACKER_SUBNET_ID", viper.GetStringMapString("packer")["subnet_id"])
	}

	if os.Getenv("PACKER_VOLUME_SIZE") == "" {
		os.Setenv("PACKER_VOLUME_SIZE", viper.GetStringMapString("packer")["volume_size"])
	}

	if os.Getenv("PACKER_VPC_ID") == "" {
		os.Setenv("PACKER_VPC_ID", viper.GetStringMapString("packer")["vpc_id"])
	}

}
