package cmd

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Vagrant build properties
type Vagrant struct {
	RunArg   string
	Template string
}

// VagrantConfig - Vagrantfile template data
type VagrantConfig struct {
	ImageName string
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

	v.Template = `# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

    config.vm.box      = "ubuntu/xenial64"
    config.vm.hostname = "{{.ImageName}}"

	config.vm.define "{{.ImageName}}"

	config.vm.synced_folder "puppet/data", "/tmp/vagrant-puppet/data"
	config.vm.synced_folder "test", "/tmp/test"

	config.vm.provider "virtualbox" do |vb|
		vb.name   = "{{.ImageName}}"
		vb.memory = 2048

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
`
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
		t = template.Must(t.Parse(v.Template))

		f, err := os.Create(vagrantfilePath)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Create Vagrantfile")
			return
		}

		t.Execute(f, VagrantConfig{ImageName: imageName})
		f.Close()
	}

	shell(fmt.Sprintf("vagrant %s", v.RunArg))
}

// Clean up build artifacts
func (v *Vagrant) Clean() {
	shell("vagrant destroy -f")
}
