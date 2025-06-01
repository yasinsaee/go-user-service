package util

import (
	"fmt"
	"time"

	ptime "github.com/yaa110/go-persian-calendar"
)

func ConvertTime(ti time.Time, lang string) string {
	timeParse := ti

	t := time.Date(timeParse.Year(), timeParse.Month(), timeParse.Day(), timeParse.Hour(), timeParse.Minute(), timeParse.Second(), 0, ptime.Iran())
	pt := ptime.New(t)

	if lang == "fa" {
		return fmt.Sprintln(pt.Day(), pt.Month().String(), pt.Year())
	}

	return fmt.Sprintln(ti.Day(), ti.Month().String(), ti.Year())
}

func GetDay(ti time.Time, lang string) string {
	timeParse := ti

	t := time.Date(timeParse.Year(), timeParse.Month(), timeParse.Day(), timeParse.Hour(), timeParse.Minute(), timeParse.Second(), 0, ptime.Iran())
	pt := ptime.New(t)

	if lang == "fa" {
		return fmt.Sprintln(pt.Day())
	}

	return fmt.Sprintln(ti.Day())
}

func GetMonth(ti time.Time, lang string) string {
	timeParse := ti

	t := time.Date(timeParse.Year(), timeParse.Month(), timeParse.Day(), timeParse.Hour(), timeParse.Minute(), timeParse.Second(), 0, ptime.Iran())
	pt := ptime.New(t)

	if lang == "fa" {
		return fmt.Sprintln(pt.Month().String())
	}

	return fmt.Sprintln(ti.Month().String(),)
}
