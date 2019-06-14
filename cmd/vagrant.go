package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Vagrant struct {
	runArg string
}

func (v *Vagrant) Run() {
	shell(fmt.Sprintf("vagrant %s", v.runArg))
}

func (v *Vagrant) Clean() {
	shell("vagrant destroy -f")
}

func (v *Vagrant) Init(imageName string) {
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

func shell(command string) {
	cmdArgs := strings.Fields(command)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
