// tools.go

/// +build ignore

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"
	"regexp"

	"github.com/gotk3/gotk3/gtk"
)

// ComboBoxTextGetAllEntries: Retrieve all ComboBoxText entries to a string slice.
func ComboBoxTextGetAllEntries(cbxEntry *gtk.ComboBoxText) (outSlice []string) {
	model, err := cbxEntry.GetModel()
	iter, ok := model.GetIterFirst()
	for ok {
		if glibValue, err := model.GetValue(iter, 0); err == nil {
			if entry, err := glibValue.GetString(); err == nil {
				outSlice = append(outSlice, entry)
				ok = model.IterNext(iter)
			}
		}
		if err != nil {
			fmt.Errorf("ComboBoxTextGetAllEntries: %s", err.Error())
		}
	}
	return outSlice
}

// Fill / Clean comboBoxText
func ComboBoxTextFill(cbxEntry *gtk.ComboBoxText, entries []string, options ...bool) {
	var prepend, removeAll bool
	switch len(options) {
	case 1:
		prepend = options[0]
	case 2:
		prepend = options[0]
		removeAll = options[1]
	}
	if !removeAll {
		for _, word := range entries {
			ComboBoxTextAddSetEntry(cbxEntry, word, prepend)
		}
		return
	}
	cbxEntry.RemoveAll()
}

// ComboBoxTextAddSetEntry: Add newEntry if not exist to ComboBoxText, Option: prepend:bool.
// Get index and set cbxText at it if already exist.
func ComboBoxTextAddSetEntry(cbxEntry *gtk.ComboBoxText, newEntry string, prepend ...bool) (existAtPos int) {
	var prependEntry bool
	var count int
	var iter *gtk.TreeIter
	var ok bool
	existAtPos = -1
	if len(prepend) > 0 {
		prependEntry = prepend[0]
	}
	model, err := cbxEntry.GetModel()
	iter, ok = model.GetIterFirst()
	for ok {
		if glibValue, err := model.GetValue(iter, 0); err == nil {
			if entry, err := glibValue.GetString(); err == nil {
				if entry == newEntry {
					existAtPos = count
					break
				}
				count++
				ok = model.IterNext(iter)
			}
		}
		if err != nil {
			fmt.Errorf("ComboBoxTextAddSetEntry: %s", err.Error())
		}
	}
	if existAtPos == -1 {
		switch {
		case prependEntry:
			cbxEntry.PrependText(newEntry)
		default:
			cbxEntry.AppendText(newEntry)
		}
	} else {
		cbxEntry.SetActiveIter(iter)
	}
	return existAtPos
}

// GetEntryText: retrieve value of an entry control.
func GetEntryText(entry *gtk.Entry) (outString string) {
	var err error
	outString, err = entry.GetText()
	if err != nil {
		fmt.Errorf("GetEntryText: %s", err.Error())
	}
	return outString
}

// SpinbuttonSetValues: Configure spin button
// min, max, value, stepIncrement, pageIncrement, pageSize
func SpinbuttonSetValues(sb *gtk.SpinButton, min, max, value int, step ...int) {
	incStep, pageIncrement, pageSize := 1, 0, 0
	switch len(step) {
	case 1:
		incStep = step[0]
	case 2:
		incStep = step[0]
		pageIncrement = step[1]
	case 3:
		incStep = step[0]
		pageIncrement = step[1]
		pageSize = step[2]
	}
	if ad, err := gtk.AdjustmentNew(float64(value), float64(min), float64(max),
		float64(incStep), float64(pageIncrement), float64(pageSize)); err == nil {
		sb.Configure(ad, 1, 0)
	} else {
		fmt.Errorf("SpinbuttonSetValues: %s", err.Error())
	}
}

// MarkupHttpClickable: Search for http adresses to be treated as clickable link
func MarkupHttpClickable(inString string) (outString string) {
	outString = inString // In case nothing is found, returned value is same as entered value
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	indexes := reg.FindAllIndex([]byte(outString), -1)

	for idx := len(indexes) - 1; idx >= 0; idx-- {
		inLeft := outString[:indexes[idx][0]]
		inRight := outString[indexes[idx][1]:]
		url := outString[indexes[idx][0]:indexes[idx][1]]
		outString = inLeft + `<a href="` + url + `">` + url + `</a>` + inRight
	}
	return
}
