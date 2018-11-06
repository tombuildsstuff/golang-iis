package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/tombuildsstuff/golang-iis/iis/helpers"
)

func (c Client) Run(commands string) (*string, *string, error) {
	rInt := helpers.RandomInt()
	filename := fmt.Sprintf("command-%d.ps1", rInt)
	err := ioutil.WriteFile(filename, []byte(commands), os.FileMode(0700))
	if err != nil {
		return nil, nil, fmt.Errorf("Error writing command file: %+v", err)
	}

	var stderr bytes.Buffer
	var stdout bytes.Buffer

	// TODO: we could remove the need for a file by running these commands via WinRM, maybe?
	cmd := exec.Command("powershell.exe", "-ExecutionPolicy", "Bypass", "-NoLogo", "-NonInteractive", "-NoProfile", "-File", filename)

	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	defer os.Remove(filename)

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("Error starting: %+v", err)
	}

	if err := cmd.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Error waiting: %+v", err)
	}

	stdOutStr := stdout.String()
	stdErrStr := stderr.String()

	return &stdOutStr, &stdErrStr, nil
}
