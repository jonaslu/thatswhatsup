package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func assert(testcase string, expected string, actual string) {
	if expected != actual {
		fmt.Println("Tescase " + testcase + " failed. Expected " + expected + " actual " + actual)
	}
}

func compile() string {
	content, err := ioutil.ReadFile("resources/compile-unit.s")

	if err != nil {
		log.Fatal(err)
	}

	runcode := strings.Replace(string(content), "[insert]", "movl $42, %eax", 1)

	return runcode
}

func writeRunCode(runcode string) {
	file, err := os.Create("output/compile-unit.s")

	if err != nil {
		log.Fatal(err)
	}

	file.WriteString(runcode)
}

func makeRunCodeBinary() {
	gccBinaryCmd := exec.Command("gcc", "resources/run-assembly.c", "output/compile-unit.s", "-o", "output/run-assembly")

	_, err := gccBinaryCmd.Output()

	if err != nil {
		log.Fatal(err)
	}
}

func runTests() {
	compiledBinaryOutput := exec.Command("output/run-assembly")

	result, err := compiledBinaryOutput.Output()

	if err != nil {
		log.Fatal(err)
	}

	assert("Output to equal 42", "42", string(result))
}

func main() {
	runcode := compile()
	writeRunCode(runcode)
	makeRunCodeBinary()
	runTests()
}
