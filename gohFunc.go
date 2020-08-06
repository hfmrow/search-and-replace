// gohFunc.go

/*
	Source file auto-generated on Thu, 06 Aug 2020 20:25:48 using Gotk3ObjHandler v1.5 ©2018-20 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-20 H.F.M - Search And Replace v1.8 github.com/hfmrow/sAndReplace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

/*******************************************************/
/* Functions declarations, used to initialize objects */
/*****************************************************/
// newBuilder: initialise builder with glade xml string
func newBuilder(varPath interface{}) (err error) {
	var Gtk3Interface []byte
	if Gtk3Interface, err = getBytesFromVarAsset(varPath); err == nil {
		if mainObjects.mainUiBuilder, err = gtk.BuilderNew(); err == nil {
			err = mainObjects.mainUiBuilder.AddFromString(string(Gtk3Interface))
		}
	}
	return err
}

// loadObject: Load GtkObject to be transtyped ...
func loadObject(name string) (newObj glib.IObject) {
	var err error
	if newObj, err = mainObjects.mainUiBuilder.GetObject(name); err != nil {
		fmt.Printf("Unable to load %s object, maybe it was deleted from the Glade file ... : %s\n%s\n",
			name, err.Error(),
			fmt.Sprint("An update with GOH may avoid this issue."))
		os.Exit(1)
	}
	return newObj
}

// WindowDestroy: is the triggered handler when closing/destroying the gui window.
func windowDestroy() {
	/* Doing something before quit. Put something here ... */
	if err := mainOptions.Write(); err != nil { /* Update mainOptions with values of gtk conctrols and write to file */
		fmt.Printf("%s\n%v\n", "Writing options error.", err)
	}
	gtk.MainQuit()
}

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

/***************************************/
/* Embedded data conversion functions */
/* Used to make variable content     */
/* available in go-source           */
/***********************************/
// getBytesFromVarAsset: Get []byte representation from file or asset, depending on type
func getBytesFromVarAsset(varPath interface{}) (outBytes []byte, err error) {
	switch v := varPath.(type) {
	case string:
		return ioutil.ReadFile(v)
	case []uint8:
		return v, err
	}
	return
}

// HexToBytes: Convert Gzip Hex to []byte used for embedded binary in source code
func HexToBytes(varPath string, gzipData []byte) (outByte []byte) {
	r, err := gzip.NewReader(bytes.NewBuffer(gzipData))
	if err == nil {
		var bBuffer bytes.Buffer
		if _, err = io.Copy(&bBuffer, r); err == nil {
			if err = r.Close(); err == nil {
				return bBuffer.Bytes()
			}
		}
	} else {
		fmt.Printf("An error occurred while reading: %s\n%s\n", varPath, err.Error())
	}
	return
}

/*******************************/
/* Simplified files Functions */
/*****************************/
// Make temporary directory
func tempMake(prefix string) (dir string) {
	var err error
	if dir, err = ioutil.TempDir("", prefix+"-"); err != nil {
		log.Fatal(err)
	}
	return dir + string(os.PathSeparator)
}

// Retrieve current realpath and options filename
func getAbsRealPath() (absoluteRealPath, optFilename string) {
	absoluteBaseName, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	absoluteRealPath, base := filepath.Split(absoluteBaseName)
	splited := strings.Split(base, ".")
	if len(splited) == 1 {
		optFilename = filepath.Join(absoluteRealPath, base, ".") + ".opt"
	} else {
		splited = splited[:len(splited)-1]
		optFilename = filepath.Join(absoluteRealPath, strings.Join(splited, ".")+".opt")
	}
	return
}

// Used as fake function for signals section
func blankNotify() {}
