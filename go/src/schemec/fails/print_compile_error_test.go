package fails

import (
	"testing"
)

func assertIntValue(t *testing.T, result int, expected int) {
	if result != expected {
		t.Error(result, expected)
	}
}

func assertStringValue(t *testing.T, result string, expected string) {
	if result != expected {
		t.Error(result, expected)
	}
}

func TestStartErrorOnFirstRow(t *testing.T) {
	program := "(+ a 1)"
	startingErrorRowContextIndex := getStartingRowErrorContext(program, SourceMarker{3, 4})

	assertIntValue(t, startingErrorRowContextIndex, 0)
}

func TestStartErrorOnSecondRow(t *testing.T) {
	program := "\n(+ a 1)"
	startingErrorRowContextIndex := getStartingRowErrorContext(program, SourceMarker{4, 5})

	assertIntValue(t, startingErrorRowContextIndex, 1)
}

func TestStartErrorOnThirdRow(t *testing.T) {
	program := "\n\n(+ a 1)"
	startingErrorRowContextIndex := getStartingRowErrorContext(program, SourceMarker{7, 8})

	assertIntValue(t, startingErrorRowContextIndex, 2)
}

func TestStartErrorOnSomeRow(t *testing.T) {
	program := "\n(+ b 2)\n(+ a 1)"
	startingErrorRowContextIndex := getStartingRowErrorContext(program, SourceMarker{11, 12})

	assertIntValue(t, startingErrorRowContextIndex, 9)
}

func TestEndErrorOnFirstRow(t *testing.T) {
	program := "(+ a 1)"
	endingErrorRowContextIndex := getEndingRowErrorContext(program, SourceMarker{3, 4})

	assertIntValue(t, endingErrorRowContextIndex, 7)
}

func TestEndErrorOnLastRow(t *testing.T) {
	program := "\n\n(+ a 1)"
	endingErrorRowContextIndex := getEndingRowErrorContext(program, SourceMarker{7, 8})

	assertIntValue(t, endingErrorRowContextIndex, 9)
}

func TestEndErrorOnSomeRow(t *testing.T) {
	program := "\n(+ b 2)\n(+ a 1)"
	endingErrorRowContextIndex := getEndingRowErrorContext(program, SourceMarker{4, 5})

	assertIntValue(t, endingErrorRowContextIndex, 8)
}

func TestErrorLineCountOnFirstRow(t *testing.T) {
	program := "(+ a 1)"
	errorLineCount := getErrorLineCount(program, SourceMarker{3, 4})

	assertIntValue(t, errorLineCount, 1)
}

func TestErrorLineCountOnThirdRow(t *testing.T) {
	program := "\n\n(+ a 1)"
	errorLineCount := getErrorLineCount(program, SourceMarker{5, 6})

	assertIntValue(t, errorLineCount, 3)
}

func TestErrorLineCountOnSomeRow(t *testing.T) {
	program := "\n(+ b 2)\n(+ a 1)"
	errorLineCount := getErrorLineCount(program, SourceMarker{4, 5})

	assertIntValue(t, errorLineCount, 2)
}

func TestMessagesOnFirstRow(t *testing.T) {
	program := "(+ a 1)"
	errorLineColumnMessage, errorRowMessage, errorRowMakerMessage := getErrorMessages(program, SourceMarker{3, 4})

	assertStringValue(t, errorLineColumnMessage, "1:4")
	assertStringValue(t, errorRowMessage, program)
	assertStringValue(t, errorRowMakerMessage, "   ^")
}

func TestMessagesOnThirdRow(t *testing.T) {
	program := "\n\n(+ a 1)"

	errorLineColumnMessage, errorRowMessage, errorRowMakerMessage := getErrorMessages(program, SourceMarker{5, 6})

	assertStringValue(t, errorLineColumnMessage, "3:4")
	assertStringValue(t, errorRowMessage, "(+ a 1)")
	assertStringValue(t, errorRowMakerMessage, "   ^")
}

func TestErrorMessagesOnSomeRow(t *testing.T) {
	program := "\n(+ b 2)\n(+ a 1)"
	errorLineColumnMessage, errorRowMessage, errorRowMakerMessage := getErrorMessages(program, SourceMarker{5, 6})

	assertStringValue(t, errorLineColumnMessage, "2:5")
	assertStringValue(t, errorRowMessage, "(+ b 2)")
	assertStringValue(t, errorRowMakerMessage, "    ^")
}
