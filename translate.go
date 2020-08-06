// translate.go

// File generated on Tue, 19 Nov 2019 21:53:02 using Gotk3ObjectsTranslate v1.3 2019 H.F.M

/*
* 	This program comes with absolutely no warranty.
*	See the The MIT License (MIT) for details:
*	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// initGtkObjectsText: read translations from structure and set them to objects.
func (trans *MainTranslate) initGtkObjectsText() {
	trans.setTextToGtkObjects(&mainObjects.btnFind.Widget, "btnFind")
	trans.setTextToGtkObjects(&mainObjects.btnReplaceInClipboard.Widget, "btnReplaceInClipboard")
	trans.setTextToGtkObjects(&mainObjects.btnScan.Widget, "btnScan")
	trans.setTextToGtkObjects(&mainObjects.btnShowClipboard.Widget, "btnShowClipboard")
	trans.setTextToGtkObjects(&mainObjects.chkCaseSensitive.Widget, "chkCaseSensitive")
	trans.setTextToGtkObjects(&mainObjects.chkCharacterClass.Widget, "chkCharacterClass")
	trans.setTextToGtkObjects(&mainObjects.chkCharacterClassStrictMode.Widget, "chkCharacterClassStrictMode")
	trans.setTextToGtkObjects(&mainObjects.chkFollowSymlinkDir.Widget, "chkFollowSymlinkDir")
	trans.setTextToGtkObjects(&mainObjects.chkRegex.Widget, "chkRegex")
	trans.setTextToGtkObjects(&mainObjects.chkUseEscapeChar.Widget, "chkUseEscapeChar")
	trans.setTextToGtkObjects(&mainObjects.chkWholeWord.Widget, "chkWholeWord")
	trans.setTextToGtkObjects(&mainObjects.chkWildcard.Widget, "chkWildcard")
	trans.setTextToGtkObjects(&mainObjects.entryExtMask.Widget, "entryExtMask")
	trans.setTextToGtkObjects(&mainObjects.entryReplace.Widget, "entryReplace")
	trans.setTextToGtkObjects(&mainObjects.entrySearch.Widget, "entrySearch")
	trans.setTextToGtkObjects(&mainObjects.fileChooserBtn.Widget, "fileChooserBtn")
	trans.setTextToGtkObjects(&mainObjects.findWinCancelBtn.Widget, "findWinCancelBtn")
	trans.setTextToGtkObjects(&mainObjects.findWinChkBackUp.Widget, "findWinChkBackUp")
	trans.setTextToGtkObjects(&mainObjects.findWinChkDispForbFiles.Widget, "findWinChkDispForbFiles")
	trans.setTextToGtkObjects(&mainObjects.findWinReplaceBtn.Widget, "findWinReplaceBtn")
	trans.setTextToGtkObjects(&mainObjects.findWinTreeView.Widget, "findWinTreeView")
	trans.setTextToGtkObjects(&mainObjects.imgTop.Widget, "imgTop")
	trans.setTextToGtkObjects(&mainObjects.lblExtMask.Widget, "lblExtMask")
	trans.setTextToGtkObjects(&mainObjects.lblReplace.Widget, "lblReplace")
	trans.setTextToGtkObjects(&mainObjects.lblSearch.Widget, "lblSearch")
	trans.setTextToGtkObjects(&mainObjects.listViewFiles.Widget, "listViewFiles")
	trans.setTextToGtkObjects(&mainObjects.MainButtonOptions.Widget, "MainButtonOptions")
	trans.setTextToGtkObjects(&mainObjects.mainWinBtnClose.Widget, "mainWinBtnClose")
	trans.setTextToGtkObjects(&mainObjects.OptionButtonDone.Widget, "OptionButtonDone")
	trans.setTextToGtkObjects(&mainObjects.OptionsEntryMaxFileSize.Widget, "OptionsEntryMaxFileSize")
	trans.setTextToGtkObjects(&mainObjects.OptionsEntryMinFileSize.Widget, "OptionsEntryMinFileSize")
	trans.setTextToGtkObjects(&mainObjects.OptionsImageTop.Widget, "OptionsImageTop")
	trans.setTextToGtkObjects(&mainObjects.OptionsLabelChooseDepth.Widget, "OptionsLabelChooseDepth")
	trans.setTextToGtkObjects(&mainObjects.OptionsLabelEnableDirChooser.Widget, "OptionsLabelEnableDirChooser")
	trans.setTextToGtkObjects(&mainObjects.OptionsLabelMaxFileSize.Widget, "OptionsLabelMaxFileSize")
	trans.setTextToGtkObjects(&mainObjects.OptionsLabelMinFileSize.Widget, "OptionsLabelMinFileSize")
	trans.setTextToGtkObjects(&mainObjects.spinButtonDepth.Widget, "spinButtonDepth")
	trans.setTextToGtkObjects(&mainObjects.switchFileChooserButton.Widget, "switchFileChooserButton")
	trans.setTextToGtkObjects(&mainObjects.textWinBtnDone.Widget, "textWinBtnDone")
	trans.setTextToGtkObjects(&mainObjects.textWinChkShowModifications.Widget, "textWinChkShowModifications")
	trans.setTextToGtkObjects(&mainObjects.textWinChkSyntxHighlight.Widget, "textWinChkSyntxHighlight")
	trans.setTextToGtkObjects(&mainObjects.textWinChkWrap.Widget, "textWinChkWrap")
	trans.setTextToGtkObjects(&mainObjects.textWinComboBoxLanguage.Widget, "textWinComboBoxLanguage")
	trans.setTextToGtkObjects(&mainObjects.textWinComboBoxTextStyleChooser.Widget, "textWinComboBoxTextStyleChooser")
}
// Translations structure declaration. To be used in main application.
var translate = new(MainTranslate)

// sts: some sentences/words used in the application. Mostly used in Development mode.
// You must add there all sentences used in your application. Or not ...
// They'll be added to language file each time application started
// when "devMode" is set at true.
var sts = map[string]string{
	`regexpErr`: `Regex mistake ...`,
	`done`: `Operation done.`,
	`sbFile`: `File`,
	`savef`: `Save file`,
	`dir-rem`: `Directory does not exist. The current application directory will be used.`,
	`no`: `No`,
	`proceed`: `Are you sure you want to replace pattern in files ?`,
	`sbFiles`: `Files:`,
	`file-LinkNotExist`: `File does not exist, or symlink endpoint not found.`,
	`file-rem`: `File(s) does not exist.`,
	`status`: `Status:`,
	`dispTime`: `Display time:`,
	`file-perm`: `File permissions error.`,
	`in`: `in`,
	`allow`: `Allow`,
	`unexpected`: `An unexpected error occurred`,
	`totalModified`: `Occurrence(s) modified`,
	`titleSearchResults`: `Search results`,
	`removed`: `File has been removed before processing.`,
	`alert`: `Alert ...`,
	`openf`: `Open file`,
	`clpbrdPreview`: `Clipboard content preview.`,
	`totalOccurrences`: `Occurrence(s) found:`,
	`searchTime`: `Search time:`,
	`confirm`: `Confirmation`,
	`missing`: `Something missing`,
	`emptyCB`: `Nothing to do with an empty clipboard ...`,
	`yes`: `Yes`,
	`sbFilesSel`: `Files selected:`,
	`sbFileSel`: `File selected:`,
	`cancel`: `Cancel`,
	`nothingToSearch`: `Nothing to search ...`,
	`file`: `file(s)`,
	`ok`: `Ok`,
	`titlePreviewText`: `Preview window`,
	`forbiddenFiles`: `Some files could not be accessed ...
Unchecking "Follow symlink" may be useful.`,
	`deny`: `Deny`,
	`noFileSel`: `No selected file(s) to search in ...`,
	`notFound`: `Nothing was found ...`,
	`scanTime`: `Scan files time:`,
	`retry`: `Retry`,
}


// Translations structure with methods
type MainTranslate struct {
	// Public
	ProgInfos    progInfo
	Language     language
	Options      parsingFlags
	ObjectsCount int
	Objects      []object
	Sentences    map[string]string
	// Private
	objectsLoaded bool
}

// MainTranslateNew: Initialise new translation structure and assign language file content to GtkObjects.
// devModeActive, indicate that the new sentences must be added to previous language file.
func MainTranslateNew(filename string, devModeActive ...bool) (mt *MainTranslate) {
	var err error
	mt = new(MainTranslate)
	if err = mt.read(filename); err == nil {
		mt.initGtkObjectsText()
		if len(devModeActive) != 0 {
			if devModeActive[0] {
				mt.Sentences = sts
				err := mt.write(filename)
				if err != nil {
					fmt.Printf("%s\n%s\n", "Cannot write actual sentences to language file.", err.Error())
				}
			}
		}
	} else {
		fmt.Printf("%s\n%s\n", "Error loading language file !\nNot an error when you just creating from glade Xml or GOH project file.", err.Error())
	}
	return
}

// readFile: language file.
func (trans *MainTranslate) read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		if err = json.Unmarshal(textFileBytes, &trans); err == nil {
			trans.objectsLoaded = true
		}
	}
	return
}

// Write json datas to file
func (trans *MainTranslate) write(filename string) (err error) {
	var out bytes.Buffer
	var jsonData []byte
	if jsonData, err = json.Marshal(&trans); err == nil && trans.objectsLoaded {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), 0644)
		}
	}
	return
}

type parsingFlags struct {
	SkipLowerCase  bool
	SkipEmptyLabel bool
	SkipEmptyName  bool
	DoBackup       bool
}

type progInfo struct {
	Name              string
	Version           string
	Creat             string
	MainObjStructName string
	GladeXmlFilename  string
	TranslateFilename string
}

type language struct {
	LangNameLong string
	LangNameShrt string
	Author       string
	Date         string
	Updated      string
	Contributors []string
}

type object struct {
	Class         string
	Id            string
	Label         string
	LabelMarkup   bool
	LabelWrap     bool
	Tooltip       string
	TooltipMarkup bool
	Text          string
	Uri           string
	Comment       string
	Idx           int
}

// Define available property within objects
type propObject struct {
	Class         string
	Label         bool
	LabelMarkup   bool
	LabelWrap     bool
	Tooltip       bool
	TooltipMarkup bool
	Text          bool
	Uri           bool
}

// Property that exists for Gtk3 Object ...	(Used for Class capability)
var propPerObjects = []propObject{
	{Class: "GtkButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkToolButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkToggleButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkLabel", Label: true, LabelMarkup: true, Tooltip: true, TooltipMarkup: true, LabelWrap: true},
	{Class: "GtkSpinButton", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkEntry", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkCheckButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkProgressBar", Tooltip: true, TooltipMarkup: true, Text: true},
	{Class: "GtkSearchBar", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkImage", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkRadioButton", Label: true, LabelMarkup: false, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkComboBoxText", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkComboBox", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkLinkButton", Label: true, Tooltip: true, TooltipMarkup: true, Uri: true},
	{Class: "GtkSwitch", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkTreeView", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkFileChooserButton", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkTextView", Tooltip: true, TooltipMarkup: true},
}

// setTextToGtkObjects: read translations from structure and set them to object.
// like this: setTextToGtkObjects(&mainObjects.TransLabelHint.Widget, "TransLabelHint")
func (trans *MainTranslate) setTextToGtkObjects(obj *gtk.Widget, objectId string) {
	for _, currObject := range trans.Objects {
		if currObject.Id == objectId {
			for _, props := range propPerObjects {
				if currObject.Class == props.Class {
					if props.Label {
						obj.SetProperty("label", currObject.Label)
						if props.LabelMarkup {
							obj.SetProperty("use-markup", currObject.LabelMarkup)
							obj.SetProperty("label", strings.ReplaceAll(currObject.Label, "&", "&amp;"))
						}
					}
					if props.LabelWrap {
						obj.SetProperty("wrap", currObject.LabelWrap)
					}
					if props.Tooltip && !currObject.TooltipMarkup {
						obj.SetProperty("tooltip_text", currObject.Tooltip)
					}
					if props.Tooltip && currObject.TooltipMarkup {
						obj.SetProperty("tooltip_markup", strings.ReplaceAll(currObject.Tooltip, "&", "&amp;"))
					}
					if props.Text {
						obj.SetProperty("text", currObject.Text)
					}
					if props.Uri {
						obj.SetProperty("uri", currObject.Uri)
					}
				}
			}
		}
	}
}
