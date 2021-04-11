// translate.go

// File generated on Sat, 10 Apr 2021 22:05:31 using Gotk3 Objects Translate v1.5.2 2019-21 hfmrow

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
	trans.setTextToGtkObjects(&obj.btnFind.Widget, "btnFind")
	trans.setTextToGtkObjects(&obj.btnReplaceInClipboard.Widget, "btnReplaceInClipboard")
	trans.setTextToGtkObjects(&obj.btnScan.Widget, "btnScan")
	trans.setTextToGtkObjects(&obj.btnShowClipboard.Widget, "btnShowClipboard")
	trans.setTextToGtkObjects(&obj.chkCaseSensitive.Widget, "chkCaseSensitive")
	trans.setTextToGtkObjects(&obj.chkCharacterClass.Widget, "chkCharacterClass")
	trans.setTextToGtkObjects(&obj.chkCharacterClassStrictMode.Widget, "chkCharacterClassStrictMode")
	trans.setTextToGtkObjects(&obj.chkFollowSymlinkDir.Widget, "chkFollowSymlinkDir")
	trans.setTextToGtkObjects(&obj.chkRegex.Widget, "chkRegex")
	trans.setTextToGtkObjects(&obj.chkUseEscapeChar.Widget, "chkUseEscapeChar")
	trans.setTextToGtkObjects(&obj.chkUseEscapeCharToReplace.Widget, "chkUseEscapeCharToReplace")
	trans.setTextToGtkObjects(&obj.chkWholeWord.Widget, "chkWholeWord")
	trans.setTextToGtkObjects(&obj.chkWildcard.Widget, "chkWildcard")
	trans.setTextToGtkObjects(&obj.entryExtMask.Widget, "entryExtMask")
	trans.setTextToGtkObjects(&obj.entryReplace.Widget, "entryReplace")
	trans.setTextToGtkObjects(&obj.entrySearch.Widget, "entrySearch")
	trans.setTextToGtkObjects(&obj.fileChooserBtn.Widget, "fileChooserBtn")
	trans.setTextToGtkObjects(&obj.findWinBtnDeselect.Widget, "findWinBtnDeselect")
	trans.setTextToGtkObjects(&obj.findWinBtnInvertSel.Widget, "findWinBtnInvertSel")
	trans.setTextToGtkObjects(&obj.findWinCancelBtn.Widget, "findWinCancelBtn")
	trans.setTextToGtkObjects(&obj.findWinChkBackUp.Widget, "findWinChkBackUp")
	trans.setTextToGtkObjects(&obj.findWinChkDispForbFiles.Widget, "findWinChkDispForbFiles")
	trans.setTextToGtkObjects(&obj.findWinChkExpandAll.Widget, "findWinChkExpandAll")
	trans.setTextToGtkObjects(&obj.findWinReplaceBtn.Widget, "findWinReplaceBtn")
	trans.setTextToGtkObjects(&obj.findWinScrollWin.Widget, "findWinScrollWin")
	trans.setTextToGtkObjects(&obj.findWinTreeView.Widget, "findWinTreeView")
	trans.setTextToGtkObjects(&obj.imgTop.Widget, "imgTop")
	trans.setTextToGtkObjects(&obj.lblExtMask.Widget, "lblExtMask")
	trans.setTextToGtkObjects(&obj.lblReplace.Widget, "lblReplace")
	trans.setTextToGtkObjects(&obj.lblSearch.Widget, "lblSearch")
	trans.setTextToGtkObjects(&obj.listViewFiles.Widget, "listViewFiles")
	trans.setTextToGtkObjects(&obj.MainButtonOptions.Widget, "MainButtonOptions")
	trans.setTextToGtkObjects(&obj.mainWinBtnClose.Widget, "mainWinBtnClose")
	trans.setTextToGtkObjects(&obj.OptionButtonDone.Widget, "OptionButtonDone")
	trans.setTextToGtkObjects(&obj.OptionsEntryMaxFileSize.Widget, "OptionsEntryMaxFileSize")
	trans.setTextToGtkObjects(&obj.OptionsEntryMinFileSize.Widget, "OptionsEntryMinFileSize")
	trans.setTextToGtkObjects(&obj.OptionsImageTop.Widget, "OptionsImageTop")
	trans.setTextToGtkObjects(&obj.OptionsLabelChooseDepth.Widget, "OptionsLabelChooseDepth")
	trans.setTextToGtkObjects(&obj.OptionsLabelEnableDirChooser.Widget, "OptionsLabelEnableDirChooser")
	trans.setTextToGtkObjects(&obj.OptionsLabelMaxFileSize.Widget, "OptionsLabelMaxFileSize")
	trans.setTextToGtkObjects(&obj.OptionsLabelMinFileSize.Widget, "OptionsLabelMinFileSize")
	trans.setTextToGtkObjects(&obj.scrolledWindowTreeView.Widget, "scrolledWindowTreeView")
	trans.setTextToGtkObjects(&obj.SourceToggleButtonMapWidth.Widget, "SourceToggleButtonMapWidth")
	trans.setTextToGtkObjects(&obj.spinButtonDepth.Widget, "spinButtonDepth")
	trans.setTextToGtkObjects(&obj.statusbar.Widget, "statusbar")
	trans.setTextToGtkObjects(&obj.switchFileChooserButton.Widget, "switchFileChooserButton")
	trans.setTextToGtkObjects(&obj.textWinBtnDone.Widget, "textWinBtnDone")
	trans.setTextToGtkObjects(&obj.textWinChkShowModifications.Widget, "textWinChkShowModifications")
	trans.setTextToGtkObjects(&obj.textWinChkSyntxHighlight.Widget, "textWinChkSyntxHighlight")
	trans.setTextToGtkObjects(&obj.textWinChkWrap.Widget, "textWinChkWrap")
	trans.setTextToGtkObjects(&obj.textWinComboBoxLanguage.Widget, "textWinComboBoxLanguage")
	trans.setTextToGtkObjects(&obj.textWinComboBoxTextStyleChooser.Widget, "textWinComboBoxTextStyleChooser")
	trans.setTextToGtkObjects(&obj.View.Widget, "View")
}
// Translations structure declaration. To be used in main application.
var translate = new(MainTranslate)

// sts: some sentences/words used in the application. Mostly used in Development mode.
// You must add there all sentences used in your application. Or not ...
// They'll be added to language file each time application started
// when "devMode" is set at true.
var sts = map[string]string{
	`scanTime`: `Scan files time:`,
	`totalModified`: `Occurrence(s) modified`,
	`allow`: `Allow`,
	`file-LinkNotExist`: `File does not exist, or symlink endpoint not found.`,
	`sbFiles`: `Files:`,
	`status`: `Status:`,
	`deny`: `Deny`,
	`savef`: `Save file`,
	`regexpErr`: `Regex mistake ...`,
	`titleSearchResults`: `Search results`,
	`removed`: `File has been removed before processing.`,
	`file-rem`: `File(s) does not exist.`,
	`file`: `file(s)`,
	`noFileSel`: `No selected file(s) to search in ...`,
	`dispTime`: `Display time:`,
	`searchTime`: `Search time:`,
	`missing`: `Something missing`,
	`totalOccurrences`: `Occurrence(s) found:`,
	`in`: `in`,
	`yes`: `Yes`,
	`clpbrdPreview`: `Clipboard content preview.`,
	`no`: `No`,
	`nothingToSearch`: `Nothing to search ...`,
	`openf`: `Open file`,
	`forbiddenFiles`: `Some files could not be accessed ...
Unchecking "Follow symlink" may be useful.`,
	`notFound`: `Nothing was found ...`,
	`sbFile`: `File`,
	`alert`: `Alert ...`,
	`sbFilesSel`: `Files selected:`,
	`titlePreviewText`: `Preview window`,
	`file-perm`: `File permissions error.`,
	`done`: `Operation done.`,
	`ok`: `Ok`,
	`sbFileSel`: `File selected:`,
	`cancel`: `Cancel`,
	`dir-rem`: `Directory does not exist. The current application directory will be used.`,
	`unexpected`: `An unexpected error occurred`,
	`confirm`: `Confirmation`,
	`proceed`: `Are you sure you want to replace pattern in files ?`,
	`emptyCB`: `Nothing to do with an empty clipboard ...`,
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
	SkipLowerCase,
	SkipEmptyLabel,
	SkipEmptyName,
	DoBackup bool
}

type progInfo struct {
	Name,
	Version,
	Creat,
	MainObjStructName,
	GladeXmlFilename,
	TranslateFilename,
	ProjectRootDir,
	GohProjFile string
}

type language struct {
	LangNameLong,
	LangNameShrt,
	Author,
	Date,
	Updated string
	Contributors []string
}

type object struct {
	Class,
	Id,
	Label,
	Tooltip,
	Text,
	Uri,
	Comment string
	LabelMarkup,
	LabelWrap,
	TooltipMarkup bool
	Idx int
}

// Define available property within objects
type propObject struct {
	Class string
	Label,
	LabelMarkup,
	LabelWrap,
	Tooltip,
	TooltipMarkup,
	Text,
	Uri bool
}

// Property that exists for Gtk3 Object ...	(Used for Class capability)
var propPerObjects = []propObject{
	{Class: "GtkButton", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkMenuButton", Label: true, Tooltip: true, TooltipMarkup: true},

	// {Class: "GtkToolButton", Label: true, Tooltip: true, TooltipMarkup: true},    // Deprecated since 3.10
	// {Class: "GtkImageMenuItem", Label: true, Tooltip: true, TooltipMarkup: true}, // Deprecated since 3.10

	{Class: "GtkMenuItem", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkCheckMenuItem", Label: true, Tooltip: true, TooltipMarkup: true},
	{Class: "GtkRadioMenuItem", Label: true, Tooltip: true, TooltipMarkup: true},

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
	{Class: "GtkSourceView", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkStatusbar", Tooltip: true, TooltipMarkup: true},
	{Class: "GtkScrolledWindow", Tooltip: true, TooltipMarkup: true},
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
					break
				}
			}
			break
		}
	}
}
