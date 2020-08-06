package chromaHighlight

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

// Hfmrow default theme.
var Hfmrow = styles.Register(chroma.MustNewStyle("hfmrow", chroma.StyleEntries{
	chroma.Whitespace:     "#C0C0C0 ",
	chroma.Comment:        "#D94B00",
	chroma.CommentPreproc: "#BC7A00",
	chroma.CommentSpecial: "bg:#E0E0E0 bold #006000",

	chroma.Keyword:            "bold #000080",
	chroma.KeywordPseudo:      "nobold",
	chroma.KeywordType:        "bold #366F36",
	chroma.KeywordConstant:    "bold #3F0000",
	chroma.KeywordDeclaration: "bold #000080",

	chroma.Operator:     "bold #234567",
	chroma.OperatorWord: "bold #542111",

	chroma.NameBuiltin:   "bold #5F1860",
	chroma.NameFunction:  "bold #18455E",
	chroma.NameClass:     "bold #0000FF",
	chroma.NameNamespace: "bold #0000FF",
	chroma.NameException: "bold #D2413A",
	chroma.NameVariable:  "#366F36",
	chroma.NameConstant:  "#366F36",
	chroma.NameLabel:     "#A0A000",
	chroma.NameEntity:    "bold #999999",
	chroma.NameAttribute: "#7D9029",
	chroma.NameTag:       "bold #008000",
	chroma.NameDecorator: "#AA22FF",

	chroma.String:         "#0E6200",
	chroma.StringDoc:      "italic",
	chroma.StringInterpol: "bold #0E6200",
	chroma.StringEscape:   "bold #0E6266",
	chroma.StringRegex:    "#0E6233",
	chroma.StringSymbol:   "#19177C",
	chroma.StringOther:    "#008000",
	chroma.Number:         "#450099",

	chroma.GenericHeading:    "bold #000080",
	chroma.GenericSubheading: "bold #800080",
	chroma.GenericDeleted:    "#A00000",
	chroma.GenericInserted:   "#00A000",
	chroma.GenericError:      "#FF0000",
	chroma.GenericEmph:       "italic",
	chroma.GenericStrong:     "bold",
	chroma.GenericPrompt:     "bold #000080",
	chroma.GenericOutput:     "#888",
	chroma.GenericTraceback:  "#04D",
	chroma.GenericUnderline:  "underline",

	chroma.Error: "border:#FF0000",
}))
