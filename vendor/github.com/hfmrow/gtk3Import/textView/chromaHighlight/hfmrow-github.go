package chromaHighlight

import (
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
)

// Hfmrow default theme.
var HfmrowGithub = styles.Register(chroma.MustNewStyle("hfmrow-github", chroma.StyleEntries{
	chroma.Whitespace:     "#C0C0C0 ",
	chroma.Comment:        "#6A737D",
	chroma.CommentPreproc: "#6A737D",

	chroma.Keyword:            "#FF1010",
	chroma.KeywordPseudo:      "nobold",
	chroma.KeywordType:        "#E46209",
	chroma.KeywordConstant:    "#FF0000",
	chroma.KeywordDeclaration: "#E46209",

	chroma.Operator:     "#234567",
	chroma.OperatorWord: "#D83A49",

	chroma.NameBuiltin:   "#FF1010",
	chroma.NameFunction:  "#6138AA",
	chroma.NameClass:     "#FF1010",
	chroma.NameNamespace: "#FF1010",
	chroma.NameException: "#D2413A",
	chroma.NameVariable:  "#032F62",
	chroma.NameConstant:  "#032F62",
	chroma.NameLabel:     "#A0A000",
	chroma.NameEntity:    "#999999",
	chroma.NameAttribute: "#7D9029",
	chroma.NameTag:       "#032F62",
	chroma.NameDecorator: "#AA22FF",

	chroma.String:         "#0E6200",
	chroma.StringDoc:      "italic",
	chroma.StringInterpol: "#0E6200",
	chroma.StringEscape:   "#0E6266",
	chroma.StringRegex:    "#0E6233",
	chroma.StringSymbol:   "#19177C",
	chroma.StringOther:    "#032F62",
	chroma.Number:         "#450099",

	chroma.GenericHeading:    "#000080",
	chroma.GenericSubheading: "#800080",
	chroma.GenericDeleted:    "#A00000",
	chroma.GenericInserted:   "#00A000",
	chroma.GenericError:      "#FF0000",
	chroma.GenericEmph:       "italic",
	chroma.GenericStrong:     "nobold",
	chroma.GenericPrompt:     "#000080",
	chroma.GenericOutput:     "#888",
	chroma.GenericTraceback:  "#04D",
	chroma.GenericUnderline:  "underline",

	chroma.Error: "border:#FF0000",
}))
