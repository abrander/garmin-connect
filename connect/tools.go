package main

import (
	"fmt"
	"time"
)

func formatDate(t time.Time) string {
	if t == (time.Time{}) {
		return "never"
	}

	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
