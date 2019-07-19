// translate.go

// File generated on Fri, 19 Jul 2019 03:17:23 using Gotk3ObjectsTranslate v1.0 2019 H.F.M

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
	"os"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

// initGtkObjectsText: read translations from structure and set them to objects.
func (trans *MainTranslate) initGtkObjectsText() {
	trans.setTextToGtkObjects(&mainObjects.btnFind.Widget, "btnFind")
	trans.setTextToGtkObjects(&mainObjects.btnReplaceInClipboard.Widget, "btnReplaceInClipboard")
	trans.setTextToGtkObjects(&mainObjects.btnShowClipboard.Widget, "btnShowClipboard")
	trans.setTextToGtkObjects(&mainObjects.chkCaseSensitive.Widget, "chkCaseSensitive")
	trans.setTextToGtkObjects(&mainObjects.chkCharacterClass.Widget, "chkCharacterClass")
	trans.setTextToGtkObjects(&mainObjects.chkCharacterClassStrictMode.Widget, "chkCharacterClassStrictMode")
	trans.setTextToGtkObjects(&mainObjects.chkFollowSymlinkDir.Widget, "chkFollowSymlinkDir")
	trans.setTextToGtkObjects(&mainObjects.chkRegex.Widget, "chkRegex")
	trans.setTextToGtkObjects(&mainObjects.chkSubDir.Widget, "chkSubDir")
	trans.setTextToGtkObjects(&mainObjects.chkUseEscapeChar.Widget, "chkUseEscapeChar")
	trans.setTextToGtkObjects(&mainObjects.chkWildcard.Widget, "chkWildcard")
	trans.setTextToGtkObjects(&mainObjects.chkWoleWord.Widget, "chkWoleWord")
	trans.setTextToGtkObjects(&mainObjects.entryExtMask.Widget, "entryExtMask")
	trans.setTextToGtkObjects(&mainObjects.entryReplace1.Widget, "entryReplace1")
	trans.setTextToGtkObjects(&mainObjects.entrySearch1.Widget, "entrySearch1")
	trans.setTextToGtkObjects(&mainObjects.fileChooserBtn.Widget, "fileChooserBtn")
	trans.setTextToGtkObjects(&mainObjects.findWinCancelBtn.Widget, "findWinCancelBtn")
	trans.setTextToGtkObjects(&mainObjects.findWinChkBackUp.Widget, "findWinChkBackUp")
	trans.setTextToGtkObjects(&mainObjects.findWinLabel.Widget, "findWinLabel")
	trans.setTextToGtkObjects(&mainObjects.findWinLabel1.Widget, "findWinLabel1")
	trans.setTextToGtkObjects(&mainObjects.findWinLabel2.Widget, "findWinLabel2")
	trans.setTextToGtkObjects(&mainObjects.findWinReplaceBtn.Widget, "findWinReplaceBtn")
	trans.setTextToGtkObjects(&mainObjects.findWinTreeView.Widget, "findWinTreeView")
	trans.setTextToGtkObjects(&mainObjects.imgTop.Widget, "imgTop")
	trans.setTextToGtkObjects(&mainObjects.lblExtMask.Widget, "lblExtMask")
	trans.setTextToGtkObjects(&mainObjects.lblReplace.Widget, "lblReplace")
	trans.setTextToGtkObjects(&mainObjects.lblSearch.Widget, "lblSearch")
	trans.setTextToGtkObjects(&mainObjects.SwitchFileChooserButton.Widget, "SwitchFileChooserButton")
	trans.setTextToGtkObjects(&mainObjects.textWinBtnDone.Widget, "textWinBtnDone")
	trans.setTextToGtkObjects(&mainObjects.textWinChkShowModifications.Widget, "textWinChkShowModifications")
	trans.setTextToGtkObjects(&mainObjects.textWinChkWrap.Widget, "textWinChkWrap")
	trans.setTextToGtkObjects(&mainObjects.treeviewFiles.Widget, "treeviewFiles")
}
// Translations structure declaration. To be used in main application.
var translate = new(MainTranslate)

// sts: some sentences/words used in the application. Mostly used in Development mode.
// You must add there all sentences used in your application. Or not ...
// They'll be added to language file each time application started
// when "devMode" is set at true.
var sts = map[string]string{
	`savef`: `Save file`,
	`file-rem`: `File(s) does not exist.`,
	`retry`: `Retry`,
	`yes`: `Yes`,
	`ok`: `Ok`,
	`deny`: `Deny`,
	`no`: `No`,
	`openf`: `Open file`,
	`allow`: `Allow`,
	`cancel`: `Cancel`,
	`dir-rem`: `Directory does not exist. The current application directory will be used.`,
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
	mt = new(MainTranslate)
	if _, err := os.Stat(filename); err == nil {
		mt.read(filename)
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
		fmt.Printf("%s\n%s\n", "Error loading language file !\nNot an error when is just creating from glade Xml or GOH project file.", err.Error())
	}
	return mt
}

// readFile: language file.
func (trans *MainTranslate) read(filename string) (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		if err = json.Unmarshal(textFileBytes, &trans); err == nil {
			trans.objectsLoaded = true
		}
	}
	return err
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
	return err
}

type parsingFlags struct {
	SkipLowerCase       bool
	SkipEmptyLabel      bool
	DoBackup            bool
}

type progInfo struct {
	Name                 string
	Version              string
	Creat                string
	MainObjStructName    string
	GladeXmlFilename     string
	TranslateFilename    string
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
	Class   string
	Id      string
	Label   string
	Tooltip string
	Text    string
	Uri     string
	Wrap    bool
	Markup  bool
	Comment string
}

// Define available property within objects
type propObject struct {
	Class   string
	Label   bool
	Tooltip bool
	Markup  bool
	Text    bool
	Uri     bool
	Wrap    bool
}

// Property that exists for GObject ...	(Used for Class)
var propPerObjects = []propObject{
	{Class: "GtkButton", Label: true, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkToggleButton", Label: true, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkLabel", Label: true, Tooltip: true, Markup: true, Text: false, Uri: false, Wrap: true},
	{Class: "GtkSpinButton", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkEntry", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkCheckButton", Label: true, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkProgressBar", Label: false, Tooltip: true, Markup: true, Text: true, Uri: false},
	{Class: "GtkSearchBar", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkImage", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkRadioButton", Label: true, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkComboBoxText", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkComboBox", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkLinkButton", Label: true, Tooltip: true, Markup: true, Text: false, Uri: true},
	{Class: "GtkSwitch", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkTreeView", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
	{Class: "GtkFileChooserButton", Label: false, Tooltip: true, Markup: true, Text: false, Uri: false},
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
					}
					if props.Tooltip && !currObject.Markup {
						obj.SetProperty("tooltip_text", currObject.Tooltip)
					}
					if props.Tooltip && currObject.Markup {
						obj.SetProperty("tooltip_markup", strings.ReplaceAll(currObject.Tooltip, "&", "&amp;"))
					}
					if props.Text {
						obj.SetProperty("text", currObject.Text)
					}
					if props.Uri {
						obj.SetProperty("uri", currObject.Uri)
					}
					if props.Wrap {
						obj.SetProperty("wrap", currObject.Wrap)
					}
				}
			}
		}
	}
}
