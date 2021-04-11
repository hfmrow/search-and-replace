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
)

func getArguments() (treeviewSelectedRows []string, entrySearchText, entryReplaceText string, err error) {
	var rows [][]string

	entrySearchText = GetEntryText(obj.entrySearch)
	entryReplaceText = GetEntryText(obj.entryReplace)

	if len(entrySearchText) == 0 {
		err = errors.New(sts["nothingToSearch"])
		return
	}

	// if there is no row selected, consider there are all selected.
	if tvsList.Selection.CountSelectedRows() == 0 {
		// err = errors.New(sts["noFileSel"])
		// return
		treeviewSelectedRows, err = tvsList.StoreColToStringSl(opt.mapListStore["pathReal"])

	} else if rows, err = tvsList.GetSelectedRows(); err == nil {
		for _, row := range rows {
			treeviewSelectedRows = append(treeviewSelectedRows, row[opt.mapListStore["pathReal"]])
		}
	}
	return
}

func searchAndReplace(treeviewSelectedRows []string,
	entrySearchText, entryReplaceText string,
	removeEmptyResult, replace bool) (mainFound []glfsft.SearchAndReplaceFiles, occurrences int, err error) {

	return SearchAndReplaceInFiles(
		treeviewSelectedRows,
		entrySearchText,
		entryReplaceText,
		opt.FileMinSizeLimit, // Size limit >
		opt.FileMaxSizeLimit, // Size limit <
		obj.chkCaseSensitive.GetActive(),
		obj.chkCharacterClass.GetActive(),
		obj.chkCharacterClassStrictMode.GetActive(),
		obj.chkRegex.GetActive(),
		obj.chkWildcard.GetActive(),
		obj.chkUseEscapeChar.GetActive(),
		obj.chkUseEscapeCharToReplace.GetActive(),
		obj.chkWholeWord.GetActive(),
		replace, // DoReplace
		acceptBinary,
		removeEmptyResult)
}

func displayResults(occurrences int) (countFiles int) {
	if occurrences > 0 {
		// Prepare displaying results window
		if !alreadyPlacedFoundWin {
			x, y := obj.mainWin.GetPosition()
			obj.findWin.Move(x+opt.CascadeDepth, y+opt.CascadeDepth)
			obj.findWin.Resize(obj.mainWin.GetSize())
			alreadyPlacedFoundWin = true
		}

		countFiles = showResults(&filesFoundMulti)
		findWinTitle.Update([]string{fmt.Sprintf("%s %d %s %d %s", sts["totalOccurrences"], occurrences, sts["in"], countFiles, sts["file"])})

		BringToFront(obj.findWin)
		return
	}
	DialogMessage(obj.mainWin, "warning", sts["done"], "\n\n"+sts["notFound"], "", "Ok")
	return
}

// onTheFlySearch:
func onTheFlySearch(inTextBytes []byte, doReplace bool) (outTextBytes []byte, err error) {

	if len(inTextBytes) != 0 {
		if !fileFoundSingle.IsReadyToReplace() {
			fileFoundSingle = SearchAndReplaceNew("", []byte{}, "", "")
		}

		fileFoundSingle.Init(
			inTextBytes,
			fmt.Sprint(GetEntryText(obj.entrySearch)),
			fmt.Sprint(GetEntryText(obj.entryReplace)),
			obj.chkCaseSensitive.GetActive(),
			obj.chkCharacterClass.GetActive(),
			obj.chkCharacterClassStrictMode.GetActive(),
			obj.chkRegex.GetActive(),
			obj.chkWildcard.GetActive(),
			obj.chkUseEscapeChar.GetActive(),
			obj.chkUseEscapeCharToReplace.GetActive(),
			obj.chkWholeWord.GetActive(),
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
