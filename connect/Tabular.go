package main

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type Tabular struct {
	maxLength int
	titles    []string
	values    []Value
}

type Value struct {
	Unit  string
	Value interface{}
}

func (v Value) String() string {
	str := stringer(v.Value)

	return "\033[1m" + str + "\033[0m " + v.Unit
}

func NewTabular() *Tabular {
	return &Tabular{}
}

func (t *Tabular) AddValue(title string, value interface{}) {
	t.AddValueUnit(title, value, "")
}

func (t *Tabular) AddValueUnit(title string, value interface{}, unit string) {
	v := Value{
		Unit:  unit,
		Value: value,
	}

	t.titles = append(t.titles, title)
	t.values = append(t.values, v)

	if len(title) > t.maxLength {
		t.maxLength = len(title)
	}
}

func leftPad(in string, length int) string {
	result := ""
	inLen := utf8.RuneCountInString(in)

	for i := 0; i < length-inLen; i++ {
		result += " "
	}

	return result + in
}

func (t *Tabular) Output(writer io.Writer) {
	for i, value := range t.values {
		fmt.Fprintf(writer, "%s %s\n", leftPad(t.titles[i], t.maxLength), value.String())
	}
}
