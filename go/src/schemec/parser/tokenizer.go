package parser

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type (
	// SourceMarker contains the tokens matched starting and ending string
	// index position in the given program string
	SourceMarker struct {
		startPos int
		endPos   int
	}

	// Token represents a token in the given program string
	Token struct {
		value        string
		sourceMarker SourceMarker
	}
)

// integer boolean char comment or identifier
var re = regexp.MustCompile("0|[-]?[1-9][0-9]*|true|false|\\\\#[^\\s]|\\(|\\)|;[^\n]*|[^\\s\\(\\);]+")

func isComment(value string) bool {
	return strings.HasPrefix(value, ";")
}

func tokenize(program string) []Token {
	allIndexes := re.FindAllStringIndex(program, -1)
	allTokens := []Token{}

	for _, pos := range allIndexes {
		match := program[pos[0]:pos[1]]

		if !isComment(match) {
			allTokens = append(allTokens, Token{match, SourceMarker{pos[0], pos[1]}})
		}
	}

	return allTokens
}

type (
	// Symbol represents the ast type of anything that is not
	// an integer, char or boolean
	Symbol struct {
		value        string
		sourceMarker SourceMarker
	}

	// Integer ast type
	Integer struct {
		value        int
		sourceMarker SourceMarker
	}

	// Boolean ast type
	Boolean struct {
		value        bool
		sourceMarker SourceMarker
	}

	// Char (single character) ast type
	Char struct {
		value        rune
		sourceMarker SourceMarker
	}

	// List ast type
	List struct {
		value        []interface{}
		sourceMarker SourceMarker
	}
)

func peekToken(tokens []Token) (Token, []Token) {
	if len(tokens) == 0 {
		return Token{}, tokens
	}

	if len(tokens) == 1 {
		return tokens[0], []Token{}
	}

	return tokens[0], tokens[1:]
}

func parseList(tokens []Token) ([]Token, List, error) {
	returnList := []interface{}{}
	nextToken, _ := peekToken(tokens)
	listStartPos := nextToken.sourceMarker.startPos

	for len(tokens) > 0 {
		nextToken, remainingTokens := peekToken(tokens)

		if nextToken.value == ")" {
			return remainingTokens,
				List{returnList, SourceMarker{listStartPos, nextToken.sourceMarker.endPos}},
				nil
		}

		remainingTokens, parseResult, err := parse(tokens)

		if err != nil {
			log.Fatal(err)
		}

		returnList = append(returnList, parseResult)
		tokens = remainingTokens
	}

	// TODO Print nice error-message
	return nil, List{}, errors.New("End of list not found")
}

func parse(tokens []Token) ([]Token, interface{}, error) {
	nextToken, remainingTokens := peekToken(tokens)
	nextTokenValue := nextToken.value
	sourceMarker := nextToken.sourceMarker

	if nextTokenValue == "(" {
		return parseList(remainingTokens)
	}

	if strings.HasPrefix(nextTokenValue, "\\#") {
		return remainingTokens,
			Char{[]rune(nextTokenValue)[2:3][0], sourceMarker},
			nil
	}

	if nextTokenValue == "true" {
		return remainingTokens, Boolean{true, sourceMarker}, nil
	}

	if nextTokenValue == "false" {
		return remainingTokens, Boolean{false, sourceMarker}, nil
	}

	intValue, err := strconv.Atoi(nextTokenValue)
	if err == nil {
		return remainingTokens, Integer{intValue, sourceMarker}, nil
	}

	return remainingTokens, Symbol{nextTokenValue, sourceMarker}, nil
}

// GetAst takes a string and produces an ast
func GetAst(program string) (interface{}, error) {
	if program == "" {
		return nil, nil
	}

	tokens := tokenize(program)

	if len(tokens) == 0 {
		return nil, nil
	}

	tokens, ast, err := parse(tokens)

	if err != nil {
		return nil, err
	}

	if len(tokens) > 0 {
		// TODO Print error context (use SourceMarker)
		return nil, errors.New("Unrecognized symbol")
	}

	return ast, nil
}
