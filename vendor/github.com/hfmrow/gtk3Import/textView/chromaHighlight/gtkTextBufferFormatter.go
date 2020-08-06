// gtkTextBufferFormatter.go

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
	"html"
	"io"
	"regexp"
	"strconv"
	"strings"

	glsg "github.com/hfmrow/genLib/strings"

	"github.com/alecthomas/chroma"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

// tagProp: used to parse tag informations from chroma
type tagProp struct {
	Name  string
	Type  string
	Value interface{}
}

// tagDef: used to parse tag informations from chroma
type tagDef struct {
	Name, Prio string
	Props      []tagProp
}

// tagsRegistration: register and add tag definitions that was
// created with the syntax highlighter into the TextBuffer.TagList
func (c *ChromaHighlight) tagsRegistration(tagList string) (err error) {

	var splitted []string
	var def tagDef

	matchTagName := regexp.MustCompile(`<tag name=`)
	matchTagAttr := regexp.MustCompile(`<attr name=`)
	matchTagEnd := regexp.MustCompile(`</tag>`)
	matchEndList := regexp.MustCompile(`.*</tags>.*`)
	eol := glsg.GetTextEOL([]byte(tagList))

	lines := strings.Split(tagList, eol)
	for _, line := range lines {
		splitted = strings.Split(line, `"`)
		switch {
		case matchTagName.MatchString(line):
			def.Name = splitted[1]
			def.Prio = splitted[3]
		case matchTagAttr.MatchString(line):
			def.Props = append(def.Props, tagProp{Name: splitted[1], Type: splitted[3], Value: splitted[5]})
		case matchTagEnd.MatchString(line):
			if _, err = c.buildTag(def); err != nil && err != c.errAlreadyExist {
				return
			}
			err = nil
			def = tagDef{}
		case matchEndList.MatchString(line):
			return
		}
	}

	c.initialised = true
	return
}

// buildTag: create tag with properties and add it to buffer.
func (c *ChromaHighlight) buildTag(def tagDef) (tag *gtk.TextTag, err error) {
	// Convert GdkColor to (gchar)
	var getColor = func(v interface{}) string {
		return "#" + v.(string)[:2] + v.(string)[5:7] + v.(string)[10:12]
	}
	if tag = c.TextTagList[def.Name]; tag == nil {
		// Convert Buffer RTF to tag
		tagProp := map[string]interface{}{}
		for _, prop := range def.Props {
			switch prop.Name {
			case "style":
				tagProp["style"] = pango.STYLE_ITALIC
			case "underline":
				tagProp["underline"] = pango.UNDERLINE_SINGLE
			case "strikethrough": // Not used actually, implemented for further usage
				tagProp["strikethrough"] = true
			case "weight": // Bold
				val, _ := strconv.Atoi(prop.Value.(string))
				tagProp[prop.Name] = val
			case "foreground-gdk":
				tagProp["foreground"] = getColor(prop.Value)
			case "background-gdk":
				tagProp["background"] = getColor(prop.Value)
			}
		}
		c.TextTagList[def.Name] = c.buff.CreateTag(def.Name, tagProp)

	}
	return
}

// gtkTextBufferFormatter: is a part of "ChromaHighlight" function.
func (c *ChromaHighlight) gtkTextBufferFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) (err error) {
	var name string
	var entry chroma.StyleEntry
	style = c.chromaClearBackground(style)
	// Seek tokens
	for tkn := it(); tkn != chroma.EOF; tkn = it() {
		name = strings.ToLower(tkn.Type.String())
		entry = style.Get(tkn.Type)

		if !entry.IsZero() {
			// The tag does not already exist, so, create it
			if !c.tagDefList[name] {
				c.tagDefinition += c.buildOneTagDef(name, &entry)
			}
			// Create RTF tag
			fmt.Fprint(w, `<apply_tag name="`+name+`">`+html.EscapeString(tkn.Value)+`</apply_tag>`)
		} else {
			fmt.Fprint(w, html.EscapeString(tkn.Value))
		}
	}
	return
}

// buildOneTagDef: create definition for named tag.
func (c *ChromaHighlight) buildOneTagDef(name string, entry *chroma.StyleEntry) (outString string) {
	var attr string
	var r, g, b uint8
	var getColour = func(color chroma.Colour) string {
		r, g, b = color.Red(), color.Green(), color.Blue()
		return fmt.Sprintf("%02X%02X:%02X%02X:%02X%02X", r, r, g, g, b, b)
	}
	if entry.Bold == chroma.Yes {
		attr = `   <attr name="weight" type="gint" value="700" />`
	}
	if entry.Italic == chroma.Yes {
		if len(attr) > 0 { // add EOL if not alone
			attr += "\n"
		}
		attr += `   <attr name="style" type="PangoStyle" value="PANGO_STYLE_ITALIC" />`
	}
	if entry.Underline == chroma.Yes {
		if len(attr) > 0 { // add EOL if not alone
			attr += "\n"
		}
		attr += `   <attr name="underline" type="PangoUnderline" value="PANGO_UNDERLINE_SINGLE" />`
	}
	if entry.Colour.IsSet() {
		if len(attr) > 0 { // add EOL if not alone
			attr += "\n"
		}
		attr += `   <attr name="foreground-gdk" type="GdkColor" value="` + getColour(entry.Colour) + `" />`
	}
	if entry.Background.IsSet() {
		if len(attr) > 0 { // add EOL if not alone
			attr += "\n"
		}
		attr += `   <attr name="background-gdk" type="GdkColor" value="` + getColour(entry.Background) + `" />`
	}
	if entry.Border.IsSet() {
		if len(attr) > 0 { // add EOL if not alone
			attr += "\n"
		}
		attr += `   <attr name="background-gdk" type="GdkColor" value="` + getColour(entry.Border) + `" />`
	}
	if len(attr) > 0 {
		c.tagDefList[name] = true
		outString = "\n" + `  <tag name="` + name + `" priority="` + fmt.Sprintf("%d", len(c.tagDefList)) + `">
` + attr + `
  </tag>`
	}
	return
}
