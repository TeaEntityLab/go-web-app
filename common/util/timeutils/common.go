package timeutils

import (
	"strconv"
	"strings"
	"time"
)

const (
	TimeFormatLayoutDateYYYYMMDDByDash      = "2006-01-02"
	TimeFormatLayoutDateYYYYMMDDBySlash     = "2006/01/02"
	TimeFormatLayoutDateYYYYMMDDByDot       = "2006.01.02"
	TimeFormatLayoutDateYYYYMMDDByBackslash = "2006\\01\\02"

	TimeFormatLayoutDateYYYYMMDDHHMMSSCSVExcel = "2006-01-02 15:04:05"
)

var (
	TimeFormatLayoutDateYYYYMMDDSeries = []string{
		time.RFC3339,
		TimeFormatLayoutDateYYYYMMDDHHMMSSCSVExcel,
		TimeFormatLayoutDateYYYYMMDDByDash,
		TimeFormatLayoutDateYYYYMMDDBySlash,
		TimeFormatLayoutDateYYYYMMDDByDot,
		TimeFormatLayoutDateYYYYMMDDByBackslash,
	}
)

func GetTimestamp() int64 {
	return time.Now().UTC().Unix()
}

func GetTimestampString() string {
	return strconv.FormatInt(GetTimestamp(), 10)
}

func GetRFC3339StringForCodeGen(givenTime *time.Time) string {
	datetime := GetRFC3339String(givenTime)
	datetime = strings.Replace(datetime, ":", "_", -1)
	datetime = strings.Replace(datetime, "-", "_", -1)

	return datetime
}

func GetRFC3339String(givenTime *time.Time) string {
	if givenTime == nil {
		val := time.Now().UTC()
		givenTime = &val
	}
	datetime := givenTime.Format(time.RFC3339)
	return datetime
}

func GetDateYYYYMMDDByDashString(givenTime *time.Time) string {
	if givenTime == nil {
		val := time.Now().UTC()
		givenTime = &val
	}
	datetime := givenTime.Format(TimeFormatLayoutDateYYYYMMDDByDash)
	return datetime
}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func TruncateByDuration(t time.Time, d time.Duration) time.Time {
	return t.UTC().Truncate(d)
}

func TruncateByHour(t time.Time) time.Time {
	return TruncateByDuration(t, time.Hour)
}

func TruncateByDay(t time.Time) time.Time {
	return TruncateByDuration(t, 24*time.Hour)
}
