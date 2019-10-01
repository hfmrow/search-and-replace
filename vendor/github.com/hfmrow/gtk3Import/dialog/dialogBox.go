// dialogBox.go

/*
  Create a Dialog, who accept GtkWidgets defined into structure.
  The Structure contain all needed options to fill most of usages.

Usage:
	if tw, err := gtk.TreeViewNew(); err == nil {
		dbs := gi.DialogBoxNew(MainWindow, gtk.DIALOG_DESTROY_WITH_PARENT, tw, "test Title", "No", "Yes", "why")
		dbs.ButtonsImages = dbs.ValuesToInterfaceSlice("assets/images/Sign-cancel-20.png", "", signSelect20) // "signSelect20" is []byte of image
		dbs.ScrolledArea = true
		dbs.SetModal = false
		result=dbs.Run()
		// Do what you want with "result"
	}
*/

package gtk3Import

import (
	"fmt"
	"log"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"

	p "github.com/hfmrow/gtk3Import/pixbuff"
)

// DialogBoxStructure: Wrap a Dialog with desired count of
// buttons and widgets. The structure have defaults parameters
type DialogBoxStructure struct {
	BoxHAlign       gtk.Align
	BoxVAlign       gtk.Align
	BoxOrientation  gtk.Orientation
	BoxHExpand      bool
	BoxVExpand      bool
	WidgetExpend    bool
	WidgetFill      bool
	HSize, VSize    int
	SkipTaskbarHint bool
	KeepAbove       bool
	Flag            gtk.DialogFlags
	SetModal        bool
	ScrolledArea    bool
	Buttons         []string
	ButtonsImages   []interface{} // image representation from file or []byte, depending on type
	Widgets         []gtk.IWidget
	Title, Text     string
	MarkupLabel     bool
	LabelLineWrap   bool
	Padding         uint

	buttonsWithImages bool
	dialog            *gtk.Dialog
	label             *gtk.Label
	box               *gtk.Box
	scrolledWindow    *gtk.ScrolledWindow
	window            *gtk.Window
}

// DialogBoxNew: Create a new structure to wrap a GtkDialog including
// defaults parameters. "widget" can be "nul" to only display "text".
func DialogBoxNew(window *gtk.Window, flag gtk.DialogFlags, widget gtk.IWidget, title string, buttons ...string) *DialogBoxStructure {
	dbs := new(DialogBoxStructure)
	dbs.window = window
	dbs.Flag = flag
	dbs.BoxHAlign = gtk.ALIGN_FILL
	dbs.BoxVAlign = gtk.ALIGN_FILL
	dbs.BoxOrientation = gtk.ORIENTATION_VERTICAL
	dbs.BoxHExpand = true
	dbs.BoxVExpand = true
	dbs.HSize, dbs.VSize = 640, 480
	dbs.SkipTaskbarHint = true
	dbs.KeepAbove = true
	dbs.Flag = flag
	dbs.SetModal = false
	dbs.Title = title
	dbs.LabelLineWrap = true
	dbs.WidgetExpend = true
	dbs.WidgetFill = true
	dbs.Padding = 0
	if widget != nil {
		dbs.Widgets = append(dbs.Widgets, widget)
	}
	if len(buttons) == 0 {
		dbs.Buttons = []string{"Ok"}
	} else {
		dbs.Buttons = buttons
	}
	return dbs
}

// SliceToInterface: Convert String slice to interface, for simplify adding text rows
func (dbs *DialogBoxStructure) ValuesToInterfaceSlice(inSlice ...interface{}) (outIface []interface{}) {
	for _, value := range inSlice {
		outIface = append(outIface, value)
	}
	return
}

// Run: calling Run function return "value" < 0 for cross closed or >= 0
// corresponding to buttons order representation starting with 0 at left.
func (dbs *DialogBoxStructure) Run() (value int) {
	var btnObj *gtk.Button
	var err error
	// Build Objects

	// if dbs.dialog, err = gtk.DialogNewWithButtons(dbs.Title, dbs.window, dbs.Flag); err == nil { // is waiting merging, pull request #425
	if dbs.dialog, err = gtk.DialogNew(); err == nil {
		dbs.dialog.SetTransientFor(dbs.window)
		dbs.dialog.SetModal(dbs.SetModal)

		if dbs.label, err = gtk.LabelNew(""); err == nil {
			dbs.box, err = dbs.dialog.GetContentArea()
		}
	}
	if err != nil {
		log.Fatalf("Enable to create dialog: %s\n", err.Error())
	}

	// Markup & Label options
	dbs.label.SetSizeRequest(dbs.box.GetSizeRequest())
	dbs.label.SetLineWrap(dbs.LabelLineWrap)
	dbs.label.SetLineWrapMode(pango.WRAP_WORD)
	if dbs.MarkupLabel {
		dbs.label.SetLabel(dbs.Text)
	} else {
		dbs.label.SetMarkup(dbs.Text)
	}

	// Control
	dbs.buttonsWithImages = len(dbs.ButtonsImages) != 0
	if len(dbs.Buttons) != len(dbs.ButtonsImages) && dbs.buttonsWithImages {
		log.Fatalf("You must provide an image or an empty string for each button\nButton(s) count: %d, Image(s) count: %d\n",
			len(dbs.Buttons), len(dbs.ButtonsImages))
	}

	// Dialog options
	dbs.dialog.SetTitle(dbs.Title)
	dbs.dialog.SetSkipTaskbarHint(dbs.SkipTaskbarHint)
	dbs.dialog.SetKeepAbove(dbs.KeepAbove)

	// Box options
	dbs.box.SetHAlign(dbs.BoxHAlign)
	dbs.box.SetVAlign(dbs.BoxVAlign)
	dbs.box.SetOrientation(dbs.BoxOrientation)
	dbs.box.SetHExpand(dbs.BoxHExpand)
	dbs.box.SetVExpand(dbs.BoxVExpand)
	dbs.box.SetSizeRequest(dbs.HSize, dbs.VSize)

	// Buttons
	for idxBtn, btnLbl := range dbs.Buttons {
		if btnObj, err = dbs.dialog.AddButton(btnLbl, gtk.ResponseType(idxBtn)); err == nil {
			if dbs.buttonsWithImages && len(dbs.ButtonsImages[idxBtn].(string)) != 0 {
				p.SetButtonImage(btnObj, dbs.ButtonsImages[idxBtn])
			}
		}
	}

	// Packing
	if len(dbs.Text) != 0 {
		dbs.box.PackStart(dbs.label, dbs.WidgetExpend, dbs.WidgetFill, dbs.Padding)
	}
	if len(dbs.Widgets) != 0 {
		if dbs.ScrolledArea {
			if dbs.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil); err == nil {
				for _, wdg := range dbs.Widgets {
					dbs.scrolledWindow.Add(wdg)
				}
				dbs.box.PackStart(dbs.scrolledWindow, dbs.WidgetExpend, dbs.WidgetFill, dbs.Padding)
			}
		} else {
			for _, wdg := range dbs.Widgets {
				dbs.box.PackStart(wdg, dbs.WidgetExpend, dbs.WidgetFill, dbs.Padding)
			}
		}
	}
	dbs.box.ShowAll()
	dbs.dialog.ShowAll()
	// The show must go on
	value = int(dbs.dialog.Run())
	dbs.dialog.Destroy()
	return value
}

/**********************/
/* PixBuff functions */
/********************/

// setButtonImage: Set Icon to GtkButton objects
func SetButtonImage(object *gtk.Button, varPath interface{}, size ...int) {
	var image *gtk.Image
	inPixbuf, err := GetPixBuff(varPath, size...)
	if err == nil {
		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
			object.SetImage(image)
			object.SetAlwaysShowImage(true)
			return
		}
	}
	if err != nil && len(varPath.(string)) != 0 {
		fmt.Printf("SetButtonImage: An error occurred on image: %s\n%v\n", varPath, err.Error())
	}
}

// GetPixBuff: Get gtk.Pixbuff image representation from file or []byte, depending on type
// size: resize height keeping porportions. 0 = no change
func GetPixBuff(varPath interface{}, size ...int) (outPixbuf *gdk.Pixbuf, err error) {
	sze := 0
	if len(size) != 0 {
		sze = size[0]
	}
	switch reflect.TypeOf(varPath).String() {
	case "string":
		outPixbuf, err = gdk.PixbufNewFromFile(varPath.(string))
	case "[]uint8":
		pbLoader, err := gdk.PixbufLoaderNew()
		if err == nil {
			outPixbuf, err = pbLoader.WriteAndReturnPixbuf(varPath.([]byte))
		}
	}
	if err == nil && sze != 0 {
		newWidth, wenHeight := normalizeSize(outPixbuf.GetWidth(), outPixbuf.GetHeight(), sze, 2)
		outPixbuf, err = outPixbuf.ScaleSimple(newWidth, wenHeight, gdk.INTERP_BILINEAR)
	}
	return outPixbuf, err
}

// NormalizeSize: compute new size with kept proportions based on defined format.
// format: 0 percent, 1 reducing width, 2 reducing height
func normalizeSize(oldWidth, oldHeight, newValue, format int) (outWidth, outHeight int) {
	switch format {
	case 0: // percent
		outWidth = int((float64(oldWidth) * float64(newValue)) / 100)
		outHeight = int((float64(oldHeight) * float64(newValue)) / 100)
	case 1: // Width
		outWidth = newValue
		outHeight = int(float64(oldHeight) * (float64(newValue) / float64(oldWidth)))
	case 2: // Height
		outWidth = int(float64(oldWidth) * (float64(newValue) / float64(oldHeight)))
		outHeight = newValue
	}
	return outWidth, outHeight
}
