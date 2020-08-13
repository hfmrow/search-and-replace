// comboBox.go

/*
	©2019 H.F.M. MIT license
*/

package gtk3Import

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

// ComboBoxTextGetAllEntries: Retrieve all ComboBoxText entries to a string slice.
func ComboBoxTextGetAllEntries(cbxEntry *gtk.ComboBoxText) (outSlice []string) {
	iTreeModel, err := cbxEntry.GetModel()
	model := iTreeModel.ToTreeModel()
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
	iTreeModel, err := cbxEntry.GetModel()
	model := iTreeModel.ToTreeModel()
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

func ComboBoxTextClearAll(cbxEntry *gtk.ComboBoxText) {
	cbxEntry.PrependText("")
	cbxEntry.SetActive(0)
	cbxEntry.RemoveAll()

}
