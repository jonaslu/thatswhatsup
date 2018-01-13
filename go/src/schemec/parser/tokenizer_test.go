package parser

import (
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
