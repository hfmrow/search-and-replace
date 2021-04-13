// scanFileDir.go

package scanFileDir

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Scan dir and subdir to get symlinks with specified endpoint.
func RecurseScanSymlink(path, linkEndPoint string) (fileList []string, err error) {
	err = filepath.Walk(path, func(filePath string, fileInfos os.FileInfo, err error) error {
		if fileInfos.Mode()&os.ModeSymlink != 0 {
			realPath, err := os.Readlink(filePath)
			if err != nil {
				return err
			}
			if strings.Contains(realPath, linkEndPoint) {
				fileList = append(fileList, strings.Replace(filePath, path, "", -1))
			}
		}
		return nil
	})
	if err != nil {
		return fileList, err
	}
	return fileList, nil
}

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
// depth = -1 means infinite, depth = 0 means no sub-dir. optParams: showDir, followSymlinks as bool.
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

type ScanDirFileInfos struct {
	Filename string
	FileInfo os.FileInfo
}

// ScanDirDepthFileInfo: retrieve files in a specific directory and his sub-directory depending on depth argument.
// depth = -1 means infinite, depth = 0 means no sub-dir. optParams: showDir, followSymlinks as bool.
// return a structure that contain filename and os.FileInfo.
func ScanDirDepthFileInfo(root string, depth int, optParam ...bool) (files []ScanDirFileInfos, err error) {
	var showDirs, followSl, isDir bool
	switch len(optParam) {
	case 1:
		showDirs = optParam[0]
	case 2:
		showDirs = optParam[0]
		followSl = optParam[1]
	}
	var depthRecurse int
	var tmpFiles []ScanDirFileInfos
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
						files = append(files, ScanDirFileInfos{Filename: filepath.Join(root, file.Name()), FileInfo: file})
					}
					if depth != 0 {
						depthRecurse--
						if tmpFiles, err = ScanDirDepthFileInfo(filepath.Join(root, file.Name()), depthRecurse, showDirs, followSl); err == nil {
							files = append(files, tmpFiles...)
						} else {
							return files, checkUnrecoverableErr(err)
						}
					}
				} else { // Not a Dir
					files = append(files, ScanDirFileInfos{Filename: filepath.Join(root, file.Name()), FileInfo: file})
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

// ScanSubDir: retrieve files in a specific directory and his sub-directory.
// does not follow symlink (walk)
func ScanSubDir(root string, showDirs ...bool) (files []string, err error) {
	var listDir bool
	if len(showDirs) != 0 {
		listDir = showDirs[0]
	}
	err = filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, path)
			} else if listDir {
				files = append(files, path)
			}
			return nil
		})
	return files, err
}

// ScanDir: retrieve filenames in a specific directory
func ScanDir(root string, showDirs ...bool) (files []string, err error) {
	var listDir bool
	if len(showDirs) != 0 {
		listDir = showDirs[0]
	}
	filesInfos, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range filesInfos {
		if !file.IsDir() {
			files = append(files, filepath.Join(root, file.Name()))
		} else if listDir {
			files = append(files, filepath.Join(root, file.Name()))
		}
	}
	return files, err
}

// ScanDirFileInfo: retrieve []os.FileInfo from a defined directory
func ScanDirFileInfo(root string) (infosF []os.FileInfo, err error) {
	var fRoot *os.File
	if fRoot, err = os.Open(root); err == nil {
		defer fRoot.Close()
		infosF, err = fRoot.Readdir(-1)
	}
	return
}

// ScanDirFileInfoMask: retrieve []os.FileInfo from a defined directory
// "masks" is used to filter files that are kept.
// infosFiles, infosDirs, infosSymlink as return arguments.
func ScanDirFileInfoMask(root string, masks []string) (
	infosFiles,
	infosDirs,
	infosSymlink []os.FileInfo,
	err error) {

	var infosF []os.FileInfo
	var fRoot *os.File

	var isMatching = func(fInfo os.FileInfo) (match bool) {
		var err error
		for _, msk := range masks {
			if match, err = filepath.Match(msk, fInfo.Name()); err == nil {
				if match {
					return match
				}
			}
		}
		return false
	}
	if fRoot, err = os.Open(root); err == nil {
		defer fRoot.Close()
		infosF, err = fRoot.Readdir(-1)

		for _, fi := range infosF {
			if isMatching(fi) {
				switch {
				case fi.IsDir():
					infosDirs = append(infosDirs, fi)
				case fi.Mode()&os.ModeSymlink != 0:
					infosSymlink = append(infosSymlink, fi)
				default:
					infosFiles = append(infosFiles, fi)
				}
			}
		}
	}
	return
}
