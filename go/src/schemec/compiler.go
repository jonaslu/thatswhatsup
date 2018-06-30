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
const booleanShiftBits = 7

const charactersTag = 15
const charactersShiftBits = 8

// 00101111
const emptyListTag = 47

var spIndex = 0

const stackWordSize = 4

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

	case parser.List:
		if len(n.Value) == 0 {
			return strconv.Itoa(emptyListTag), nil
		}

		return "", errors.New("Not an immediate value")

	default:
		// !! TODO !! Fix pretty error-printing
		return "", errors.New("Not an immediate value")
	}
}

func storeImmediateRepresentationInEax(immediateRepresentation string) string {
	return "movl $" + immediateRepresentation + ", %eax"
}

func compareIfEaxContainsValue(compareValue string) []string {
	compareWithEmptyList := "cmpl $" + compareValue + ", %eax"
	zeroEax := "movl $0, %eax"
	setLowBitOfEaxToOneIfEqual := "sete %al"

	shiftUpby7Bits := "sall $7, %eax"
	setBooleanTag := "orl $31, %eax"

	return []string{compareWithEmptyList, zeroEax, setLowBitOfEaxToOneIfEqual, shiftUpby7Bits, setBooleanTag}
}

func parseList(list parser.List) ([]string, error) {
	listValues := list.Value

	if len(listValues) == 0 {
		return []string{storeImmediateRepresentationInEax(strconv.Itoa(emptyListTag))}, nil
	}

	if n, ok := listValues[0].(parser.Symbol); ok {
		switch n.Value {
		case "add1":
			firstArgumentInstructions := compileAst(listValues[1])

			addOneImmediateValue := getIntegerImmediateRepresentation(1)
			addOneToValueInEax := "addl $" + addOneImmediateValue + ", %eax"

			instructions := append(firstArgumentInstructions, addOneToValueInEax)
			return instructions, nil

		case "char->integer":
			firstArgumentInstructions := compileAst(listValues[1])

			shiftUpBy6Bits := "sall $6, %eax"
			setCharTag := "addl $" + strconv.Itoa(charactersTag) + ", %eax"

			charToIntegerInstructions := append(firstArgumentInstructions, shiftUpBy6Bits, setCharTag)
			return charToIntegerInstructions, nil

		case "integer->char":
			firstArgumentInstructions := compileAst(listValues[1])

			shiftByDown6Bits := "sarl $6, %eax"

			charToIntegerInstructions := append(firstArgumentInstructions, shiftByDown6Bits)
			return charToIntegerInstructions, nil

		case "null?":
			firstArgumentInstructions := compileAst(listValues[1])

			checkIfNullInstructions := append(firstArgumentInstructions, compareIfEaxContainsValue(strconv.Itoa(emptyListTag))...)
			return checkIfNullInstructions, nil

		case "zero?":
			firstArgumentInstructions := compileAst(listValues[1])

			zeroImmedateValue := getIntegerImmediateRepresentation(0)

			checkIfZeroInstructions := append(firstArgumentInstructions, compareIfEaxContainsValue(zeroImmedateValue)...)
			return checkIfZeroInstructions, nil

		case "not":
			firstArgumentInstructions := compileAst(listValues[1])

			shiftByDown7Bytes := "sarl $7, %eax"
			xorWithOne := "xorl $1, %eax"

			shiftUpby7Bits := "sall $7, %eax"
			setBooleanTag := "orl $31, %eax"

			notInstructions := append(firstArgumentInstructions, shiftByDown7Bytes, xorWithOne, shiftUpby7Bits, setBooleanTag)
			return notInstructions, nil

		case "+":
			spIndex = spIndex - stackWordSize
			secondArgumentInstructions := compileAst(listValues[2])
			saveEaxOnStackInstruction := "movl %eax, " + strconv.Itoa(spIndex) + "(%rsp)"

			firstArgumentInstructions := compileAst(listValues[1])
			addStackInstructionsToEax := "addl " + strconv.Itoa(spIndex) + "(%rsp), %eax"
			spIndex = spIndex + stackWordSize

			instructions := []string{}
			instructions = append(instructions, secondArgumentInstructions...)
			instructions = append(instructions, saveEaxOnStackInstruction)
			instructions = append(instructions, firstArgumentInstructions...)
			instructions = append(instructions, addStackInstructionsToEax)

			return instructions, nil
		}
	}

	// TODO Fix pretty error-printing
	return nil, errors.New("Expected symbol at position")
}

func compileAst(ast interface{}) []string {
	switch n := ast.(type) {
	case parser.List:
		instructions, err := parseList(n)

		if err != nil {
			logAndQuit(err)
		}

		return instructions

	default:
		// !! TODO !! Handle error
		immediateRepresentation, _ := getImmediateValue(ast)
		writeValue := storeImmediateRepresentationInEax(immediateRepresentation)

		return []string{writeValue}
	}
}

func compile(program string) string {
	ast, err := parser.GetAst(program)

	if err != nil {
		logAndQuit(err)
	}

	instructions := compileAst(ast)
	writeValue := strings.Join(instructions, "\n")

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
