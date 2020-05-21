package shell

import (
	"bytes"
	"errors"
	"os/exec"
)

// Execute executes the command specified in the commandName parameter
// and with the arguments specified in the args parameter
func Execute(commandName string, args ...string) (string, error) {
	cmd := exec.Command(commandName, args...)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()

	if err != nil {
		return stdOut.String(), err
	}
	if stdErr.Len() > 0 {
		return stdOut.String(), errors.New(stdErr.String())
	}
	return stdOut.String(), nil
}
