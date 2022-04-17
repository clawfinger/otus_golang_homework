package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, value := range env {
		if value.NeedRemove {
			if err := os.Unsetenv(key); err != nil {
				return 1
			}
		} else {
			if err := os.Setenv(key, value.Value); err != nil {
				return 1
			}
		}
	}

	cmdStruct := exec.Command(cmd[0], cmd[1:]...) //nolint

	cmdStruct.Stdout = os.Stdout
	cmdStruct.Stderr = os.Stderr

	if err := cmdStruct.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitErr := err.(*exec.ExitError) //nolint
			return exitErr.ExitCode()
		}
	}

	return
}
