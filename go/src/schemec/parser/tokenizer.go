package parser

import (
	"errors"
	"regexp"
	"schemec/fails"
	"strconv"
	"strings"
)

var inputProgram string

type (
	// Token represents a token in the given program string
	Token struct {
		value        string
		sourceMarker fails.SourceMarker
	}
)

// integer boolean char comment or identifier
var re = regexp.MustCompile("0|[-]?[1-9][0-9]*|true|false|#\\\\[^\\s]|\\(|\\)|;[^\n]*|[^\\s\\(\\);]+")

type (
	// Symbol represents the ast type of anything that is not
	// an integer, char or boolean
	Symbol struct {
		Value        string
		SourceMarker fails.SourceMarker
	}

	// Integer ast type
	Integer struct {
		Value        int
		SourceMarker fails.SourceMarker
	}

	// Boolean ast type
	Boolean struct {
		Value        bool
		SourceMarker fails.SourceMarker
	}

	// Char (single character) ast type
	Char struct {
		Value        rune
		SourceMarker fails.SourceMarker
	}

	// List ast type
	List struct {
		Value        []interface{}
		SourceMarker fails.SourceMarker
	}
)

type HasSourceMarker interface {
	GetSourceMarker() fails.SourceMarker
}

func (s Symbol) GetSourceMarker() fails.SourceMarker {
	return s.SourceMarker
}

func (i Integer) GetSourceMarker() fails.SourceMarker {
	return i.SourceMarker
}

func (b Boolean) GetSourceMarker() fails.SourceMarker {
	return b.SourceMarker
}

func (c Char) GetSourceMarker() fails.SourceMarker {
	return c.SourceMarker
}

func (l List) GetSourceMarker() fails.SourceMarker {
	return l.SourceMarker
}

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
	listStartPos := nextToken.sourceMarker.StartPos

	for len(tokens) > 0 {
		nextToken, remainingTokens := peekToken(tokens)

		if nextToken.value == ")" {
			return remainingTokens, List{returnList, fails.SourceMarker{listStartPos, nextToken.sourceMarker.EndPos}}, nil
		}

		remainingTokens, parseResult := parse(tokens)

		returnList = append(returnList, parseResult)
		tokens = remainingTokens
	}

	return nil, List{}, errors.New("Cannot find end of list")
}

func parse(tokens []Token) ([]Token, interface{}) {
	nextToken, remainingTokens := peekToken(tokens)
	nextTokenValue := nextToken.value
	sourceMarker := nextToken.sourceMarker

	if nextTokenValue == "(" {
		if remainingTokens, parsedList, err := parseList(remainingTokens); err != nil {
			panic(fails.CompileFail(inputProgram, err.Error(), nextToken.sourceMarker))
		} else {
			return remainingTokens, parsedList
		}
	}

	if strings.HasPrefix(nextTokenValue, "#\\") {
		return remainingTokens,
			Char{[]rune(nextTokenValue)[2:3][0], sourceMarker}
	}

	if nextTokenValue == "true" {
		return remainingTokens, Boolean{true, sourceMarker}
	}

	if nextTokenValue == "false" {
		return remainingTokens, Boolean{false, sourceMarker}
	}

	intValue, err := strconv.Atoi(nextTokenValue)
	if err == nil {
		return remainingTokens, Integer{intValue, sourceMarker}
	}

	return remainingTokens, Symbol{nextTokenValue, sourceMarker}
}

func isComment(value string) bool {
	return strings.HasPrefix(value, ";")
}

func tokenize(program string) []Token {
	allIndexes := re.FindAllStringIndex(program, -1)
	allTokens := []Token{}

	for _, pos := range allIndexes {
		match := program[pos[0]:pos[1]]

		if !isComment(match) {
			allTokens = append(allTokens, Token{match, fails.SourceMarker{pos[0], pos[1]}})
		}
	}

	return allTokens
}

// GetAst takes a string and produces an ast
func GetAst(program string) interface{} {
	inputProgram = program

	if program == "" {
		return nil
	}

	tokens := tokenize(program)

	if len(tokens) == 0 {
		return nil
	}

	tokens, ast := parse(tokens)

	if len(tokens) > 0 {
		panic(fails.CompileFail(inputProgram, "Unrecognized symbol", tokens[0].sourceMarker))
	}

	return ast
}
