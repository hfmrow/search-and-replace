package fileText

import (
	"io/ioutil"
	"os"
	"regexp"
)

// replaceInFile: allow using regexp in argument.
func ReplaceInFile(filename, search, replace string, doBackup ...bool) (found bool, err error) {
	var findReg = regexp.MustCompile(`(` + regexp.QuoteMeta(search) + `)`)
	var data []byte
	data, err = ioutil.ReadFile(filename)
	if err == nil {
		if len(doBackup) != 0 {
			if doBackup[0] {
				if err = os.Rename(filename, filename+"~"); err != nil {
					return found, err
				}
			}
		}
		if !findReg.Match(data) {
			return false, err
		}
		found = true
		err = ioutil.WriteFile(filename, findReg.ReplaceAll(data, []byte(replace)), 0644)
	}
	return found, err
}

// Write string to file low lvl format with append possibility.
func WriteTextFile(filename, data string, appendIfExist ...bool) error {
	var apnd bool
	var file *os.File
	var err error
	if len(appendIfExist) != 0 {
		apnd = appendIfExist[0]
	}
	// open file using READ & WRITE permission
	if apnd { // append
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0660)
		defer file.Close()
		if err != nil {
			return err
		}
	} else { // Overwrite
		file, err = os.Create(filename)
		defer file.Close()
		if err != nil {
			return err
		}
	}
	// write some text to file
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	// save changes
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}
