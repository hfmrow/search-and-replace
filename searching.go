// searching.go

/*
	Source file auto-generated on Thu, 17 Oct 2019 22:43:58 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
	This software use gotk3 that is licensed under the ISC License:
	https://github.com/gotk3/gotk3/blob/master/LICENSE

	Copyright ©2018-19 H.F.M - Search And Replace
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package main

import (
	"errors"
	"fmt"

	glfsft "github.com/hfmrow/genLib/files/fileText"

	gidg "github.com/hfmrow/gtk3Import/dialog"
	gitl "github.com/hfmrow/gtk3Import/tools"
)

func getArguments() (treeviewSelectedRows []string, entrySearchText, entryReplaceText string, err error) {
	var rows [][]string

	entrySearchText = gitl.GetEntryText(mainObjects.entrySearch)
	entryReplaceText = gitl.GetEntryText(mainObjects.entryReplace)

	if len(entrySearchText) == 0 {
		err = errors.New(sts["nothingToSearch"])
		return
	}

	if tvsList.Selection.CountSelectedRows() == 0 {
		err = errors.New(sts["noFileSel"])
		return
	}

	if rows, err = tvsList.GetSelectedRows(); err == nil {
		for _, row := range rows {
			treeviewSelectedRows = append(treeviewSelectedRows, row[3])
		}
	}
	return
}

func searchAndReplace(treeviewSelectedRows []string,
	entrySearchText, entryReplaceText string,
	removeEmptyResult, replace bool) (mainFound []glfsft.SearchAndReplaceFiles, occurrences int, err error) {

	return glfsft.SearchAndReplaceInFiles(
		treeviewSelectedRows,
		entrySearchText,
		entryReplaceText,
		mainOptions.FileMinSizeLimit, // Size limit >
		mainOptions.FileMaxSizeLimit, // Size limit <
		mainObjects.chkCaseSensitive.GetActive(),
		mainObjects.chkCharacterClass.GetActive(),
		mainObjects.chkCharacterClassStrictMode.GetActive(),
		mainObjects.chkRegex.GetActive(),
		mainObjects.chkWildcard.GetActive(),
		mainObjects.chkUseEscapeChar.GetActive(),
		mainObjects.chkUseEscapeCharToReplace.GetActive(),
		mainObjects.chkWholeWord.GetActive(),
		replace, // DoReplace
		replace, // DoSave
		mainOptions.MakeBackup,
		acceptBinary,
		removeEmptyResult)
}

func displayResults(occurrences int) (countFiles int) {
	if occurrences > 0 {
		// Prepare displaying results window
		if !alreadyPlacedFoundWin {
			x, y := mainObjects.mainWin.GetPosition()
			mainObjects.findWin.Move(x+mainOptions.CascadeDepth, y+mainOptions.CascadeDepth)
			mainObjects.findWin.Resize(mainObjects.mainWin.GetSize())
			alreadyPlacedFoundWin = true
		}

		countFiles = showResults(&filesFoundMulti)
		findWinTitle.Update([]string{fmt.Sprintf("%s %d %s %d %s", sts["totalOccurrences"], occurrences, sts["in"], countFiles, sts["file"])})

		BringToFront(mainObjects.findWin)
		return
	}
	gidg.DialogMessage(mainObjects.mainWin, "warning", sts["done"], "\n\n"+sts["notFound"], "", "Ok")
	return
}

// onTheFlySearch:
func onTheFlySearch(inTextBytes []byte, doReplace bool) (outTextBytes []byte, err error) {

	if len(inTextBytes) != 0 {
		if !fileFoundSingle.ReadyToReplace() {
			fileFoundSingle = glfsft.SearchAndReplaceNew([]byte{}, "", "")
		}

		fileFoundSingle.Init(
			inTextBytes,
			fmt.Sprint(gitl.GetEntryText(mainObjects.entrySearch)),
			fmt.Sprint(gitl.GetEntryText(mainObjects.entryReplace)),
			mainObjects.chkCaseSensitive.GetActive(),
			mainObjects.chkCharacterClass.GetActive(),
			mainObjects.chkCharacterClassStrictMode.GetActive(),
			mainObjects.chkRegex.GetActive(),
			mainObjects.chkWildcard.GetActive(),
			mainObjects.chkUseEscapeChar.GetActive(),
			mainObjects.chkUseEscapeCharToReplace.GetActive(),
			mainObjects.chkWholeWord.GetActive(),
			doReplace)

		// Return text after replacement (i.e: clipboard replace)
		if err = fileFoundSingle.SearchAndReplace(); err == nil {
			outTextBytes = fileFoundSingle.TextBytes

			if doReplace {
				// Get back original TextBytes, we don't want a real replacement ...
				fileFoundSingle.TextBytes = inTextBytes
			}
		}
	}
	return
}
