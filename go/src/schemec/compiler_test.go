package main

import (
	"os/exec"
	"testing"
)

func captureBinaryOutput(binaryFilePath string) string {

	compiledBinaryOutput := exec.Command(binaryFilePath)

	result, err := compiledBinaryOutput.Output()

	if err != nil {
		logAndQuit(err)
	}

	return string(result)
}

func runTest(t *testing.T, program string, expected string) {
	binaryFilePath := Compile(program)
	result := captureBinaryOutput(binaryFilePath)

	if result != expected {
		t.Error(result, expected)
	}
}

func TestCompileIntegers(t *testing.T) {
	runTest(t, "42", "42")
	runTest(t, "666", "666")
	runTest(t, "0", "0")
}

func TestCompileBooleans(t *testing.T) {
	runTest(t, "true", "true")
	runTest(t, "false", "false")
}

func TestCompileEmptyList(t *testing.T) {
	runTest(t, "()", "()")
}

func TestCompileChar(t *testing.T) {
	runTest(t, "#\\a", "#\\a")
	runTest(t, "#\\1", "#\\1")
	runTest(t, "#\\€", "#\\€")
}

func TestCompileAddOperator(t *testing.T) {
	runTest(t, "(add1 0)", "1")
	runTest(t, "(add1 1)", "2")
	runTest(t, "(add1 665)", "666")
}

func TestCompileCharToIntegerOperator(t *testing.T) {
	runTest(t, "(char->integer 65)", "#\\A")
	runTest(t, "(char->integer 97)", "#\\a")
}

func TestCompileIntegerToCharOperator(t *testing.T) {
	runTest(t, "(integer->char #\\A)", "65")
	runTest(t, "(integer->char #\\a)", "97")
}

func TestCompileCheckIfNull(t *testing.T) {
	runTest(t, "(null? ())", "true")
	runTest(t, "(null? 1)", "false")
}

func TestCompileCheckIfZero(t *testing.T) {
	runTest(t, "(zero? 0)", "true")
	runTest(t, "(zero? 1)", "false")
}
