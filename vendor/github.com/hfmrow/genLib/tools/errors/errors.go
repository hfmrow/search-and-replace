// errors.go

package errors

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

type errorInf struct {
	fn string // function
	f  string // file
}

// Check: Display errors in conveiniant way
// bool, "options" : [0] = exit on error
func Check(err error, options ...bool) (errFound bool) {
	var stacked []errorInf
	var outStrErr string
	var brk bool

	switch len(options) {
	case 1:
		brk = options[0]
	}
	if err != nil {
		errFound = true
		stack := strings.Split(string(debug.Stack()), "\n")
		for errIdx := 5; errIdx < len(stack)-1; errIdx++ {
			stacked = append(stacked, errorInf{fn: stack[errIdx], f: strings.TrimSpace(stack[errIdx+1])})
			errIdx++
		}
		baseMessage := strings.Split(err.Error(), "\n\n")
		for _, mess := range baseMessage {
			outStrErr += fmt.Sprintf("[%s]\n", mess)
		}
		for errIdx := 1; errIdx < len(stacked); errIdx++ {
			outStrErr += fmt.Sprintf("[%s]*[%s]\n", strings.SplitN(stacked[errIdx].fn, "(", 2)[0], stacked[errIdx].f)
		}
	}
	fmt.Print(outStrErr)
	if brk {
		os.Exit(1)
	}
	return
}

// Reduce filepath length.
// func makeDispPath(path string, length int) (outPath string) {
// 	splited := strings.Split(path, string(os.PathSeparator))
// 	if len(splited) > length+1 {
// 		return "..." + string(os.PathSeparator) + filepath.Join(splited[len(splited)-length:]...)
// 	}
// 	return path
// }
