// textViewRowNumber.go

// Source file auto-generated on Tue, 20 Aug 2019 18:38:23 using Gotk3ObjHandler v1.3.6 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

/*
* Sync Textview scroll with numbers
 */

package gtk3Import

import (
	"errors"
	"fmt"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

var err error

/*
* Main structure
 */
type TextViewRowNumber struct {
	LinesCount           int
	alreadyTagged        bool
	tagIt                bool
	signalHdlChanged     glib.SignalHandle
	signalHdlScrollEvent glib.SignalHandle
	textViewNumbersRef   *gtk.TextView
	textViewTextRef      *gtk.TextView
	buffNumb             *gtk.TextBuffer
	buffText             *gtk.TextBuffer
	scrolledNumbersRef   *gtk.ScrolledWindow
	scrolledTextRef      *gtk.ScrolledWindow
}

/*
* Tags stuff
 */
type TextTagProps map[string]interface{}

type textTagList map[string]TextTagProps

/*
* Functions
 */
// TextViewNumberedSetting: Synchronize two texview, first containing lines numbers, second contain text.
func TextViewRowNumberNew(textViewNumbers, textViewText *gtk.TextView,
	scrolledNumbers, scrolledText *gtk.ScrolledWindow, applyTags ...bool) (tvrn *TextViewRowNumber, err error) {
	tvrn = new(TextViewRowNumber)
	err = tvrn.StructSetup(textViewNumbers, textViewText, scrolledNumbers, scrolledText, applyTags...)
	return
}

func (tvrn *TextViewRowNumber) StructSetup(textViewNumbers, textViewText *gtk.TextView,
	scrolledNumbers, scrolledText *gtk.ScrolledWindow, applyTags ...bool) (err error) {

	tvrn.textViewNumbersRef, tvrn.textViewTextRef = textViewNumbers, textViewText
	tvrn.scrolledNumbersRef, tvrn.scrolledTextRef = scrolledNumbers, scrolledText

	if len(applyTags) > 0 {
		tvrn.tagIt = applyTags[0]
	}

	if !tvrn.alreadyTagged {
		if err = tvrn.DetachBuffers(); err == nil {
			tvrn.AttachBuffers()
		} else {
			return errors.New(fmt.Sprintf("Could not get buffers: %s", err.Error()))
		}
		// Prepare Text textview
		tvrn.textViewTextRef.SetProperty("left-margin", 4)
		tvrn.textViewTextRef.SetProperty("right-margin", 3)
		tvrn.textViewTextRef.SetProperty("editable", true)

		tvrn.scrolledTextRef.SetShadowType(gtk.SHADOW_NONE)
		tvrn.scrolledTextRef.SetProperty("hexpand", true)
		tvrn.scrolledTextRef.SetProperty("vexpand", true)
		tvrn.scrolledTextRef.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

		// Prepare Numbers textview
		tvrn.textViewNumbersRef.SetSensitive(false)
		tvrn.textViewNumbersRef.SetProperty("left-margin", 2)
		tvrn.textViewNumbersRef.SetProperty("right-margin", 3)
		tvrn.textViewNumbersRef.SetProperty("editable", false)
		tvrn.textViewNumbersRef.SetProperty("justification", gtk.JUSTIFY_RIGHT)

		tvrn.scrolledNumbersRef.SetSensitive(false)
		tvrn.scrolledNumbersRef.SetShadowType(gtk.SHADOW_NONE)
		tvrn.scrolledNumbersRef.SetProperty("window-placement", gtk.CORNER_TOP_RIGHT)
		tvrn.scrolledNumbersRef.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
		tvrn.scrolledNumbersRef.SetProperty("halign", gtk.ALIGN_START)
		// Add tags to TagTable if requested
		tvrn.createTags()
		// Connect signals to make interactive, synchronize scroll
		tvrn.scrolled() // To synch both scrolledWin before start
		if tvrn.signalHdlScrollEvent, err = tvrn.scrolledTextRef.Connect("scroll-event", tvrn.scrolled); err == nil {
			tvrn.signalHdlChanged, err = tvrn.buffText.Connect("changed", tvrn.UpdateTextViewNumbers)
		}

		if err != nil {
			return errors.New(fmt.Sprintf("Could not add signal handler: %s", err.Error()))
		}
		tvrn.alreadyTagged = true
	}
	// Add/update numbers to designed textview.
	tvrn.UpdateTextViewNumbers()

	// Apply some tags ...
	// buffNumb.ApplyTagByName("blue_foreground", buffNumb.GetStartIter(), buffNumb.GetEndIter())
	return
}

// Build and store lines numbers to his textview...
func (tvrn *TextViewRowNumber) UpdateTextViewNumbers() {
	linesCount := tvrn.buffText.GetLineCount()
	if linesCount != tvrn.LinesCount {
		if linesCount < tvrn.LinesCount {
			tvrn.LinesCount = 0
			tvrn.buffNumb.Delete(tvrn.buffNumb.GetStartIter(), tvrn.buffNumb.GetEndIter()) // Clear TextBuffer
		}
		// Stop signal handler.
		tvrn.buffText.HandlerBlock(tvrn.signalHdlChanged)
		tvrn.scrolledTextRef.HandlerBlock(tvrn.signalHdlScrollEvent)

		for idx := tvrn.LinesCount + 1; idx <= linesCount; idx++ {
			tvrn.buffNumb.Insert(tvrn.buffNumb.GetEndIter(), fmt.Sprintf("%d\n", idx))
		}
		// Restart signal handler.
		tvrn.buffText.HandlerUnblock(tvrn.signalHdlChanged)
		tvrn.scrolledTextRef.HandlerUnblock(tvrn.signalHdlScrollEvent)
		tvrn.LinesCount = linesCount
	}
}

func (tvrn *TextViewRowNumber) ScrollToLine(line int) {
	TextViewScrollToLine(tvrn.textViewTextRef, line)
	// if line > 0 {
	// 	iter := tvrn.buffText.GetIterAtLine(line)
	// 	iter1 := tvrn.buffText.GetIterAtLine(line - 1)
	// 	tvrn.buffText.PlaceCursor(iter)
	// 	for gtk.EventsPending() {
	// 		gtk.MainIterationDo(false)
	// 	}
	// 	tvrn.textViewTextRef.ScrollToIter(iter, 0.0, true, 0.5, 0.5)

	// 	tvrn.buffText.SelectRange(iter, iter1) // HighLight current line.
	// }
}

func (tvrn *TextViewRowNumber) DetachBuffers() (err error) {
	if tvrn.buffNumb, err = tvrn.textViewNumbersRef.GetBuffer(); err == nil {
		if tvrn.buffText, err = tvrn.textViewTextRef.GetBuffer(); err == nil {
			tvrn.buffNumb.Ref()
			tvrn.buffText.Ref()
			tvrn.textViewNumbersRef.SetBuffer(nil)
			tvrn.textViewTextRef.SetBuffer(nil)
		}
	}
	return
}

func (tvrn *TextViewRowNumber) AttachBuffers() {
	tvrn.textViewNumbersRef.SetBuffer(tvrn.buffNumb)
	tvrn.textViewTextRef.SetBuffer(tvrn.buffText)
	tvrn.buffNumb.Unref()
	tvrn.buffText.Unref()
}

// createTags: Creating tags elements
func (tvrn *TextViewRowNumber) createTags() {
	var tagList = textTagList{
		"normal":           TextTagProps{"style": pango.STYLE_NORMAL},
		"oblique":          TextTagProps{"style": pango.STYLE_OBLIQUE},
		"italic":           TextTagProps{"style": pango.STYLE_ITALIC},
		"bold":             TextTagProps{"weight": pango.WEIGHT_BOLD},
		"underline":        TextTagProps{"underline": pango.UNDERLINE_SINGLE},
		"double_underline": TextTagProps{"underline": pango.UNDERLINE_DOUBLE},
		"strikethrough":    TextTagProps{"strikethrough": true},
		"heading":          TextTagProps{"weight": pango.WEIGHT_BOLD, "size": 15 * pango.SCALE},
		"big":              TextTagProps{"size": 20 * pango.SCALE},
		// Can't use pango variables, ... cannot convert type ...
		"xx-small":          TextTagProps{"scale": 0.5787037037037}, /* pango.SCALE_XX_SMALL */
		"x-small":           TextTagProps{"scale": 0.6444444444444}, /* pango.SCALE_X_SMALL */
		"small":             TextTagProps{"scale": 0.8333333333333}, /* pango.SCALE_SMALL */
		"medium":            TextTagProps{"scale": 1.0},             /* pango.SCALE_MEDIUM */
		"large":             TextTagProps{"scale": 1.2},             /* pango.SCALE_LARGE */
		"x-large":           TextTagProps{"scale": 1.4399999999999}, /* pango.SCALE_X_LARGE */
		"xx-large":          TextTagProps{"scale": 1.728},           /* pango.SCALE_XX_LARGE */
		"superscript":       TextTagProps{"rise": 10 * pango.SCALE, "size": 8 * pango.SCALE},
		"subscript":         TextTagProps{"rise": -10 * pango.SCALE, "size": 8 * pango.SCALE},
		"monospace":         TextTagProps{"family": "monospace"},
		"blue_foreground":   TextTagProps{"foreground": "#3050FF"},
		"yellow_background": TextTagProps{"background": "#FFFF99"},
	}
	for tagName, prop := range tagList {
		tvrn.buffText.CreateTag(tagName, prop)
	}
}

// scrolled: Synchronize ScrollWindow to scroll as together
func (tvrn *TextViewRowNumber) scrolled() {
	tvrn.scrolledNumbersRef.SetVAdjustment(tvrn.scrolledTextRef.GetVAdjustment())
}

// textViewScrollToLine: Scroll to line and highligth it
func TextViewScrollToLine(textView *gtk.TextView, line int) {
	var err error
	if line > 0 {
		if buf, err := textView.GetBuffer(); err == nil {
			iter := buf.GetIterAtLine(line)
			iter1 := buf.GetIterAtLine(line - 1)
			buf.PlaceCursor(iter)
			for gtk.EventsPending() {
				gtk.MainIterationDo(false)
			}
			textView.ScrollToIter(iter, 0.0, true, 0.5, 0.5)

			buf.SelectRange(iter, iter1) // HighLight current line.
		}
	}
	if err != nil {
		fmt.Println(err.Error())
	}
}
