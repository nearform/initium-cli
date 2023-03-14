package logger

import (
	"fmt"
	"os"
	"io"

	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)

type LoggerUtil struct {
}

func (ls LoggerUtil) PrintError(message string, err error) error {
    fmt.Println(message, err)
    os.Exit(1)

    return nil
}

func (ls LoggerUtil) PrintInfo(message string) {
    fmt.Println(message)
}

func (ls LoggerUtil) PrintStream(reader io.Reader) error {
    termFd, isTerm := term.GetFdInfo(os.Stdout)
    err := jsonmessage.DisplayJSONMessagesStream(reader, os.Stdout, termFd, isTerm, nil)
    if (err != nil) {
        ls.PrintError("Failed to display logs", err)
    }

    return nil
}