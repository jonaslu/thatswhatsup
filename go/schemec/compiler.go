package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func assert(testcase string, expected string, actual string) bool {
	if expected != actual {
		fmt.Println("Tescase " + testcase + " failed. Expected " + expected + " actual " + actual)

		return false
	}

	return true
}

func compile(program string) string {
	content, err := ioutil.ReadFile("resources/compile-unit.s")

	if err != nil {
		log.Fatal(err)
	}

	runcode := strings.Replace(string(content), "[insert]", "movl $"+program+", %eax", 1)

	return runcode
}

func writeAssemblyFile(runcode string) string {
	assemblyFilePath := "output/compile-unit.s"
	file, err := os.Create(assemblyFilePath)

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(runcode)

	return assemblyFilePath
}

func makeRunCodeBinary(assemblyOutputFile string) string {
	binaryName := "output/run-assembly"
	gccBinaryCmd := exec.Command("gcc", "resources/run-assembly.c", assemblyOutputFile, "-o", binaryName)

	_, err := gccBinaryCmd.Output()

	if err != nil {
		log.Fatal(err)
	}

	return binaryName
}

func captureBinaryOutput(binaryFilePath string) string {
	compiledBinaryOutput := exec.Command(binaryFilePath)

	result, err := compiledBinaryOutput.Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

func runTest(testname string, program string, output string) bool {
	assembly := compile(program)
	assemblyFilePath := writeAssemblyFile(assembly)
	binaryFilePath := makeRunCodeBinary(assemblyFilePath)
	result := captureBinaryOutput(binaryFilePath)

	return assert("Output to equal 42", output, result)
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
