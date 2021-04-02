package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// ExecCommand: execute shell command and return stdout on
// success or stderr with non nil error on fail.
func ExecCommand(cmd []string) (string, error) {
	var stdout, stderr bytes.Buffer
	var err error
	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stdout = &stdout
	execCmd.Stderr = &stderr
	err = execCmd.Run()
	if err != nil {
		if strings.Contains(stderr.String(), "already exists. Exiting.") {
			err = os.ErrExist
		}
		return fmt.Sprintf(
			"stdout:\n%s\nstderr:\n%s\n",
			stdout.String(),
			stderr.String()), err
	}
	return stdout.String(), nil
}

// ExecCommandProgress: execute command and retrieve in realtime
// the last 'lineCount' lines from the stdout.
// NOTE: the 'progress' function must contain waiting operation.
func ExecCommandProgress(cmd []string, lineCount int, progress func(lines []string)) error {
	var err, errLoop error
	var stdout io.ReadCloser
	var stderr bytes.Buffer
	// build command
	execCmd := exec.Command(cmd[0], cmd[1:]...)
	execCmd.Stderr = &stderr
	// stdout redirection
	if stdout, err = execCmd.StdoutPipe(); err == nil {
		// Assignate to a reader
		buf := bufio.NewReader(stdout)
		if err = execCmd.Start(); err != nil {
			return err
		}
		// main loop that retrieve 'lineCount' lines on each
		// occurrence. Reader position is kept to continue reading
		// at next line on the next call until EOF or error occurs.
		lines := make([]string, lineCount)
		var lineBytes []byte
		for {
			for i := 0; i < lineCount; i++ {
				lineBytes, _, errLoop = buf.ReadLine()
				lines[i] = string(lineBytes)
			}
			if errLoop != nil {
				break
			}
			progress(lines)
		}
	}
	err = execCmd.Wait()
	if err != nil {
		return fmt.Errorf("%s", stderr.String())
	}
	return nil
}

// GetTerminalOut: retrieve terminal output via 'f' func.
func GetTerminalOut(f func(item string), cmd ...string) error {
	out, err := ExecCommand(cmd)
	tmpOut := strings.Split(string(out), "\n")
	for _, val := range tmpOut {
		if len(val) > 0 {
			f(val)
		}
	}
	return err
}

// ExecCommand: launch commandline application with arguments
// return terminal output and error.
// func ExecCommand(infos string, cmds ...string) (outTerm []byte, err error) {
// 	execCmd := exec.Command(cmds[0], cmds[1:]...)
// 	execCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
// 	outTerm, err = execCmd.CombinedOutput()
// 	if err != nil {
// 		err = errors.New(
// 			fmt.Sprintf("[%s][%s]\nCommand failed: %v\nTerminal:\n%v",
// 				infos,
// 				strings.Join(cmds, " "),
// 				err,
// 				string(outTerm)))
// 		return
// 	}
// 	return
// }

// // GetTerminalOut: retrieve terminal output via 'f' func.
// func GetTerminalOut(f func(item string), cmd ...string) error {
// 	out, err := ExecCommand(cmd[0], cmd...)
// 	tmpOut := strings.Split(string(out), "\n")
// 	for _, val := range tmpOut {
// 		if len(val) > 0 {
// 			f(val)
// 		}
// 	}
// 	return err
// }

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
// func GetEnvVar(envVar string) string {
// 	return os.Getenv(envVar)
// }

// GetUser: retrieve realUser and currentUser.
// func GetUser() (realUser, currentUser *user.User, err error) {
// 	if currentUser, err = user.Current(); err == nil {
// 		realUser, err = user.Lookup(os.Getenv("SUDO_USER"))
// 	}
// 	return realUser, currentUser, err
// }
