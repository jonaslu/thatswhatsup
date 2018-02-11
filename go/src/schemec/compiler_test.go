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
}
