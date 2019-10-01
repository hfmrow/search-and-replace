// cssWdgScnLoad.go

/*
	Load or read from data css style for an object(widget) or for entire screen.
*/

package gtk3Import

import (
	"fmt"
	"io/ioutil"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

// CssWidgetLoad: Load or read from data and apply css to
// widget if it's provided. Apply to screen otherwise.
func CssWdgScnLoad(filename string, wdgt ...*gtk.Widget) {
	var err error
	var cssProv *gtk.CssProvider

	if bytes, err := ioutil.ReadFile(filename); err == nil {
		filename = string(bytes)
	}

	if cssProv, err = gtk.CssProviderNew(); err == nil {

		if err = cssProv.LoadFromData(filename); err == nil {
			if len(wdgt) == 0 {
				var screen *gdk.Screen
				if screen, err = gdk.ScreenGetDefault(); err == nil {
					gtk.AddProviderForScreen(screen, cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
				}
			} else {
				var styleContext *gtk.StyleContext
				if styleContext, err = wdgt[0].GetStyleContext(); err == nil {
					styleContext.AddProvider(cssProv, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
				}
			}
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}

/*
 #MainTextViewEditorText * {
	color: shade (#332211, 1.06);
	background-color: #fefefe;
	opacity: 0.99;
}

 #MainTextViewEditorText text selection {
	background-color: #aaddff;
	color:shade (#332211, 1.06);
}

 #MainTextViewEditorNumbers * {
	color: shade (#0033ff, 1.06);
	background-color: #eeeeee;
	opacity: 0.99;
}

 #MainTextViewEditorNumbers text {
	color: shade (#0022ff, 1.06);
	background-color: #eeeeee;
	opacity: 0.99;
}
*/
