// pixbuffByte.go

/*
	Source file auto-generated on Sun, 20 Oct 2019 13:50:31 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	©2019 H.F.M

	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package gtk3Import

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

/*************************************************/
/* Images functions, used to initialize objects */
/* You can use it to load your own embedded    */
/* images, icons ...                          */
/*********************************************/

// SetPic: Assign image to an Object depending on type, accept filename or []byte.
// options: 1- size (int), 2- enable animation (bool)
// ie:
// - Load a gif animated image and specify using animation (true), resizing not allowed with animations.
//     SetPict(GtkImage, "linearProgressHorzBlue.gif", 0, true)
// - Resize to 32 pixels height, keep porportions & assign image to GtkButton.
//     SetPict(GtkButton, "stillImage.png", 32)
// - With default size, resizing not allowed for GtkSpinButton.
//     SetPict(GtkSpinButton, "stillImage.png")
func SetPict(iObject, varPath interface{}, options ...interface{}) {
	var err error
	var sze int
	var inPixbuf *gdk.Pixbuf
	var inPixbufAnimation *gdk.PixbufAnimation
	var image *gtk.Image
	var isAnim bool
	pos := gtk.ENTRY_ICON_PRIMARY
	// Options parsing
	var getSizeOrPos = func() {
		switch s := options[0].(type) {
		case int:
			sze = s
		case string:
			if s == "right" {
				pos = gtk.ENTRY_ICON_SECONDARY
			}
		}
	}
	switch len(options) {
	case 1:
		getSizeOrPos()
	case 2: //
		getSizeOrPos()
		isAnim = options[1].(bool)
	}
	if isAnim { // PixbufAnimation case
		if inPixbufAnimation, err = GetPixBufAnimation(varPath); err == nil {
			image, err = gtk.ImageNewFromAnimation(inPixbufAnimation)
		}
	} else { // Pixbuf case
		if inPixbuf, err = GetPixBuf(varPath, sze); err == nil {
			image, err = gtk.ImageNewFromPixbuf(inPixbuf)
		}
	}
	if err != nil {
		if _, err = os.Stat(varPath.(string)); !os.IsNotExist(err) {
			log.Fatalf("SetPict: %v\n%s\n", varPath, err.Error())
			return
		}
	}
	// Objects parsing
	if image != nil {
		switch object := iObject.(type) {
		case *gtk.Image: // Set Image to GtkImage
			if isAnim {
				object.SetFromAnimation(inPixbufAnimation)
			} else {
				object.SetFromPixbuf(inPixbuf)
			}
		case *gtk.Window: // Set Icon to GtkWindow, No Animate.
			object.SetIcon(inPixbuf)
		case *gtk.Button: // Set Image to GtkButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)
		case *gtk.ToolButton: // Set Image to GtkToolButton
			object.SetIconWidget(image)
		case *gtk.ToggleButton: // Set Image to GtkToggleButton
			object.SetImage(image)
			object.SetAlwaysShowImage(true)
		case *gtk.SpinButton: // Set Icon to GtkSpinButton. options[0] = "left" or "right", No resize, No Animate.
			object.SetIconFromPixbuf(pos, inPixbuf)
		case *gtk.Box: // Add Image to GtkBox
			object.Add(image)
		}
	}
	return
}

// GetPixBuf: Get gdk.PixBuf from filename or []byte, depending on type
// size: resize height keeping porportions. 0 = no change
func GetPixBuf(varPath interface{}, size ...int) (outPixbuf *gdk.Pixbuf, err error) {
	var pixbufLoader *gdk.PixbufLoader
	sze := 0
	if len(size) != 0 {
		sze = size[0]
	}
	switch varPath.(type) {
	case string:
		outPixbuf, err = gdk.PixbufNewFromFile(varPath.(string))
	case []uint8:
		if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
			outPixbuf, err = pixbufLoader.WriteAndReturnPixbuf(varPath.([]byte))
		}
	}
	if err == nil && sze != 0 {
		newWidth, newHeight := NormalizeSize(outPixbuf.GetWidth(), outPixbuf.GetHeight(), sze, 2)
		outPixbuf, err = outPixbuf.ScaleSimple(newWidth, newHeight, gdk.INTERP_BILINEAR)
	}
	return
}

// GetPixBufAnimation: Get gdk.PixBufAnimation from filename or []byte, depending on type
func GetPixBufAnimation(varPath interface{}) (outPixbufAnimation *gdk.PixbufAnimation, err error) {
	var pixbufLoader *gdk.PixbufLoader
	switch varPath.(type) {
	case string:
		outPixbufAnimation, err = gdk.PixbufAnimationNewFromFile(varPath.(string))
	case []uint8:
		if pixbufLoader, err = gdk.PixbufLoaderNew(); err == nil {
			outPixbufAnimation, err = pixbufLoader.WriteAndReturnPixbufAnimation(varPath.([]byte))
		}
	}
	return
}

// NormalizeSize: compute new size with kept proportions based on defined format.
// format: 0 percent, 1 reducing width, 2 reducing height
func NormalizeSize(oldWidth, oldHeight, newValue, format int) (outWidth, outHeight int) {
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
	return
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
	if outPixbuf, err = GetPixBuf(varPath); err == nil {
		if width != outPixbuf.GetWidth() || height != outPixbuf.GetHeight() {
			outPixbuf, err = outPixbuf.ScaleSimple(width, height, interpolation)
		}
	}
	return
}

// RotateImage: Rotate by 90,180,270 degres and get gdk.PixBuf image representation from file or []byte, depending on type
func RotateImage(varPath interface{}, angle gdk.PixbufRotation) (outPixbuf *gdk.Pixbuf, err error) {
	if outPixbuf, err = GetPixBuf(varPath); err == nil {
		switch angle {
		case 90:
			outPixbuf, err = outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_COUNTERCLOCKWISE)
		case 180:
			outPixbuf, err = outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_UPSIDEDOWN)
		case 270:
			outPixbuf, err = outPixbuf.RotateSimple(gdk.PIXBUF_ROTATE_CLOCKWISE)
		default:
			return nil, errors.New("Rotation not allowed: " + fmt.Sprintf("%d", angle))
		}
	}
	return
}

// FlipImage: Get Flipped gdk.PixBuf image representation from file or []byte, depending on type
func FlipImage(varPath interface{}, horizontal bool) (outPixbuf *gdk.Pixbuf, err error) {
	if outPixbuf, err = GetPixBuf(varPath); err == nil {
		outPixbuf, err = outPixbuf.Flip(horizontal)
	}
	return
}

// /*
// 	Old functions, not used in new programs written and
// 	will be deleted after updating all other programs
// */
// // setImage: Set Image to GtkImage objects
// func SetImage(object *gtk.Image, varPath interface{}, size ...int) {
// 	if inPixbuf, err := GetPixBuf(varPath, size...); err == nil {
// 		object.SetFromPixbuf(inPixbuf)
// 		return
// 	} else if len(varPath.(string)) != 0 {
// 		fmt.Printf("SetImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setWinIcon: Set Icon to GtkWindow objects
// func SetWinIcon(object *gtk.Window, varPath interface{}, size ...int) {
// 	if inPixbuf, err := GetPixBuf(varPath, size...); err == nil {
// 		object.SetIcon(inPixbuf)
// 	} else if len(varPath.(string)) != 0 {
// 		fmt.Printf("SetWinIcon: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setButtonImage: Set Icon to GtkButton objects
// func SetButtonImage(object *gtk.Button, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			object.SetImage(image)
// 			object.SetAlwaysShowImage(true)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("SetButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setToolButtonImage: Set Icon to GtkToolButton objects
// func SetToolButtonImage(object *gtk.ToolButton, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			object.SetIconWidget(image)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("setToolButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setToggleButtonImage: Set Icon to GtkToggleButton objects
// func SetToggleButtonImage(object *gtk.ToggleButton, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			object.SetImage(image)
// 			object.SetAlwaysShowImage(true)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("SetToggleButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // SetSpinButtonImage: Set Icon to GtkSpinButton objects. Position = "left" or "right"
// func SetSpinButtonImage(object *gtk.SpinButton, varPath interface{}, position ...string) {
// 	var inPixbuf *gdk.Pixbuf
// 	var err error
// 	pos := gtk.ENTRY_ICON_PRIMARY
// 	if len(position) > 0 {
// 		if position[0] == "right" {
// 			pos = gtk.ENTRY_ICON_SECONDARY
// 		}
// 	}
// 	if inPixbuf, err = GetPixBuf(varPath); err == nil {
// 		object.SetIconFromPixbuf(pos, inPixbuf)
// 		return
// 	} else if len(varPath.(string)) != 0 {
// 		fmt.Printf("SetSpinButtonImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }

// // setBoxImage:  Set Image to GtkBox objects
// func SetBoxImage(object *gtk.Box, varPath interface{}, size ...int) {
// 	var image *gtk.Image
// 	inPixbuf, err := GetPixBuf(varPath, size...)
// 	if err == nil {
// 		if image, err = gtk.ImageNewFromPixbuf(inPixbuf); err == nil {
// 			image.Show()
// 			object.Add(image)
// 			return
// 		}
// 	}
// 	if err != nil && len(varPath.(string)) != 0 {
// 		fmt.Printf("setBoxImage: An error occurred on image: %v\n%s\n", varPath, err.Error())
// 	}
// }
