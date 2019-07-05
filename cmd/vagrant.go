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
	v.ImageName = imageName
}

// Run Vagrant image build
func (v *Vagrant) Run() {
	out, err := exec.Command("vagrant", "status", imageName).Output()

	if err != nil || strings.Contains(string(out), "not created (virtualbox)") {
		v.RunArg = "up"
	} else if strings.Contains(string(out), "running (virtualbox)") {
		v.RunArg = "provision"
	} else {
		fmt.Println("Unknown Vagrant machine state")
		os.Exit(1)
	}

	if v.RunArg == "up" && viper.IsSet("driver.vagrant.memory") {
		v.Memory = viper.GetInt("driver.vagrant.memory")
	} else {
		v.Memory = 1024
	}

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
	shell(fmt.Sprintf("vagrant destroy -f %s", v.ImageName))
}

// Test image configuration
func (v *Vagrant) Test() {
	user := "vagrant"
	host := "127.0.0.1"
	port := 2222
	keyFile := fmt.Sprintf(".vagrant/machines/%s/virtualbox/private_key", v.ImageName)

	fmt.Print("InSpec version: ")
	shell("inspec version")
	shell(fmt.Sprintf("inspec exec test/image/%s --target ssh://%s@%s:%d --key-files %s --sudo --chef-license=accept-silent", v.ImageName, user, host, port, keyFile))
}
