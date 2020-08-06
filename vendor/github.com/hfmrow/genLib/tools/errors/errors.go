// errors.go

package errors

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

// Check: Display errors in convenient way with stack display
// "options" set to true -> exit on error.
// NOTICE: exit option must not be used if a "defer" function
// is initiated, otherwise, defer will never be applied !
func Check(err error, options ...bool) (isError bool) {
	if err != nil {
		type errorInf struct {
			fn string // function
			f  string // file
		}
		var stacked []errorInf
		var outStrErr string
		isError = true
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
		fmt.Print(outStrErr)
		if len(options) > 0 {
			if options[0] {
				os.Exit(1)
			}
		}
	}
	return
}
