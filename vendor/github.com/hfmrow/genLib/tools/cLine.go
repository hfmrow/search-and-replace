package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
)

// ExecCommand: launch commandline application with arguments
// return terminal output and error.
func ExecCommand(commands string, args ...string) (output []byte, err error) {
	output, err = exec.Command(commands, args...).CombinedOutput()
	if len(output) > 0 {
		if err != nil {
			err = errors.New(fmt.Sprintf("Issue:\n%s\nTerminal:\n%s\n", err.Error(), string(output)))
		} else {
			err = errors.New(fmt.Sprintf("Terminal:\n%s\n", string(output)))
		}
	}
	return
}

// CheckCmd: Check for command if exist
func CheckCmd(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}
	return true
}

// Get input from commandline stdin
func GetStdin(ask string) (input string) {
	fmt.Print(ask + ": ")
	fmt.Scanln(&input)
	return input
}

// Print a colored line to stdout
// i.e:	PrintColor(196, "text to be displayed: ")	// Display colored text without line feed
//		PrintColor(49, "New name", "\n")			// Add some text with another color and printing line feed.
func PrintColored(col int, inStr ...interface{}) {
	var outStr string
	var out io.Writer = os.Stdout
	var buf = &bytes.Buffer{}
	out = buf
	for _, str := range inStr[1:] {
		outStr += str.(string)
	}
	fmt.Fprintf(out, "\x1b[38;5;%dm%s\x1b[0m%s", col, inStr[0], outStr)
	fmt.Printf(buf.String())
}

// GetEnvVar: retrieve an environment variable value.
func GetEnvVar(envVar string) string {
	return os.Getenv(envVar)
}

// GetUser: retrieve realUser and currentUser.
func GetUser() (realUser, currentUser *user.User, err error) {
	if currentUser, err = user.Current(); err == nil {
		realUser, err = user.Lookup(os.Getenv("SUDO_USER"))
	}
	return realUser, currentUser, err
}
