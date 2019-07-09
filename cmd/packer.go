package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Packer type build properties
type Packer struct {
	amiName string
}

// Configure Packer build properties
func (p *Packer) Configure() {
	p.amiName = fmt.Sprintf("%s-%s", imageName, os.Getenv("PACKER_BUILD_VERSION"))

	os.Setenv("PACKER_IMAGE_NAME", imageName)
	os.Setenv("PACKER_AMI_NAME", p.amiName)

	if os.Getenv("PACKER_AWS_REGION") == "" {
		if os.Getenv("AWS_REGION") != "" {
			os.Setenv("PACKER_AWS_REGION", os.Getenv("AWS_REGION"))
		} else {
			os.Setenv("PACKER_AWS_REGION", viper.GetString("driver.packer.aws_region"))
		}
	}

	if os.Getenv("PACKER_AWS_ACCESS_KEY_ID") == "" {
		os.Setenv("PACKER_AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
		os.Setenv("PACKER_AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	}

	if os.Getenv("PACKER_INSTANCE_TYPE") == "" {
		os.Setenv("PACKER_INSTANCE_TYPE", viper.GetString("driver.packer.instance_type"))
	}

	if os.Getenv("PACKER_SECURITY_GROUP_ID") == "" {
		os.Setenv("PACKER_SECURITY_GROUP_ID", viper.GetString("driver.packer.security_group_id"))
	}

	if os.Getenv("PACKER_SOURCE_AMI") == "" {
		os.Setenv("PACKER_SOURCE_AMI", viper.GetString("driver.packer.source_ami"))
	}

	if os.Getenv("PACKER_SPOT_PRICE") == "" {
		os.Setenv("PACKER_SPOT_PRICE", viper.GetString("driver.packer.spot_price"))
	}

	if os.Getenv("PACKER_SUBNET_ID") == "" {
		os.Setenv("PACKER_SUBNET_ID", viper.GetString("driver.packer.subnet_id"))
	}

	if os.Getenv("PACKER_VOLUME_SIZE") == "" {
		os.Setenv("PACKER_VOLUME_SIZE", viper.GetString("driver.packer.volume_size"))
	}

	if os.Getenv("PACKER_VPC_ID") == "" {
		os.Setenv("PACKER_VPC_ID", viper.GetString("driver.packer.vpc_id"))
	}

}

// Run Packer image build
func (p *Packer) Run() {
	log.WithFields(log.Fields{
		"\nami_name":          os.Getenv("PACKER_AMI_NAME"),
		"\ninstance_type":     os.Getenv("PACKER_INSTANCE_TYPE"),
		"\nregion":            os.Getenv("PACKER_AWS_REGION"),
		"\nsecurity_group_id": os.Getenv("PACKER_SECURITY_GROUP_ID"),
		"\nsource_ami":        os.Getenv("PACKER_SOURCE_AMI"),
		"\nspot_price":        os.Getenv("PACKER_SPOT_PRICE"),
		"\nsubnet_id":         os.Getenv("PACKER_SUBNET_ID"),
		"\nvolume_size":       os.Getenv("PACKER_VOLUME_SIZE"),
		"\nvpc_id":            os.Getenv("PACKER_VPC_ID"),
	}).Info("Packer build properties")

	shell("packer", "build", "packer.json")
}

// Destroy up build artifacts
func (p *Packer) Destroy() {
	svc := ec2.New(session.New(&aws.Config{Region: aws.String(os.Getenv("PACKER_AWS_REGION"))}))

	describeResult, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{&p.amiName},
			},
			{
				Name:   aws.String("name"),
				Values: []*string{&p.amiName},
			},
		},
		Owners: []*string{aws.String("self")},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}

		return
	}

	if len(describeResult.Images) > 0 {
		fmt.Printf("Deregistering AMI => %s\n", *describeResult.Images[0].ImageId)
		_, err := svc.DeregisterImage(&ec2.DeregisterImageInput{ImageId: describeResult.Images[0].ImageId})

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}

			return
		}
	} else {
		log.Error(fmt.Sprintf("Could not find target AMI => %s", p.amiName))
	}
}

// Test image configuration
func (p *Packer) Test() {
	// TODO
}
