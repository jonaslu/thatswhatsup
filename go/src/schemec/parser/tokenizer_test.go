package parser

import (
	"fmt"
	"reflect"
	"schemec/fails"
	"strings"
	"testing"
)

func testTokenizeHelper(t *testing.T, program string, matches ...Token) {
	result := tokenize(program)

	for num, match := range matches {
		if result[num] != match {
			t.Error(result[num])
		}
	}
}

func TestTokenizeBoolean(t *testing.T) {
	testTokenizeHelper(t, "true", Token{"true", fails.SourceMarker{0, 4}})
	testTokenizeHelper(t, "false", Token{"false", fails.SourceMarker{0, 5}})
	testTokenizeHelper(t, "true false",
		Token{"true", fails.SourceMarker{0, 4}},
		Token{"false", fails.SourceMarker{5, 10}})
}

func TestTokenizeList(t *testing.T) {
	testTokenizeHelper(t, "(", Token{"(", fails.SourceMarker{0, 1}})
	testTokenizeHelper(t, ")", Token{")", fails.SourceMarker{0, 1}})
	testTokenizeHelper(t, "(   )",
		Token{"(", fails.SourceMarker{0, 1}},
		Token{")", fails.SourceMarker{4, 5}})
}

// TODO Test on end of list not found

func TestTokenizeChar(t *testing.T) {
	testTokenizeHelper(t, "\\#m", Token{"\\#m", fails.SourceMarker{0, 3}})
	testTokenizeHelper(t, "\\#A", Token{"\\#A", fails.SourceMarker{0, 3}})
	testTokenizeHelper(t, "\\#0", Token{"\\#0", fails.SourceMarker{0, 3}})
}

func TestTokenizeInteger(t *testing.T) {
	testTokenizeHelper(t, "1234", Token{"1234", fails.SourceMarker{0, 4}})
	testTokenizeHelper(t, "1", Token{"1", fails.SourceMarker{0, 1}})
	testTokenizeHelper(t, "0", Token{"0", fails.SourceMarker{0, 1}})
	testTokenizeHelper(t, "-1", Token{"-1", fails.SourceMarker{0, 2}})
}

func TestTokenizeComment(t *testing.T) {
	result := tokenize(";; Whatever")

	if len(result) > 0 {
		t.Error(result)
	}

	testTokenizeHelper(t, "+ 1 2 ;; whatever\n+ 3 4",
		Token{"+", fails.SourceMarker{0, 1}},
		Token{"1", fails.SourceMarker{2, 3}})
}

func TestTokenizeAssorted(t *testing.T) {
	testTokenizeHelper(t, "false true 0 \\#m",
		Token{"false", fails.SourceMarker{0, 5}},
		Token{"true", fails.SourceMarker{6, 10}},
		Token{"0", fails.SourceMarker{11, 12}},
		Token{"\\#m", fails.SourceMarker{13, 16}})

	testTokenizeHelper(t, "(if true ;;  Whatever\n",
		Token{"(", fails.SourceMarker{0, 1}},
		Token{"if", fails.SourceMarker{1, 3}},
		Token{"true", fails.SourceMarker{4, 8}})
}

func testParseHelper(t *testing.T, tokens []Token, match interface{}) {
	remainingTokens, result := parse(tokens)

	if len(remainingTokens) > 0 {
		t.Error("Tokens remaining", remainingTokens)
	}

	if result == nil {
		if match != nil {
			t.Error(result, match)
		}

		return
	}

	if reflect.TypeOf(result).String() == "parser.List" {
		if !reflect.DeepEqual(result, match) {
			t.Error(result, match)
		}

		return
	}

	if result != match {
		t.Error(result, match)
	}
}

func makeTokenizeReturnValue(values ...string) []Token {
	returnValue := []Token{}

	for _, value := range values {
		returnValue = append(returnValue, Token{value, fails.SourceMarker{0, 0}})
	}

	return returnValue
}

func TestGetAstForIntegers(t *testing.T) {
	singlePositiveDigitToken := makeTokenizeReturnValue("1")
	testParseHelper(t, singlePositiveDigitToken, Integer{1, singlePositiveDigitToken[0].sourceMarker})

	singleNegativeDigitToken := makeTokenizeReturnValue("-1")
	testParseHelper(t, singleNegativeDigitToken, Integer{-1, singleNegativeDigitToken[0].sourceMarker})

	positiveDigitsToken := makeTokenizeReturnValue("1234")
	testParseHelper(t, positiveDigitsToken, Integer{1234, positiveDigitsToken[0].sourceMarker})

	negativeDigitsToken := makeTokenizeReturnValue("-1234")
	testParseHelper(t, negativeDigitsToken, Integer{-1234, negativeDigitsToken[0].sourceMarker})
}

func TestGetAstForBooleans(t *testing.T) {
	trueToken := makeTokenizeReturnValue("true")
	testParseHelper(t, trueToken, Boolean{true, trueToken[0].sourceMarker})

	falseToken := makeTokenizeReturnValue("false")
	testParseHelper(t, falseToken, Boolean{false, falseToken[0].sourceMarker})
}

func TestGetAstForChars(t *testing.T) {
	charToken := makeTokenizeReturnValue("#\\m")
	testParseHelper(t, charToken, Char{'m', charToken[0].sourceMarker})
}

func makeListParseReturnValue(listMembers ...interface{}) List {
	returnValue := []interface{}{}

	for _, value := range listMembers {
		if intVal, ok := value.(int); ok {
			returnValue = append(returnValue, Integer{intVal, fails.SourceMarker{0, 0}})
		}

		if strVal, ok := value.(string); ok {
			returnValue = append(returnValue, Symbol{strVal, fails.SourceMarker{0, 0}})
		}

		if reflect.TypeOf(value).String() == "parser.List" {
			returnValue = append(returnValue, value)
		}
	}

	return List{returnValue, fails.SourceMarker{0, 0}}
}

func TestGetAstForLists(t *testing.T) {
	emptyListTokens := makeTokenizeReturnValue("(", ")")

	emptyList := List{[]interface{}{}, fails.SourceMarker{0, 0}}
	testParseHelper(t, emptyListTokens, emptyList)

	oneTwoThreeList := makeTokenizeReturnValue("(", "1", "2", "3", ")")
	testParseHelper(t, oneTwoThreeList, makeListParseReturnValue(1, 2, 3))

	addOneTwoList := makeTokenizeReturnValue("(", "add", "2", "3", ")")
	testParseHelper(t, addOneTwoList, makeListParseReturnValue("add", 2, 3))

	nestedLists := makeTokenizeReturnValue("(", "(", "+", "2", "3", ")", "(", "3", "a", ")", ")")
	testParseHelper(t, nestedLists,
		makeListParseReturnValue(
			makeListParseReturnValue("+", 2, 3),
			makeListParseReturnValue(3, "a")),
	)
}

func TestGetAstForSymbols(t *testing.T) {
	testParseHelper(t, makeTokenizeReturnValue("a"), Symbol{"a", fails.SourceMarker{0, 0}})
	testParseHelper(t, makeTokenizeReturnValue("Aaa"), Symbol{"Aaa", fails.SourceMarker{0, 0}})
}

func testGetAstFail(t *testing.T, program string, expectedError string) {
	defer func() {
		if r := recover().(error); r != nil {
			fmt.Println(program, r.Error(), expectedError)
			if !strings.Contains(r.Error(), expectedError) {
				t.Error("Failing testcase invalid", r.Error(), expectedError)
			}
		} else {
			t.Error("Program expected to fail", expectedError)
		}
	}()

	GetAst(program)
}

func TestUnknownSymbol(t *testing.T) {
	testGetAstFail(t, "1.0", "Unrecognized symbol")
}

func TestEndOfListNotFound(t *testing.T) {
	testGetAstFail(t, "(   add1 1", "Cannot find end of list")
}
