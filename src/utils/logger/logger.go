package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

// Print error message
// Exit after print
func PrintError(message string, err error) error {
	fmt.Println(message, err)
	os.Exit(1)

	return nil
}

// Print message to standard output
func PrintInfo(message string) {
	fmt.Println(message)
}

// Take the reader object and output the stream of messages
func PrintStream(reader io.Reader) error {
	termFd, isTerm := term.GetFdInfo(os.Stdout)
	err := jsonmessage.DisplayJSONMessagesStream(reader, os.Stdout, termFd, isTerm, nil)
	if err != nil {
		PrintError("Failed with", err)
	}

	return nil
}
