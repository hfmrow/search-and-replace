// filesFunc.go

/*
	Source file auto-generated on Sat, 19 Oct 2019 01:48:39 using Gotk3ObjHandler v1.3.9 ©2018-19 H.F.M
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
	"os"
	"path/filepath"

	glfssf "github.com/hfmrow/genLib/files/scanFileDir"
)

/*****************************/
/* Scan directory functions */
/***************************/
// isSymlinkDir: File is a symlinked directory ?
func isSymlinkDir(slRoot string, slStat os.FileInfo, followSymlinkDir bool) (slIsDir bool, err error) {
	var fName string
	if slStat.Mode()&os.ModeSymlink != 0 && followSymlinkDir {
		if fName, err = os.Readlink(filepath.Join(slRoot, slStat.Name())); err == nil {
			if slStat, err = os.Stat(fName); err == nil {
				if slStat.IsDir() {
					return true, nil
				}
			}
		}
	}
	return false, err
}

// checkUnrecoverableErr: Error come from file not exist or file permission ?
func checkUnrecoverableErr(err error) error {
	if err != nil {
		if !(os.IsPermission(err) || os.IsNotExist(err)) {
			return errors.New(fmt.Sprintf("Got error: %s\n", err.Error()))
		}
	}
	return nil
}

// ScanDirDepth: retrieve files in a specific directory and his sub-directory depending on depth argument.
// depth = -1 mean infinite, depth = 0 mean no sub-dir. optParams: showDir, followSymlinks as bool.
func ScanDirDepth(root string, depth int, optParam ...bool) (files []string, err error) {
	var showDirs, followSl, isDir bool
	switch len(optParam) {
	case 1:
		showDirs = optParam[0]
	case 2:
		showDirs = optParam[0]
		followSl = optParam[1]
	}
	var depthRecurse int
	var tmpFiles []string
	var newFi []os.FileInfo
	var fRoot *os.File
	// Starting scannig directory
	if fRoot, err = os.Open(root); err == nil {
		defer fRoot.Close()
		if newFi, err = fRoot.Readdir(-1); err == nil {
			for _, file := range newFi {
				depthRecurse = depth
				if isDir, err = isSymlinkDir(root, file, followSl); err != nil {
					return files, checkUnrecoverableErr(err)
				}
				if isDir || file.IsDir() { // It's a Dir
					if showDirs {
						files = append(files, filepath.Join(root, file.Name()))
					}
					if depth != 0 {
						depthRecurse--
						if tmpFiles, err = ScanDirDepth(filepath.Join(root, file.Name()), depthRecurse, showDirs, followSl); err == nil {
							files = append(files, tmpFiles...)
						} else {
							return files, checkUnrecoverableErr(err)
						}
					}
				} else { // Not a Dir
					files = append(files, filepath.Join(root, file.Name()))
				}
			}
		} else {
			return files, checkUnrecoverableErr(err)
		}
	} else {
		return files, checkUnrecoverableErr(err)
	}
	return files, err
}

// scanForSubDir: In case where a display refresh is requiered from an existing files/dirs list.
func scanForSubDir(inFilesList []string) (err error) {
	var stat os.FileInfo
	var isDir bool
	var filesList []glfssf.ScanDirFileInfos

	toDispFileList = toDispFileList[:0]

	for idx := len(inFilesList) - 1; idx > -1; idx-- {
		file := inFilesList[idx]
		if len(file) != 0 {

			if stat, err = os.Lstat(file); err != nil {
				return err
			} else if isDir, err = isSymlinkDir(file, stat,
				mainObjects.chkFollowSymlinkDir.GetActive()); !(os.IsPermission(err) || os.IsNotExist(err)) && err != nil {
				return err
			}

			if isDir || stat.IsDir() {
				if filesList, err = glfssf.ScanDirDepthFileInfo(file, mainObjects.spinButtonDepth.GetValueAsInt(),
					false, mainObjects.chkFollowSymlinkDir.GetActive()); os.IsPermission(err) || os.IsNotExist(err) || err == nil {
					toDispFileList = append(toDispFileList, filesList...)
				} else {
					return err
				}
			} else {
				toDispFileList = append(toDispFileList, glfssf.ScanDirFileInfos{Filename: file, FileInfo: stat})
			}
		}
	}
	return err
}
