package main

import (
	"os/exec"
	"strings"
	"testing"
)

func runTest(t *testing.T, program string) (string, error) {
	binaryFilePath := Compile(program)

	compiledBinaryOutput := exec.Command(binaryFilePath)
	result, err := compiledBinaryOutput.Output()

	return string(result), err
}

func runSuccess(t *testing.T, program string, expected string) {
	result, err := runTest(t, program)

	if err != nil {
		t.Error("Program failed", err)
	}

	if result != expected {
		t.Error("Test failed", result, expected)
	}
}

func runFail(t *testing.T, program string, expectedError string) {
	defer func() {
		if r := recover().(error); r != nil {
			if !strings.Contains(r.Error(), expectedError) {
				t.Error("Failing testcase invalid", r.Error(), expectedError)
			}
		} else {
			t.Error("Program expected to fail", expectedError)
		}
	}()

	runTest(t, program)
}

func TestCompileIntegers(t *testing.T) {
	runSuccess(t, "42", "42")
	runSuccess(t, "666", "666")
	runSuccess(t, "0", "0")
}

func TestCompileBooleans(t *testing.T) {
	runSuccess(t, "true", "true")
	runSuccess(t, "false", "false")
}

func TestCompileEmptyList(t *testing.T) {
	runSuccess(t, "()", "()")
}

func TestCompileChar(t *testing.T) {
	runSuccess(t, "#\\a", "#\\a")
	runSuccess(t, "#\\1", "#\\1")
	runSuccess(t, "#\\€", "#\\€")
}

func TestCompileAddOperator(t *testing.T) {
	runSuccess(t, "(add1 0)", "1")
	runSuccess(t, "(add1 1)", "2")
	runSuccess(t, "(add1 665)", "666")
}

func TestCompileCharToIntegerOperator(t *testing.T) {
	runSuccess(t, "(char->integer 65)", "#\\A")
	runSuccess(t, "(char->integer 97)", "#\\a")
}

func TestCompileIntegerToCharOperator(t *testing.T) {
	runSuccess(t, "(integer->char #\\A)", "65")
	runSuccess(t, "(integer->char #\\a)", "97")
}

func TestCompileCheckIfNull(t *testing.T) {
	runSuccess(t, "(null? ())", "true")
	runSuccess(t, "(null? 1)", "false")
}

func TestCompileCheckIfZero(t *testing.T) {
	runSuccess(t, "(zero? 0)", "true")
	runSuccess(t, "(zero? 1)", "false")
}

func TestCompileNot(t *testing.T) {
	runSuccess(t, "(not false)", "true")
	runSuccess(t, "(not true)", "false")
}

func TestCompileAddNumbers(t *testing.T) {
	runSuccess(t, "(+ 1 1)", "2")
}

func TestCompileAddAdd1Numbers(t *testing.T) {
	runSuccess(t, "(+ 1 (add1 1))", "3")
}

func TestCompileAddRecursiveNumbers(t *testing.T) {
	runSuccess(t, "(+ 1 (+ 1 1))", "3")
}

func TestExpectedSymbolAtFirstInListPosition(t *testing.T) {
	runFail(t, "(1)", "Expected a function first in list")
}

func TestUnknownFunctionFirstInListPosition(t *testing.T) {
	runFail(t, "(yggdrasil)", "Unknown function")
}

func TestUnknownVariable(t *testing.T) {
	runFail(t, "yggdrasil", "Unknown variable")
}
