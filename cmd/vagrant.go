package cmd

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Vagrant build properties
type Vagrant struct {
	RunArg    string
	ImageName string
	Memory    int
}

// Configure Vagrant build properties
func (v *Vagrant) Configure() {
	out, err := exec.Command("vagrant", "status", imageName).Output()

	if err != nil || strings.Contains(string(out), "not created (virtualbox)") {
		v.RunArg = "up"
	} else if strings.Contains(string(out), "running (virtualbox)") {
		v.RunArg = "provision"
	} else {
		fmt.Println("Unknown Vagrant machine state")
		os.Exit(1)
	}

	v.ImageName = imageName

	if viper.IsSet("driver.vagrant.memory") {
		v.Memory = viper.GetInt("driver.vagrant.memory")
	} else {
		v.Memory = 1024
	}
}

// Run Vagrant image build
func (v *Vagrant) Run() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	vagrantfilePath := fmt.Sprintf("%s/Vagrantfile", cwd)

	if _, err := os.Stat(vagrantfilePath); os.IsNotExist(err) {
		t := template.New("Vagrantfile")
		t = template.Must(t.Parse(`# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

	config.vm.box      = "ubuntu/xenial64"
	config.vm.hostname = "{{.ImageName}}"

	config.vm.define "{{.ImageName}}"

	config.vm.synced_folder "puppet/data", "/tmp/vagrant-puppet/data"
	config.vm.synced_folder "test", "/tmp/test"

	config.vm.provider "virtualbox" do |vb|
		vb.name   = "{{.ImageName}}"
		vb.memory = {{.Memory}}

		vb.customize ["modifyvm", :id, "--cpuexecutioncap", "50"]
	end

	config.vm.provision "shell", path: "bin/bootstrap.sh"

	config.vm.provision "puppet" do |puppet|
		puppet.manifests_path    = "puppet/manifests"
		puppet.manifest_file     = "site.pp"
		puppet.hiera_config_path = "puppet/hiera.yaml"
		puppet.module_path       = ["puppet/site", "puppet/modules"]
		puppet.options           = "--verbose"
	
		puppet.facter = {
			"image" => "{{.ImageName}}"
		}
	end

	config.vm.provision "shell" do |s|
		s.path = "bin/test.sh"
		s.env = {
			"IMAGE_NAME" => "{{.ImageName}}"
		}
	end
end
`))

		f, err := os.Create(vagrantfilePath)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Create Vagrantfile")
			return
		}

		t.Execute(f, v)
		f.Close()
	}

	shell(fmt.Sprintf("vagrant %s", v.RunArg))
}

// Destroy up build artifacts
func (v *Vagrant) Destroy() {
	shell("vagrant destroy -f")
}
