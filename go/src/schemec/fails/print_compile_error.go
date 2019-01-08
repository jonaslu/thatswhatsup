package fails

import (
	"fmt"
	"strconv"
	"strings"
)

// SourceMarker contains the tokens matched starting and ending string
// index position in the given program string.
// This is then used to
type SourceMarker struct {
	StartPos int
	EndPos   int
}

type Fail struct {
	program      string
	errorMessage string
	sourceMarker SourceMarker
}

func (e Fail) Error() string {
	errorLineColumnMessage, errorRowMessage, errorRowMarkerMessage := getErrorMessages(e.program, e.sourceMarker)

	errorMessage := fmt.Sprintf("Compile error on row:column %s\n%s\n%s\n\n", errorLineColumnMessage, errorRowMessage, errorRowMarkerMessage)
	errorMessage = errorMessage + fmt.Sprintf("%s\n", e.errorMessage)

	return errorMessage
}

func getErrorRowMarkerMessage(errorColumnIndex int) string {
	return strings.Repeat(" ", errorColumnIndex) + "^"
}

func getErrorRowMessage(program string, startingErrorRowIndex int, endOfErrorRowIndex int) string {
	return program[startingErrorRowIndex:endOfErrorRowIndex]
}

func getErrorLineColumnMessage(errorRowCount int, errorColumnIndex int) string {
	return strconv.Itoa(errorRowCount) + ":" + strconv.Itoa(errorColumnIndex)
}

func getStartingRowErrorContext(program string, sourceMarker SourceMarker) int {
	contextBeforeError := program[:sourceMarker.StartPos]
	startingErrorRowContextIndex := strings.LastIndex(contextBeforeError, "\n") + 1

	return startingErrorRowContextIndex
}

func getEndingRowErrorContext(program string, sourceMarker SourceMarker) int {
	contextAfterError := program[sourceMarker.EndPos:]
	endingErrorRowContextIndex := strings.Index(contextAfterError, "\n")

	if endingErrorRowContextIndex == -1 {
		endingErrorRowContextIndex = len(program)
	} else {
		endingErrorRowContextIndex = endingErrorRowContextIndex + sourceMarker.EndPos
	}
	return endingErrorRowContextIndex
}

func getErrorLineCount(program string, sourceMarker SourceMarker) int {
	programBeforeError := program[:sourceMarker.StartPos]
	errorRowCount := strings.Count(programBeforeError, "\n") + 1

	return errorRowCount
}

func getErrorMessages(program string, sourceMarker SourceMarker) (string, string, string) {
	startingErrorRowIndex := getStartingRowErrorContext(program, sourceMarker)
	endOfErrorRowIndex := getEndingRowErrorContext(program, sourceMarker)
	errorRowCount := getErrorLineCount(program, sourceMarker)

	errorColumnIndex := sourceMarker.StartPos - startingErrorRowIndex

	errorLineColumnMessage := getErrorLineColumnMessage(errorRowCount, errorColumnIndex+1)
	errorRowMessage := getErrorRowMessage(program, startingErrorRowIndex, endOfErrorRowIndex)
	errorRowMarkerMessage := getErrorRowMarkerMessage(errorColumnIndex)

	return errorLineColumnMessage, errorRowMessage, errorRowMarkerMessage
}

// CompileFail returns a struct that implements the Error() interface
// and is used in panic() when the compilation failsz
func CompileFail(program string, errorMessage string, sourceMarker SourceMarker) Fail {
	return Fail{program, errorMessage, sourceMarker}
}
