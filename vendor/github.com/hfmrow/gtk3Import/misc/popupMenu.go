// popupMenu.go

/*
	©2019 H.F.M. MIT license

	A simple builder for popup menu that may handle icons.

i.e:
	var popupMenu *PopupMenu // global declaration

	func initPopup() {
		popupMenu = PopupMenuNew()
		popupMenu.WithIcons = true
		popupMenu.LMB = false
		popupMenu.IconSize = 18 // if wanted, must be set before adding menu
		popupMenu.AddItem("_small", func() { assignTagToolButton("small") }, image1)
		popupMenu.AddSeparator()
		popupMenu.AddItem("_medium", func() { assignTagToolButton("medium") }, image2)
		popupMenu.AddItem("_large", func() { assignTagToolButton("large") }, image3)
		popupMenu.MenuBuild()
	}

Signal:
	obj.Connect("button-press-event", ObjectButtonPressEvent)
Callback:
	func TreeViewFoundButtonPressEvent(obj interface{}, event *gdk.Event) bool {
		popupMenu.CheckRMB(event)
		return false // Propagate event
	}

- May be used to append items to an existing gtk.Menu using: AppendToExistingMenu() method
  useful for textview with his existing context menu.
i.e:
Signal:
	TextView.Connect("populate-popup", popupTextViewPopulateMenu)

Callback:
	// popupTextViewPopulateMenu: Append some items to contextmenu of the TextView
	func popupTextViewPopulateMenu(txtView *gtk.TextView, popup *gtk.Widget) {
		// Convert gtk.Widget to gtk.Menu object
		pop := &gtk.Menu{gtk.MenuShell{gtk.Container{*popup}}}
		// create new menuitems
		popMenuTextView = PopupMenuNew()
		popMenuTextView.WithIcons = true
		popMenuTextView.AddSeparator()
		popMenuTextView.AddItem("Open _directory", func() { openDir() }, image1)
		popMenuTextView.AddItem("Open _file", func() { openFile() }, image2)
		// Append them to the existing menu
		popMenuTextView.AppendToExistingMenu(pop)
	}
*/

package gtk3Import

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	gipf "github.com/hfmrow/gtk3Import/pixbuff"
)

// PopupMenu: Structure that hold popup menu options and methods.
// A simple builder for popup menu that may handle icons.
type PopupMenu struct {
	Menu            *gtk.Menu
	WithIcons       bool // Adding icon or not
	IconsSize       int
	LMB             bool          // left mouse button instead of right
	CurrentUserData []interface{} // Hold current userData on RMB click (from caller object)

	items      []*gtk.MenuItem
	separators []*gtk.SeparatorMenuItem
}

// PopupMenuNew: return a new PopupMenu structure
func PopupMenuNew() (pop *PopupMenu) {
	pop = new(PopupMenu)
	pop.IconsSize = 14
	return
}

// CheckRMB: Check whether an event comes from the right button of the
// mouse and display the popup if it is the case at the mouse position.
func (pop *PopupMenu) CheckRMB(event *gdk.Event, userData ...interface{}) bool {
	eventButton := gdk.EventButtonNewFromEvent(event)
	if uint(eventButton.Button()) == pop.mouseBtn() {
		pop.CurrentUserData = userData
		pop.Menu.PopupAtPointer(event)
		return true
	}
	return false
}

// CheckRMBFromTreeView: May be directly used as callback function for
// TreeView' "button-press-event" signal, considere to initialize the
// popup menu before setting this function as callback. Otherwise, the
// call will generate error "nil pointer ..."
func (pop *PopupMenu) CheckRMBFromTreeView(tv *gtk.TreeView, event *gdk.Event) bool {
	if selection, err := tv.GetSelection(); err == nil {
		if selected := selection.CountSelectedRows(); selected > 0 {
			eventButton := gdk.EventButtonNewFromEvent(event)
			if uint(eventButton.Button()) == pop.mouseBtn() {
				// pop.CurrentUserData = userData
				pop.Menu.PopupAtPointer(event)
				if selected > 1 {
					return true
				}
			}
		}
	}
	return false
}

// AddItem: Add items to menu.
func (pop *PopupMenu) AddItem(lbl string, activateFunction interface{},
	icon ...interface{}) (err error) {

	var menuItem *gtk.MenuItem
	var pixbuf *gdk.Pixbuf

	if len(icon) != 0 {
		// The function below is a part of personal gotk3 library that
		// allow to load image with some facilities. May handle
		// filename or embedded binary data (hex/zip compressed).
		// pixbuf, err = gdk.PixbufNewFromFile(filename)
		pixbuf, err = gipf.GetPixBuf(icon[0], pop.IconsSize)
	}

	if pop.WithIcons {
		menuItem, err = pop.menuItemNewWithImage(lbl, pixbuf)
	} else {
		menuItem, err = gtk.MenuItemNewWithMnemonic(lbl)
	}
	// Handle the "activate" signal from the related item.
	if err == nil {
		menuItem.Connect("activate", activateFunction.(func()))
		pop.items = append(pop.items, menuItem)
		pop.separators = append(pop.separators, nil)
	}
	return err
}

// AddSeparator: Add separator to menu.
func (pop *PopupMenu) AddSeparator() (err error) {
	if separatorItem, err := gtk.SeparatorMenuItemNew(); err == nil {
		pop.items = append(pop.items, nil)
		pop.separators = append(pop.separators, separatorItem)
	}
	return err
}

// MenuBuild: Build popupmenu.
func (pop *PopupMenu) MenuBuild() *gtk.Menu {
	var err error
	if pop.Menu, err = gtk.MenuNew(); err == nil {
		for idx, menuItem := range pop.items {
			if pop.separators[idx] != nil {
				pop.Menu.Append(pop.separators[idx])
			} else {
				pop.Menu.Append(menuItem)
			}
		}
		pop.Menu.Connect("move-focus", func(menu *gtk.Menu, event *gdk.Event) {
			pop.Menu.Hide()
			fmt.Println(menu.GetVisible())
		})

		pop.Menu.ShowAll()
	} else {
		log.Println("Popup menu creation error !")
		return nil
	}
	return pop.Menu
}

// AppendToExistingMenu: append "MenuItems" to an existing "*gtk.Menu"
// Useful when you want to just add some entries to the context menu that
// already exist in a gtk.TextView or gtk.Entry by using "populate-popup"
// signal.
func (pop *PopupMenu) AppendToExistingMenu(menu *gtk.Menu) *gtk.Menu {
	for idx, menuItem := range pop.items {
		if pop.separators[idx] != nil {
			menu.Append(pop.separators[idx])
		} else {
			menu.Append(menuItem)
		}
	}
	menu.ShowAll()
	return menu
}

// menuItemNewWithImage: Build an item containing an image.
func (pop *PopupMenu) menuItemNewWithImage(label string,
	pixbuf *gdk.Pixbuf) (menuItem *gtk.MenuItem, err error) {
	var image *gtk.Image
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	if err == nil {
		if image, err = gtk.ImageNewFromPixbuf(pixbuf); err == nil {
			label, err := gtk.LabelNewWithMnemonic(label)
			if err == nil {
				menuItem, err = gtk.MenuItemNew()
				if err == nil {
					label.SetHAlign(gtk.ALIGN_START)
					box.Add(image)
					box.PackEnd(label, true, true, 8)
					box.SetHAlign(gtk.ALIGN_START)
					menuItem.Container.Add(box)
					menuItem.ShowAll()
				}
			}
		}
	}
	return menuItem, err
}

// mouseBtn: get uint value of specified button to match
func (pop *PopupMenu) mouseBtn() uint {
	if pop.LMB {
		return 1 // LMB
	}
	return 3 // RMB
}
