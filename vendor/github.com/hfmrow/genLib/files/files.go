// files.go

package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	glss "github.com/hfmrow/genLib/slices"
)

var TmpCount int

// CompareFilesContent:
func CompareFilesContent(filename1, filename2, filenameOutDiff string) (diff []string) {
	var slice1, slice2 []string
	var err error
	var data1, data2 []byte
	if data1, err = ReadFile(filename1); err == nil {
		slice1 = strings.Split(string(data1), "\n")
		if data2, err = ReadFile(filename2); err == nil {
			slice2 = strings.Split(string(data2), "\n")
			for _, f1 := range slice1 {
				if !glss.IsExistSl(slice2, f1) {
					diff = append(diff, f1)
				}
			}
		}
	}
	WriteFile(filenameOutDiff, []byte(strings.Join(diff, "\n")))
	return diff
}

// writeFile: with file backup capability
func WriteFile(filename string, datas []byte, doBackup ...bool) (err error) {
	if len(doBackup) != 0 {
		if doBackup[0] {
			if _, err = os.Stat(filename); !os.IsNotExist(err) {
				if err = os.Rename(filename, filename+"~"); err != nil {
					return err
				}
			}
		}
	}
	return ioutil.WriteFile(filename, datas, os.ModePerm)
}

// copyFile:
// func CopyFile(inFile, outFile string, doBackup ...bool) (err error) {
// 	var inBytes []byte
// 	if inBytes, err = ioutil.ReadFile(inFile); err == nil {
// 		if len(doBackup) != 0 {
// 			if doBackup[0] {
// 				if err = os.Rename(outFile, outFile+"~"); err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		err = ioutil.WriteFile(outFile, inBytes, os.ModePerm)
// 	}
// 	return err
// }

// ReadFile:
func ReadFile(filename string) (data []byte, err error) {
	return ioutil.ReadFile(filename)
}

// renameProjectFiles: Mass rename function
// func RenameListFiles(fromFileList, toFileList []string) (err error) {
// 	for idx, file := range fromFileList {
// 		if file != toFileList[idx] {
// 			if err = os.Rename(file, toFileList[idx]); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return err
// }

// GetFileBytesString: Retrieve 'from' 'to' bytes from file in string format.
func ReadFileToStrBytes(filename string, from, to int) (outString string) {
	var WriteBytesString = func(p []byte) (data string) {
		const lowerHex = "0123456789abcdef"
		if len(p) == 0 {
			return data
		}
		buf := []byte(`\x00`)
		var b byte
		for _, b = range p {
			buf[2] = lowerHex[b/16]
			buf[3] = lowerHex[b%16]
			data += string(buf)
		}
		return data
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	buff := make([]byte, to-from)
	_, err = file.ReadAt(buff, int64(from))
	if err != nil {
		fmt.Println(err)
	}
	return WriteBytesString(buff)
}

// // FindDir retrieve file in a specific directory with more options.
// func FindDir(dir, mask string, returnedStrSlice *[][]string, scanSub, showDir, followSymlinkDir bool) error {
// 	var fName, time, size string
// 	// Remove unwanted os path separator if exist
// 	//	dir = strings.TrimSuffix(dir, string(os.PathSeparator))

// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		return err
// 	}
// 	for _, file := range files {
// 		fName = filepath.Join(dir, file.Name())
// 		if followSymlinkDir { // Check for symlink ..
// 			file, err = os.Lstat(fName)
// 			if err != nil {
// 				return err
// 			}
// 			if file.Mode()&os.ModeSymlink != 0 { // Is a symlink ?
// 				fName, err := os.Readlink(fName) // Then read it...
// 				if err != nil {
// 					return err
// 				}
// 				file, err = os.Stat(fName) // Get symlink infos.
// 				if err != nil {
// 					return err
// 				}
// 				fName = filepath.Join(dir, file.Name())
// 			}
// 		}
// 		// Recursive play if it's a directory
// 		if file.IsDir() && scanSub {
// 			tmpFileList := new([][]string)
// 			err = FindDir(fName, mask, tmpFileList, scanSub, true, followSymlinkDir)
// 			*returnedStrSlice = append(*returnedStrSlice, *tmpFileList...)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		// get information to be displayed.
// 		size = fmt.Sprintf("%s", humanize.Bytes(uint64(file.Size())))
// 		time = fmt.Sprintf("%s.", humanize.Time(file.ModTime()))
// 		// Check for ext matching well.
// 		ok, err := filepath.Match(mask, file.Name())
// 		if err != nil {
// 			return err
// 		}
// 		if ok {
// 			if showDir { // Limit display directories if requested
// 				*returnedStrSlice = append(*returnedStrSlice, []string{file.Name(), size, time, fName})
// 			} else {
// 				if !file.IsDir() {
// 					*returnedStrSlice = append(*returnedStrSlice, []string{file.Name(), size, time, fName})
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// File struct (SplitFilePath)
type Filepath struct {
	Absolute             string
	Relative             string
	Path                 string
	Base                 string
	BaseNoExt            string
	Ext                  string
	ExecFullName         string
	RealPath             string
	RealName             string
	OutputNewExt         string
	OutputAppendFilename string
	OsSeparator          string
	IsDir                bool
	SymLink              bool
	SymLinkTo            string
}

// Split full filename into path, ext, name, ... optionally add suffix before original extension or change extension
// Relative: SplitFilepath("wanted relative path", fullpath).Relative
// Absolute: SplitFilepath("relative path", fullpath).Absolute
func SplitFilepath(filename string, newExt ...string) Filepath {
	var dir, link bool
	var f = Filepath{}
	var newExtension, dot, addToFilename string
	if len(newExt) != 0 {
		addToFilename = newExt[0]
		if !strings.Contains(newExt[0], ".") {
			dot = "."
		}
		newExtension = dot + newExt[0]
	}
	// IsDir
	fileInfos, err := os.Lstat(filename)
	if err == nil {
		dir = (fileInfos.Mode()&os.ModeDir != 0)
		link = (fileInfos.Mode()&os.ModeSymlink != 0)
		f.IsDir = dir
		if link {
			// IsLink
			f.SymLink = link
			// Symlink endpoint
			f.SymLinkTo, _ = os.Readlink(filename)
			// Symlink and Directory
			ls, err := os.Lstat(f.SymLinkTo)
			if err == nil {
				f.IsDir = (ls.Mode()&os.ModeDir != 0)
			}
		}
	}
	// Absolute
	f.Absolute, _ = filepath.Abs(filename)
	// Relative - Use the optional argument to set as basepath ...
	f.Relative, _ = filepath.Rel(newExtension, filename)
	// OsSep
	f.OsSeparator = string(os.PathSeparator)
	// Path
	if f.Path = filepath.Dir(filename); f.Path == "." {
		f.Path = ""
	}
	// Base
	f.Base = filepath.Base(filename)
	// Ext
	f.Ext = filepath.Ext(filename)
	// BaseNoExt
	splited := strings.Split(f.Base, ".")
	length := len(splited)
	if length == 1 {
		f.BaseNoExt = f.Base

	} else {
		if f.Base[:1] == "." { // Case of hidden file starting with dot
			f.Ext = ""
			f.BaseNoExt = f.Base
		} else {
			splited = splited[:length-1]
			f.BaseNoExt = strings.Join(splited, ".")
		}
	}
	// ExecFullName
	f.ExecFullName, _ = os.Executable()
	// RealPath
	realPathName, _ := filepath.EvalSymlinks(filename)
	if f.RealPath = filepath.Dir(realPathName); f.RealPath == "." {
		f.RealPath = ""
	}
	// RealName
	if f.RealName = filepath.Base(realPathName); f.RealName == "." {
		f.RealName = ""
	}
	// OutNewExt
	if f.Path == "" {
		f.OutputNewExt = f.BaseNoExt + newExtension
	} else {
		f.OutputNewExt = f.Path + f.OsSeparator + f.BaseNoExt + newExtension
	}
	// OutputAppendFilename
	if f.Path == "" {
		f.OutputAppendFilename = f.BaseNoExt + addToFilename + f.Ext
	} else {
		f.OutputAppendFilename = f.Path + f.OsSeparator + f.BaseNoExt + addToFilename + f.Ext
	}
	return f
}
