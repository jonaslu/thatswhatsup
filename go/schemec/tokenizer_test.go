package main

import (
	"testing"
)

func testSplit(t *testing.T, program string, matches ...string) {
	result := split(program)

	for num, match := range matches {
		if result[num] != match {
			t.Error(result[num])
		}
	}
}

func TestSplitOnBoolean(t *testing.T) {
	testSplit(t, "true", "true")
	testSplit(t, "false", "false")
	testSplit(t, "true false", "true", "false")
}

func TestSplitOnList(t *testing.T) {
	testSplit(t, "(", "(")
	testSplit(t, ")", ")")
	testSplit(t, "(   )", "(", ")")
}

func TestSplitOnChar(t *testing.T) {
	testSplit(t, "\\#m", "\\#m")
	testSplit(t, "\\#A", "\\#A")
	testSplit(t, "\\#0", "\\#0")
}

func TestSplitOnInteger(t *testing.T) {
	testSplit(t, "1234", "1234")
	testSplit(t, "1", "1")
	testSplit(t, "0", "0")
	testSplit(t, "-1", "-1")
}

func TestSplitOnComment(t *testing.T) {
	testSplit(t, ";;  Whatever\n", ";;  Whatever")
}

func TestSplitOnAssorted(t *testing.T) {
	testSplit(t, "false true 0 \\#m", "false", "true", "0", "\\#m")
	testSplit(t, "(if true ;;  Whatever\n", "(", "if", "true", ";;  Whatever")
}
