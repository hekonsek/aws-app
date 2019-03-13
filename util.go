package awsom

import (
	"bytes"
	"fmt"
	"github.com/Pallinder/sillyname-go"
	"github.com/go-errors/errors"
	"io"
	"os"
	"strings"
)

const UnixExitCodeGeneralError = 1

func RandomName() string {
	lowerCased := strings.ToLower(sillyname.GenerateStupidName())
	return strings.Replace(lowerCased, " ", "", -1)
}

func CliError(err error) {
	fmt.Printf("Something went wrong: %s", err)
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
