package parser

import (
	"reflect"
	"testing"
)

func testTokenize(t *testing.T, program string, matches ...string) {
	result := tokenize(program)

	for num, match := range matches {
		if result[num] != match {
			t.Error(result[num])
		}
	}
}

func TestTokenizeOnBoolean(t *testing.T) {
	testTokenize(t, "true", "true")
	testTokenize(t, "false", "false")
	testTokenize(t, "true false", "true", "false")
}

func TestTokenizeOnList(t *testing.T) {
	testTokenize(t, "(", "(")
	testTokenize(t, ")", ")")
	testTokenize(t, "(   )", "(", ")")
}

func TestTokenizeOnChar(t *testing.T) {
	testTokenize(t, "\\#m", "\\#m")
	testTokenize(t, "\\#A", "\\#A")
	testTokenize(t, "\\#0", "\\#0")
}

func TestTokenizeOnInteger(t *testing.T) {
	testTokenize(t, "1234", "1234")
	testTokenize(t, "1", "1")
	testTokenize(t, "0", "0")
	testTokenize(t, "-1", "-1")
}

func TestTokenizeOnComment(t *testing.T) {
	testTokenize(t, ";;  Whatever\n", ";;  Whatever")
}

func TestTokenizeOnAssorted(t *testing.T) {
	testTokenize(t, "false true 0 \\#m", "false", "true", "0", "\\#m")
	testTokenize(t, "(if true ;;  Whatever\n", "(", "if", "true", ";;  Whatever")
}

func testGetAst(t *testing.T, program string, match interface{}) {
	result, error := GetAst(program)

	if error != nil {
		if error.Error() != match {
			t.Error(error)
		}

		return
	}

	if result == nil {
		if match != nil {
			t.Error(result, match)
		}

		return
	}

	if reflect.TypeOf(result).String() == "[]interface {}" {
		if !reflect.DeepEqual(result, match) {
			t.Error(result, match)
		}

		return
	}

	if result != match {
		t.Error(result, match)
	}
}

func TestGetAstForIntegers(t *testing.T) {
	testGetAst(t, "1", 1)
	testGetAst(t, "1234", 1234)
	testGetAst(t, "-1234", -1234)
}

func TestGetAstForBooleans(t *testing.T) {
	testGetAst(t, "true", true)
	testGetAst(t, "false", false)
}

func TestGetAstForChars(t *testing.T) {
	testGetAst(t, "\\#m", rune('m'))
}

func TestGetAstForLists(t *testing.T) {
	testGetAst(t, "()", []interface{}{})
	testGetAst(t, "(1 2 3)", []interface{}{1, 2, 3})
	testGetAst(t, "(add 2 3)", []interface{}{Symbol{"add"}, 2, 3})
	testGetAst(t, "(+ 2 3)", []interface{}{Symbol{"+"}, 2, 3})
	testGetAst(t, "((+ 2 3) (+ 3 a))", []interface{}{[]interface{}{Symbol{"+"}, 2, 3}, []interface{}{Symbol{"+"}, 3, Symbol{"a"}}})
}

func TestGetAstForComments(t *testing.T) {
	testGetAst(t, ";; all comments baby", nil)
	testGetAst(t, "((+ 2 3) (+ 3 a) ;; yahrrr\n)", []interface{}{[]interface{}{Symbol{"+"}, 2, 3}, []interface{}{Symbol{"+"}, 3, Symbol{"a"}}})
}

func TestGetAstForSymbols(t *testing.T) {
	testGetAst(t, "a", Symbol{"a"})
	testGetAst(t, "aab", Symbol{"aab"})
	testGetAst(t, "AAb", Symbol{"AAb"})
}

func TestNotAllowedSymbols(t *testing.T) {
	testGetAst(t, "1aa", "Unrecognized symbol")
}

func TestGetAstMixedCases(t *testing.T) {
	testGetAst(t, "", nil)
	testGetAst(t, "   (+    2 3) ;; comment", []interface{}{Symbol{"+"}, 2, 3})
	testGetAst(t, "   ((+    2 3) ;; comment", "End of list not found")
}
