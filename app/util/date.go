package util

import "time"

const dateTimeFormat = "2006-01-02 15:04:05"
const dateFormat = "2006-01-02"
const timeFormat = "15:04:05"

func NewDateHelper() *dateHelper {
	return &dateHelper{}
}

type dateHelper struct{}

func (n *dateHelper) Format(c time.Time, f string) string {
	return c.Format(f)
}

func (n *dateHelper) FormatDateTime(c time.Time) string {
	return c.Format(dateTimeFormat)
}

func (n *dateHelper) FormatDate(c time.Time) string {
	return c.Format(dateFormat)
}

func (n *dateHelper) FormatTime(c time.Time) string {
	return c.Format(timeFormat)
}

func (n *dateHelper) StrToDateTime(s string) time.Time {
	dateTime, _ := time.Parse(dateTimeFormat, s)
	return dateTime
}
