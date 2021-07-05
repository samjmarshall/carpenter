package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Packer type build properties
type Packer struct {
	AMIName      string
	InstanceType string
	SubnetID     string
	AwsConfig    *aws.Config
}

// Configure Packer build properties
func (p *Packer) Configure(imageName string) {
	if os.Getenv("PACKER_BUILD_VERSION") == "" {
		p.AMIName = imageName
	} else {
		p.AMIName = fmt.Sprintf("%s-%s", imageName, os.Getenv("PACKER_BUILD_VERSION"))
	}

	os.Setenv("PACKER_AMI_NAME", p.AMIName)

	if os.Getenv("PACKER_AWS_REGION") == "" {
		if os.Getenv("AWS_REGION") != "" {
			os.Setenv("PACKER_AWS_REGION", os.Getenv("AWS_REGION"))
		} else {
			os.Setenv("PACKER_AWS_REGION", viper.GetString("image.driver.packer.aws_region"))
		}
	}

	if os.Getenv("PACKER_AWS_ACCESS_KEY_ID") == "" {
		os.Setenv("PACKER_AWS_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
		os.Setenv("PACKER_AWS_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))
	}

	p.AwsConfig = &aws.Config{
		Region:      aws.String(os.Getenv("PACKER_AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("PACKER_AWS_ACCESS_KEY_ID"), os.Getenv("PACKER_AWS_SECRET_ACCESS_KEY"), ""),
	}

	os.Setenv("PACKER_BUILD_NAME", p.AMIName)

	os.Setenv("PACKER_IMAGE_LAYERS", strings.Join(getLayers(p.AMIName), ","))
	os.Setenv("PACKER_INSPEC_LOCATIONS", inspecLocations(p.AMIName))

	if os.Getenv("PACKER_INSTANCE_TYPE") == "" {
		os.Setenv("PACKER_INSTANCE_TYPE", viper.GetString("image.driver.packer.instance_type"))
	}

	p.InstanceType = os.Getenv("PACKER_INSTANCE_TYPE")

	if os.Getenv("PACKER_SECURITY_GROUP_ID") == "" {
		os.Setenv("PACKER_SECURITY_GROUP_ID", viper.GetString("image.driver.packer.security_group_id"))
	}

	if os.Getenv("PACKER_SOURCE_AMI") == "" {
		if viper.IsSet("image.driver.packer.source_ami_filter") && viper.IsSet("image.driver.packer.source_ami_owner") {
			var filters []*ec2.Filter

			err := viper.UnmarshalKey("image.driver.packer.source_ami_filter", &filters)
			if err != nil {
				log.Error("Unable to unmarshal source AMI filters")
				os.Exit(1)
			}

			os.Setenv("PACKER_SOURCE_AMI", p.getLatestAMI(filters, viper.GetString("image.driver.packer.source_ami_owner")))
		} else {
			os.Setenv("PACKER_SOURCE_AMI", viper.GetString("image.driver.packer.source_ami"))
		}
	}

	if os.Getenv("PACKER_SUBNET_ID") == "" {
		os.Setenv("PACKER_SUBNET_ID", viper.GetString("image.driver.packer.subnet_id"))
	}

	p.SubnetID = os.Getenv("PACKER_SUBNET_ID")

	if os.Getenv("PACKER_SPOT_PRICE") == "" {
		if viper.IsSet("image.driver.packer.spot_price") {
			os.Setenv("PACKER_SPOT_PRICE", viper.GetString("image.driver.packer.spot_price"))
		} else {
			spotPrice := p.getSpotPrice()
			os.Setenv("PACKER_SPOT_PRICE", spotPrice)
		}
	}

	if os.Getenv("PACKER_VOLUME_SIZE") == "" {
		os.Setenv("PACKER_VOLUME_SIZE", viper.GetString("image.driver.packer.volume_size"))
	}

	if os.Getenv("PACKER_VPC_ID") == "" {
		os.Setenv("PACKER_VPC_ID", viper.GetString("image.driver.packer.vpc_id"))
	}

}

// Run Packer image build
func (p *Packer) Run() {
	log.WithFields(log.Fields{
		"amiName":         os.Getenv("PACKER_AMI_NAME"),
		"awsRegion":       os.Getenv("PACKER_AWS_REGION"),
		"instanceType":    os.Getenv("PACKER_INSTANCE_TYPE"),
		"securityGroupId": os.Getenv("PACKER_SECURITY_GROUP_ID"),
		"sourceAmi":       os.Getenv("PACKER_SOURCE_AMI"),
		"spotPrice":       os.Getenv("PACKER_SPOT_PRICE"),
		"subnetId":        os.Getenv("PACKER_SUBNET_ID"),
		"volumeSize":      os.Getenv("PACKER_VOLUME_SIZE"),
		"vpcId":           os.Getenv("PACKER_VPC_ID"),
	}).Info("Packer build properties")

	shell("packer", "build", "packer.json")
}

// Destroy up build artifacts
func (p *Packer) Destroy() {
	svc := ec2.New(session.New(p.AwsConfig))

	result, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("name"),
				Values: []*string{&p.AMIName},
			},
		},
		Owners: []*string{aws.String("self")},
	})

	if err != nil {
		log.Error("Unable to perform operation 'ec2:DescribeImages'")
		handleAWSError(err)
		return
	}

	if len(result.Images) > 0 {
		log.WithFields(log.Fields{"imageId": *result.Images[0].ImageId}).Info("Deregistering AMI")
		_, err := svc.DeregisterImage(&ec2.DeregisterImageInput{ImageId: result.Images[0].ImageId})

		if err != nil {
			log.Error("Unable to perform operation 'ec2:DeregisterImage'")
			handleAWSError(err)
			return
		}
	} else {
		log.WithFields(log.Fields{"amiName": p.AMIName}).Error("Could not find target AMI")
	}
}

// Test image configuration
func (p *Packer) Test() {
	// TODO
}

func (p *Packer) getSpotPrice() string {
	svc := ec2.New(session.Must(session.NewSession(p.AwsConfig)))

	subnetsResult, err := svc.DescribeSubnets(&ec2.DescribeSubnetsInput{
		SubnetIds: []*string{&p.SubnetID},
	})

	if err != nil {
		log.Warn("Unable to perform operation 'ec2:DescribeSubnets'")
		handleAWSError(err)
		return ""
	}

	historyResult, err := svc.DescribeSpotPriceHistory(&ec2.DescribeSpotPriceHistoryInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-type"),
				Values: []*string{&p.InstanceType},
			},
		},
		ProductDescriptions: []*string{aws.String("Linux/UNIX")},
		AvailabilityZone:    subnetsResult.Subnets[0].AvailabilityZone,
	})

	if err != nil {
		log.Warn("Unable to perform operation 'ec2:DescribeSpotPriceHistory'")
		handleAWSError(err)
		return ""
	}

	f, err := strconv.ParseFloat(*historyResult.SpotPriceHistory[0].SpotPrice, 64)

	if err != nil {
		log.Error(err.Error())
		return ""
	}

	return fmt.Sprintf("%.6f", (f*10/100)+f) // Set spot price 10% above current market price.
}

func (p *Packer) getLatestAMI(filters []*ec2.Filter, owner string) string {
	svc := ec2.New(session.Must(session.NewSession(p.AwsConfig)))
	result, err := svc.DescribeImages(&ec2.DescribeImagesInput{
		Filters: filters,
		Owners:  []*string{aws.String(owner)},
	})

	if err != nil {
		log.Error("Unable to perform operation 'ec2:DescribeImages'")
		handleAWSError(err)
		os.Exit(1)
	}

	if len(result.Images) > 1 {
		sort.Slice(result.Images, func(i, j int) bool {
			itime, _ := time.Parse(time.RFC3339, aws.StringValue(result.Images[i].CreationDate))
			jtime, _ := time.Parse(time.RFC3339, aws.StringValue(result.Images[j].CreationDate))
			return itime.Unix() > jtime.Unix()
		})
	}

	return *result.Images[0].ImageId
}

func handleAWSError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		default:
			log.WithFields(log.Fields{
				"code":    aerr.Code(),
				"message": aerr.Message(),
			}).Error(aerr.Error())
		}
	} else {
		log.Error(err.Error())
	}
}
