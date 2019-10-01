// pixbuffByte.go

// Source file auto-generated on Thu, 21 Feb 2019 00:54:07 using Gotk3ObjHandler v1.0 ©2019 H.F.M

/*
	©2019 H.F.M

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*************************************************/
/* Images functions, used to initialize objects */
/***********************************************/
// setImage: Set Image to GtkImage objects
func SetImage(object *gtk.Image, varPath interface{}, size ...int) {
	if inPixbuf, err := GetPixBuff(varPath, size...); err == nil {
		object.SetFromPixbuf(inPixbuf)
		return
	} else if len(varPath.(string)) != 0 {
		fmt.Printf("SetImage: An error occurred on image: %s\n%v\n", varPath, err.Error())
	}
}

// setWinIcon: Set Icon to GtkWindow objects
func SetWinIcon(object *gtk.Window, varPath interface{}, size ...int) {
	if inPixbuf, err := GetPixBuff(varPath, size...); err == nil {
		object.SetIcon(inPixbuf)
	} else if len(varPath.(string)) != 0 {
		fmt.Printf("SetWinIcon: An error occurred on image: %s\n%v\n", varPath, err.Error())
	}
}

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

// setToolButtonImage: Set Icon to GtkToolButton objects
func SetToolButtonImage(object *gtk.ToolButton, varPath interface{}, size ...int) {
	var image *gtk.Image
	inPixbuf, err := GetPixBuff(varPath, size...)
	if err == nil {
		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
			object.SetIconWidget(image)
			return
		}
	}
	if err != nil && len(varPath.(string)) != 0 {
		fmt.Printf("setToolButtonImage: An error occurred on image: %s\n%v\n", varPath, err.Error())
	}
}

// setToggleButtonImage: Set Icon to GtkToggleButton objects
func SetToggleButtonImage(object *gtk.ToggleButton, varPath interface{}, size ...int) {
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
		fmt.Printf("SetToggleButtonImage: An error occurred on image: %s\n%v\n", varPath, err.Error())
	}
}

// SetSpinButtonImage: Set Icon to GtkSpinButton objects. Position = "left" or "right"
func SetSpinButtonImage(object *gtk.SpinButton, varPath interface{}, position ...string) {
	var inPixbuf *gdk.Pixbuf
	var err error
	pos := gtk.ENTRY_ICON_PRIMARY
	if len(position) > 0 {
		if position[0] == "right" {
			pos = gtk.ENTRY_ICON_SECONDARY
		}
	}
	if inPixbuf, err = GetPixBuff(varPath); err == nil {
		object.SetIconFromPixbuf(pos, inPixbuf)
		return
	} else if len(varPath.(string)) != 0 {
		fmt.Printf("SetSpinButtonImage: An error occurred on image: %s\n%v\n", varPath, err.Error())
	}
}

// setBoxImage:  Set Image to GtkBox objects
func SetBoxImage(object *gtk.Box, varPath interface{}, size ...int) {
	var image *gtk.Image
	inPixbuf, err := GetPixBuff(varPath, size...)
	if err == nil {
		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
			image.Show()
			object.Add(image)
			return
		}
	}
	if err != nil && len(varPath.(string)) != 0 {
		fmt.Printf("setBoxImage: An error occurred on image: %s\n%v\n", varPath, err.Error())
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

// ResizeImage: Get Resized gtk.Pixbuff image representation from file or []byte, depending on type
// interp: 0 GDK_INTERP_NEAREST, 1 GDK_INTERP_TILES, 2 GDK_INTERP_BILINEAR (default), 3 GDK_INTERP_HYPER.
func ResizeImage(varPath interface{}, width, height int, interp ...int) (outPixbuf *gdk.Pixbuf, err error) {
	interpolation := gdk.INTERP_BILINEAR
	if len(interp) != 0 {
		switch interp[0] {
		case 0:
			interpolation = gdk.INTERP_NEAREST
		case 1:
			interpolation = gdk.INTERP_TILES
		case 3:
			interpolation = gdk.INTERP_HYPER
		}
	}
	if outPixbuf, err = GetPixBuff(varPath); err == nil {
		if width != outPixbuf.GetWidth() || height != outPixbuf.GetHeight() {
			return outPixbuf.ScaleSimple(width, height, interpolation)
		}
	}
	return nil, err
}

// RotateImage: Rotate by 90,180,270 degres and get gtk.Pixbuff image representation from file or []byte, depending on type
func RotateImage(varPath interface{}, angle gdk.PixbufRotation) (outPixbuf *gdk.Pixbuf, err error) {
	if outPixbuf, err = GetPixBuff(varPath); err == nil {
		switch angle {
		case 90:
			return outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_COUNTERCLOCKWISE)
		case 180:
			return outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_UPSIDEDOWN)
		case 270:
			return outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_CLOCKWISE)
		default:
			return nil, errors.New("Rotation options not allowed: " + fmt.Sprintf("%d", angle))
		}
	}
	return outPixbuf, err
}

// FlipImage: Get Flipped gtk.Pixbuff image representation from file or []byte, depending on type
func FlipImage(varPath interface{}, horizontal bool) (outPixbuf *gdk.Pixbuf, err error) {
	if outPixbuf, err = GetPixBuff(varPath); err == nil {
		return outPixbuf.Flip(horizontal)
	}
	return nil, err
}
