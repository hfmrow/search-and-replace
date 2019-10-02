// dialogError.go

package gtk3Import

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	gler "github.com/hfmrow/genLib/tools/errors"
)

// DialogError: Display error messaged dialog returning true in error case.
// options: devMode, forceExit bool
func DialogError(window *gtk.Window, title, text string, err error, options ...bool) bool {
	var devMode, forceExit = true, false
	switch {
	case len(options) == 1:
		devMode = options[0]
	case len(options) == 2:
		devMode = options[0]
		forceExit = options[1]
	case len(options) > 2:
		fmt.Printf("Wrong , arguments number, %v\n", options)
		os.Exit(1)
	}
	if gler.Check(err) {
		if devMode {
			if DialogMessage(
				window,
				"error",
				title,
				fmt.Sprintf("\n\n"+text+":\n%s", err.Error()),
				"",
				"Stop",
				"Continue") == 0 {
				os.Exit(1)
			}
		} else {
			DialogMessage(
				window,
				"error",
				title,
				fmt.Sprintf("\n\n"+text+":\n%s", err.Error()),
				"",
				"Ok")
			if forceExit {
				os.Exit(1)
			}
		}
		return true
	}
	return false
}
