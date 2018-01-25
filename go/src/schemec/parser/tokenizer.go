package parser

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// integer boolean char comment or identifier
var re = regexp.MustCompile("0|[-]?[1-9][0-9]*|true|false|\\\\#[^\\s]|\\(|\\)|;[^\n]*|[^\\s\\(\\);]+")

func tokenize(program string) []string {
	return re.FindAllString(program, -1)
}

type (
	// Symbol symbol ast value
	Symbol struct {
		value string
	}
)

func peekToken(tokens []string) (string, []string) {
	if len(tokens) == 0 {
		return "", tokens
	}

	if len(tokens) == 1 {
		return tokens[0], []string{}
	}

	return tokens[0], tokens[1:]
}

func parseList(tokens []string) ([]string, []interface{}, error) {
	returnList := []interface{}{}

	for len(tokens) > 0 {
		nextToken, remainingTokens := peekToken(tokens)

		if nextToken == ")" {
			return remainingTokens, returnList, nil
		}

		remainingTokens, parseResult, err := parse(tokens)

		if err != nil {
			log.Fatal(err)
		}

		returnList = append(returnList, parseResult)
		tokens = remainingTokens
	}

	return nil, nil, errors.New("End of list not found")
}

func parse(tokens []string) ([]string, interface{}, error) {
	nextToken, remainingTokens := peekToken(tokens)

	if nextToken == "(" {
		return parseList(remainingTokens)
	}

	if strings.HasPrefix(nextToken, "\\#") {
		// A rune - represents a char
		return remainingTokens, []rune(nextToken)[2:3][0], nil
	}

	if nextToken == "true" {
		return remainingTokens, true, nil
	}

	if nextToken == "false" {
		return remainingTokens, false, nil
	}

	intValue, err := strconv.Atoi(nextToken)
	if err == nil {
		return remainingTokens, intValue, nil
	}

	return remainingTokens, Symbol{nextToken}, nil
}

func filter(strings []string, shouldInclude func(string) bool) []string {
	result := []string{}
	for _, str := range strings {
		if shouldInclude(str) {
			result = append(result, str)
		}
	}
	return result
}

// GetAst takes a string and produces an ast
func GetAst(program string) (interface{}, error) {
	if program == "" {
		return nil, nil
	}

	tokens := tokenize(program)
	tokens = filter(tokens, func(v string) bool { return !strings.HasPrefix(v, ";") })

	if len(tokens) == 0 {
		return nil, nil
	}

	tokens, ast, err := parse(tokens)

	if err != nil {
		return nil, err
	}

	if len(tokens) > 0 {
		return nil, errors.New("Unrecognized symbol")
	}

	return ast, nil
}
