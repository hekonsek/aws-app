package awsom

import (
	"bytes"
	"fmt"
	"github.com/go-errors/errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

const UnixExitCodeGeneralError = 1

func CliError(err error) {
	fmt.Printf("Something went wrong: %s", err)
}

func ExitOnCliError(err error) {
	if err != nil {
		CliError(err)
		os.Exit(UnixExitCodeGeneralError)
	}
}

func CliCapture(handler func() error) (string, error) {
	readStdOut := os.Stdout
	pipeIn, pipeOut, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = pipeOut

	err = handler()
	if err != nil {
		return "", errors.WrapPrefix(err, "error from handler", 0)
	}

	channelOut := make(chan string)
	go func() {
		var buffer bytes.Buffer
		io.Copy(&buffer, pipeIn)
		channelOut <- buffer.String()
	}()

	err = pipeOut.Close()
	if err != nil {
		return "", err
	}
	os.Stdout = readStdOut
	out := <-channelOut

	return out, nil
}

type Exec struct {
	Command    string
	WorkingDir string
}

func (exe Exec) Run() ([]string, error) {
	parsedCommand := strings.Split(exe.Command, " ")
	cmd := exec.Command(parsedCommand[0], parsedCommand[1:]...)
	cmd.Dir = exe.WorkingDir
	stdoutStderr, err := cmd.CombinedOutput()
	outputLines := strings.Split(string(stdoutStderr), "\n")
	return outputLines[0 : len(outputLines)-1], err
}
