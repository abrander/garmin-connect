package main

import (
	"fmt"
)

// nzf is a type that will print "-" instead of 0.0 when used as a stringer.
type nzf float64

func (nzf nzf) String() string {
	if nzf != 0.0 {
		return fmt.Sprintf("%.01f", nzf)
	}

	return "-"
}
