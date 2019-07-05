package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func shell(command ...string) {
	cmd := exec.Command(command[0], command[1:len(command)]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
