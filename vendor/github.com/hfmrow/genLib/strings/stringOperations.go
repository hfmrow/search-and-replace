// stringOperations.go

package strings

import (
	"crypto/md5"
	"errors"
	"fmt"
	"html"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LowercaseAtFirst: true if 1st char is lowercase
func LowercaseAtFirst(inString string) bool {
	if len(inString) != 0 {
		charType, _ := regexp.Compile("[[:lower:]]")
		return charType.MatchString(inString[:1])
	}
	return true
}

// ToCamel: Turn string into camel case
func ToCamel(inString string, lowerAtFirst ...bool) (outString string) {
	var laf bool
	if len(lowerAtFirst) != 0 {
		laf = lowerAtFirst[0]
	}
	nonAlNum := regexp.MustCompile(`[[:punct:][:space:]]`)
	tmpString := nonAlNum.Split(inString, -1)

	for idx, word := range tmpString {
		if laf && idx < 1 {
			outString += strings.ToLower(word)
		} else {
			outString += strings.Title(word)
		}
	}
	return outString
}

// GenFileName: Generate a randomized file name
func GenFileName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(r.Int63n(time.Now().UnixNano())))))
}

// RemoveNonAlNum: Remove all non alpha-numeric char
func RemoveNonAlNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// RemoveNonNum: Remove all non numeric char
func RemoveNonNum(inString string) string {
	nonAlNum := regexp.MustCompile(`[[:punct:][:alpha:]]`)
	return nonAlNum.ReplaceAllString(inString, "")
}

// ReplaceSpace: replace all [[:space::]] with underscore "_"
func ReplaceSpace(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, "_")
}

// RemoveSpace: remove all [[:space::]]
func RemoveSpace(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:space:]]`)
	return spaceRegex.ReplaceAllString(inString, "")
}

// ReplacePunct: replace all [[:punct::]] with underscore "_"
func ReplacePunct(inString string) string {
	spaceRegex := regexp.MustCompile(`[[:punct:]]`)
	return spaceRegex.ReplaceAllString(inString, "_")
}

// SplitNumeric: Split and keep all numeric values in a string
func SplitNumeric(inString string) (outText []string, err error) {
	toSplit := regexp.MustCompile(`[[:alpha:][:punct:]]`)
	spaceSepared := string(toSplit.ReplaceAll([]byte(inString), []byte(" ")))
	spaceSepared, err = TrimSpace(spaceSepared, "-c")
	if err != nil {
		return outText, err
	}
	outText = strings.Split(spaceSepared, " ")
	return outText, err
}

// ReducePath: Reduce path length by preserving count element from the end
func TruncatePath(fullpath string, count ...int) (reduced string) {
	elemCnt := 2
	if len(count) != 0 {
		elemCnt = count[0]
	}
	splited := strings.Split(fullpath, string(os.PathSeparator))
	if len(splited) > elemCnt+1 {
		return "..." + string(os.PathSeparator) + filepath.Join(splited[len(splited)-elemCnt:]...)
	}
	return fullpath
}

// TrimSpace: Some multiple way to trim strings. cmds is optionnal or accept multiples args
func TrimSpace(inputString string, cmds ...string) (newstring string, err error) {

	osForbiden := regexp.MustCompile(`[<>:"/\\|?*]`)
	remInside := regexp.MustCompile(`[\s\p{Zs}]{2,}`)    //	to match 2 or more whitespace symbols inside a string
	remInsideNoTab := regexp.MustCompile(`[\p{Zs}]{2,}`) //	(preserve \t) to match 2 or more space symbols inside a string

	if len(cmds) != 0 {
		for _, command := range cmds {
			switch command {
			case "+h": //	Escape html
				inputString = html.EscapeString(inputString)
			case "-h": //	UnEscape html
				inputString = html.UnescapeString(inputString)
			case "+e": //	Escape specials chars
				inputString = fmt.Sprintf("%q", inputString)
			case "-e": //	Un-Escape specials chars
				tmpString, err := strconv.Unquote(`"` + inputString + `"`)
				if err != nil {
					return inputString, err
				}
				inputString = tmpString
			case "-w": //	Change all illegals chars (for path in linux and windows) into "-"
				inputString = osForbiden.ReplaceAllString(inputString, "-")
			case "+w": //	clean all illegals chars (for path in linux and windows)
				inputString = osForbiden.ReplaceAllString(inputString, "")
			case "-c": //	Trim [[:space:]] and clean multi [[:space:]] inside
				inputString = strings.TrimSpace(remInside.ReplaceAllString(inputString, " "))
			case "-ct": //	Trim [[:space:]] and clean multi [[:space:]] inside (preserve TAB)
				inputString = strings.Trim(remInsideNoTab.ReplaceAllString(inputString, " "), " ")
			case "-s": //	To match 2 or more whitespace leading/ending/inside a string (include \t, \n)
				inputString = strings.Join(strings.Fields(inputString), " ")
			case "-&": //	Replace ampersand CHAR with ampersand HTML code
				inputString = strings.Replace(inputString, "&", "&amp;", -1)
			case "+&": //	Replace ampersand HTML code with ampersand CHAR
				inputString = strings.Replace(inputString, "&amp;", "&", -1)
			default:
				return inputString, errors.New("TrimSpace, " + command + ", does not exist")
			}
		}
	}
	return inputString, nil
}
