// treeView.go

/*
	©2019 H.F.M
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This library allow to facilitate Treeview operations.
	Can manage ListView and TreeView, only one of them at a time.

	i.e:
		func exampleTreeViewStructure() {
			var err error
			var tvs *gi.TreeViewStructure
			var storeSlice [][]interface{}
			var parentIter *gtk.TreeIter

			if tw, err := gtk.TreeViewNew(); err == nil { // Create TreeView. You can use existing one.
				if tvs, err = gi.TreeViewStructureNew(tw, false, false); err == nil { // Create Structure
					tvs.AddColumn("", "active", true, false, false, false, false) // With his columns
					tvs.AddColumn("Category", "markup", true, false, false, false, true)
					tvs.StoreSetup(new(gtk.TreeStore)) // Setup structure with desired TreeModel

					tvs.StoreDetach()        // Free TreeStore from TreeView while fill it. (useful for very large entries)
					for j := 0; j < 3; j++ { // Fill with parent nodes
						parentIter, _ = tvs.AddRow(nil, tvs.ColValuesIfaceToIfaceSlice(false, fmt.Sprintf("Parent %d", j)))

						for i := 0; i < 3; i++ { // Fill parents with childs nodes
							tvs.AddRow(parentIter, tvs.ColValuesIfaceToIfaceSlice(false, fmt.Sprintf("entry %d", i)))
						}
					}
					tvs.StoreAttach() // Say to TreeView that it get his StoreModel right now
				}
			}
			// Retrieve raw values with paths [][]interface{}. Can be done as [][]string too, and [][]interface{} without path.
			if err == nil {
				if storeSlice, err = tvs.StoreToIfaceSliceWithPaths(); err == nil {
					fmt.Println(storeSlice)
				}
			}
			if err != nil {
				log.Fatal(err)
			}
		}
*/

package gtk3Import

import (
	"errors"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// Notice: All options, functions, if they're needed, must be set before starting "StoreSetup" function.
// Otherwise, you can modify all of them at run time using Gtk3 objects
// (TreeView, ListStore, TreeStore, Columns, and so on). You can access it using the main structure.
type TreeViewStructure struct {
	store                gtk.ITreeModel // Used to determine wich TreeModel we work with.
	Model                *gtk.TreeModel // Actual TreeModel. Used in some functions to avoid use of (switch ... case) type selection.
	TreeView             *gtk.TreeView
	ListStore            *gtk.ListStore
	TreeStore            *gtk.TreeStore
	Selection            *gtk.TreeSelection
	MultiSelection       bool
	ActivateSingleClick  bool
	Modified             bool
	Columns              []column
	SelectionChangedFunc func() // Function to call when the selection has (possibly) changed.
	ForEachFunc          func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool
}

type column struct {
	Name           string
	Editable       bool
	ReadOnly       bool
	Sortable       bool
	Resizable      bool
	Expand         bool
	Visible        bool
	Attribute      string // "text", "markup", "pixbuf", "progress", "spinner", "active" (toggle button)
	Column         *gtk.TreeViewColumn
	EditTextFunc   func(cellRendererText *gtk.CellRendererText, path, text string, col int) // "text"
	EditActiveFunc func(cellRendererToggle *gtk.CellRendererToggle, path string, col int)   // "active" (toggle button)
}

// Create a new treeview structure (*TreeViewStructure)
func TreeViewStructureNew(treeView *gtk.TreeView, multiselection, activateSingleClick bool) (tvs *TreeViewStructure, err error) {
	tvs = new(TreeViewStructure)
	// Store data
	tvs.ActivateSingleClick = activateSingleClick
	tvs.MultiSelection = multiselection
	tvs.TreeView = treeView
	tvs.ClearAll()
	if tvs.Selection, err = tvs.TreeView.GetSelection(); err != nil {
		return nil, errors.New(fmt.Sprintln("Unable to get gtk.TreeSelection: ", err))
	}
	return
}

// StoreSetup: Configure the TreeView columns and build the *gtk.ListStore or
// *gtk.TreeStore object. The "store" argument must be *gtk.ListStore or
// *gtk.TreeStore to indicate with what kind of TreeModel we are working ...
// i.e: StoreSetup(new(gtk.TreeStore)), configure struct to work with a TreeStore.
func (tvs *TreeViewStructure) StoreSetup(store gtk.ITreeModel) (err error) {
	var colTypeSl []glib.Type
	var tmpColType glib.Type

	tvs.store = store
	// Set options
	tvs.TreeView.SetActivateOnSingleClick(tvs.ActivateSingleClick)
	if tvs.MultiSelection {
		tvs.Selection.SetMode(gtk.SELECTION_MULTIPLE)
	}
	// Build columns and his (default) edit function according to his type.
	for colIdx, _ := range tvs.Columns {
		if tmpColType, err = tvs.insertColumn(colIdx); err != nil {
			return errors.New(fmt.Sprintf("Unable to insert column nb %d: %s", colIdx, err.Error()))
		} else {
			colTypeSl = append(colTypeSl, tmpColType)
		}
	}
	if err == nil {
		if err = tvs.buildStore(colTypeSl); err == nil {
			// For replacment (In some cases) of typed methode assignation (switch ... case)
			switch store.(type) {
			case *gtk.ListStore:
				tvs.Model = &tvs.ListStore.TreeModel
			case *gtk.TreeStore:
				tvs.Model = &tvs.TreeStore.TreeModel
			}
		}
	}
	return err
}

// StoreDetach: Unlink "TreeModel" from TreeView. Useful when lot of rows must be inserted.
// After insertion, StoreAttach() must be used to restore the link with the treeview.
// tips: must be used before gtk.ListStore.Clear().
func (tvs *TreeViewStructure) StoreDetach() {
	if tvs.store != nil {
		tvs.Model.Ref()
		tvs.TreeView.SetModel(nil)
	}
}

// StoreAttach: To use after data insertion to restore the link with TreeView.
func (tvs *TreeViewStructure) StoreAttach() {
	if tvs.store != nil {
		tvs.TreeView.SetModel(tvs.Model)
		tvs.Model.Unref()
	}
}

// RemoveColumns: Remove column from MainStructure and TreeView.
func (tvs *TreeViewStructure) RemoveColumn(col int) (columnCount int) {
	columnCount = tvs.TreeView.RemoveColumn(tvs.Columns[col].Column)
	tvs.Columns = append(tvs.Columns[:col], tvs.Columns[col+1:]...)
	tvs.Modified = true
	return
}

// InsertColumn: Insert new column to MainStructure.
func (tvs *TreeViewStructure) InsertColumn(name, attribute string, pos int, editable, readOnly,
	sortable, resizable, expand, visible bool) {
	newCol := []column{{Name: name, Attribute: attribute, Editable: editable, ReadOnly: readOnly,
		Sortable: sortable, Resizable: resizable, Expand: expand, Visible: visible}}

	tvs.Columns = append(tvs.Columns[:pos], append(newCol, tvs.Columns[pos:]...)...)
}

// AddColumn: Add new column to MainStructure.
// attribute may be: text, pixbuf, progress, spinner, toggle
func (tvs *TreeViewStructure) AddColumn(name, attribute string, editable, readOnly,
	sortable, resizable, expand, visible bool) {
	col := column{Name: name, Attribute: attribute, Editable: editable, ReadOnly: readOnly,
		Sortable: sortable, Resizable: resizable, Expand: expand, Visible: visible}

	tvs.Columns = append(tvs.Columns, col)
}

// RemoveColumns: Clear TreeView columns and ListStore or TreeStore object.
// Depending on provided object type in "store" variable.
func (tvs *TreeViewStructure) ClearAll( /**/ ) (err error) {
	if tvs.TreeView != nil && len(tvs.Columns) != 0 {
		tvs.TreeView.SetModel(nil)
		for idx := len(tvs.Columns) - 1; idx > -1; idx-- {
			tvs.RemoveColumn(idx) // Remove columns (from last to first)
		}
		switch tvs.store.(type) {
		case *gtk.ListStore:
			if tvs.ListStore != nil {
				tvs.ListStore.Clear()
				tvs.ListStore.Unref()
			}
		case *gtk.TreeStore:
			if tvs.TreeStore != nil {
				tvs.TreeStore.Clear()
				tvs.TreeStore.Unref()
			}
		}
	}
	return
}

// insertColumn: Insert column at defined position
func (tvs *TreeViewStructure) insertColumn(colIdx int) (colType glib.Type, err error) {
	// renderCell: Set cellRenderer type and column options
	var renderCell = func(cellRenderer gtk.ICellRenderer, colIdx int) (err error) {
		var column *gtk.TreeViewColumn
		if column, err = gtk.TreeViewColumnNewWithAttribute(tvs.Columns[colIdx].Name,
			cellRenderer, tvs.Columns[colIdx].Attribute, colIdx); err == nil {
			tvs.Columns[colIdx].Column = column
			column.SetExpand(tvs.Columns[colIdx].Expand)
			column.SetResizable(tvs.Columns[colIdx].Resizable) // Set column resizable
			column.SetVisible(tvs.Columns[colIdx].Visible)
			if tvs.Columns[colIdx].Sortable {
				column.SetSortColumnID(colIdx) // Set column sortable
			}
			// tvs.TreeView.AppendColumn(column)
			tvs.TreeView.InsertColumn(column, colIdx)
		}
		return err
	}
	atrribute := tvs.Columns[colIdx].Attribute
	switch {
	case atrribute == "active": // "toggle"
		var cellRenderer *gtk.CellRendererToggle
		if cellRenderer, err = gtk.CellRendererToggleNew(); err == nil {
			// Define an edit function if not previously given
			if tvs.Columns[colIdx].EditActiveFunc == nil {
				tvs.Columns[colIdx].EditActiveFunc = func(cellRendererToggle *gtk.CellRendererToggle, path string, col int) {
					if !tvs.Columns[col].ReadOnly {
						var iter /*parentIter,*/ *gtk.TreeIter
						var childIter = new(gtk.TreeIter)
						var gValue *glib.Value
						var goValue interface{}
						var ok bool
						var changeValues = func(childIter *gtk.TreeIter, col int, goValue interface{}) (ok bool, err error) {
							switch tvs.store.(type) {
							case *gtk.ListStore:
								if err = tvs.ListStore.SetValue(childIter, col, !goValue.(bool)); err == nil {
									tvs.Modified = true
								}
							case *gtk.TreeStore:
								if err = tvs.TreeStore.SetValue(childIter, col, !goValue.(bool)); err == nil {
									tvs.Modified = true
								}
							}
							if err == nil {
								return tvs.Model.IterNext(childIter), nil
							}
							return ok, err
						}
						if iter, err = tvs.Model.GetIterFromString(path); err == nil {
							if gValue, err = tvs.Model.GetValue(iter, col); err == nil {
								if goValue, err = gValue.GoValue(); err == nil {
									if tvs.Model.IterHasChild(iter) { // Change state of all childs if there exists
										ok = tvs.Model.IterChildren(iter, childIter)
										for ok {
											if ok, err = changeValues(childIter, col, goValue); err != nil {
												log.Fatalf("Unable to edit (toggle) cell col %d, path %s: %s\n", col, path, err.Error())
											}

										}
									}
								}
							}
							_, err = changeValues(iter, col, goValue)
						}
						if err != nil {
							log.Fatalf("Unable to edit (toggle) cell col %d, path %s: %s\n", col, path, err.Error())
						}
					}
				}
			}
			if err == nil {
				if _, err = cellRenderer.Connect("toggled", tvs.Columns[colIdx].EditActiveFunc, colIdx); err == nil {
					if err = renderCell(cellRenderer, colIdx); err == nil {
						colType = glib.TYPE_BOOLEAN
					}
				}
			}
		}
	case atrribute == "spinner":
		var cellRenderer *gtk.CellRendererSpinner
		if cellRenderer, err = gtk.CellRendererSpinnerNew(); err == nil {
			cellRenderer.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRenderer, colIdx); err == nil {
				colType = glib.TYPE_FLOAT
			}
		}
	case atrribute == "progress":
		var cellRenderer *gtk.CellRendererProgress
		if cellRenderer, err = gtk.CellRendererProgressNew(); err == nil {
			cellRenderer.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRenderer, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case atrribute == "pixbuf":
		var cellRenderer *gtk.CellRendererPixbuf
		if cellRenderer, err = gtk.CellRendererPixbufNew(); err == nil {
			if err = renderCell(cellRenderer, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case atrribute == "text" || atrribute == "markup":
		var cellRenderer *gtk.CellRendererText
		cellRenderer, err = gtk.CellRendererTextNew()
		cellRenderer.SetProperty("editable", tvs.Columns[colIdx].Editable)
		// Define an edit function if not previously given
		if tvs.Columns[colIdx].EditTextFunc == nil {
			tvs.Columns[colIdx].EditTextFunc = func(cellRendererText *gtk.CellRendererText, path, text string, col int) {
				if !tvs.Columns[col].ReadOnly {
					var iter *gtk.TreeIter
					if iter, err = tvs.Model.GetIterFromString(path); err == nil {
						switch tvs.store.(type) {
						case *gtk.ListStore:
							if err = tvs.ListStore.SetValue(iter, col, text); err == nil {
								tvs.Modified = true
							}
						case *gtk.TreeStore:
							if err = tvs.TreeStore.SetValue(iter, col, text); err == nil {
								tvs.Modified = true
							}
						}
					}
					if err != nil {
						log.Fatalf("Unable to edit (text) cell col %d, path %s, text %s: %s\n", col, path, text, err.Error())
					}
				}
			}
		}
		if err == nil {
			if _, err = cellRenderer.Connect("edited", tvs.Columns[colIdx].EditTextFunc, colIdx); err == nil {
				if err = renderCell(cellRenderer, colIdx); err == nil {
					colType = glib.TYPE_STRING
				}
			}
		}
	default:
		err = errors.New(fmt.Sprintf("Error on setting attribute: %s is not implemented or inexistent.\n", tvs.Columns[colIdx].Attribute))
	}
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to Renderer Cell: %s\n", err.Error()))
	}
	return colType, err
}

// buildStore: Build ListStore or TreeStore object. Depending on provided object type in "store" variable.
func (tvs *TreeViewStructure) buildStore(colTypeSl []glib.Type) (err error) {
	switch tvs.store.(type) {
	case *gtk.ListStore:
		// Creating a ListStore. This is what holds the data that will be shown on our TreeView.
		if tvs.ListStore, err = gtk.ListStoreNew(colTypeSl...); err != nil {
			return errors.New(fmt.Sprintf("Unable to create ListStore: %s\n", err.Error()))
		}
		tvs.TreeView.SetModel(tvs.ListStore)
	case *gtk.TreeStore:
		// Creating a TreeStore. This is what holds the data that will be shown on our TreeView.
		if tvs.TreeStore, err = gtk.TreeStoreNew(colTypeSl...); err != nil {
			return errors.New(fmt.Sprintf("Unable to create TreeStore: %s\n", err.Error()))
		}
		tvs.TreeView.SetModel(tvs.TreeStore)
	}
	// Emitted whenever the selection has (possibly) changed.
	if tvs.SelectionChangedFunc != nil {
		_, err = tvs.Selection.Connect("changed", tvs.SelectionChangedFunc)
	}
	return err
}

// CountRows: Return the number of rows in treeview.
func (tvs *TreeViewStructure) CountRows() (count int) {
	var err error
	var forEachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool {
		count++
		return false
	}
	switch tvs.store.(type) {
	case *gtk.ListStore:
		err = tvs.ListStore.ForEach(forEachFunc)
	case *gtk.TreeStore:
		err = tvs.TreeStore.ForEach(forEachFunc)
	}
	if err != nil {
		fmt.Printf("Unable to retrieve the number of rows: %s", err.Error())
	}
	return
}

// AddRow: Append a row to the Store (defined by type of "store" variable).
// "parent" is useless for ListStore, if its set to nil on TreeStore, it will create a new parent
func (tvs *TreeViewStructure) AddRow(parent *gtk.TreeIter, row []interface{}) (iter *gtk.TreeIter, err error) {
	var colIdx []int
	switch tvs.store.(type) {
	case *gtk.ListStore:
		for idx, _ := range row {
			colIdx = append(colIdx, idx)
		}
		iter = tvs.ListStore.Append()                               // Get an iterator for a new row at the end of the ListStore
		if err = tvs.ListStore.Set(iter, colIdx, row); err != nil { // Set the contents of the row that the iterator represents
			return
		}
	case *gtk.TreeStore:
		iter = tvs.TreeStore.Append(parent) // Get an iterator for a new row at the end of the TreeStore or under parent if not nil
		for col, value := range row {
			if err = tvs.TreeStore.SetValue(iter, col, value); err != nil { // Set the contents of the row that the iterator represents
				return
			}
		}
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Unable to add row: %s", err.Error()))
	}
	tvs.Modified = true
	return iter, err
}

// AddRow: Insert a row after/before iter to "store": ListStore/Treestore. Parent may be nil for Liststore.
func (tvs *TreeViewStructure) InsertRow(inIter, parent *gtk.TreeIter, row []interface{}, before ...bool) (iter *gtk.TreeIter, err error) {
	var tmpBefore bool
	var colIdx []int
	// var path *gtk.TreePath
	if len(before) != 0 {
		tmpBefore = before[0]
	}
	for idx, _ := range row {
		colIdx = append(colIdx, idx)
	}
	switch tvs.store.(type) {
	case *gtk.ListStore:
		if tmpBefore { // Get the insertion iter
			iter = tvs.ListStore.InsertBefore(inIter)
		} else {
			iter = tvs.ListStore.InsertAfter(inIter)
		}
		err = tvs.ListStore.Set(iter, colIdx, row)
	case *gtk.TreeStore:
		if tmpBefore { // Get the insertion iter
			iter = tvs.TreeStore.InsertBefore(parent, inIter)
		} else {
			iter = tvs.TreeStore.InsertAfter(parent, inIter)
		}
		err = tvs.TreeStore.SetValue(iter, colIdx[0], row[0])
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Unable to insert row: %s", err.Error()))
	}
	tvs.Modified = true
	return iter, err
}

// DuplicateRow: Copy a row after iter to the listStore
func (tvs *TreeViewStructure) DuplicateRow(inIter, parent *gtk.TreeIter) (iter *gtk.TreeIter, err error) {
	var glibValue *glib.Value
	var goValue interface{}
	switch tvs.store.(type) {
	case *gtk.ListStore:
		iter = tvs.ListStore.InsertAfter(inIter)
		for colIdx, _ := range tvs.Columns {
			if glibValue, err = tvs.ListStore.GetValue(inIter, colIdx); err == nil {
				if goValue, err = glibValue.GoValue(); err == nil {
					err = tvs.ListStore.SetValue(iter, colIdx, goValue)
				}
			}
		}
	case *gtk.TreeStore:
		iter = tvs.TreeStore.InsertAfter(parent, inIter)
		for colIdx, _ := range tvs.Columns {
			if glibValue, err = tvs.TreeStore.GetValue(inIter, colIdx); err == nil {
				if goValue, err = glibValue.GoValue(); err == nil {
					err = tvs.TreeStore.SetValue(iter, colIdx, goValue)
				}
			}
		}
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintln("Unable to duplicating row: %s", err.Error()))
	}
	tvs.Modified = true
	tvs.ItersUnselect(inIter)
	tvs.ItersSelect(iter)
	return iter, err
}

// RemoveSelectedRows: Delete entries from selected iters or from given iters.
func (tvs *TreeViewStructure) RemoveSelectedRows(iters ...*gtk.TreeIter) (err error) {
	var ok bool
	if len(iters) == 0 {
		if iters, err = tvs.GetSelectedIters(); err != nil {
			return err
		}
	}
	tvs.StoreDetach()
	switch tvs.store.(type) {
	case *gtk.ListStore:
		for idx := len(iters) - 1; idx > -1; idx-- {
			ok = tvs.ListStore.Remove(iters[idx])
		}
	case *gtk.TreeStore:
		for idx := len(iters) - 1; idx > -1; idx-- {
			ok = tvs.TreeStore.Remove(iters[idx])
		}
	}
	tvs.StoreAttach()
	if ok {
		tvs.Modified = true
	}
	return err
}

// GetSelectedRows: Retrieve rows from selected iters or from given iters.
func (tvs *TreeViewStructure) GetSelectedRows(iters ...*gtk.TreeIter) (outSlice [][]string, err error) {
	var tmpSlice []string

	if len(iters) == 0 {
		if iters, err = tvs.GetSelectedIters(); err != nil {
			return outSlice, err
		}
	}
	for _, iter := range iters {
		if tmpSlice, err = tvs.GetRow(iter); err == nil {
			outSlice = append(outSlice, tmpSlice)
			tvs.Selection.SelectIter(iter) // To keep iters selected
		}
	}
	return outSlice, err
}

// GetSelectedIters: retreve list of selected iters
func (tvs *TreeViewStructure) GetSelectedIters() (iters []*gtk.TreeIter, err error) {
	var iter *gtk.TreeIter
	var getIters = func(glist *glib.List) (iters []*gtk.TreeIter, err error) {
		for row := glist; row != nil; row = row.Next() {
			if iter, err = tvs.Model.GetIter(row.Data().(*gtk.TreePath)); err == nil {
				iters = append(iters, iter)
			}
		}
		return
	}
	switch tvs.store.(type) {
	case *gtk.ListStore:
		glist := tvs.Selection.GetSelectedRows(tvs.ListStore)
		iters, err = getIters(glist)
	case *gtk.TreeStore: // TODO make it to run through parents/child's nodes
		glist := tvs.Selection.GetSelectedRows(tvs.TreeStore)
		iters, err = getIters(glist)
	}
	if err != nil {
		return iters, errors.New(fmt.Sprintln("Unable to retrieve selected iters: ", err))
	}
	return iters, err
}

// ItersSelect:
func (tvs *TreeViewStructure) ItersSelect(iters ...*gtk.TreeIter) {
	for _, iter := range iters {
		if !tvs.Selection.IterIsSelected(iter) {
			tvs.Selection.SelectIter(iter)
		}
	}
}

// ItersUnselect:
func (tvs *TreeViewStructure) ItersUnselect(iters ...*gtk.TreeIter) {
	for _, iter := range iters {
		if tvs.Selection.IterIsSelected(iter) {
			tvs.Selection.UnselectIter(iter)
		}
	}
}

// ItersSelectRange:
func (tvs *TreeViewStructure) ItersSelectRange(startIter, endIter *gtk.TreeIter) (err error) {
	var startPath, endPath *gtk.TreePath

	if startPath, err = tvs.Model.GetPath(startIter); err == nil {
		if endPath, err = tvs.Model.GetPath(endIter); err == nil {
			tvs.Selection.SelectRange(startPath, endPath)
		}
	}
	return err
}

// ScrollToCell: "column" argument set to nul, mean column 0,
func (tvs *TreeViewStructure) ScrollToCell(path *gtk.TreePath, column *gtk.TreeViewColumn, align bool, xalign, yalign float32) {
	if column == nil {
		column = tvs.Columns[0].Column
	}
	tvs.TreeView.ScrollToCell(path, column, align, xalign, yalign)
}

// ScrollToIter:
func (tvs *TreeViewStructure) ScrollToIter(iter *gtk.TreeIter) (err error) {
	var path *gtk.TreePath
	if path, err = tvs.Model.GetPath(iter); err == nil {
		tvs.ScrollToCell(path, nil, true, 0.5, 0.5)
	}
	return
}

// getRow: Get row from iter as []string
func (tvs *TreeViewStructure) GetRow(iter *gtk.TreeIter) (outSlice []string, err error) {
	var glibValue *glib.Value
	var valueString string

	for colIdx := 0; colIdx < len(tvs.Columns); colIdx++ {
		if glibValue, err = tvs.Model.GetValue(iter, colIdx); err == nil {
			if valueString, err = tvs.getStringCellValueByType(glibValue); err == nil {
				outSlice = append(outSlice, valueString)
			}
		}
		if err != nil {
			break
		}
	}
	return outSlice, err
}

// getCellValueByType: Retrieve cell value and convert it to string based on his type
func (tvs *TreeViewStructure) getStringCellValueByType(glibValue *glib.Value) (valueString string, err error) {
	var actualType glib.Type
	var valueIface interface{}

	if actualType, _, err = glibValue.Type(); err == nil {
		switch actualType {
		case glib.TYPE_STRING: // Strings
			valueString, err = glibValue.GetString()
		case glib.TYPE_BOOLEAN: // Boolean
			if valueIface, err = glibValue.GoValue(); err == nil {
				if valueIface.(bool) {
					valueString = "true"
				} else {
					valueString = "false"
				}
			}
		// case glib.TYPE_BOXED:
		default:
			err = errors.New(fmt.Sprintf("Type %s: not yet implemented\n", tvs.getStringGlibType(actualType)))
		}
	}
	return
}

// GetRowIface: Get row from iter as []interface{}
func (tvs *TreeViewStructure) GetRowIface(iter *gtk.TreeIter) (outIface []interface{}, err error) {
	var glibValue *glib.Value
	var value interface{}

	for colIdx := 0; colIdx < len(tvs.Columns); colIdx++ {
		if glibValue, err = tvs.Model.GetValue(iter, colIdx); err == nil {
			if value, err = glibValue.GoValue(); err == nil {
				outIface = append(outIface, value)
			}
		}
		if err != nil {
			break
		}
	}
	return outIface, err
}

// StoreToSlice: Retrieve all the rows values from a "store" as [][]string
func (tvs *TreeViewStructure) StoreToStringSlice() (outSlice [][]string, err error) {
	var tmpSlice []string
	// Foreach Function
	var foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool {
		if tmpSlice, err = tvs.GetRow(iter); err == nil {
			outSlice = append(outSlice, tmpSlice)
		} else {
			return true
		}
		return false
	}
	// Gathering columns names
	for _, col := range tvs.Columns {
		tmpSlice = append(tmpSlice, col.Column.GetTitle())
	}
	outSlice = append(outSlice, tmpSlice)
	// Retrieve values
	tvs.Model.ForEach(foreachFunc)

	return outSlice, err
}

// StoreToIface: Retrieve all the rows values from a "store" as [][]interface{}
func (tvs *TreeViewStructure) StoreToIfaceSlice() (outIface [][]interface{}, err error) {
	var tmpIface []interface{}
	// Foreach Function
	var retrieveValuesForeachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool {
		if tmpIface, err = tvs.GetRowIface(iter); err == nil {
			outIface = append(outIface, tmpIface)
		} else {
			return true
		}
		return false
	}
	// Gathering columns names
	for _, col := range tvs.Columns { // Gathering of columns names
		tmpIface = append(tmpIface, col.Column.GetTitle())
	}
	outIface = append(outIface, tmpIface)
	// Retrieve values
	tvs.Model.ForEach(retrieveValuesForeachFunc)

	return outIface, err
}

// StoreToIface: Retrieve all the rows values from a "store" as [][]interface{}
// The path to the iter is at start: [[0] [Parent 0]] [[0 0] [entry 0]] [[0 1] [entry 1]]
// [[1] [Parent 1]] [[1 0] [entry 0]] [[1 1] [entry 1]] [[1 2] [entry 2]] ...
// Thiw function does not retrieve the columns names !
func (tvs *TreeViewStructure) StoreToIfaceSliceWithPaths() (outIface [][]interface{}, err error) {
	var tmpIface []interface{}
	// Foreach Function
	var retrieveValuesForeachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool {
		if tmpIface, err = tvs.GetRowIface(iter); err == nil {
			tmpIface = append([]interface{}{path.GetIndices()}, tmpIface)
			outIface = append(outIface, tmpIface)
		} else {
			return true
		}
		return false
	}
	// Retrieve values
	tvs.Model.ForEach(retrieveValuesForeachFunc)

	return outIface, err
}

// ColValuesStringSliceToIfaceSlice: Convert string list to []interface, for simplify adding text rows
func (tvs *TreeViewStructure) ColValuesStringSliceToIfaceSlice(inSlice ...string) (outIface []interface{}) {
	for _, value := range inSlice {
		outIface = append(outIface, value)
	}
	return
}

// ColValuesIfaceToIfaceSlice: Convert interface list to []interface, for simplify adding text rows
func (tvs *TreeViewStructure) ColValuesIfaceToIfaceSlice(inSlice ...interface{}) (outIface []interface{}) {
	for _, value := range inSlice {
		outIface = append(outIface, value)
	}
	return
}

// func (tvs *TreeViewStructure) getDecendants(iter *gtk.TreeIter) (descendants [][]interface{}, err error) {
// 	// var parentPath, path *gtk.TreePath
// 	var iterDesc *gtk.TreeIter
// 	var ok bool = true
// 	var rowIface []interface{}

// 	ok = tvs.TreeStore.IterChildren(iter, iterDesc)
// 	for ok {
// 		rowIface, err = tvs.GetRowIface(iterDesc)
// 		descendants = append(descendants, rowIface)
// 		ok = tvs.TreeStore.IterNext(iterDesc)
// 	}
// 	return
// }

// TEST function
func (tvs *TreeViewStructure) selectRange(start, end *gtk.TreeIter) (err error) {
	var startPath, endPath *gtk.TreePath
	if startPath, err = tvs.ListStore.GetPath(start); err == nil {
		if endPath, err = tvs.ListStore.GetPath(end); err == nil {
			tvs.Selection.SelectRange(startPath, endPath)
		}
	}
	return err
}

// TEST function
func (tvs *TreeViewStructure) pathSelected(start *gtk.TreeIter) (err error) {
	var startPath *gtk.TreePath
	if startPath, err = tvs.ListStore.GetPath(start); err == nil {
		fmt.Println("iter", tvs.Selection.IterIsSelected(start))
		fmt.Println("path", tvs.Selection.PathIsSelected(startPath))
	}
	return err
}

// TEST function
func (tvs *TreeViewStructure) forEach() {
	var err error
	var model *gtk.TreeModel
	var ipath *gtk.TreePath
	var foreachFunc gtk.TreeModelForeachFunc
	foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool {
		if ipath, err = model.GetPath(iter); err == nil {
			fmt.Printf("path: %s, iter: %s\n", path.String(), ipath.String())
		} else {
			fmt.Println("error occured inside func: " + err.Error())
			return true
		}
		return false
	}
	if model, err = tvs.TreeView.GetModel(); err == nil {
		err = model.ForEach(foreachFunc)
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

// TEST function
func (tvs *TreeViewStructure) idx() {
	var err error
	// var model *gtk.TreeModel
	var path, cpypath *gtk.TreePath
	if path, err = gtk.TreePathNewFirst(); err == nil {
		fmt.Printf("path: %s\n", path.String())
		path.AppendIndex(3)
		fmt.Printf("depth: %d\n", path.GetDepth())
		path.PrependIndex(6)
		fmt.Printf("depth to copy: %d:%s\n", path.GetDepth(), path.String())
		if cpypath, err = path.Copy(); err == nil {
			fmt.Printf("copied: %d:%s\n", cpypath.GetDepth(), cpypath.String())
			fmt.Printf("compared: %d\n", cpypath.Compare(cpypath))
			cpypath.Next()
			fmt.Printf("next: :%s\n", cpypath.String())
			cpypath.Prev()
			fmt.Printf("prev: :%s\n", cpypath.String())
			cpypath.Up()
			fmt.Printf("up: :%s\n", cpypath.String())
			cpypath.Down()
			fmt.Printf("down: :%s\n", cpypath.String())
			fmt.Printf("IsAncestor: :%v\n", cpypath.IsAncestor(path))
			fmt.Printf("IsDescendant: :%v\n", cpypath.IsDescendant(path))
			if path, err = gtk.TreePathNewFromIndicesv([]int{2, 3, 4, 7, 8}); err == nil {
				fmt.Printf("new indices: %d:%s\n", path.GetDepth(), path.String())
			}
		}
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

// TEST function
func (tvs *TreeViewStructure) indices() {
	var err error
	var model *gtk.TreeModel
	var ipath, jpath *gtk.TreePath
	var foreachFunc gtk.TreeModelForeachFunc
	foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool {
		if ipath, err = model.GetPath(iter); err == nil {
			indices := ipath.GetIndices()
			jpath, _ = gtk.TreePathNewFromIndicesv(indices)
			indices1 := jpath.GetIndices()
			fmt.Printf("indices %v -> pathString: %v -> indices %v\n", indices, jpath.String(), indices1)
		} else {
			fmt.Println("error occured inside func: " + err.Error())
			return true
		}
		return false
	}
	if model, err = tvs.TreeView.GetModel(); err == nil {
		err = model.ForEach(foreachFunc)
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

// TEST function
// func (tvs *TreeViewStructure) getColsNames() (err error) {
// 	var glist *glib.List
// 	if glist, err = tvs.TreeView.GetColumns(); err == nil {
// 		for l := glist; l != nil; l = l.Next() {
// 			col := l.Data().(*gtk.TreeViewColumn)
// 			fmt.Println(col.GetTitle())
// 		}
// 	}
// 	if err != nil {
// 		err = errors.New("error occured while reading cols names: " + err.Error())
// 	}
// 	return err
// }

var glibType = map[int]string{
	0:  "glib.TYPE_INVALID",
	4:  "glib.TYPE_NONE",
	8:  "glib.TYPE_INTERFACE",
	12: "glib.TYPE_CHAR",
	16: "glib.TYPE_UCHAR",
	20: "glib.TYPE_BOOLEAN",
	24: "glib.TYPE_INT",
	28: "glib.TYPE_UINT",
	32: "glib.TYPE_LONG",
	36: "glib.TYPE_ULONG",
	40: "glib.TYPE_INT64",
	44: "glib.TYPE_UINT64",
	48: "glib.TYPE_ENUM",
	52: "glib.TYPE_FLAGS",
	56: "glib.TYPE_FLOAT",
	60: "glib.TYPE_DOUBLE",
	64: "glib.TYPE_STRING",
	68: "glib.TYPE_POINTER",
	72: "glib.TYPE_BOXED",
	76: "glib.TYPE_PARAM",
	80: "glib.TYPE_OBJECT",
	84: "glib.TYPE_VARIANT",
}

// getStringGlibType: Retrieve the string of glib value type.
func (tvs *TreeViewStructure) getStringGlibType(t glib.Type) string {
	for val, str := range glibType {
		if val == int(t) {
			return str
		}
	}
	return "Unnowen type"
}
