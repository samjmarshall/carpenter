package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Vagrant build properties
type Vagrant struct {
	runArg string
}

// Configure Vagrant build properties
func (v *Vagrant) Configure() {
	os.Setenv("VAGRANT_IMAGE_NAME", imageName)

	out, err := exec.Command("vagrant", "status", imageName).Output()

	if err != nil || strings.Contains(string(out), "not created (virtualbox)") {
		v.runArg = "up"
	} else if strings.Contains(string(out), "running (virtualbox)") {
		v.runArg = "provision"
	} else {
		fmt.Println("Unknown Vagrant machine state")
		os.Exit(1)
	}
}

// Run Vagrant image build
func (v *Vagrant) Run() {
	shell(fmt.Sprintf("vagrant %s", v.runArg))
}

// Clean up build artifacts
func (v *Vagrant) Clean() {
	shell("vagrant destroy -f")
}
