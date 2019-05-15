package main

import (
	"fmt"
	"strconv"
	"time"
)

func formatDate(t time.Time) string {
	if t == (time.Time{}) {
		return "-"
	}

	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func stringer(value interface{}) string {
	stringer, ok := value.(fmt.Stringer)
	if ok {
		return stringer.String()
	}

	str := ""
	switch v := value.(type) {
	case string:
		str = v
	case int, int64:
		str = fmt.Sprintf("%d", v)
	case float64:
		str = strconv.FormatFloat(v, 'f', 1, 64)
	case bool:
		if v {
			str = gotIt
		}
	default:
		panic(fmt.Sprintf("no idea what to do about %T:%v", value, value))
	}

	return str
}

func sliceStringer(values []interface{}) []string {
	ret := make([]string, len(values))

	for i, value := range values {
		ret[i] = stringer(value)
	}

	return ret
}

func hoursAndMinutes(dur time.Duration) string {
	if dur == 0 {
		return "-"
	}

	if dur < 60*time.Minute {
		m := dur.Truncate(time.Minute)

		return fmt.Sprintf("%dm", m/time.Minute)
	}

	h := dur.Truncate(time.Hour)
	m := (dur - h).Truncate(time.Minute)

	h /= time.Hour
	m /= time.Minute

	return fmt.Sprintf("%dh%dm", h, m)
}
