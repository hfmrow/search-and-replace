// popupMenu.go

/*
	©2019 H.F.M. MIT license

	Make popup menu menu
*/

/* * * * * *
* Usage: *
* * * * *
func initPopup() {
	mainOptions.popupMenu = new(PopupMenu)
	mainOptions.popupMenu.WithIcons = false
	mainOptions.popupMenu.LMB = false
	mainOptions.popupMenu.PopupAddItem("_small", func() { assignTagToolButton("small") }, "")
	mainOptions.popupMenu.PopupAddSeparator()
	mainOptions.popupMenu.PopupAddItem("_medium", func() { assignTagToolButton("medium") }, "")
	mainOptions.popupMenu.PopupAddSeparator()
	mainOptions.popupMenu.PopupAddItem("_large", func() { assignTagToolButton("large") }, "")
	mainOptions.popupMenu.PopupMenuBuild()
}

SIGNAL:

	obj.Connect("button-press-event", ObjectButtonReleaseEvent)

func ObjectButtonReleaseEvent(obj interface{}, event *gdk.Event) bool {
	return mainOptions.PopupCheckRMB(event)
}
*/

package gtk3Import

import (
	"log"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"

	p "github.com/hfmrow/gtk3Import/pixbuff"
)

type PopupMenu struct {
	Menu       *gtk.Menu
	WithIcons  bool // Adding icon or not
	LMB        bool // left mouse button instead of right
	items      []*gtk.MenuItem
	separators []*gtk.SeparatorMenuItem
}

// PopupCheckRMB: Check if event come from right mouse button, and display popup if it is.
func (pop *PopupMenu) PopupCheckRMB(event *gdk.Event) bool {
	mouseBTN := uint(3)
	if pop.LMB {
		mouseBTN = 1
	}
	eventButton := gdk.EventButtonNewFromEvent(event)
	if eventButton.Button() == mouseBTN {
		pop.Menu.PopupAtPointer(event)
		return true // return true to stop propagate event.
	}
	return false // return false to propagate event.
}

// PopupAddItem: Add items to menu.
func (pop *PopupMenu) PopupAddItem(lbl string, activateFunction interface{}, icon ...interface{}) (err error) {
	var menuItem *gtk.MenuItem
	var image interface{}
	if len(icon) != 0 {
		image = icon[0]
	}
	if pop.WithIcons {
		menuItem, err = pop.menuItemNewWithImage(lbl, image)
	} else {
		menuItem, err = gtk.MenuItemNewWithMnemonic(lbl)
	}
	if err == nil {
		menuItem.Connect("activate", activateFunction.(func()))
		pop.items = append(pop.items, menuItem)
		pop.separators = append(pop.separators, nil)
	}
	return err
}

// PopupAddSeparator: Add separator to menu.
func (pop *PopupMenu) PopupAddSeparator() (err error) {
	if separatorItem, err := gtk.SeparatorMenuItemNew(); err == nil {
		pop.items = append(pop.items, nil)
		pop.separators = append(pop.separators, separatorItem)
	}
	return err
}

// PopupMenuBuild: Build popupmenu.
func (pop *PopupMenu) PopupMenuBuild() *gtk.Menu {
	var err error
	if pop.Menu, err = gtk.MenuNew(); err == nil {
		for idx, menuItem := range pop.items {
			if pop.separators[idx] != nil {
				pop.Menu.Append(pop.separators[idx])
			} else {
				pop.Menu.Append(menuItem)
			}
		}
		pop.Menu.ShowAll()
	} else {
		log.Println("Popup menu creation error !")
		return nil
	}
	return pop.Menu
}

// menuItemNewWithImage: Build an item containing an image.
func (pop *PopupMenu) menuItemNewWithImage(label string, icon interface{}) (menuItem *gtk.MenuItem, err error) {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	if err == nil {
		image, err := gtk.ImageNew()
		if err == nil {
			p.SetImage(image, icon, 14)
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
