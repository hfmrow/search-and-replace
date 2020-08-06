// chromaHighlight.go

/*
	This library use:
	- gotk3 that is licensed under the ISC License:
	  https://github.com/gotk3/gotk3/blob/master/LICENSE

	- Chroma — A general purpose syntax highlighter in pure Go, under the MIT License:
	  https://github.com/alecthomas/chroma/LICENSE

	Copyright ©2019 H.F.M gotk3_chroma_syntax_highlighter library "https://github/hfmrow"
	This library comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php

	- Information: It is strongly recommended to detach the buffer before filling it,
	(using the same method as for a TreeView with its model when it is filled)

	i.e:	"sigHdlBufTxt" is a "glib.SignalHandle", stored on signal "connect" method

	// DetachBuffers: and block emitting "changed" signal
	func DetachBuffers(TextView *gtk.TextView) (err error) {
		if BuffTxt, err = TextView.GetBuffer(); err == nil {
			BuffTxt.HandlerBlock(sigHdlBufTxt) // Use it if you have a callback on "change" signal
			BuffTxt.Ref()
			TextView.SetBuffer(nil)
		}
		return
	}

	// AttachBuffers: and unblock emitting "changed" signal
	func AttachBuffers(TextView *gtk.TextView, BuffTxt *gtk.TextBuffer) {
		tvn.TextView.SetBuffer(BuffTxt)
		BuffTxt.Unref()
		BuffTxt.HandlerUnblock(sigHdlBufTxt) // Use it if you have a callback on "change" signal
	}
*/

package chromaHighlight

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	glco "github.com/hfmrow/genLib/crypto"
	glfs "github.com/hfmrow/genLib/files"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/gotk3/gotk3/gtk"
)

// ChromaHighlight: structure, see below for informations
type ChromaHighlight struct {
	Styles      []string
	Lexers      []string
	Output      []byte
	TextTagList map[string]*gtk.TextTag // Used to store list of used tags in textBuffer
	Formatter   int

	formatter        string
	tagDefList       map[string]bool // Used to create tags definition to RTF textBuffer format
	preExistsTagList map[string]bool // store tags that exists before using highlighter
	intStyles        map[string]bool // internal usage to check if exists
	intLexers        map[string]bool // internal usage to check if exists

	tagDefinition string
	textTagTable  *gtk.TextTagTable
	buff          *gtk.TextBuffer
	regBG, regFG  *regexp.Regexp
	priority      int
	defaultLexerName,
	defaultStyleName,
	lexerName,
	styleName,
	lastMd5 string
	md5SizeAnalyze  int
	initialised     bool // TagList has been initialyzed
	errAlreadyExist error

	// chroma style, used by reset background function.
	chromaStyle *chroma.Style
}

// TODO rewrite with new bench ... !!! (1) - (0) - (2)
// ChromaHighlightNew: Syntax highlighter using Chroma
// syntax highlighter, "github.com/alecthomas/chroma"
// there is three way to highlight using (0), (1), (2) formatters:
// (0)- Use the Tags insetion method, all is done in
// one step and TextBuffer is directly filled after the TextTag
// creation process. No need to use ToTextBuff() method !
// (1)- Used as default, is a three steps method, the first pass
// collect the visual informations, the second pass compile them
// and the 3rd pass build the TextTags and display to textBuffer
// using TextBuffer rich text format' implementation, adventage
// of this method, is that you can save the result to a file that
// it will be reused later. You must use ToTextBuff() method to
// display result in TextView, or get the content of the Output
// variable' struct.
// (2)- The pango version use TextBuffer. InsertMarkup function,
// is a two steps method and need more processing time, as previous,
// you can save result for further usage, and also use output in
// all component that have Markup capabilities ("tooltips", "label",
// "treeview", ...). You must use ToTextBuff() method to display
// result in TextView, or the Output variable' struct.
// Each of them give the same visual result. Bench say
// approximatly: (0) 6.92 x faster than (2), (1) 10,46 x faster
// than (2). With a Golang' source file 290K with "hfmrow" style
// using an old 4*core celeron 2.2ghz, the processing time take:
// (0) -> 4.863sec, (1) -> 3.219sec, (2) -> 33.687sec.
func ChromaHighlightNew(buff *gtk.TextBuffer, formatter ...int) (c *ChromaHighlight, err error) {

	c = new(ChromaHighlight)

	// Get styles & languages lists
	c.Lexers = lexers.Names(false)
	c.Styles = styles.Names()
	c.intStyles = make(map[string]bool)
	c.intLexers = make(map[string]bool)
	for _, name := range c.Lexers {
		c.intLexers[name] = true
	}
	for _, name := range c.Styles {
		c.intStyles[name] = true
	}
	// Set default values
	c.defaultLexerName = "Go"
	c.defaultStyleName = "pygments"

	c.buff = buff
	c.errAlreadyExist = errors.New("Already exist !")

	// formatter option
	if len(formatter) > 0 {
		switch formatter[0] {
		case 1:
			c.formatter = "gtkTextBuffer"
			c.Formatter = 1
		case 2:
			c.formatter = "pango"
			c.Formatter = 2
		default:
			c.formatter = "gtkDirectToTextBuffer"
			c.Formatter = 0
		}
	} else {
		c.formatter = "gtkTextBuffer"
	}
	return
}

// init: initialize the ChromaHighlight structure
// & registering formatters.
func (c *ChromaHighlight) init() (err error) {

	// Option handling registering formatters
	switch c.formatter {
	case "gtkDirectToTextBuffer":
		formatters.Register("gtkDirectToTextBuffer", chroma.FormatterFunc(c.gtkDirectToTextBufferFormatter))
	case "gtkTextBuffer":
		formatters.Register("gtkTextBuffer", chroma.FormatterFunc(c.gtkTextBufferFormatter))
	case "pango":
		formatters.Register("pango", chroma.FormatterFunc(c.pangoFormatter))
	}

	// Used to parse GdkColor
	c.regBG = regexp.MustCompile(`bg:#[a-fA-F|0-9]{6}`)
	c.regFG = regexp.MustCompile(`#[a-fA-F|0-9]{6}`)

	c.RemoveTags()

	c.md5SizeAnalyze = 1024 // Set to 0 mean there is no limit

	c.textTagTable, err = c.buff.GetTagTable()
	return
}

// Highlight: Doing the job and Let there be more light ...
// like a pig on the wings.
func (c *ChromaHighlight) Highlight(inputString, lexerName, styleName string) (err error) {
	var yes bool

	// don't care if error occure, it will be
	// handled by the caller.
	if yes, err = c.sameAsPrevious(&inputString, lexerName, styleName); yes {
		return
	}

	c.initialised = false
	if err = c.init(); err != nil {
		return
	}

	// Option handling pre-formatting
	switch c.formatter {
	case "gtkTextBuffer":
		// Clear tags definitions
		c.tagDefinition = ""
	case "gtkDirectToTextBuffer":
		// Clear TextBuffer
		c.buff.Delete(c.buff.GetStartIter(), c.buff.GetEndIter())
	}

	var buff = new(bytes.Buffer)
	var bufioWriter = bufio.NewWriter(buff)
	if err = c.highlight(
		bufioWriter,
		inputString,
		c.lexerName,
		c.formatter,
		c.styleName); err != nil {
		return
	}
	bufioWriter.Flush()

	// Option handling post-formatting
	switch c.formatter {
	case "gtkTextBuffer":
		err = c.tagsRegistration(c.tagDefinition)

		mime := "GTKTEXTBUFFERCONTENTS-0001"
		header := "<text_view_markup>\n <tags>"
		endTags := "\n </tags>\n"
		startText := "<text>"
		footer := "</text>\n</text_view_markup>"

		defPart := header + c.tagDefinition + endTags + startText + string(buff.Bytes()) + footer
		c.Output = []byte(mime + string(glfs.SizeToBytes(uint32(len(defPart)))) + defPart)
		// Register new tags

	case "gtkDirectToTextBuffer":
		// Nothing to do since the output is only
		// recieved by the given TextBuffer
	case "pango":
		c.Output = buff.Bytes()
	}

	return
}

// ToTextBuff: output to TextBuffer. unified version.
func (c *ChromaHighlight) ToTextBuff() (err error) {
	// Clear TextBuffer
	c.buff.Delete(c.buff.GetStartIter(), c.buff.GetEndIter())

	// Option handling for display
	switch c.formatter {
	case "gtkTextBuffer":
		// Deserialize data
		tags := c.buff.RegisterDeserializeTagset("")
		_, err = c.buff.Deserialize(c.buff, tags, c.buff.GetStartIter(), c.Output)
	case "gtkDirectToTextBuffer":
		c.initialised = true
	case "pango":
		c.buff.InsertMarkup(c.buff.GetStartIter(), string(c.Output))
	}
	return
}

// ToFile: output to file, cannot be used with (0).
func (c *ChromaHighlight) ToFile(filename string) (err error) {
	return ioutil.WriteFile(filename, c.Output, os.ModePerm)
}

// Initialised: return initialisation state
func (c *ChromaHighlight) Initialised() bool {
	return c.initialised
}

// highlight: Determine lexer, formatter and style
// to proceed with.
func (c *ChromaHighlight) highlight(w io.Writer, source, lexer, formatter, style string) (err error) {
	if l := lexers.Get(strings.ToLower(lexer)); l != nil {
		l = chroma.Coalesce(l)
		if f := formatters.Get(formatter); f != nil {
			if c.chromaStyle = styles.Get(style); c.chromaStyle != nil {
				if it, err := l.Tokenise(nil, source); err == nil {
					return f.Format(w, c.chromaStyle, it)
				} else {
					err = errors.New(fmt.Sprintf("styles.Get: %s\n", err.Error()))
				}
			} else {
				err = errors.New("styles.Get")
			}
		} else {
			err = errors.New("formatters.Get")
		}
	} else {
		err = errors.New("lexers.Get")
	}
	return
}

// Clear the background colour.
func (c *ChromaHighlight) chromaClearBackground(style *chroma.Style) *chroma.Style {
	builder := style.Builder()
	bg := builder.Get(chroma.Background)
	bg.Background = 0
	bg.NoInherit = true
	builder.AddEntry(chroma.Background, bg)
	style, _ = builder.Build()
	return style
}

// sameAsPrevious: used to find out if the text has
// changed since the last initialization, used to
// build md5 too. In fact, avoid to uselessly redoing
// that have been already done.
func (c *ChromaHighlight) sameAsPrevious(inputString *string, lexerName, styleName string) (yes bool, err error) {
	var sameAsPreviousText = func() bool {
		sze := len(*inputString) - 1
		// Ther is no size limit if "md5SizeAnalyze" is set to 0
		if c.md5SizeAnalyze > 0 {
			if sze > c.md5SizeAnalyze {
				// limit analyzed data size as specified in structure
				sze = c.md5SizeAnalyze
			}
		}
		currentMd5 := glco.Md5String((*inputString)[:sze])
		if c.lastMd5 == currentMd5 {
			return true
		}
		c.lastMd5 = currentMd5
		return false
	}

	// To not do again if the parameters are the same
	// as the previous ones
	if c.initialised &&
		sameAsPreviousText() &&
		c.lexerName == lexerName &&
		c.styleName == styleName {
		switch c.formatter {
		case "gtkTextBuffer": // don't rebuild all, just register tags again.
			return true, c.tagsRegistration(c.tagDefinition)
		case "gtkDirectToTextBuffer": // tags need to be rebuild anyway.
			return false, nil
		}
		return true, nil
	}
	c.lexerName = lexerName
	c.styleName = styleName

	// build md5.
	sameAsPreviousText()

	// Safe use whether no valid language or style is
	// present, set the default values.
	if !c.intLexers[lexerName] {
		c.lexerName = c.defaultLexerName
	}
	if !c.intStyles[styleName] {
		c.styleName = c.defaultStyleName
	}
	return false, nil
}

/*
	First old version
*/

// // ColorTtyToPango: Convert console tty colors code to Pango
// // markup style. Usually used to highlight source code with:
// // "github.com/alecthomas/chroma"
// // memo:		quick.Highlight(os.Stdout,
// //					someSourceCode,
// //					"go",
// //					"terminal16m",
// //					"monokai")
// func ColorTtyToPango(src *[]byte) (err error) {

// 	var (
// 		rgbColor   string
// 		colors     []string
// 		ttyColor   []byte
// 		rd, vd, bd int
// 		replaceTty *regexp.Regexp
// 	)

// 	escape := string([]byte{0x1b})
// 	color := regexp.MustCompile(`(` + escape + `\[\d{1,3};\d{1,3};\d{1,3};\d{1,3};\d{1,3}m)`)
// 	end := regexp.MustCompile(`(` + escape + `\[0m)`)

// 	// Convert ttyColor to hex form "#AABBCC"
// 	var buidRGB = func(ttyColor []byte) (out string, err error) {
// 		// var outColor string
// 		colors = strings.Split(string(ttyColor), ";")
// 		if rd, err = strconv.Atoi(colors[2]); err == nil {
// 			if vd, err = strconv.Atoi(colors[3]); err == nil {
// 				if bd, err = strconv.Atoi(strings.TrimSuffix(colors[4], "m")); err == nil {
// 					out = `#` + fmt.Sprintf("%02X", rd) + fmt.Sprintf("%02X", vd) + fmt.Sprintf("%02X", bd)
// 				}
// 			}
// 		}

// 		return
// 	}

// 	// Convert ends
// 	*src = end.ReplaceAll(*src, []byte(`</span>`))

// 	// Convert colors
// 	ok := true
// 	for ok {
// 		ttyColor = color.Find(*src)
// 		if len(ttyColor) > 0 {
// 			replaceTty = regexp.MustCompile(`(` + regexp.QuoteMeta(string(ttyColor)) + `)`)
// 			if rgbColor, err = buidRGB(ttyColor); err != nil {
// 				return
// 			}
// 			*src = replaceTty.ReplaceAll(*src, []byte(`<span foreground="`+rgbColor+`">`))
// 			continue
// 		}
// 		ok = false
// 	}
// 	return
// }
