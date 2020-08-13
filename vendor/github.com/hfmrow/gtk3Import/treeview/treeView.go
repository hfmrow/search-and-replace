// treeView.go

/*
	Source file auto-generated on Tue, 12 Nov 2019 22:14:11 using Gotk3ObjHandler v1.5 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - TreeView library
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This library allow you to facilitate Treeview operations.
	Can manage ListView and TreeView, only one of them at a time.
	This lib handle every kind of column type as defined in the gtk3 development manual.
	Some conveignant functions are available to manage entries, columns, values, rows ...

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

// This library allow you to facilitate Treeview operations.
// Can manage ListView and TreeView, only one of them at a time.
// This lib handle every kind of column type as defined in the gtk3
// development manual. Some conveignant functions are available to
// manage entries, columns, values, rows ...
// Notice: All options, functions, if they're needed, must be set
// before starting "StoreSetup" function. Otherwise, you can modify
// all of them at runtime using Gtk3 objects (TreeView, ListStore,
// TreeStore, Columns, and so on). You can access it using the main
// structure.
type TreeViewStructure struct {

	// Actual TreeModel. Used in some functions to avoid use of
	// (switch ... case) type selection.
	Model *gtk.TreeModel

	// Direct acces to implicated objects
	TreeView  *gtk.TreeView
	ListStore *gtk.ListStore
	TreeStore *gtk.TreeStore
	Selection *gtk.TreeSelection

	// Basic options
	MultiSelection      bool
	ActivateSingleClick bool

	// The model has been modified?
	Modified bool

	// All columns option are available throught this structure
	Columns []column

	// TODO check out this ...
	// When "HasTooltip" is true, this function is launched,
	// Case of use: display tooltip according to rows currently hovered.
	// returned "bool" mean display or not the tooltip.
	CallbackTooltipFunc func(iter *gtk.TreeIter, path *gtk.TreePath, column *gtk.TreeViewColumn, tooltip *gtk.Tooltip) bool
	HasTooltip          bool

	// Function to call when the selection has (possibly) changed.
	SelectionChangedFunc func()

	// Used in gtk.Model or gtk.TreeSelection ForEach functions
	ModelForEachFunc     func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) bool
	SelectionForEachFunc func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{})

	// Used to substract from Y coordinates when using tooltip
	headerHeight int
	// Used to determine wich TreeModel we work with.
	store gtk.ITreeModel
}

type column struct {
	Name      string
	Editable  bool
	ReadOnly  bool
	Sortable  bool
	Resizable bool
	Expand    bool
	Visible   bool
	// Attributes with layout: "text", "markup", "pixbuf", "progress", "spinner", "active" (toggle button)
	// Attributes without layout: "pointer", "integer", "uint64", "int64"
	Attribute string

	// Direct access to the GtkTreeViewColumn
	Column  *gtk.TreeViewColumn
	ColType glib.Type

	// There is some default function defined for cell edition, normally, you don't have to define it yourself.
	// But in the case where you need specific operations, you can build you own edition function.
	EditTextFunc   func(cellRendererText *gtk.CellRendererText, path, text string, col int) // "text"
	EditActiveFunc func(cellRendererToggle *gtk.CellRendererToggle, path string, col int)   // "active" (toggle button)
}

// GetHeaderButton: Retrieve the button that assigned to the column header.
// May be used after "StoreSetup()" method.
func (col *column) GetHeaderButton() (button *gtk.Button, err error) {
	// var iWdgt *gtk.IWidget
	button = new(gtk.Button)
	if iWdgt, err := col.Column.GetButton(); err == nil {
		// Wrap button
		wdgt := iWdgt.ToWidget()
		actionable := &gtk.Actionable{wdgt.Object}
		button = &gtk.Button{gtk.Bin{gtk.Container{*wdgt}}, actionable}
	}
	return
}

/*********************\
*    Setup Functions   *
* Func that applied to *
* buils and handle     *
* treeview & models    *
***********************/

// Create a new treeview structure.
func TreeViewStructureNew(treeView *gtk.TreeView, multiselection, activateSingleClick bool) (tvs *TreeViewStructure, err error) {
	tvs = new(TreeViewStructure)
	// Store data
	tvs.ActivateSingleClick = activateSingleClick
	tvs.MultiSelection = multiselection
	tvs.TreeView = treeView
	tvs.ClearAll()
	tvs.HasTooltip = true
	if tvs.Selection, err = tvs.TreeView.GetSelection(); err != nil {
		return nil, fmt.Errorf("Unable to get gtk.TreeSelection: %s\n", err.Error())
	}
	return
}

// StoreSetup: Configure the TreeView columns and build the gtk.ListStore or
// gtk.TreeStore object. The "store" argument must be *gtk.ListStore or
// *gtk.TreeStore to indicate which kind of TreeModel we are working with ...
// i.e:
//   StoreSetup(new(gtk.TreeStore)), configure struct to work with a TreeStore.
//   StoreSetup(new(gtk.ListStore)), configure struct to work with a ListStore.
func (tvs *TreeViewStructure) StoreSetup(store gtk.ITreeModel) (err error) {
	var colTypeSl []glib.Type
	var tmpColType glib.Type

	tvs.store = store
	tvs.headerHeight = -1

	// Removing existing columns if exists ...
	for idx := int(tvs.TreeView.GetNColumns()) - 1; idx > -1; idx-- {
		tvs.TreeView.RemoveColumn(tvs.TreeView.GetColumn(idx))
	}

	// Tooltip setup
	if tvs.HasTooltip {
		tvs.TreeView.SetProperty("has-tooltip", true)
		tvs.TreeView.Connect("query-tooltip", tvs.treeViewQueryTooltip)
	} else {
		tvs.TreeView.SetProperty("has-tooltip", false)
	}

	// Set options
	tvs.TreeView.SetActivateOnSingleClick(tvs.ActivateSingleClick)
	if tvs.MultiSelection {
		tvs.Selection.SetMode(gtk.SELECTION_MULTIPLE)
	}

	// Build columns and his (default) edit function according to his type.
	for colIdx, _ := range tvs.Columns {
		if tmpColType, err = tvs.insertColumn(colIdx); err != nil {
			return fmt.Errorf("Unable to insert column nb %d: %s\n", colIdx, err.Error())
		} else {
			tvs.Columns[colIdx].ColType = tmpColType
			colTypeSl = append(colTypeSl, tmpColType)
		}
	}
	return tvs.buildStore(colTypeSl)
}

/*******************************\
*    Struct columns Functions    *
* Funct that applied to columns  *
* handled by the main structure  *
* this is the step before cols   *
* integration to the treeview    *
*********************************/

// RemoveColumn: Remove column from MainStructure and TreeView.
func (tvs *TreeViewStructure) RemoveColumn(col int) (columnCount int) {
	columnCount = tvs.TreeView.RemoveColumn(tvs.Columns[col].Column)
	tvs.Columns = append(tvs.Columns[:col], tvs.Columns[col+1:]...)
	tvs.Modified = true
	return
}

// InsertColumn: Insert new column to MainStructure, StoreSetup method must be called after.
func (tvs *TreeViewStructure) InsertColumn(name, attribute string, pos int, editable, readOnly,
	sortable, resizable, expand, visible bool) {
	newCol := []column{{Name: name, Attribute: attribute, Editable: editable, ReadOnly: readOnly,
		Sortable: sortable, Resizable: resizable, Expand: expand, Visible: visible}}

	tvs.Columns = append(tvs.Columns[:pos], append(newCol, tvs.Columns[pos:]...)...)
	tvs.Modified = true
}

// AddColumn: Adds a single new column to MainStructure.
// attribute may be: text, markup, pixbuf, progress, spinner, active ...
// see above for complete list.
func (tvs *TreeViewStructure) AddColumn(name, attribute string, editable, readOnly,
	sortable, resizable, expand, visible bool) {

	col := column{Name: name, Attribute: attribute, Editable: editable, ReadOnly: readOnly,
		Sortable: sortable, Resizable: resizable, Expand: expand, Visible: visible}
	tvs.Columns = append(tvs.Columns, col)
	tvs.Modified = true
}

// AddColumns: Adds several new columns to MainStructure.
// attribute may be: text, markup, pixbuf, progress, spinner, active
func (tvs *TreeViewStructure) AddColumns(nameAndAttribute [][]string, editable, readOnly,
	sortable, resizable, expand, visible bool) {
	for _, inCol := range nameAndAttribute {
		tvs.AddColumn(inCol[0], inCol[1], editable, readOnly, sortable, resizable, expand, visible)
	}
	tvs.Modified = true
}

/**************************\
*    Building Functions     *
* Make to create treemodel  *
* and handle the differant  *
* kind of columns with      *
* predefined edit functions *
****************************/

// insertColumn: Insert column at defined position
func (tvs *TreeViewStructure) insertColumn(colIdx int) (colType glib.Type, err error) {
	// renderCell: Set cellRenderer type and column options
	var renderCell = func(cellRenderer gtk.ICellRenderer, colIdx int) (err error) {
		var column *gtk.TreeViewColumn
		if column, err = gtk.TreeViewColumnNewWithAttribute(tvs.Columns[colIdx].Name,
			cellRenderer, tvs.Columns[colIdx].Attribute, colIdx); err == nil {
			tvs.Columns[colIdx].Column = column                // store column object to main struct.
			column.SetExpand(tvs.Columns[colIdx].Expand)       // Set expand option
			column.SetResizable(tvs.Columns[colIdx].Resizable) // Set resizable option
			column.SetVisible(tvs.Columns[colIdx].Visible)     // Set visible option
			if tvs.Columns[colIdx].Sortable {                  // Set sortable option
				column.SetSortColumnID(colIdx)
			}
			tvs.TreeView.InsertColumn(column, colIdx)
		}
		return err
	}
	attribute := tvs.Columns[colIdx].Attribute
	switch {
	case attribute == "active": // "toggle"
		var cellRenderer *gtk.CellRendererToggle
		if cellRenderer, err = gtk.CellRendererToggleNew(); err == nil {

			// An edit function may be user defined before structure initialisation,
			// if not, this one, is set as the default edit function for checkboxes.
			if tvs.Columns[colIdx].EditActiveFunc == nil {
				tvs.Columns[colIdx].EditActiveFunc = func(cellRendererToggle *gtk.CellRendererToggle, path string, col int) {
					if !tvs.Columns[col].ReadOnly {
						var iter *gtk.TreeIter
						var anotherIter = new(gtk.TreeIter)
						var goValue, goValueTmp interface{}
						var ok bool
						if iter, err = tvs.Model.GetIterFromString(path); err == nil {
							goValue = tvs.GetColValue(iter, col)
							// Switch the value of main iter
							if err = tvs.SetColValue(iter, col, !goValue.(bool)); err == nil {
								if tvs.Model.IterHasChild(iter) {
									// Parent: change state of all childs if they exists
									ok = tvs.changeTreeStateBool(iter, col, !goValue.(bool))
								}
								if ok = tvs.Model.IterParent(anotherIter, iter); ok {
									// Check each child and define his parent as true or false:
									// Whether all childs are set to true, the parent(s) must be set to true.
									// If at least one child is set to false, parent(s) must be set to false.
									ok = tvs.Model.IterChildren(anotherIter, iter)
									for ok {
										if goValueTmp = tvs.GetColValue(iter, col); !goValueTmp.(bool) {
											break
										}
										ok = tvs.Model.IterNext(iter)
									}
									ok = anotherIter != nil
									for ok {
										iter = anotherIter
										if err = tvs.SetColValue(iter, col, goValueTmp); err == nil {
											anotherIter = new(gtk.TreeIter) // Need to be initialised each time (RTFM).
											ok = tvs.Model.IterParent(anotherIter, iter)
										}
									}
								}
							}
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
	case attribute == "spinner":
		var cellRenderer *gtk.CellRendererSpinner
		if cellRenderer, err = gtk.CellRendererSpinnerNew(); err == nil {
			cellRenderer.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRenderer, colIdx); err == nil {
				colType = glib.TYPE_FLOAT
			}
		}
	case attribute == "progress":
		var cellRenderer *gtk.CellRendererProgress
		if cellRenderer, err = gtk.CellRendererProgressNew(); err == nil {
			cellRenderer.SetProperty("editable", tvs.Columns[colIdx].Editable)
			if err = renderCell(cellRenderer, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case attribute == "pixbuf":
		var cellRenderer *gtk.CellRendererPixbuf
		if cellRenderer, err = gtk.CellRendererPixbufNew(); err == nil {
			if err = renderCell(cellRenderer, colIdx); err == nil {
				colType = glib.TYPE_OBJECT
			}
		}
	case attribute == "text" || attribute == "markup":
		var cellRenderer *gtk.CellRendererText
		cellRenderer, err = gtk.CellRendererTextNew()
		cellRenderer.SetProperty("editable", tvs.Columns[colIdx].Editable)

		// An edit function may be user-defined before structure initialisation,
		// if not, this one, is set as the default edit function for text cells.
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
	case attribute == "pointer": // Pointer
		colType = glib.TYPE_POINTER
	case attribute == "integer": // INT
		colType = glib.TYPE_INT
	case attribute == "uint64": // UINT64
		colType = glib.TYPE_UINT64
	case attribute == "int64": // INT64
		colType = glib.TYPE_INT64
	default:
		err = fmt.Errorf("Error on setting attribute: %s is not implemented or inexistent.\n", tvs.Columns[colIdx].Attribute)
	}
	if err != nil {
		err = fmt.Errorf("Unable to make Renderer Cell: %s\n", err.Error())
	} else {
		// Add type to columns structure.
		tvs.Columns[colIdx].ColType = colType
	}
	return colType, err
}

// buildStore: Build ListStore or TreeStore object. Depending on provided
// object type in "store" variable.
func (tvs *TreeViewStructure) buildStore(colTypeSl []glib.Type) (err error) {

	switch tvs.store.(type) {
	case *gtk.ListStore: // Create the ListStore.
		if tvs.ListStore, err = gtk.ListStoreNew(colTypeSl...); err != nil {
			return fmt.Errorf("Unable to create ListStore: %s\n", err.Error())
		}
		tvs.Model = &tvs.ListStore.TreeModel
		tvs.TreeView.SetModel(tvs.ListStore)

	case *gtk.TreeStore: // Create the TreeStore.
		if tvs.TreeStore, err = gtk.TreeStoreNew(colTypeSl...); err != nil {
			return fmt.Errorf("Unable to create TreeStore: %s\n", err.Error())
		}
		tvs.Model = &tvs.TreeStore.TreeModel
		tvs.TreeView.SetModel(tvs.TreeStore)
	}
	// Emitted whenever the selection has (possibly, RTFM) changed.
	if tvs.SelectionChangedFunc != nil { // link to callback function if exists.
		_, err = tvs.Selection.Connect("changed", tvs.SelectionChangedFunc)
	}
	return err
}

/**********************\
*    Path Functions     *
* Funct that            *
* applied to Path       *
************************/

// // ScrollToCell: "column" argument set to nul, mean column 0,
// func (tvs *TreeViewStructure) scrollToCell(path *gtk.TreePath, column *gtk.TreeViewColumn, align bool, xalign, yalign float32) {
// 	if column == nil {
// 		column = tvs.Columns[0].Column
// 	}
// 	tvs.TreeView.ScrollToCell(path, column, align, xalign, yalign)
// }

/**********************\
*    Iters Functions    *
* Funct that            *
* applied to Iters      *
************************/

// GetSelectedIters: retrieve list of selected iters,
// return nil whether nothing selected.
func (tvs *TreeViewStructure) GetSelectedIters() (iters []*gtk.TreeIter) {
	iters = make([]*gtk.TreeIter, tvs.Selection.CountSelectedRows())
	var count int
	// tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData interface{}) {
	tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) {
		iters[count] = iter
		count++
	})
	if len(iters) == 0 {
		iters = nil
	}
	return
}

// GetSelectedPaths: retrieve list of selected paths,
// return nil whether nothing selected.
func (tvs *TreeViewStructure) GetSelectedPaths() (paths []*gtk.TreePath) {
	paths = make([]*gtk.TreePath, tvs.Selection.CountSelectedRows())
	var count int
	tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) {
		paths[count] = path
		count++
	})
	if len(paths) == 0 {
		paths = nil
	}
	return
}

// ItersSelect: Select provided Iters.
func (tvs *TreeViewStructure) ItersSelect(iters ...*gtk.TreeIter) {
	for _, iter := range iters {
		if !tvs.Selection.IterIsSelected(iter) {
			tvs.Selection.SelectIter(iter)
		}
	}
}

// ItersUnselectAll: Unselect all selected iters.
func (tvs *TreeViewStructure) ItersUnselectAll() {
	tvs.Selection.UnselectAll()
}

// ItersUnselect: Unselect provided Iters.
func (tvs *TreeViewStructure) ItersUnselect(iters ...*gtk.TreeIter) {
	for _, iter := range iters {
		if tvs.Selection.IterIsSelected(iter) {
			tvs.Selection.UnselectIter(iter)
		}
	}
}

// ItersSelectRange: Select range between start and end iters.
func (tvs *TreeViewStructure) ItersSelectRange(startIter, endIter *gtk.TreeIter) (err error) {
	var startPath, endPath *gtk.TreePath

	if startPath, err = tvs.Model.GetPath(startIter); err == nil {
		if endPath, err = tvs.Model.GetPath(endIter); err == nil {
			tvs.Selection.SelectRange(startPath, endPath)
		}
	}
	return err
}

// ScrollToIter: scroll to iter, pointing to the column if it has been specified.
func (tvs *TreeViewStructure) IterScrollTo(iter *gtk.TreeIter, column ...int) (err error) {
	var path *gtk.TreePath
	var colNb int
	if len(column) > 0 {
		colNb = column[0]
	}
	if col := tvs.TreeView.GetColumn(colNb); col != nil {
		if path, err = tvs.Model.GetPath(iter); err == nil {
			if path != nil {
				tvs.TreeView.ScrollToCell(path, col, true, 0.5, 0.5)
			} else {
				err = fmt.Errorf("IterScrollTo: Unable to get path from iter\n")
			}
		}
	} else {
		err = fmt.Errorf("IterScrollTo: Unable to get column %d\n", colNb)
	}
	return
}

/******************************\
*    Cols Functions             *
* Funct that applied to Columns *
********************************/

// GetColValue: Get value from iter of specific column as interface type.
func (tvs *TreeViewStructure) GetColValue(iter *gtk.TreeIter, col int) (value interface{}) {
	var gValue *glib.Value
	var err error
	if gValue, err = tvs.Model.GetValue(iter, col); err == nil {
		if value, err = gValue.GoValue(); err == nil {
			return
		}
	}
	log.Fatalf("GetColValue: %s\n", err.Error())
	return
}

// SetColValue: Set value to iter for a specific column as interface type.
func (tvs *TreeViewStructure) SetColValue(iter *gtk.TreeIter, col int, goValue interface{}) (err error) {
	switch tvs.store.(type) {
	case *gtk.ListStore:
		err = tvs.ListStore.SetValue(iter, col, goValue)
	case *gtk.TreeStore:
		err = tvs.TreeStore.SetValue(iter, col, goValue)
	}
	if err != nil {
		log.Fatalf("SetColValue: %s\n", err.Error())
		return
	}
	tvs.Modified = true
	return
}

// GetColValueFromPath: Get value from path of specific column as interface type.
// Note: should be used only if there is no other choice, prefer using iter to get values.
func (tvs *TreeViewStructure) GetColValuePath(path *gtk.TreePath, col int) (value interface{}) {
	var err error
	var iter *gtk.TreeIter

	switch tvs.store.(type) {
	case *gtk.ListStore:
		iter, err = tvs.ListStore.GetIter(path)
	case *gtk.TreeStore:
		iter, err = tvs.TreeStore.GetIter(path)
	}
	if err != nil {
		log.Fatalf("GetColValuePath: unable to get iter from path: %s\n", err.Error())
		return
	}
	return tvs.GetColValue(iter, col)
}

// SetColValue: Set value to path for a specific column as interface type.
// Note: should be used only if there is no other choice, prefer using iter to set values.
func (tvs *TreeViewStructure) SetColValuePath(path *gtk.TreePath, col int, goValue interface{}) (err error) {
	var iter *gtk.TreeIter
	switch tvs.store.(type) {
	case *gtk.ListStore:
		if iter, err = tvs.ListStore.GetIter(path); err == nil {
			err = tvs.ListStore.SetValue(iter, col, goValue)
		}
	case *gtk.TreeStore:
		if iter, err = tvs.TreeStore.GetIter(path); err == nil {
			err = tvs.TreeStore.SetValue(iter, col, goValue)
		}
	}
	if err != nil {
		return
	}
	tvs.Modified = true
	return
}

/**************************\
*    Rows Functions         *
* Func that applied to rows *
****************************/

// CountRows: Return the number of rows in treeview.
func (tvs *TreeViewStructure) CountRows() int {
	return tvs.Model.IterNChildren(nil)
}

// GetRowNbIter: Return the row number handled by the given iter,
// without any depth consideration.
func (tvs *TreeViewStructure) GetRowNbIter(iter *gtk.TreeIter) int {
	path, err := tvs.Model.GetPath(iter)
	if err != nil {
		fmt.Printf("Unable to get row number: %s\n", err.Error())
		return -1
	}
	ind := path.GetIndices()
	return ind[len(ind)-1:][0]
}

// AddRow: Append a row to the Store (defined by type of "store" variable).
// "parent" is useless for ListStore, if its set to nil on TreeStore,
// it will create a new parent
func (tvs *TreeViewStructure) AddRow(parent *gtk.TreeIter, row ...interface{}) (iter *gtk.TreeIter, err error) {
	return tvs.InsertRow(parent, -1, row...)
}

// InsertRow: Insert a row to the Store (defined by type of "store" variable).
// "parent" is useless for ListStore, if its set to nil on TreeStore,
// it will create a new parent otherwise the new iter will be a child of it.
// "insertPos" indicate row number for insertion, set to -1 mean append at the end.
func (tvs *TreeViewStructure) InsertRow(parent *gtk.TreeIter, insertPos int, row ...interface{}) (iter *gtk.TreeIter, err error) {

	iter = new(gtk.TreeIter)

	var colIdx = make([]int, len(row))
	for idx := 0; idx < len(row); idx++ {
		colIdx[idx] = idx
	}

	switch tvs.store.(type) {
	case *gtk.ListStore:
		err = tvs.ListStore.InsertWithValues(iter, insertPos, colIdx, row)
	case *gtk.TreeStore:
		err = tvs.TreeStore.InsertWithValues(iter, parent, insertPos, colIdx, row)
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to add row %d: %s\n", insertPos, err.Error())
	}
	tvs.Modified = true
	return
}

// InsertRowAtIter: Insert a row after/before iter to "store": ListStore/Treestore.
// Parent may be nil for Liststore.
func (tvs *TreeViewStructure) InsertRowAtIterN(inIter, parent *gtk.TreeIter, row []interface{}, before ...bool) (iter *gtk.TreeIter, err error) {
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
		return nil, errors.New(fmt.Sprintf("Unable to insert row: %s\n", err.Error()))
	}
	tvs.Modified = true
	return iter, err
}

// InsertRowAtIter: Insert a row after/before iter to "store": ListStore/Treestore.
// Parent may be nil for Liststore.
func (tvs *TreeViewStructure) InsertRowAtIter(inIter, parent *gtk.TreeIter, row ...interface{}) (iter *gtk.TreeIter, err error) {
	if iter, err = tvs.InsertRow(parent, tvs.GetRowNbIter(inIter)+1, row...); err == nil {
		tvs.Modified = true
	}
	return iter, err
}

// TODO rewrite to be much faster
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
		return nil, fmt.Errorf("Unable to duplicating row: %s\n", err.Error())
	}
	tvs.Modified = true
	tvs.ItersUnselect(inIter)
	tvs.ItersSelect(iter)
	return iter, err
}

// RemoveSelectedRows: Delete entries from selected iters or from given iters.
func (tvs *TreeViewStructure) RemoveSelectedRows(iters ...*gtk.TreeIter) (removed int, err error) {
	var ok bool
	if len(iters) == 0 {
		iters = tvs.GetSelectedIters()
	}
	for idx := len(iters) - 1; idx > -1; idx-- {
		if ok = tvs.RemoveRow(iters[idx]); !ok {
			err = fmt.Errorf("Unable to remove selected row: %d\n", idx)
			break
		}
		removed++
	}
	if ok {
		tvs.Modified = true
	}
	return
}

// GetSelectedRows: Get entries from selected iters as [][]string.
func (tvs *TreeViewStructure) GetSelectedRows() (outSlice [][]string, err error) {
	outSlice = make([][]string, tvs.Selection.CountSelectedRows())
	var count int
	tvs.Selection.SelectedForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) {
		if outSlice[count], err = tvs.GetRow(iter); err != nil {
			err = fmt.Errorf("Unable to get selected row: %d\n", count)
		}
		count++
	})
	return
}

// RemoveRow: Unified remove iter function
func (tvs *TreeViewStructure) RemoveRow(iter *gtk.TreeIter) (ok bool) {
	switch tvs.store.(type) {
	case *gtk.ListStore:
		ok = tvs.ListStore.Remove(iter)
	case *gtk.TreeStore:
		ok = tvs.TreeStore.Remove(iter)
	}
	if ok {
		tvs.Modified = true
	}
	return ok
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

/*****************************\
*    Convenient Functions     *
* Designed to make life easier *
*******************************/

// GetColumns: Retieve columns available in the current TreeView.
func (tvs *TreeViewStructure) GetColumns() (out []*gtk.TreeViewColumn) {
	glibList := tvs.TreeView.GetColumns()
	glibList.Foreach(func(item interface{}) {
		out = append(out, item.(*gtk.TreeViewColumn))
	})
	return
}

// StoreDetach: Unlink "TreeModel" from TreeView. Useful when lot of rows need
// be inserted. After insertion, StoreAttach() must be used to restore the link
// with the treeview. tips: must be used before ListStore/TreeStore.Clear().
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

// Clear: Clear the current used Model:
// unified version of gtk.TreeStore.Clear() or gtk.ListStore.Clear()
func (tvs *TreeViewStructure) Clear() {
	switch tvs.store.(type) {
	case *gtk.ListStore:
		tvs.ListStore.Clear()
	case *gtk.TreeStore:
		tvs.TreeStore.Clear()
	}
}

// ClearAll: Clear TreeView's columns, ListStore / TreeStore object.
// Depending on provided object type into the "store" variable.
// To reuse structure, you must execute StoreSetup() again after
// added new columns.
func (tvs *TreeViewStructure) ClearAll() (err error) {
	if tvs.TreeView != nil {
		// Removing existing columns if exists ...
		for idx := int(tvs.TreeView.GetNColumns()) - 1; idx > -1; idx-- {
			tvs.TreeView.RemoveColumn(tvs.TreeView.GetColumn(idx))
		}
		tvs.Columns = tvs.Columns[:0]
		tvs.TreeView.SetModel(nil)
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
		tvs.Modified = false
	}
	return
}

// StoreToSlice: Retrieve all the rows values from a "store" as [][]string
func (tvs *TreeViewStructure) StoreToStringSlice() (outSlice [][]string, err error) {
	var tmpSlice []string
	// Foreach Function
	var foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
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
	var retrieveValuesForeachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
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

// GetTree: get selected and unselected items.
// contained by the treestore
func (tvs *TreeViewStructure) GetTree() (checked, unChecked []string, err error) {
	var row []interface{}

	tvs.TreeStore.ForEach(func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
		row, err = tvs.GetRowIface(iter)

		if row[0].(bool) && len(row[2].(string)) > 0 {
			checked = append(checked, row[2].(string))
		} else {
			unChecked = append(unChecked, row[2].(string))
		}
		return false
	})
	return checked, unChecked, err
}

// AddTree: This function, add a full tree to a TreeStore, childs will be added to
// parent' tree if exists. Each treeview entry handle a checkbox and a name only. i.e:
// CASE 1: "pathSplitted" = []string{"github.com","hfmrow","gtk3Import","pango"},
// will add all nodes for this tree to column "pathCol" and a checkbox with
// "stateDefault" at column nbr "toggleCol".
// CASE 21: "pathSplitted" = []string{"true","github.com","hfmrow","gtk3Import","pango"},
// will add all nodes for this tree to column "pathCol" and a checkbox with
// "true" at column nbr "toggleCol".
// The returned variable "outIter", point to the iter of the last entry.
func (tvs *TreeViewStructure) AddTree(toggleCol, pathCol int, stateDefault bool, iFace ...interface{}) (outIter *gtk.TreeIter, err error) {
	var iterate, childsCount int
	var tmpIter *gtk.TreeIter
	var ok bool
	var value string
	var pathSplitted []interface{}

	// addItem: delete all entries and erase whole golang environment ... No, just kidding.
	// It only does what the name says.
	var addItem = func(toAdd string, iter *gtk.TreeIter) (tmpIter *gtk.TreeIter) {
		var err error
		tmpIter = tvs.TreeStore.Append(iter)
		if err = tvs.TreeStore.SetValue(tmpIter, toggleCol, stateDefault); err == nil {
			err = tvs.TreeStore.SetValue(tmpIter, pathCol, toAdd)
		}
		if err != nil {
			log.Fatalf("Unable to addItem: %s\n", err.Error())
			tmpIter = nil
		}
		return
	}

	// Check current entry type and do what we need with values (checked, notChecked or undefined)
	var getCurrentState = func() /*(outBool bool)*/ {
		switch iFace[0].(type) {
		case bool:
			stateDefault = iFace[0].(bool) // Get bool value for checkbox.
			pathSplitted = iFace[1:]       // Get full splitted path ingnoring 1st element.
		default:
			pathSplitted = iFace // Get full splitted path, checkbox value set as default.
		}
	}

	// findCreatFirstParent: Add or find first parent that match "name" and return it's iter.
	var findCreatFirstParent = func(name string) (tmpIter *gtk.TreeIter, ok bool, err error) {
		if len(name) > 0 {
			tmpIter, ok = tvs.TreeStore.GetIterFirst()
			for ok {
				value := tvs.GetColValue(tmpIter, pathCol).(string)
				if value == name {
					return
				}
				ok = tvs.TreeStore.IterNext(tmpIter)
			} // Nothing found, then create it
			tmpIter = addItem(pathSplitted[iterate].(string), outIter)
			ok = tmpIter != nil
		} else {
			err = errors.New("Could not proceed with empty parent.")
		}
		return
	}
	// searchMatch: Walk trought iter to retrieve the one matching "toMatch".
	var searchMatch = func(toMatch string, outIter *gtk.TreeIter) (childIter *gtk.TreeIter, ok bool) {
		childIter = new(gtk.TreeIter)
		ok = tvs.TreeStore.IterChildren(outIter, childIter)
		for ok {
			value = tvs.GetColValue(childIter, pathCol).(string)
			if value == toMatch {
				return
			}
			ok = tvs.Model.IterNext(childIter)
		}
		return
	}

	// Get toggle state
	getCurrentState()

	// parse "dirPath" entry into treestore.
	if outIter, ok, err = findCreatFirstParent(pathSplitted[0].(string)); err == nil {
		for ok {
			value = tvs.GetColValue(outIter, pathCol).(string)
			if value == pathSplitted[iterate] {
				childsCount = tvs.Model.IterNChildren(outIter)
				if childsCount > 0 {
					iterate++
					if iterate >= len(pathSplitted) {
						break
					}
					if tmpIter, ok = searchMatch(pathSplitted[iterate].(string), outIter); !ok {
						outIter = addItem(pathSplitted[iterate].(string), outIter)
						ok = outIter != nil
					} else {
						outIter = tmpIter
					}
					continue
				} else {
					iterate++
					if iterate >= len(pathSplitted) {
						break
					}
					outIter = addItem(pathSplitted[iterate].(string), outIter)
					ok = outIter != nil
				}
			} else {
				if tmpIter, ok = searchMatch(pathSplitted[iterate].(string), outIter); ok {
					outIter = addItem(pathSplitted[iterate].(string), tmpIter)
				}
			}
		}
	}
	return
}

// func (tvs *TreeViewStructure) AddTree(toggleCol, pathCol int, stateDefault bool, pathSplitted ...interface{}) (outIter *gtk.TreeIter, err error) {
// 	var iterate, childsCount int
// 	var tmpIter *gtk.TreeIter
// 	var ok bool
// 	var value string

// 	// addItem: delete all entries and erase whole golang environment ... No, just kidding.
// 	// It only does what the name says.
// 	var addItem = func(toAdd string, iter *gtk.TreeIter) (tmpIter *gtk.TreeIter) {
// 		var err error
// 		tmpIter = tvs.TreeStore.Append(iter)
// 		if err = tvs.TreeStore.SetValue(tmpIter, toggleCol, stateDefault); err == nil {
// 			err = tvs.TreeStore.SetValue(tmpIter, pathCol, toAdd)
// 		}
// 		if err != nil {
// 			log.Fatalf("Unable to addItem: %s\n", err.Error())
// 			tmpIter = nil
// 		}
// 		return
// 	}

// 	// Check current entry stat (checked, notChecked or undefined)
// 	var getCurrentState = func() bool {
// 		switch pathSplitted[0].(string) {
// 		case "true":
// 			pathSplitted = pathSplitted[1:]
// 			return true
// 		case "false":
// 			pathSplitted = pathSplitted[1:]
// 			return false
// 		}
// 		return stateDefault
// 	}

// 	// findCreatFirstParent: Add or find first parent that match "name" and return it's iter.
// 	var findCreatFirstParent = func(name string) (tmpIter *gtk.TreeIter, ok bool, err error) {
// 		stateDefault = getCurrentState()
// 		name = pathSplitted[0].(string)
// 		if len(name) > 0 {
// 			tmpIter, ok = tvs.TreeStore.GetIterFirst()
// 			for ok {
// 				value := tvs.GetColValue(tmpIter, pathCol).(string)
// 				if value == name {
// 					return
// 				}
// 				ok = tvs.TreeStore.IterNext(tmpIter)
// 			} // Nothing found, then create it
// 			tmpIter = addItem(pathSplitted[iterate].(string), outIter)
// 			ok = tmpIter != nil
// 		} else {
// 			err = errors.New("Could not proceed with empty parent.")
// 		}
// 		return
// 	}
// 	// searchMatch: Walk trought iter to retrieve the one matching "toMatch".
// 	var searchMatch = func(toMatch string, outIter *gtk.TreeIter) (childIter *gtk.TreeIter, ok bool) {
// 		childIter = new(gtk.TreeIter)
// 		ok = tvs.TreeStore.IterChildren(outIter, childIter)
// 		for ok {
// 			value = tvs.GetColValue(childIter, pathCol).(string)
// 			if value == toMatch {
// 				return
// 			}
// 			ok = tvs.Model.IterNext(childIter)
// 		}
// 		return
// 	}

// 	// parse "dirPath" entry into treestore.
// 	if outIter, ok, err = findCreatFirstParent(pathSplitted[0].(string)); err == nil {
// 		for ok {
// 			stateDefault = getCurrentState()
// 			value = tvs.GetColValue(outIter, pathCol).(string)
// 			if value == pathSplitted[iterate] {
// 				childsCount = tvs.Model.IterNChildren(outIter)
// 				if childsCount > 0 {
// 					iterate++
// 					if iterate >= len(pathSplitted) {
// 						break
// 					}
// 					if tmpIter, ok = searchMatch(pathSplitted[iterate].(string), outIter); !ok {
// 						outIter = addItem(pathSplitted[iterate].(string), outIter)
// 						ok = outIter != nil
// 					} else {
// 						outIter = tmpIter
// 					}
// 					continue
// 				} else {
// 					iterate++
// 					if iterate >= len(pathSplitted) {
// 						break
// 					}
// 					outIter = addItem(pathSplitted[iterate].(string), outIter)
// 					ok = outIter != nil
// 				}
// 			} else {
// 				if tmpIter, ok = searchMatch(pathSplitted[iterate].(string), outIter); ok {
// 					outIter = addItem(pathSplitted[iterate].(string), tmpIter)
// 				}
// 			}
// 		}
// 	}
// 	return
// }

// ColValuesStringSliceToIfaceSlice: Convert string list to []interface, for simplify adding text rows
func (tvs *TreeViewStructure) ColValuesStringSliceToIfaceSlice(inSlice ...string) (outIface []interface{}) {
	outIface = make([]interface{}, len(inSlice))
	for idx, data := range inSlice {
		outIface[idx] = data
	}
	return
}

// // ColValuesIfaceToIfaceSlice: Convert interface list to []interface, for simplify adding text rows
// func (tvs *TreeViewStructure) ColValuesIfaceToIfaceSlice(inSlice ...interface{}) (outIface []interface{}) {
// 	outIface = make([]interface{}, len(inSlice))
// 	for idx, data := range inSlice {
// 		outIface[idx] = data
// 	}
// 	return
// }

// glibType:  glib value type List structure.
var glibType = map[glib.Type]string{
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

// // getStringGlibType: Retrieve the string of glib value type.
// func (tvs *TreeViewStructure) getStringGlibType(t glib.Type) string {
// 	for val, str := range glibType {
// 		if val == int(t) {
// 			return str
// 		}
// 	}
// 	return "Unnowen type"
// }

/************************\
*    Helpers Functions    *
* Made to simplify some   *
* more complex functions  *
**************************/

// treeViewQueryTooltip: function to display tooltip according to rows currently hovered
func (tvs *TreeViewStructure) treeViewQueryTooltip(tw *gtk.TreeView, x, y int, KeyboardMode bool, tooltip *gtk.Tooltip) bool {
	if tvs.CountRows() > 0 && tvs.HasTooltip && tvs.CallbackTooltipFunc != nil {
		var (
			path    *gtk.TreePath
			column  *gtk.TreeViewColumn
			isBlank bool
		)
		// we must substract header height to "y" position to get the correct path.
		if path, column, _, _, isBlank = tvs.TreeView.IsBlankAtPos(x, y-tvs.getHeaderHeight()); !isBlank {
			if iter, err := tvs.Model.GetIter(path); err == nil {
				return tvs.CallbackTooltipFunc(iter, path, column, tooltip)
			} else {
				log.Printf("treeViewQueryTooltip:GetIter: %s\n", err.Error())
			}
		}
	}
	return false
}

// getHeaderHeight: Used to get height of header to use with [TreeView.IsBlankAtPos],
// It is needed to decrease y pos by height of cells to get a correct path value.
func (tvs *TreeViewStructure) getHeaderHeight() (height int) {
	if tvs.headerHeight < 0 { // That mean that is launched only at the first call.
		for gtk.EventsPending() {
			gtk.MainIteration() // Wait for pending events (until the widget is redrawn)
		}
		// Getting header height
		backupVisibleHeader := tvs.TreeView.GetHeadersVisible()
		tvs.TreeView.SetHeadersVisible(true)
		withHeader, _ := tvs.TreeView.GetPreferredHeight()
		tvs.TreeView.SetHeadersVisible(false)
		withoutHeader, _ := tvs.TreeView.GetPreferredHeight()
		tvs.TreeView.SetHeadersVisible(backupVisibleHeader)
		tvs.headerHeight = withHeader - withoutHeader
	}
	return tvs.headerHeight
}

// getCellValueByType: Retrieve cell value and convert it to string based on
// his type. Used by GetRow func.
func (tvs *TreeViewStructure) getStringCellValueByType(glibValue *glib.Value) (valueString string, err error) {
	var actualType glib.Type
	var valueIface interface{}

	if actualType, _, err = glibValue.Type(); err == nil {
		switch actualType {

		// Strings
		case glib.TYPE_STRING:
			valueString, err = glibValue.GetString()

			// Numeric values
		case glib.TYPE_INT64, glib.TYPE_UINT64, glib.TYPE_INT:
			if valueIface, err = glibValue.GoValue(); err == nil {
				switch val := valueIface.(type) {
				case int, uint64, int64:
					valueString = fmt.Sprintf("%d", val)
				}
			}

			// Pointer, just say that it's what's it ...
		case glib.TYPE_POINTER:
			valueString = "pointer"

			// Boolean
		case glib.TYPE_BOOLEAN:
			if valueIface, err = glibValue.GoValue(); err == nil {
				if valueIface.(bool) {
					valueString = "true"
				} else {
					valueString = "false"
				}
			}

			// Need to be implemented
		default:
			err = fmt.Errorf("getStringCellValueByType: Type %s, not yet implemented\n", glibType[actualType])
		}
	}
	return
}

// changeTreeStateBool: Modify the state of the entire tree starting at parent.
func (tvs *TreeViewStructure) changeTreeStateBool(parent *gtk.TreeIter, col int, goValue interface{}) (ok bool) {
	var err error
	childIter := new(gtk.TreeIter)
	ok = parent != nil
	for ok {
		if err = tvs.SetColValue(parent, col, goValue); err == nil {
			if ok = tvs.Model.IterHasChild(parent); ok {
				if ok = tvs.Model.IterChildren(parent, childIter); ok {
					for ok {
						if ok = tvs.changeTreeStateBool(childIter, col, goValue); !ok {
							ok = tvs.Model.IterNext(childIter)
						}
					}
				}
			}
		}
	}
	return
}

/********************\
*    TEST Functions   *
* Not designed to be  *
* used as it !!       *
**********************/

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

func (tvs *TreeViewStructure) selectRange(start, end *gtk.TreeIter) (err error) {
	var startPath, endPath *gtk.TreePath
	if startPath, err = tvs.ListStore.GetPath(start); err == nil {
		if endPath, err = tvs.ListStore.GetPath(end); err == nil {
			tvs.Selection.SelectRange(startPath, endPath)
		}
	}
	return err
}

func (tvs *TreeViewStructure) pathSelected(start *gtk.TreeIter) (err error) {
	var startPath *gtk.TreePath
	if startPath, err = tvs.ListStore.GetPath(start); err == nil {
		fmt.Println("iter", tvs.Selection.IterIsSelected(start))
		fmt.Println("path", tvs.Selection.PathIsSelected(startPath))
	}
	return err
}

func (tvs *TreeViewStructure) forEach() {
	var err error
	var model gtk.ITreeModel
	var ipath *gtk.TreePath
	var foreachFunc gtk.TreeModelForeachFunc
	foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
		if ipath, err = model.GetPath(iter); err == nil {
			fmt.Printf("path: %s, iter: %s\n", path.String(), ipath.String())
		} else {
			fmt.Println("error occured inside func: " + err.Error())
			return true
		}
		return false
	}
	if model, err = tvs.TreeView.GetModel(); err == nil {
		model.ToTreeModel().ForEach(foreachFunc)
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

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

func (tvs *TreeViewStructure) indices() {
	var err error
	var model gtk.ITreeModel
	var ipath, jpath *gtk.TreePath
	var foreachFunc gtk.TreeModelForeachFunc
	foreachFunc = func(model *gtk.TreeModel, path *gtk.TreePath, iter *gtk.TreeIter, userData ...interface{}) bool {
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
		model.ToTreeModel().ForEach(foreachFunc)
	}
	if err != nil {
		fmt.Println("error occured outside func: " + err.Error())
	}
}

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
