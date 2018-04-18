package main

import (
	"errors"
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

// 0011111
const booleanTag = 31
const booleanShiftBits = 8

const charactersTag = 15
const charactersShiftBits = 8

// 00101111
const emptyListTag = 47

func getIntegerImmediateRepresentation(integerValue int) string {
	return strconv.Itoa(integerValue << 2)
}

func getImmediateValue(ast interface{}) (string, error) {
	switch n := ast.(type) {

	case parser.Integer:
		integerValue := n.Value
		return getIntegerImmediateRepresentation(integerValue), nil

	case parser.Boolean:
		if n.Value {
			return strconv.Itoa(1<<booleanShiftBits + booleanTag), nil
		}

		return strconv.Itoa(booleanTag), nil

	case parser.Char:
		return strconv.Itoa(int(n.Value)<<charactersShiftBits + charactersTag), nil

	default:
		// !! TODO !! Fix pretty error-printing
		return "", errors.New("Not an immediate value")
	}
}

func storeImmediateRepresentationInEax(immediateRepresentation string) string {
	return "movl $" + immediateRepresentation + ", %eax"
}

func parseList(list parser.List) ([]string, error) {
	listValues := list.Value

	if len(listValues) == 0 {
		return []string{storeImmediateRepresentationInEax(strconv.Itoa(emptyListTag))}, nil
	}

	if n, ok := listValues[0].(parser.Symbol); ok {
		switch n.Value {
		case "add1":
			// !! TODO !! Handle error
			immediateValue, _ := getImmediateValue(listValues[1])
			immediateValueStoredInEax := storeImmediateRepresentationInEax(immediateValue)

			addOneImmediateValue := getIntegerImmediateRepresentation(1)
			addOneToValueInEax := "addl $" + addOneImmediateValue + ", %eax"

			addOneToValueInstructions := []string{immediateValueStoredInEax, addOneToValueInEax}
			return addOneToValueInstructions, nil

		case "char->integer":
			immediateValue, _ := getImmediateValue(listValues[1])
			immediateValueStoredInEax := storeImmediateRepresentationInEax(immediateValue)

			shiftBy6Bits := "sall $6, %eax"
			setCharTag := "addl $" + strconv.Itoa(charactersTag) + ", %eax"

			charToIntegerInstructions := []string{immediateValueStoredInEax, shiftBy6Bits, setCharTag}
			return charToIntegerInstructions, nil
		}
	}

	// TODO Fix pretty error-printing
	return nil, errors.New("Expected symbol at position")
}

func compile(program string) string {
	ast, err := parser.GetAst(program)

	if err != nil {
		logAndQuit(err)
	}

	var writeValue string

	switch n := ast.(type) {
	case parser.List:
		var instructions []string
		instructions, err = parseList(n)
		writeValue = strings.Join(instructions, "\n")

	default:
		// !! TODO !! Handle error
		immediateRepresentation, _ := getImmediateValue(ast)
		writeValue = storeImmediateRepresentationInEax(immediateRepresentation)
	}

	content, err := ioutil.ReadFile("resources/compile-unit.s")

	if err != nil {
		logAndQuit(err)
	}

	fmt.Println(writeValue)

	runcode := strings.Replace(string(content), "[insert]", writeValue, 1)

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

// Compile takes a program as a string and returns the filepath
// of a runnable binary with the program compiled
func Compile(program string) string {
	assembly := compile(program)
	assemblyFilePath := writeAssemblyFile(assembly)
	binaryFilePath := makeRunCodeBinary(assemblyFilePath)

	return binaryFilePath
}
