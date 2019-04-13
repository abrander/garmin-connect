package main

import (
	"fmt"
	"strconv"
	"time"
)

func formatDate(t time.Time) string {
	if t == (time.Time{}) {
		return "never"
	}

	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func stringer(value interface{}) string {
	stringer, ok := value.(fmt.Stringer)
	if ok {
		return stringer.String()
	}

	str := ""
	switch value.(type) {
	case string:
		str = value.(string)
	case int, int64:
		str = fmt.Sprintf("%d", value)
	case float64:
		str = strconv.FormatFloat(value.(float64), 'f', 1, 64)
	case bool:
		if value.(bool) {
			str = gotIt
		}
	default:
		panic(fmt.Sprintf("no idea what to do about %T:%v", value, value))
	}

	return str
}

func sliceStringer(values []interface{}) []string {
	ret := make([]string, len(values), len(values))

	for i, value := range values {
		ret[i] = stringer(value)
	}

	return ret
}
