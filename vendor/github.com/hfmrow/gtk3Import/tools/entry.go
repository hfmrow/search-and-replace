// entry.go

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// GetExtEntry: Sanitize extension entries
func GetExtEntry(e *gtk.Entry, separator string) (out []string) {
	tmpOut := strings.Split(strings.TrimSpace(GetEntryText(e)), separator)
	if len(tmpOut) > 0 {
		for _, ext := range tmpOut {
			if len(ext) > 0 {
				tmp := strings.Split(strings.TrimSpace(ext), ".")
				var tmp1 []string
				for _, s := range tmp {
					if len(s) == 0 {
						s = "*"
					}
					tmp1 = append(tmp1, s)
				}
				out = append(out, strings.Join(tmp1, "."))
			}
		}
	}
	return
}

// SetExtEntry: Sanitize extension entries
func SetExtEntry(e *gtk.Entry, separator string, in []string) {
	if len(in) > 0 {
		e.SetText(strings.Join(in, separator+" "))
	} else {
		e.SetText("")
	}
}

// GetEntryText: retrieve value of an entry control.
func GetEntryText(entry *gtk.Entry) (outString string) {
	var err error
	if outString, err = entry.GetText(); err != nil {
		fmt.Errorf("GetEntryText: %s", err.Error())
	}
	return
}

// GetEntryTextAsInt: retrieve value of an entry control as integer
func GetEntryTextAsInt(entry *gtk.Entry) (outint int) {
	var err error
	var outString string
	if outString, err = entry.GetText(); err == nil {
		if outint, err = strconv.Atoi(outString); err == nil {
			return
		}
	}
	if err != nil {
		fmt.Errorf("GetEntryTextAsInt: %s", err.Error())
	}
	return
}
