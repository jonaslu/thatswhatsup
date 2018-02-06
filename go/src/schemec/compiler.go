package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"schemec/parser"
	"strconv"
	"strings"
)

func logAndQuit(err error) {
	fmt.Printf("%+v", err)

	os.Exit(1)
}

func compile(program string) string {
	ast, err := parser.GetAst(program)

	if err != nil {
		logAndQuit(err)
	}

	var writeValue string

	switch n := ast.(type) {
	case parser.Integer:
		integerValue := n.Value
		writeValue = strconv.Itoa(integerValue << 2)

	case parser.Boolean:
		if n.Value {
			writeValue = strconv.Itoa(1<<8 + 31)
		} else {
			writeValue = strconv.Itoa(31)
		}

	case parser.Symbol:

	case parser.List:
	}

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

func Compile(program string) string {
	assembly := compile(program)
	assemblyFilePath := writeAssemblyFile(assembly)
	binaryFilePath := makeRunCodeBinary(assemblyFilePath)

	return binaryFilePath
}
