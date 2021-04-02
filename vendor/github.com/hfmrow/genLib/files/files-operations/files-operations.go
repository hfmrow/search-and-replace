// files-operations.go

/*
	Copyright ©2020 H.F.M - File operations v1.0 https://github.com/hfmrow
	This program comes with absolutely no warranty. See the The MIT License (MIT) for details:
	https://opensource.org/licenses/mit-license.php
*/

package fileOp

/*
#include <stdio.h>
#include <dirent.h>

int ldir(char *root)
{
    DIR *folder;
    struct dirent *entry;
    int files = 0;
// test
    folder = opendir(root);
    if(folder == NULL)
    {
        perror("Unable to read directory");
        return(1);
    }

    while( (entry=readdir(folder)) )
    {
        files++;
        printf("File %3d: %s\n",
                files,
                entry->d_name
              );
    }

    closedir(folder);

    return(0);
}
*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

type FilesOpStruct struct {
	Masks []string

	ExcludeMasks,
	AppendToFiles,
	Backup bool

	Files    []FileDetails
	FilesRaw []string

	SrceDir string

	SortResult,
	RealUserOwner bool

	Owner *OwnerStruct
	Perms *PermsStruct
}

type FileDetails struct {
	BaseDir,
	FileRel,
	LinkTarget,
	FilenameFull,
	Error string

	Mode os.FileMode

	IsDir,
	IsSymlink bool

	Size int64

	info os.FileInfo
}

// restoreRight: NOT used
func (fos *FilesOpStruct) restoreRight(osPathname string) {

	fi, _ := os.Stat(osPathname)

	os.Chmod(osPathname, os.ModePerm&0)

	switch {
	case fi.IsDir():

		os.Chmod(osPathname, fos.Perms.Dir)

	default:

		os.Chmod(osPathname, fos.Perms.File)
	}

}

// FilesOpStructNew: Create a new 'FilesOpStruct' that hold some
// useful file/dir functions. The goal of it is to facilitate file
// interactions, can create copy, files/dir perserving
// owner/permissions. Files can be stored with all information
// needed to recreate them. This structure has been desi
func FilesOpStructNew() (fos *FilesOpStruct, err error) {

	fos = new(FilesOpStruct)
	return fos, fos.init()
}

// init: Initialize permissions and current/real user informations.
func (fos *FilesOpStruct) init() (err error) {
	fos.Perms = PermsStructNew()
	if fos.Owner, err = OwnerStructNew(); err != nil {
		fmt.Println("Error collecting real user information.")
	}
	return
}

// SetOwner: Set filename user owner, uid and gid are defined in
// main structure, fos.RealUserOwner must be enabled. The
// 'callbeforefunc' permit to run an operation just before chown
// command (like saving png, jpg ...).
func (fos *FilesOpStruct) SetOwner(path string, callbeforefunc ...interface{}) (err error) {

	if fos.RealUserOwner {
		if len(callbeforefunc) > 0 {
			callbeforefunc[0].(func())()
		}
		err = os.Chown(path, fos.Owner.Uid, fos.Owner.Gid)
	}
	return
}

func (fos *FilesOpStruct) checkForBackup(filename string) (err error) {

	if fos.Backup {
		if _, err = os.Stat(filename); err == nil {
			err = os.Rename(filename, filename+"~")
		}
	}
	return
}

// WriteFile: can set owner (if requested).
func (fos *FilesOpStruct) WriteFile(filename string, data []byte, perm os.FileMode) (err error) {

	if err = fos.checkForBackup(filename); err == nil {

		if err = ioutil.WriteFile(filename, data, perm); err == nil {
			err = fos.SetOwner(filename)
		}
	}
	return
}

// Copy: file/dir tree preserving permissions & owner (if requested).
func (fos *FilesOpStruct) CopyFile(src, dest string) (err error) {

	var info os.FileInfo
	var f, s *os.File

	if info, err = os.Lstat(src); err == nil {

		if err = fos.MkdirAll(filepath.Dir(dest), fos.Perms.Dir); err == nil {

			if err = fos.checkForBackup(dest); err == nil {

				if f, err = os.Create(dest); err == nil {
					defer f.Close()

					if err = os.Chmod(f.Name(), info.Mode()); err == nil {

						if s, err = os.Open(src); err == nil {
							defer s.Close()

							_, err = io.Copy(f, s)
							fos.SetOwner(dest)
						}
					}
				}
			}
		}
	}
	return err
}

// FileDetailsNew: Make a new 'FileDetails' Structure
func (fos *FilesOpStruct) FileDetailsNew(path string, stat os.FileInfo) (fd *FileDetails, err error) {

	rel, err := filepath.Rel(fos.SrceDir, path)
	if err != nil {
		return nil, err
	}

	if rel == "." {
		return nil, nil
	}

	var linkTarget string
	isLink := stat.Mode()&os.ModeSymlink == os.ModeSymlink
	if isLink {
		if linkTarget, err = os.Readlink(path); err != nil {
			return nil, err
		}
	}

	fd = new(FileDetails)
	fd = &FileDetails{
		FilenameFull: path,
		BaseDir:      fos.SrceDir,
		FileRel:      rel,
		Mode:         stat.Mode(),
		Size:         stat.Size(),
		IsDir:        stat.IsDir(),
		IsSymlink:    isLink,
		LinkTarget:   linkTarget,
		info:         stat}

	return fd, nil
}

// func asPtrAndLength(s string) (*C.char, int) {
// 	addr := &s
// 	hdr := (*reflect.StringHeader)(unsafe.Pointer(addr))

// 	p := (*C.char)(unsafe.Pointer(hdr.Data))
// 	n := hdr.Len

// 	// reflect.StringHeader stores the Data field as a uintptr, not
// 	// a pointer, so ensure that the string remains reachable until
// 	// the uintptr is converted.
// 	runtime.KeepAlive(addr)

// 	return p, n
// }

// readDir:
func (fos *FilesOpStruct) readDir(path string) (filesDetails []FileDetails, err error) {

	var (
		fd        *FileDetails
		tmpFD     []FileDetails
		filesInfo []os.FileInfo
	)

	fmt.Println("ldir:", C.ldir(C.CString(path)))

	return

	if filesInfo, err = ioutil.ReadDir(path); err == nil {
		for _, fi := range filesInfo {

			if fd, err = fos.FileDetailsNew(filepath.Join(path, fi.Name()), fi); err != nil {
				fd.Error = fmt.Sprintf("%v", err)
			}
			if fd == nil {
				continue
			}
			filesDetails = append(filesDetails, *fd)

			if fi.IsDir() {

				if tmpFD, err = fos.readDir(fd.FilenameFull); err != nil {
					log.Printf("%v", err)
				}
				filesDetails = append(filesDetails, tmpFD...)
			}
		}
	}
	return
}

// GetFilesDetailsReadDir: Retrieve files details starting at "root",
// they will be stored as 'FileDetails' structure in fos.Files slice.
// unlike the other function, this one does not take into account
// the options defined, Mask and sort methods are not available.
func (fos *FilesOpStruct) GetFilesDetailsReadDir(root string) (err error) {

	fos.SrceDir = root

	fos.Files, err = fos.readDir(root)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	return
}

// GetFilesDetails: Retrieve files details starting at "root", they
// will be stored as 'FileDetails' structure in fos.Files slice.
func (fos *FilesOpStruct) GetFilesDetails(root string) (err error) {
	var ok bool
	var stat os.FileInfo

	// Clear files entries
	if !fos.AppendToFiles {
		fos.Files = fos.Files[:0]
	}

	var checkFile = func(inFileDetails *FileDetails) (err error) {

		if inFileDetails != nil {

			// Check for pattern exclusion.
			if len(fos.Masks) != 0 {
				for _, mask := range fos.Masks {
					if ok, err = filepath.Match(
						mask,
						filepath.Base(inFileDetails.FilenameFull)); err != nil {

						return
					}

					switch {
					case ok && !fos.ExcludeMasks:
						fos.Files = append(fos.Files, *inFileDetails)
					case !ok && fos.ExcludeMasks:
						fos.Files = append(fos.Files, *inFileDetails)
					}
				}
			} else {
				fos.Files = append(fos.Files, *inFileDetails)
			}
		}
		return
	}

	fos.SrceDir = root

	if stat, err = os.Lstat(fos.SrceDir); err == nil {
		switch {
		case stat.IsDir():

			if err = filepath.Walk(fos.SrceDir, func(path string, fInfo os.FileInfo, err error) error {

				switch {
				case err != nil:
					return err
				case fInfo == nil:
					return fmt.Errorf("Error stat: %v", path)
				}

				f, err := fos.FileDetailsNew(path, fInfo)
				if err != nil {
					return err
				}
				return checkFile(f)

			}); err != nil {
				return err
			}

		default:

			fos.SrceDir = filepath.Dir(root)
			f, err := fos.FileDetailsNew(root, stat)
			if err != nil {
				return err
			}
			err = checkFile(f)
			if err != nil {
				return err
			}
		}
	}
	return
}

// MkdirAll: This is a modified version of original golang function,
// this one set owner for each created directories.
func (fos *FilesOpStruct) MkdirAll(path string, perm os.FileMode) error {

	dir, err := os.Stat(path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &os.PathError{"mkdir", path, syscall.ENOTDIR}
	}

	i := len(path) // Skip trailing path separator.
	for i > 0 && os.IsPathSeparator(path[i-1]) {
		i--
	}

	j := i // Scan backward over element.
	for j > 0 && !os.IsPathSeparator(path[j-1]) {
		j--
	}

	if j > 1 {
		err = fos.MkdirAll(path[:j-1], perm)
		if err != nil {
			return err
		}
	}

	if err = os.Mkdir(path, perm); err != nil {

		dir, err1 := os.Lstat(path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}

	// Change dir owner if required
	if err = fos.SetOwner(path); err != nil {
		return err
	}

	return nil
}

// Sort: Descending as default behavior
func (fos *FilesOpStruct) Sort(ascend ...bool) {

	var ascending bool
	if len(ascend) > 0 {
		ascending = ascend[0]
	}

	switch ascending {

	case false:
		sort.SliceStable(fos.Files,
			func(i, j int) bool {
				return strings.ToLower(fos.Files[i].FilenameFull) <
					strings.ToLower(fos.Files[j].FilenameFull)
			})

	case true:
		sort.SliceStable(fos.Files,
			func(i, j int) bool {
				return strings.ToLower(fos.Files[i].FilenameFull) >
					strings.ToLower(fos.Files[j].FilenameFull)
			})
	}
}

// Read FilesOpStruct from file
func (fos *FilesOpStruct) Read(filename string) (err error) {

	var textFileBytes []byte

	if textFileBytes, err = ioutil.ReadFile(filename); err == nil {
		err = json.Unmarshal(textFileBytes, &fos)
	}
	if err != nil {
		fmt.Printf("Enable to read 'FilesOpStruct' file: %s\n", err.Error())
	}
	return
}

// Write FilesOpStruct to file
func (fos *FilesOpStruct) Write(filename string) (err error) {

	var (
		jsonData []byte
		out      bytes.Buffer
	)

	if jsonData, err = json.Marshal(&fos); err == nil {
		if err = json.Indent(&out, jsonData, "", "\t"); err == nil {
			err = ioutil.WriteFile(filename, out.Bytes(), os.ModePerm)
		}
	}
	return err
}
