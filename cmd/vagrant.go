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
	AwsRegion   string
	Box         string
	Cwd         string
	ImageName   string
	Memory      int
	Provisioner string
	Running     bool
	Tester      string
}

// Configure Vagrant build properties
func (v *Vagrant) Configure(imageName string) {
	v.ImageName = imageName
	v.AwsRegion = os.Getenv("AWS_REGION")

	out, err := exec.Command("vagrant", "status", v.ImageName).Output()

	if strings.Contains(string(out), "not created") || err != nil {
		v.Running = false
	} else if strings.Contains(string(out), "running") {
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

	if viper.IsSet("image.tester") {
		v.Tester = viper.GetString("image.tester")
	}
}

// Run Vagrant image build
func (v *Vagrant) Run() {
	if viper.IsSet("image.driver.vagrant.memory") {
		v.Memory = viper.GetInt("image.driver.vagrant.memory")
	} else {
		v.Memory = 1024
	}

	if viper.IsSet("image.driver.vagrant.box") {
		v.Box = viper.GetString("image.driver.vagrant.box")
	} else {
		log.Error("Vagrant driver 'box' is not set")
		return
	}

	if viper.IsSet("image.provisioner") {
		v.Provisioner = viper.GetString("image.provisioner")
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

	config.vm.provider "virtualbox" do |vb|
		vb.name   = "{{.ImageName}}"
		vb.memory = {{.Memory}}

		vb.customize ["modifyvm", :id, "--cpuexecutioncap", "50"]
	end

	{{if eq .Provisioner "puppet"}}# Puppet apply
	config.vm.synced_folder "image/puppet", "/tmp/puppet"
	config.vm.provision "shell", inline: <<-SCRIPT
if [ ! -f /etc/apt/sources.list.d/puppet6.list ]; then
	wget -q https://apt.puppetlabs.com/puppet6-release-xenial.deb
	dpkg -i puppet6-release-xenial.deb
	apt-get update
	apt-get install -y puppet-agent git
	/opt/puppetlabs/puppet/bin/gem install r10k --no-document
fi

cd /tmp/puppet
[ ! -d modules ] && /opt/puppetlabs/puppet/bin/r10k puppetfile install
cp facts.yaml /opt/puppetlabs/facter/facts.d/facts.yaml
/opt/puppetlabs/bin/puppet apply manifests --modulepath=site:modules --hiera_config=hiera.yaml --verbose
SCRIPT{{end}}

	{{if eq .Tester "inspec"}}# InSpec test
	config.vm.synced_folder "image/inspec", "/tmp/inspec"
	config.vm.provision "shell", inline: "curl -sSL https://omnitruck.chef.io/install.sh | CI=true bash -s -- -P inspec"{{end}}

	# Upgrade all system packages
	config.vm.provision "shell", inline: "apt-get upgrade -y"
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
		shell("vagrant", "up", "--install-provider")
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
            sudo inspec exec %s --no-distinct-exit --no-create-lockfile --chef-license=accept-silent`, inspecLocations(v.ImageName)))
	}
}
