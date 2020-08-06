package tools

import (
	"regexp"
	"strings"
	"time"
)

// UrlGet: find into sentence the url available part.
func UrlsGet(inString string) []string {
	reg := regexp.MustCompile(`(http|https|ftp|ftps)\:\/\/[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,3}(\/\S*)?`)
	return reg.FindAllString(inString, -1)
}

type timeStamp struct {
	Year          string
	YearCopyRight string
	Month         string
	MonthWord     string
	Day           string
	DayWord       string
	Date          string
	Time          string
	Full          string
	FullNum       string
}

// Get current timestamp
func TimeStamp() *timeStamp {
	ts := new(timeStamp)
	timed := time.Now()
	regD := regexp.MustCompile("([^[:digit:]])")
	regA := regexp.MustCompile("([^[:alpha:]])")
	splitedNum := regD.Split(timed.Format(time.RFC3339), -1)
	splitedWrd := regA.Split(timed.Format(time.RFC850), -1)
	ts.Year = splitedNum[0]
	ts.Month = splitedNum[1]
	ts.Day = splitedNum[2]
	ts.Time = splitedNum[3] + `:` + splitedNum[4] + `:` + splitedNum[5]
	ts.DayWord = splitedWrd[0]
	ts.MonthWord = splitedWrd[5]
	ts.YearCopyRight = `©` + ts.Year
	ts.Full = strings.Join(strings.Split(timed.Format(time.RFC1123), " ")[:5], " ")

	nonAlNum := regexp.MustCompile(`[[:punct:][:alpha:]]`)
	date := nonAlNum.ReplaceAllString(time.Now().Format(time.RFC3339), "")[:14]
	ts.FullNum = date[:8] + "-" + date[8:]
	return ts
}
