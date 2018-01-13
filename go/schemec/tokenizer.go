package main

import (
	"regexp"
)

// Identifier, boolean (true or false), #\m <- char or integer
var re = regexp.MustCompile("0|[-]?[1-9][0-9]*|true|false|\\\\#[^\\s]|[a-zA-Z][a-zA-Z0-9]*|\\(|\\)|;[^\n]*")

func split(program string) []string {
	return re.FindAllString(program, -1)
}
