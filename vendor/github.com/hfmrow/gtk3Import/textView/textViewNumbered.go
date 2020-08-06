// textViewNumbered.go

/*
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2019 H.F.M - TextView Numbered  library "https//github/hfmrow"
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	This source code is a part of a personal library using gotk3.

	Allow using a scrolled TextView with line numbers (left sided) with some formatting text
	controls

	There is a limitation regarding some types of text containing non-printable characters
	(some fonts print Utf-8 codes instead of character), does not work properly, the height
	of the line numbers is not correctly synchronized with the text. These characters print
	an oversized entry in the textView column, so the numbers are shifted relative to the
	lines of text.
*/

package textView

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	glco "github.com/hfmrow/genLib/crypto"

	gimc "github.com/hfmrow/gtk3Import/misc"
	gitvtt "github.com/hfmrow/gtk3Import/textView/textTag"
)

var count int

// Allow the use of a TextView with line numbers (left sided) with scrolling
// capabilities and some formatting text controls
type TextViewNumbered struct {
	TextView    *gtk.TextView
	BuffTxt     *gtk.TextBuffer
	ShowNumbers bool
	Editable    bool

	NumFgCol        string
	NumBgCol        string
	TxtFgCol        string
	TxtBgCol        string
	SelFgCol        string
	SelBgCol        string
	ColorBgRangeSet string

	FontSze    int
	FontFamily string

	BufferChangeCallbackFunc func()

	tVNum           *gtk.TextView
	tVNumWidth      int
	lastTVNumWidth  int
	textViewLMargin int

	scrTxt      *gtk.ScrolledWindow
	container   gtk.Container
	hAdj        *gtk.Adjustment
	vAdj        *gtk.Adjustment
	hValue      float64
	vValue      float64
	currentline int

	buffNum *gtk.TextBuffer

	lastLineCount int
	sigHdlBufTxt  glib.SignalHandle
	// sigHdlscrTxt  glib.SignalHandle
	lastMd5          string
	colorBgRangeName string
}

// TextViewNumberedNew: Create and initialize a new TextViewNumbered structure
func TextViewNumberedNew(container gtk.Container) (tvn *TextViewNumbered, err error) {
	tvn = new(TextViewNumbered)
	err = tvn.Init(container)
	return
}

// Init: Initialize a TextViewNumbered structure
func (tvn *TextViewNumbered) Init(container gtk.Container) (err error) {
	tvn.container = container
	tvn.textViewLMargin = 5
	tvn.ShowNumbers = true
	tvn.Editable = true
	tvn.lastLineCount = 0
	tvn.BufferChangeCallbackFunc = func() {}

	// Default colors for numbers
	tvn.NumFgCol = "#3322FF"
	tvn.NumBgCol = "#EEEEEE"
	tvn.TxtFgCol = "#331111"
	tvn.TxtBgCol = "#F8F8F8"
	tvn.SelFgCol = "#152727"
	tvn.SelBgCol = "#CBEBEB"
	tvn.FontSze = 12
	tvn.FontFamily = `"Liberation Mono", sans-serif`
	tvn.colorBgRangeName = "highlightBCColorSoftLightGreen"
	tvn.ColorBgRangeSet = "#E6FFE6" // lightgreen

	// Setup scrolled window
	if tvn.scrTxt, err = gtk.ScrolledWindowNew(nil, nil); err == nil {
		tvn.scrTxt.SetProperty("visible", true)
		tvn.scrTxt.SetProperty("can-focus", true)
		tvn.scrTxt.SetProperty("hexpand", true)
		tvn.scrTxt.SetProperty("vexpand", true)
		tvn.scrTxt.SetProperty("overlay-scrolling", true)
		tvn.scrTxt.SetShadowType(gtk.SHADOW_NONE)
		tvn.scrTxt.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_ALWAYS)
		tvn.hAdj = tvn.scrTxt.GetHAdjustment()
		tvn.vAdj = tvn.scrTxt.GetVAdjustment()
		// Setup text
		if tvn.TextView, err = gtk.TextViewNew(); err == nil {
			tvn.TextView.SetProperty("visible", true)
			tvn.TextView.SetProperty("can-focus", true)
			tvn.TextView.SetProperty("hexpand", true)
			tvn.TextView.SetProperty("vexpand", true)
			tvn.TextView.SetProperty("left-margin", tvn.textViewLMargin)
			tvn.TextView.SetProperty("right-margin", 3)
			tvn.TextView.SetProperty("editable", true)
			tvn.scrTxt.Add(tvn.TextView)
			// Setup numbers
			if tvn.tVNum, err = gtk.TextViewNew(); err == nil {
				tvn.tVNum.SetSensitive(false)
				tvn.tVNum.SetProperty("visible", true)
				tvn.tVNum.SetProperty("can-focus", false)
				tvn.tVNum.SetProperty("hexpand", false)
				tvn.tVNum.SetProperty("vexpand", true)
				tvn.tVNum.SetProperty("left-margin", 2)
				tvn.tVNum.SetProperty("right-margin", 3)
				tvn.tVNum.SetProperty("editable", false)
				tvn.tVNum.SetProperty("justification", gtk.JUSTIFY_RIGHT)
				// Setup viewport
				tvn.TextView.AddChildInWindow(tvn.tVNum, gtk.TEXT_WINDOW_LEFT, 0, 0)
				tvn.container.Add(tvn.scrTxt)

				if tvn.buffNum, err = tvn.tVNum.GetBuffer(); err == nil {
					if tvn.BuffTxt, err = tvn.TextView.GetBuffer(); err == nil {
						tvn.lastLineCount = tvn.BuffTxt.GetLineCount()

						// tvn.TextView.Connect("realize", tvn.updateNumbers)
						tvn.sigHdlBufTxt, _ = tvn.BuffTxt.Connect("changed", tvn.updateNumbers)
						tvn.doCss()
					}
				}
			}
		}
	}
	return
}

// ResetMd5: Used to determine if the text must be displayed
// again, ext has been changed.
func (tvn *TextViewNumbered) ResetMd5() {
	tvn.lastMd5 = ""
}

// WaitForEventPending: Wait for pending events (until the widget is redrawn)
func (tvn *TextViewNumbered) WaitForEventPending() {
	for gtk.EventsPending() {
		gtk.MainIteration()
	}
}

// // scrolled: used to asjust left margin of the text if its necessary
// func (tvn *TextViewNumbered) scrolled() {
// 	if tvn.tVNumWidth != tvn.lastTVNumWidth {
// 		tvn.autoSetMargin()
// 		tvn.lastTVNumWidth = tvn.tVNumWidth
// 	}
// }

// Update: refresh TextView.
func (tvn *TextViewNumbered) Update() {
	tvn.TextView.SetEditable(tvn.Editable)
	tvn.doCss()
	tvn.updateNumbers()
}

// FromFile: opens, loads the text file into the TextView and scrolls to
// the position if it has been specified.
func (tvn *TextViewNumbered) FromFile(filename string, position ...int) (err error) {
	var data []byte
	pos := 1
	if len(position) > 0 {
		pos = position[0]
	}
	if data, err = ioutil.ReadFile(filename); err == nil {
		tvn.SetText(string(data))
		tvn.ScrollToLine(pos)
	}
	return
}

// StoreScrollPos: store position in he scroll widget.
// Usage:
//	StoreScrollPos()
//	line = GetCurrentLineNb()
//	 ... doing some adding to textview (assuming you understand adding at the end of textview ...)
//	 ... to preserve the "line" variable target position.
//	SetCurrentLineNb(line)
//	RestoreScrollPos()
// in practice, you can add text with a minimum latency time.
func (tvn *TextViewNumbered) StoreScrollPos() {
	tvn.hValue = tvn.hAdj.GetValue()
	tvn.vValue = tvn.vAdj.GetValue()
}

// RestoreScrollPos: Restore position to scroll widget.
func (tvn *TextViewNumbered) RestoreScrollPos() {
	tvn.hAdj.SetValue(tvn.hValue)
	tvn.vAdj.SetValue(tvn.vValue)
}

// Clear: Create new buffers (txt & num) and assign it.
// If you have made some signals callback, you must do it again.
func (tvn *TextViewNumbered) Reset() {
	tvn.BuffTxt, _ = gtk.TextBufferNew(nil)
	tvn.buffNum, _ = gtk.TextBufferNew(nil)
	tvn.lastLineCount = tvn.BuffTxt.GetLineCount()
	tvn.sigHdlBufTxt, _ = tvn.BuffTxt.Connect("changed", tvn.updateNumbers)
}

// Clear: Delete buffers content (txt & num)
func (tvn *TextViewNumbered) Clear() {
	tvn.BuffTxt.Delete(tvn.BuffTxt.GetStartIter(), tvn.BuffTxt.GetEndIter())
	tvn.buffNum.Delete(tvn.buffNum.GetStartIter(), tvn.buffNum.GetEndIter())
	tvn.lastLineCount = tvn.BuffTxt.GetLineCount()
}

// SetTextMarkup: and scroll to line nb if given.
func (tvn *TextViewNumbered) SetTextMarkup(text string, line ...int) {
	tvn.setTextFunc(&text, func() { tvn.BuffTxt.InsertMarkup(tvn.BuffTxt.GetStartIter(), text) }, line...)
}

// SetText: and scroll to line nb if given.
func (tvn *TextViewNumbered) SetText(text string, line ...int) {
	tvn.setTextFunc(&text, func() { tvn.BuffTxt.SetText(text) }, line...)
}

// setTextFunc:
func (tvn *TextViewNumbered) setTextFunc(text *string, setTextFct func(), line ...int) {
	// Avoid writing twice the same text to TextView.
	currMd5 := glco.Md5String(*text)
	if tvn.lastMd5 != currMd5 {
		tvn.detachBuffers()
		tvn.Clear()
		setTextFct()
		tvn.lastMd5 = currMd5
		tvn.attachBuffers()
	}
	if len(line) > 0 {
		tvn.ScrollToLine(line[0])
	}

	// glib.IdleAdd(func() {
	tvn.BufferChangeCallbackFunc() // TODO not really useful ?
	// })
}

// BufferDetach: Unlink buffer from TextView, if you hav to fill it
// with a large amount of data, it's useful to doing that.
func (tvn *TextViewNumbered) BufferDetach() {
	tvn.detachBuffers()
}

// BufferAttach: to execute after the previous one.
func (tvn *TextViewNumbered) BufferAttach() {
	tvn.attachBuffers()
}

// GetText: retrieve text from TextBuffer.
func (tvn *TextViewNumbered) GetText() (text string) {
	var err error
	if text, err = tvn.BuffTxt.GetText(tvn.BuffTxt.GetStartIter(), tvn.BuffTxt.GetEndIter(), true); err != nil {
		log.Fatalf("GetText: %s\n", err.Error)
	}
	return
}

// GetLineNumber: Get current line number at the cursor position
func (tvn *TextViewNumbered) GetCurrentLineNb() int {
	iter := tvn.BuffTxt.GetIterAtMark(tvn.BuffTxt.GetInsert())
	if iter != tvn.BuffTxt.GetEndIter() {
		return iter.GetLine()
	}
	return -1
}

// SetCursorAtLine: Set current line number & place cursor on it.
func (tvn *TextViewNumbered) SetCursorAtLine(line int) (iter *gtk.TextIter) {
	iter = tvn.BuffTxt.GetIterAtLine(line)
	tvn.BuffTxt.PlaceCursor(iter)
	tvn.TextView.GrabFocus()
	return
}

// SelectRange: select Lines
func (tvn *TextViewNumbered) SelectRange(startLine, endLine int) {
	if startLine <= endLine && startLine > 0 {

		startIter := tvn.BuffTxt.GetIterAtLine(startLine)
		endIter := tvn.BuffTxt.GetIterAtOffset(tvn.BuffTxt.GetIterAtLine(endLine).GetOffset() - 1)
		tvn.BuffTxt.SelectRange(startIter, endIter)
	}
}

// ColorBgRange: apply colored background to lines range.
func (tvn *TextViewNumbered) ColorBgRange(startLine, endLine int) {
	gitvtt.TagRemoveIfExists(tvn.BuffTxt, tvn.colorBgRangeName)
	if startLine <= endLine && startLine > 0 {
		tag := gitvtt.TagCreateIfNotExists(
			tvn.BuffTxt, tvn.colorBgRangeName,
			map[string]interface{}{"background": tvn.ColorBgRangeSet})

		startIter := tvn.BuffTxt.GetIterAtLine(startLine)
		endIter := tvn.BuffTxt.GetIterAtOffset(tvn.BuffTxt.GetIterAtLine(endLine).GetOffset() - 1)
		tvn.BuffTxt.ApplyTag(tag, startIter, endIter)
	}
}

// ScrollToLine: Scroll to line
func (tvn *TextViewNumbered) ScrollToLine(line int) (iter *gtk.TextIter) {
	if line > 0 && line < tvn.BuffTxt.GetLineCount() {

		for gtk.EventsPending() {
			gtk.MainIteration() // Wait for pending events (until the widget is redrawn)
		}

		// iter = tvn.BuffTxt.GetIterAtLine(line)
		// tvn.BuffTxt.PlaceCursor(iter)
		// mark := tvn.BuffTxt.GetInsert()
		// tvn.TextView.ScrollToMark(mark, 0.0, true, 0.5, 0.5)

		glib.IdleAdd(func() {
			var count int
			glib.TimeoutAdd(uint(64), func() bool {
				count++
				tvn.TextView.ScrollToIter(tvn.BuffTxt.GetIterAtLine(line), 0.0, true, 0.5, 0.5)
				return count <= 5
			})
		})
	}
	return
}

// updateNumbers: Build and store lines numbers to his textview...
func (tvn *TextViewNumbered) updateNumbers() {

	if tvn.TextView.GetWrapMode() == gtk.WRAP_NONE && tvn.ShowNumbers {
		if !tvn.tVNum.GetVisible() {
			tvn.tVNum.SetVisible(true)
			tvn.autoSetMargin()
		}
		linesCount := tvn.BuffTxt.GetLineCount()
		if linesCount != tvn.lastLineCount-1 {

			// detachBuffer of numbers column
			tvn.buffNum.Ref()
			tvn.tVNum.SetBuffer(nil)

			if linesCount < tvn.lastLineCount {
				// Remove some numbers
				tvn.lastLineCount = linesCount
				tvn.buffNum.Delete(tvn.buffNum.GetIterAtLine(tvn.lastLineCount), tvn.buffNum.GetEndIter())
			} else {
				// Add numbers
				var tmpLines string
				for idx := tvn.lastLineCount; idx <= linesCount; idx++ {
					tmpLines += strconv.Itoa(idx) + "\n"
				}
				tvn.buffNum.Insert(tvn.buffNum.GetEndIter(), tmpLines)
			}

			// attachBuffer of numbers column
			tvn.tVNum.SetBuffer(tvn.buffNum)
			tvn.buffNum.Unref()

			tvn.autoSetMargin()
			tvn.lastLineCount = linesCount + 1
		}
		return
	}
	// If we have a Wrap method set to the textview, never display he column numbers
	tvn.TextView.SetLeftMargin(tvn.textViewLMargin)
	tvn.tVNum.SetVisible(false)
}

// autoSetMargin: used to align left margin of the text view with numbers column width
func (tvn *TextViewNumbered) autoSetMargin() {
	var count int

	// fmt.Println("autoSetMargin")

	glib.IdleAdd(func() {

		glib.TimeoutAdd(uint(64), func() bool {
			count++
			tvn.tVNumWidth = tvn.tVNum.GetVisibleRect().GetWidth()
			// First display, to avoid having a large left margin on very large text files.
			if tvn.tVNumWidth > 150 {
				tvn.tVNumWidth = 55 // Default value if we haven't other possibility.
			}
			tvn.TextView.SetLeftMargin(tvn.textViewLMargin + tvn.tVNumWidth)
			if tvn.TextView.GetLeftMargin() == tvn.textViewLMargin+tvn.tVNumWidth ||
				count > 47 { // keep working during some seconds if there is no change
				return false
			}
			return true
		})
	})
}

// doColored: Add some colors to the numbers column
func (tvn *TextViewNumbered) doCss() {

	tvn.TextView.SetName("tvtxt")
	gimc.CssWdgScnLoad(`
#tvtxt {
	font: normal normal normal normal `+fmt.Sprintf("%d", tvn.FontSze)+`px `+tvn.FontFamily+`;
}`, &tvn.TextView.Widget)

	gimc.CssWdgScnLoad(`
#tvtxt text {
	color: `+tvn.TxtFgCol+`;
	background-color: `+tvn.TxtBgCol+`;	
}
.view text selection {
	background-color: `+tvn.SelBgCol+`;
	color: `+tvn.SelFgCol+`;
}`, &tvn.TextView.Widget)

	tvn.tVNum.SetName("tvnum")
	gimc.CssWdgScnLoad(`
#tvnum text {
	color: `+tvn.NumFgCol+`;
	background-color: `+tvn.NumBgCol+`;
}`, &tvn.tVNum.Widget)
}

// detachBuffers: and block emitting "changed" signal
func (tvn *TextViewNumbered) detachBuffers() (err error) {
	if tvn.BuffTxt, err = tvn.TextView.GetBuffer(); err == nil {
		tvn.BuffTxt.HandlerBlock(tvn.sigHdlBufTxt)
		tvn.BuffTxt.Ref()
		tvn.TextView.SetBuffer(nil)
	}
	return
}

// attachBuffers: and unblock emitting "changed" signal
func (tvn *TextViewNumbered) attachBuffers() {
	tvn.TextView.SetBuffer(tvn.BuffTxt)
	tvn.BuffTxt.Unref()
	tvn.Update()
	tvn.BuffTxt.HandlerUnblock(tvn.sigHdlBufTxt)
}
