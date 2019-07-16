package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Vagrant build properties
type Vagrant struct {
	Running     bool
	ImageName   string
	Provisioner string
	Tester      string
	Box         string
	Memory      int
	Cwd         string
}

// Configure Vagrant build properties
func (v *Vagrant) Configure() {
	v.ImageName = imageName

	if viper.IsSet("tester") {
		v.Tester = viper.GetString("tester")
	}

	out, err := exec.Command("vagrant", "status", imageName).Output()

	if err != nil || strings.Contains(string(out), "not created (virtualbox)") {
		v.Running = false
	} else if strings.Contains(string(out), "running (virtualbox)") {
		v.Running = true
	} else {
		fmt.Println("Unknown Vagrant machine state")
		os.Exit(1)
	}

	v.Cwd, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Run Vagrant image build
func (v *Vagrant) Run() {
	if viper.IsSet("driver.vagrant.memory") {
		v.Memory = viper.GetInt("driver.vagrant.memory")
	} else {
		v.Memory = 1024
	}

	if viper.IsSet("driver.vagrant.box") {
		v.Box = viper.GetString("driver.vagrant.box")
	} else {
		log.Error("Vagrant driver 'box' is not set")
		return
	}

	if viper.IsSet("provisioner") {
		v.Provisioner = viper.GetString("provisioner")
	}

	vagrantfilePath := fmt.Sprintf("%s/Vagrantfile", v.Cwd)

	if _, err := os.Stat(vagrantfilePath); os.IsNotExist(err) {
		t := template.New("Vagrantfile")
		t = template.Must(t.Parse(`# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|

	config.vm.box      = "{{.Box}}"
	config.vm.hostname = "{{.ImageName}}"

	config.vm.define "{{.ImageName}}"

	{{if eq .Provisioner "puppet"}}config.vm.synced_folder "puppet/data", "/tmp/vagrant-puppet/data"{{end}}
	config.vm.synced_folder "test", "/tmp/test"

	config.vm.provider "virtualbox" do |vb|
		vb.name   = "{{.ImageName}}"
		vb.memory = {{.Memory}}

		vb.customize ["modifyvm", :id, "--cpuexecutioncap", "50"]
	end

	{{if eq .Provisioner "puppet"}}config.vm.provision "shell", inline: <<-SCRIPT
if [ ! -f /etc/apt/sources.list.d/puppet6.list ]; then
		wget https://apt.puppetlabs.com/puppet6-release-xenial.deb
		sudo dpkg -i puppet6-release-xenial.deb
		sudo apt-get update
		sudo apt-get install puppet-agent
fi
SCRIPT

	config.vm.provision "puppet" do |puppet|
		puppet.manifests_path    = "puppet/manifests"
		puppet.manifest_file     = "site.pp"
		puppet.hiera_config_path = "puppet/hiera.yaml"
		puppet.module_path       = ["puppet/site", "puppet/modules"]
		puppet.options           = "--verbose"
	
		puppet.facter = {
			"image_type" => "ami",
			"image"      => "{{.ImageName}}",
		}
	end{{end}}

	{{if eq .Tester "inspec"}}config.vm.provision "shell", inline: 'curl -L https://omnitruck.chef.io/install.sh | sudo bash -s -- -P inspec -s once'{{end}}

	config.vm.provision "shell", inline: 'sudo apt-get -y upgrade'

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

	if v.Running {
		shell("vagrant", "provision")
	} else {
		shell("vagrant", "up")
	}
}

// Destroy up build artifacts
func (v *Vagrant) Destroy() {
	if v.Running {
		shell("vagrant", "destroy", "-f", v.ImageName)
	}

	os.RemoveAll(".vagrant/")
	os.Remove("Vagrantfile")
}

// Test image configuration
func (v *Vagrant) Test() {
	// Run InSpec
	switch v.Tester {
	case "inspec":
		shell("vagrant", "ssh", "-c", fmt.Sprintf(`echo "Inspec version: $(sudo inspec version)";
			sudo inspec vendor /tmp/test/image/%s --overwrite --chef-license=accept-silent;
			sudo inspec exec /tmp/test/image/%s`, v.ImageName, v.ImageName))
	}
}
