// pangoFormatter.go

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

	"github.com/alecthomas/chroma"
)

// pangoFormatter: is a part of "ChromaHighlight" library
// This is the Pango version, wich not use tags functionality
// but only Pango markup style.
func (c *ChromaHighlight) pangoFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
	var closer string
	style = c.chromaClearBackground(style)
	for tkn := it(); tkn != chroma.EOF; tkn = it() {

		entry := style.Get(tkn.Type)
		if entry.IsZero() {
			fmt.Fprint(w, html.EscapeString(tkn.Value))
		} else {
			if entry.Bold == chroma.Yes {
				fmt.Fprint(w, `<b>`)
				closer = `</b>`
			}
			if entry.Underline == chroma.Yes {
				fmt.Fprint(w, `<u>`)
				closer = `</u>` + closer
			}
			if entry.Italic == chroma.Yes {
				fmt.Fprint(w, `<i>`)
				closer = `</i>` + closer
			}
			if entry.Colour.IsSet() {
				fmt.Fprint(w, `<span foreground="`+
					fmt.Sprintf("#%02X%02X%02X",
						entry.Colour.Red(),
						entry.Colour.Green(),
						entry.Colour.Blue())+`">`)
				closer = `</span>` + closer
			}
			if entry.Background.IsSet() {
				fmt.Fprint(w, `<span background="`+
					fmt.Sprintf("#%02X%02X%02X",
						entry.Background.Red(),
						entry.Background.Green(),
						entry.Background.Blue())+`">`)
				closer = `</span>` + closer
			}
			if entry.Border.IsSet() {
				fmt.Fprint(w, `<span background="`+
					fmt.Sprintf("#%02X%02X%02X",
						entry.Border.Red(),
						entry.Border.Green(),
						entry.Border.Blue())+`">`)
				closer = `</span>` + closer
			}
			fmt.Fprint(w, html.EscapeString(tkn.Value)+closer)
		}
		closer = ""
	}
	return nil
}

// // pangoFormatter: is a part of "ChromaHighlight" library
// // This is the Pango version, wich not use tags functionality
// // but only Pango markup style.
// func (c *ChromaHighlight) pangoFormatter(w io.Writer, style *chroma.Style, it chroma.Iterator) error {
// 	var closer, out string
// 	style = c.chromaClearBackground(style)
// 	for tkn := it(); tkn != chroma.EOF; tkn = it() {

// 		entry := style.Get(tkn.Type)
// 		if !entry.IsZero() {
// 			if entry.Bold == chroma.Yes {
// 				out = `<b>`
// 				closer = `</b>`
// 			}
// 			if entry.Underline == chroma.Yes {
// 				out += `<u>`
// 				closer = `</u>` + closer
// 			}
// 			if entry.Italic == chroma.Yes {
// 				out += `<i>`
// 				closer = `</i>` + closer
// 			}
// 			if entry.Colour.IsSet() {
// 				out += `<span foreground="` +
// 					fmt.Sprintf("#%02X%02X%02X",
// 						entry.Colour.Red(),
// 						entry.Colour.Green(),
// 						entry.Colour.Blue()) + `">`
// 				closer = `</span>` + closer
// 			}
// 			if entry.Background.IsSet() {
// 				out += `<span background="` +
// 					fmt.Sprintf("#%02X%02X%02X",
// 						entry.Background.Red(),
// 						entry.Background.Green(),
// 						entry.Background.Blue()) + `">`
// 				closer = `</span>` + closer
// 			}
// 			if entry.Border.IsSet() {
// 				out += `<span background="` +
// 					fmt.Sprintf("#%02X%02X%02X",
// 						entry.Border.Red(),
// 						entry.Border.Green(),
// 						entry.Border.Blue()) + `">`
// 				closer = `</span>` + closer
// 			}
// 			fmt.Fprint(w, out)
// 		}
// 		fmt.Fprint(w, html.EscapeString(tkn.Value))
// 		if !entry.IsZero() {
// 			fmt.Fprint(w, closer)
// 		}
// 		closer, out = "", ""
// 	}
// 	return nil
// }
