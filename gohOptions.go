// gohOptions.go

// Source file auto-generated on Fri, 19 Jul 2019 03:46:10 using Gotk3ObjHandler v1.3 ©2019 H.F.M

/*
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	gi "github.com/hfmrow/sAndReplace/gtk3Import"
)

// App infos
var Name = "SearchAndReplace"
var Vers = "v1.7.3"
var Descr = "Search and replace text sequences in one or more files.\nThe clipboard is also included inside replacing functionality."
var Creat = "H.F.M"
var YearCreat = "2018"
var LicenseShort = "This program comes with absolutely no warranty.\nSee the The MIT License (MIT) for details:\nhttps://opensource.org/licenses/mit-license.php"
var LicenseAbrv = "License (MIT)"
var Repository = "github.com/hfmrow/sAndReplace"

// Vars declarations
var absoluteRealPath, optFilename = getAbsRealPath()
var mainOptions *MainOpt
var err error
var tempDir string

var devMode bool
var applyChanges, acceptBinary, cmdLineArg bool
var treeviewSelectedRows []string
var filesCount int
var filesSelected int
var textWinTextToShowBytes []byte

type MainOpt struct {
	/* Public, will be saved and restored */
	MainWinWidth      int
	MainWinHeight     int
	LanguageFilename  string
	AboutOptions      *gi.AboutInfos
	CaseSensitive     bool
	UseRegex          bool
	WholeWord         bool
	Mask              string
	ScanSubDir        bool
	FollowSymlinkDir  bool
	MakeBackup        bool
	UseCharacterClass bool
	UseWildcard       bool
	WrapText          bool

	ChildWindowsKeepAbove bool
	Directory             string
	TreeViewColumns       [][]string
	TreeStoreColumns      [][]string

	TxtAbout                    string
	TxtSearchResult             string
	TxtTextWin                  string
	TxtSomethingMissing         string
	TxtAlert                    string
	TxtFileRemovedBeforeProcess string
	TxtNothingFound             string
	TxtNoFileToSearch           string
	TxtNothingToSearch          string
	TextRegexError              string
	TxtClipboardEmpty           string
	ExtMask                     []string
	/* Private, will NOT be saved */
	currentInFilesList []string
}

// Main options initialisation
func (opt *MainOpt) Init() {
	opt.AboutOptions = new(gi.AboutInfos)

	opt.MainWinWidth = 800
	opt.MainWinHeight = 600
	opt.LanguageFilename = "assets/lang/eng.lang"
	opt.CaseSensitive = false
	opt.Directory = ""
	opt.FollowSymlinkDir = false
	opt.Mask = "*"
	opt.ScanSubDir = false
	opt.TreeViewColumns = [][]string{{"Name", "text"}, {"Size", "text"}, {"Date", "text"}, {"Path", "text"}}
	opt.TreeStoreColumns = [][]string{{"File", "markup"}}
	opt.UseRegex = false
	opt.WholeWord = false
	opt.MakeBackup = true
	opt.UseCharacterClass = false
	opt.UseWildcard = false
	opt.WrapText = false
	opt.ChildWindowsKeepAbove = true
	// Titles win
	opt.TxtSearchResult = "Search result"
	opt.TxtTextWin = "Text editor"
	opt.TxtAbout = "About"
	// Dialog texts
	opt.TxtSomethingMissing = "Something missing"
	opt.TxtAlert = "Alert ..."
	opt.TxtFileRemovedBeforeProcess = "file has been removed before processing."
	opt.TxtNothingFound = "Nothing was found ... "
	opt.TxtNoFileToSearch = "No file(s) to search in ... "
	opt.TxtNothingToSearch = "Nothing to search ... "
	opt.TextRegexError = "Regex mistake ... "
	opt.TxtClipboardEmpty = "Nothing to do with an empty clipboard ..."
}

// Variables -> Objects.
func (opt *MainOpt) UpdateObjects() {
	mainObjects.mainWin.Resize(opt.MainWinWidth, opt.MainWinHeight)
	mainObjects.fileChooserBtn.FileChooser.SetCurrentFolder(opt.Directory)
	mainObjects.chkCaseSensitive.ToggleButton.SetActive(opt.CaseSensitive)
	mainObjects.chkFollowSymlinkDir.ToggleButton.SetActive(opt.FollowSymlinkDir)
	mainObjects.chkSubDir.ToggleButton.SetActive(opt.ScanSubDir)
	mainObjects.chkRegex.ToggleButton.SetActive(opt.UseRegex)
	mainObjects.chkWoleWord.ToggleButton.SetActive(opt.WholeWord)
	mainObjects.entryExtMask.SetText(opt.Mask)
	mainObjects.chkWildcard.SetActive(opt.UseWildcard)
	mainObjects.chkCharacterClass.SetActive(opt.UseCharacterClass)
	mainObjects.textWinChkWrap.SetActive(opt.WrapText)
	mainObjects.findWinChkBackUp.SetActive(mainOptions.MakeBackup)
}

// Objects -> Variables.
func (opt *MainOpt) UpdateOptions() {
	opt.MainWinWidth, opt.MainWinHeight = mainObjects.mainWin.GetSize()
	opt.Directory = mainObjects.fileChooserBtn.FileChooser.GetFilename()
	opt.CaseSensitive = mainObjects.chkCaseSensitive.ToggleButton.GetActive()
	opt.FollowSymlinkDir = mainObjects.chkFollowSymlinkDir.ToggleButton.GetActive()
	opt.ScanSubDir = mainObjects.chkSubDir.ToggleButton.GetActive()
	opt.UseRegex = mainObjects.chkRegex.ToggleButton.GetActive()
	opt.WholeWord = mainObjects.chkWoleWord.ToggleButton.GetActive()
	opt.UseWildcard = mainObjects.chkWildcard.GetActive()
	opt.UseCharacterClass = mainObjects.chkCharacterClass.GetActive()
	opt.WrapText = mainObjects.textWinChkWrap.GetActive()
	opt.MakeBackup = mainObjects.findWinChkBackUp.GetActive()
	opt.Mask = getEntryText(mainObjects.entryExtMask)
}

// Read Options from file
func (opt *MainOpt) Read() (err error) {
	var textFileBytes []byte
	if textFileBytes, err = ioutil.ReadFile(optFilename); err == nil {
		err = json.Unmarshal(textFileBytes, &opt)
	}
	return err
}

// Write Options to file
func (opt *MainOpt) Write() (err error) {
	var out bytes.Buffer
	var jsonData []byte
	opt.UpdateOptions()
	if jsonData, err = json.Marshal(&opt); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(optFilename, out.Bytes(), 0644)
		}
	}
	return err
}
