package main

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type Table struct {
	columnsMax []int
	header     []string
	rows       [][]string
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) AddHeader(titles ...string) {
	t.header = titles
	t.columnsMax = make([]int, len(t.header))
	for i, title := range t.header {
		t.columnsMax[i] = utf8.RuneCountInString(title)
	}
}

func (t *Table) AddRow(columns ...interface{}) {
	cols := sliceStringer(columns)

	if len(columns) != len(t.header) {
		panic("worng number of columns")
	}

	t.rows = append(t.rows, cols)

	for i, col := range cols {
		l := utf8.RuneCountInString(col)

		if t.columnsMax[i] < l {
			t.columnsMax[i] = l
		}
	}
}

func rightPad(in string, length int) string {
	result := in
	inLen := utf8.RuneCountInString(in)

	for i := 0; i < length-inLen; i++ {
		result += " "
	}

	return result
}

func (t *Table) outputLine(w io.Writer, columns []string) {
	line := ""

	for i, column := range columns {
		line += rightPad(column, t.columnsMax[i]) + " "
	}

	fmt.Fprintf(w, "%s\n", line)
}

func (t *Table) outputHeader(w io.Writer, columns []string) {
	line := ""

	for i, column := range columns {
		line += "\033[1m" + rightPad(column, t.columnsMax[i]) + "\033[0m "
	}

	fmt.Fprintf(w, "%s\n", line)
}

func (t *Table) Output(writer io.Writer) {
	t.outputHeader(writer, t.header)
	for _, row := range t.rows {
		t.outputLine(writer, row)
	}
}
