package parser

import (
	"regexp"
)

// integer boolean char identifier or comment
var re = regexp.MustCompile("0|[-]?[1-9][0-9]*|true|false|\\\\#[^\\s]|[a-zA-Z][a-zA-Z0-9]*|\\(|\\)|;[^\n]*")

func tokenize(program string) []string {
	return re.FindAllString(program, -1)
}
