// gtkDirectToTextBufferFormatter.go

/*
	This library use:
	- gotk3 that is licensed under the ISC License:
	  https://github.com/gotk3/gotk3/blob/master/LICENSE

	- Chroma — A general purpose syntax highlighter in pure Go, under the MIT License:
	  https://github.com/alecthomas/chroma/LICENSE

	Copyright ©2019 H.F.M gotk3_chroma_syntax_highlighter library "https://github/hfmrow"
	This library comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package chromaHighlight

import (
	"fmt"
	"io"
	"strings"

	"github.com/gotk3/gotk3/gtk"

	"github.com/alecthomas/chroma"
	"github.com/gotk3/gotk3/pango"
)

// gtkTextBufferFormatter: is a part of "ChromaHighlight" function.
func (c *ChromaHighlight) gtkDirectToTextBufferFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) (err error) {
	style = c.chromaClearBackground(style)
	var tag *gtk.TextTag
	for tkn := it(); tkn != chroma.EOF; tkn = it() {
		entry := style.Get(tkn.Type)
		tokenName := strings.ToLower(tkn.Type.String())
		if !entry.IsZero() {
			tag = c.buildOneTagDefDirectToTextBuffer(tokenName, &entry)
			// Insert RTF tag direct to TextBuffer
			c.buff.InsertWithTag(c.buff.GetIterAtMark(c.buff.GetInsert()), tkn.Value, tag)
		} else {
			// No highlight, just insert normal text
			c.buff.Insert(c.buff.GetIterAtMark(c.buff.GetInsert()), tkn.Value)
		}
	}
	c.initialised = true
	return
}

// buildOneTagDefDirectToBuff: create tag and apply to TextBuffer
func (c *ChromaHighlight) buildOneTagDefDirectToTextBuffer(name string, styleEntry *chroma.StyleEntry) (tag *gtk.TextTag) {
	tagProp := map[string]interface{}{}
	if styleEntry.Bold == chroma.Yes {
		tagProp["weight"] = pango.WEIGHT_BOLD
	}
	if styleEntry.Italic == chroma.Yes {
		tagProp["style"] = pango.STYLE_ITALIC
	}
	if styleEntry.Underline == chroma.Yes {
		tagProp["underline"] = pango.UNDERLINE_SINGLE
	}
	if styleEntry.Colour.IsSet() {
		tagProp["foreground"] = fmt.Sprintf("#%02X%02X%02X",
			styleEntry.Colour.Red(),
			styleEntry.Colour.Green(),
			styleEntry.Colour.Blue())
	}
	if styleEntry.Background.IsSet() {
		tagProp["background"] = fmt.Sprintf("#%02X%02X%02X",
			styleEntry.Background.Red(),
			styleEntry.Background.Green(),
			styleEntry.Background.Blue())
	}
	if styleEntry.Border.IsSet() {
		tagProp["background"] = fmt.Sprintf("#%02X%02X%02X",
			styleEntry.Border.Red(),
			styleEntry.Border.Green(),
			styleEntry.Border.Blue())
	}
	tag = c.createTag(name, tagProp)
	// store it for further usage
	if !c.preExistsTagList[name] { // check wether did not exist before.
		if c.TextTagList[name] == nil {
			c.TextTagList[name] = tag // add wether not actually exist.
		}
	}
	return
}
