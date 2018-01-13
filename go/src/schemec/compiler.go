package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func logAndQuit(err error) {
	fmt.Printf("%+v", err)

	os.Exit(1)
}

func assert(testcase string, expected string, actual string) bool {
	if expected != actual {
		fmt.Println("Tescase " + testcase + " failed. Expected " + expected + " actual " + actual)

		return false
	}

	return true
}

func compile(program string) string {
	integerValue, err := strconv.Atoi(program)

	if err != nil {
		logAndQuit(err)
	}

	writeValue := strconv.Itoa(integerValue << 2)

	content, err := ioutil.ReadFile("resources/compile-unit.s")

	if err != nil {
		logAndQuit(err)
	}

	runcode := strings.Replace(string(content), "[insert]", "movl $"+writeValue+", %eax", 1)

	return runcode
}

func writeAssemblyFile(runcode string) string {
	assemblyFilePath := "output/compile-unit.s"
	file, err := os.Create(assemblyFilePath)

	if err != nil {
		logAndQuit(err)
	}

	file.WriteString(runcode)

	return assemblyFilePath
}

func makeRunCodeBinary(assemblyOutputFile string) string {
	binaryName := "output/run-assembly"
	gccBinaryCmd := exec.Command("gcc", "resources/run-assembly.c", assemblyOutputFile, "-o", binaryName)

	_, err := gccBinaryCmd.Output()

	if err != nil {
		logAndQuit(err)
	}

	return binaryName
}

func captureBinaryOutput(binaryFilePath string) string {
	compiledBinaryOutput := exec.Command(binaryFilePath)

	result, err := compiledBinaryOutput.Output()

	if err != nil {
		logAndQuit(err)
	}

	return string(result)
}

func runTest(testname string, program string, output string) bool {
	assembly := compile(program)
	assemblyFilePath := writeAssemblyFile(assembly)
	binaryFilePath := makeRunCodeBinary(assemblyFilePath)
	result := captureBinaryOutput(binaryFilePath)

	return assert(testname, output, result)
}

func runIntegerTests() {
	fmt.Println("Running integer tests")

	testsPassed := true

	testsPassed = testsPassed && runTest("Outout to equal 42", "42", "42")
	testsPassed = testsPassed && runTest("Outout to equal 666", "666", "666")
	testsPassed = testsPassed && runTest("Outout to equal 0", "0", "0")

	if testsPassed {
		fmt.Println("OK")
	} else {
		fmt.Println("Fail")
	}
}

func main() {
	runIntegerTests()
}
