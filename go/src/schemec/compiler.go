package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"schemec/fails"
	"schemec/parser"
	"strconv"
	"strings"
)

// 0011111
const booleanTag = 31
const booleanShiftBits = 7

const charactersTag = 15
const charactersShiftBits = 8

// 00101111
const emptyListTag = 47

var spIndex = -4

const stackWordSize = -4

var inputProgram string
var topEnvironment map[string]StackVariable

/* StackVariable
 * is used in the environment to store
 * a value in the stack (for use later in the
 * function call).
 */
type StackVariable struct {
	symbol  parser.Symbol
	spIndex int
}

func getIntegerImmediateRepresentation(integerValue int) string {
	return strconv.Itoa(integerValue << 2)
}

/* Returns the immediate representation (value combined with the
 * type tag) as a string. Often used to store the value
 * in eax to perform some calculation.
 */
func getImmediateValue(ast interface{}) string {
	switch n := ast.(type) {

	case parser.Integer:
		integerValue := n.Value
		return getIntegerImmediateRepresentation(integerValue)

	case parser.Boolean:
		if n.Value {
			return strconv.Itoa(1<<booleanShiftBits + booleanTag)
		}

		return strconv.Itoa(booleanTag)

	case parser.Char:
		return strconv.Itoa(int(n.Value)<<charactersShiftBits + charactersTag)

	case parser.List:
		if len(n.Value) == 0 {
			return strconv.Itoa(emptyListTag)
		}

		// !! TODO !! Is this needed? How do we get here?
		panic(fails.CompileFail(inputProgram, "Not an immediate value", n.SourceMarker))
	}

	// !! TODO !! If we get here it's an error in the compiler
	// since it's an added parser.Type missed in the case
	fmt.Printf("Compiler error, unknown type %v, %T", ast, ast)
	os.Exit(1)

	return ""
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

/* If the list is empty, simply returns the immediate representation
 * for a list. If not empty checks the function call (in the first
 * position of the list and switches on known built in forms.
 */
func parseList(list parser.List, environment map[string]StackVariable) []string {
	listValues := list.Value

	if len(listValues) == 0 {
		return []string{storeImmediateRepresentationInEax(strconv.Itoa(emptyListTag))}
	}

	firstInList := listValues[0]

	if n, ok := firstInList.(parser.Symbol); ok {
		switch n.Value {
		case "add1":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			addOneImmediateValue := getIntegerImmediateRepresentation(1)
			addOneToValueInEax := "addl $" + addOneImmediateValue + ", %eax"

			instructions := append(firstArgumentInstructions, addOneToValueInEax)
			return instructions

		case "char->integer":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			shiftUpBy6Bits := "sall $6, %eax"
			setCharTag := "addl $" + strconv.Itoa(charactersTag) + ", %eax"

			charToIntegerInstructions := append(firstArgumentInstructions, shiftUpBy6Bits, setCharTag)
			return charToIntegerInstructions

		case "integer->char":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			shiftByDown6Bits := "sarl $6, %eax"

			integerToCharInstructions := append(firstArgumentInstructions, shiftByDown6Bits)
			return integerToCharInstructions

		case "null?":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			checkIfNullInstructions := append(firstArgumentInstructions, compareIfEaxContainsValue(strconv.Itoa(emptyListTag))...)
			return checkIfNullInstructions

		case "zero?":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			zeroImmedateValue := getIntegerImmediateRepresentation(0)

			checkIfZeroInstructions := append(firstArgumentInstructions, compareIfEaxContainsValue(zeroImmedateValue)...)
			return checkIfZeroInstructions

		case "not":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			shiftByDown7Bytes := "sarl $7, %eax"
			xorWithOne := "xorl $1, %eax"

			shiftUpby7Bits := "sall $7, %eax"
			setBooleanTag := "orl $31, %eax"

			notInstructions := append(firstArgumentInstructions, shiftByDown7Bytes, xorWithOne, shiftUpby7Bits, setBooleanTag)
			return notInstructions

		case "+":
			// spIndex = spIndex - stackWordSize
			secondArgumentInstructions := compileAst(listValues[2], environment)
			saveEaxOnStackInstruction := "movl %eax, " + strconv.Itoa(spIndex) + "(%rsp)"

			firstArgumentInstructions := compileAst(listValues[1], environment)
			addStackInstructionsToEax := "addl " + strconv.Itoa(spIndex) + "(%rsp), %eax"
			spIndex = spIndex + stackWordSize

			instructions := []string{}
			instructions = append(instructions, secondArgumentInstructions...)
			instructions = append(instructions, saveEaxOnStackInstruction)
			instructions = append(instructions, firstArgumentInstructions...)
			instructions = append(instructions, addStackInstructionsToEax)

			return instructions

		// case "-":
		// case "*":
		// case "=":
		// case "char=?":
		// (let ((a 1) (b 2)) body-expr body-expr)
		case "let":
			return compileLet(list, environment)

		default:
			panic(fails.CompileFail(inputProgram, "Unknown function", n.SourceMarker))
		}
	}

	panic(fails.CompileFail(inputProgram, "Expected a function first in list", firstInList.(parser.HasSourceMarker).GetSourceMarker()))
}

/* The main meat in this file - takes the parsed ast and
 * and any environment and switches on the ast.
 * If it's a list we parse it as a list (with function calls).
 * If not we try to find a symbol in the environment matching the ast
 * and if not it's an immediate value (Char, Boolean, Integer)
 * and we return the value.
 */
func compileAst(ast interface{}, environment map[string]StackVariable) []string {
	switch n := ast.(type) {
	case parser.List:
		instructions := parseList(n, environment)

		return instructions

	case parser.Symbol:
		stackVariable, variableFound := environment[n.Value]

		if !variableFound {
			panic(fails.CompileFail(inputProgram, "Unknown variable", n.SourceMarker))
		}

		fetchStackPosToEax := "movl " + strconv.Itoa(stackVariable.spIndex) + "(%rsp), %eax"
		return []string{fetchStackPosToEax}

	default:
		immediateRepresentation := getImmediateValue(ast)
		writeValue := storeImmediateRepresentationInEax(immediateRepresentation)

		return []string{writeValue}
	}
}

func compile(program string) string {
	ast := parser.GetAst(program)

	topEnvironment = map[string]StackVariable{}

	instructions := compileAst(ast, topEnvironment)
	writeValue := strings.Join(instructions, "\n")

	content, err := ioutil.ReadFile("resources/compile-unit.s")

	if err != nil {
		panic(err)
	}

	runcode := strings.Replace(string(content), "[insert]", writeValue, 1)

	return runcode
}

func writeAssemblyFile(runcode string) string {
	assemblyFilePath := "output/compile-unit.s"
	file, err := os.Create(assemblyFilePath)

	if err != nil {
		panic(err)
	}

	file.WriteString(runcode)

	return assemblyFilePath
}

func makeRunCodeBinary(assemblyOutputFile string) string {
	binaryName := "output/run-assembly"
	gccBinaryCmd := exec.Command("gcc", "resources/run-assembly.c", assemblyOutputFile, "-o", binaryName)

	_, err := gccBinaryCmd.Output()

	if err != nil {
		panic(err)
	}

	return binaryName
}

// Compile takes a program as a string and returns the filepath
// of a runnable binary with the program compiled
func Compile(program string) string {
	inputProgram = program

	assembly := compile(program)
	assemblyFilePath := writeAssemblyFile(assembly)
	binaryFilePath := makeRunCodeBinary(assemblyFilePath)

	return binaryFilePath
}
