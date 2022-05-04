package main

import (
	"testing"
)

var testStr = []struct {
	name   string
	input  string
	output string
}{
	{name: "1", input: "a4bc2d5e", output: "aaaabccddddde"},
	{name: "2", input: "abcd", output: "abcd"},
	{name: "3", input: "", output: ""},
	{name: "4", input: "a5A", output: "aaaaaA"},
	{name: "5", input: "a5", output: "aaaaa"},
	{name: "6", input: "qwe\\4\\5", output: "qwe45"},
	{name: "7", input: "qwe\\45", output: "qwe44444"},
	{name: "8", input: "qwe\\\\5", output: "qwe\\\\\\\\\\"},
	{name: "9", input: "45", output: ""},
}

func TestUnpack(t *testing.T) {
	for _, v := range testStr {
		t.Run(v.name, func(t *testing.T) {
			if v.output != Unpack(v.input) {
				t.Error("Error", v.input)
			}
		})
	}
}
