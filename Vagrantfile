# -*- mode: ruby -*-
# vi: set ft=ruby :

NODE_NAME = "test"

Vagrant.configure("2") do |config|

  config.vm.box      = "ubuntu/xenial64"
  config.vm.hostname = NODE_NAME

  config.vm.define NODE_NAME

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  # config.vm.synced_folder "../data", "/vagrant_data"

  config.vm.provider :virtualbox do |vb|
    vb.name   = NODE_NAME
    vb.memory = "2048"
    vb.customize ["modifyvm", :id, "--cpuexecutioncap", "50"]
  end

  config.vm.provision :shell, inline: <<-SHELL
    wget https://apt.puppetlabs.com/puppet6-release-xenial.deb
    sudo dpkg -i puppet6-release-xenial.deb
    apt-get update
    apt-get install puppet-agent
  SHELL
end
