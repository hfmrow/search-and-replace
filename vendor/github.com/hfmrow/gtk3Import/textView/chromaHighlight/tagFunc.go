// tagFunc.go

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

import "github.com/gotk3/gotk3/gtk"

// RemoveTags: reset all tags and maps
func (c *ChromaHighlight) RemoveTags() {
	for tagName, _ := range c.TextTagList {
		c.removeExistingTag(tagName)
	}
	c.initTagsMaps()
}

// initTagsMaps: Initialyze/reset tags maps
func (c *ChromaHighlight) initTagsMaps() {
	c.TextTagList = make(map[string]*gtk.TextTag)
	c.tagDefList = make(map[string]bool)
	c.preExistsTagList = make(map[string]bool)
}

// createTag: create tag with properties and add it to buffer.
// Check wether the tag already exist in this case, return it.
func (c *ChromaHighlight) createTag(tagName string, props map[string]interface{}) (tag *gtk.TextTag) {
	if tag = c.TextTagList[tagName]; tag == nil {
		tag = c.buff.CreateTag(tagName, props) // add tag & properties
	}
	return
}

// removeExistingTag: from buffer & TextTagTable if exists.
func (c *ChromaHighlight) removeExistingTag(tagName string) {
	if tag, ok := c.lookupExistingTag(tagName); ok {
		c.buff.RemoveTag(tag,
			c.buff.GetStartIter(),
			c.buff.GetEndIter())
		c.textTagTable.Remove(tag)
	}
	return
}

// LookupExistingTag:
func (c *ChromaHighlight) lookupExistingTag(tagName string) (tag *gtk.TextTag, ok bool) {
	if tag, _ := c.textTagTable.Lookup(tagName); tag != nil { // Check wether tag exists
		return tag, true
	}
	return
}
