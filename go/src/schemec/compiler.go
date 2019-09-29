package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"schemec/compilerutils"
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

var spIndex = -8

const stackWordSize = -8

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

func getBooleanImmediateRepresentation(booleanValue bool) string {
	if booleanValue {
		return strconv.Itoa(1<<booleanShiftBits + booleanTag)
	}

	return strconv.Itoa(booleanTag)
}

func getIntegerImmediateRepresentation(integerValue int) string {
	return strconv.Itoa(integerValue << 2)
}

/* Returns the immediate representation (value combined with the
 * type tag) as a string. Often used to store the value
 * in rax to perform some calculation.
 */
func getImmediateValue(ast interface{}) string {
	switch n := ast.(type) {

	case parser.Integer:
		integerValue := n.Value
		return getIntegerImmediateRepresentation(integerValue)

	case parser.Boolean:
		return getBooleanImmediateRepresentation(n.Value)

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

func storeImmediateRepresentationInRax(immediateRepresentation string) string {
	return "movq $" + immediateRepresentation + ", %rax"
}

func compareIfRaxContainsValue(compareValue string) []string {
	compareWithEmptyList := "cmpq $" + compareValue + ", %rax"
	zeroRax := "movq $0, %rax"
	setLowBitOfRaxToOneIfEqual := "sete %al"

	shiftUpby7Bits := "salq $7, %rax"
	setBooleanTag := "orq $31, %rax"

	return []string{compareWithEmptyList, zeroRax, setLowBitOfRaxToOneIfEqual, shiftUpby7Bits, setBooleanTag}
}

/* If the list is empty, simply returns the immediate representation
 * for a list. If not empty checks the function call (in the first
 * position of the list and switches on known built in forms.
 */
func parseList(list parser.List, environment map[string]StackVariable) []string {
	listValues := list.Value

	if len(listValues) == 0 {
		return []string{storeImmediateRepresentationInRax(strconv.Itoa(emptyListTag))}
	}

	firstInList := listValues[0]

	if n, ok := firstInList.(parser.Symbol); ok {
		switch n.Value {
		case "add1":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			addOneImmediateValue := getIntegerImmediateRepresentation(1)
			addOneToValueInRax := "addq $" + addOneImmediateValue + ", %rax"

			instructions := append(firstArgumentInstructions, addOneToValueInRax)
			return instructions

		case "char->integer":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			shiftUpBy6Bits := "salq $6, %rax"
			setCharTag := "addq $" + strconv.Itoa(charactersTag) + ", %rax"

			charToIntegerInstructions := append(firstArgumentInstructions, shiftUpBy6Bits, setCharTag)
			return charToIntegerInstructions

		case "integer->char":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			shiftByDown6Bits := "sarq $6, %rax"

			integerToCharInstructions := append(firstArgumentInstructions, shiftByDown6Bits)
			return integerToCharInstructions

		case "null?":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			checkIfNullInstructions := append(firstArgumentInstructions, compareIfRaxContainsValue(strconv.Itoa(emptyListTag))...)
			return checkIfNullInstructions

		case "zero?":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			zeroImmedateValue := getIntegerImmediateRepresentation(0)

			checkIfZeroInstructions := append(firstArgumentInstructions, compareIfRaxContainsValue(zeroImmedateValue)...)
			return checkIfZeroInstructions

		case "not":
			firstArgumentInstructions := compileAst(listValues[1], environment)

			shiftByDown7Bytes := "sarq $7, %rax"
			xorWithOne := "xorq $1, %rax"

			shiftUpby7Bits := "salq $7, %rax"
			setBooleanTag := "orq $31, %rax"

			notInstructions := append(firstArgumentInstructions, shiftByDown7Bytes, xorWithOne, shiftUpby7Bits, setBooleanTag)
			return notInstructions

		case "+":
			// spIndex = spIndex - stackWordSize
			secondArgumentInstructions := compileAst(listValues[2], environment)
			saveRaxOnStackInstruction := "movq %rax, " + strconv.Itoa(spIndex) + "(%rsp)"

			firstArgumentInstructions := compileAst(listValues[1], environment)
			addStackInstructionsToRax := "addq " + strconv.Itoa(spIndex) + "(%rsp), %rax"
			spIndex = spIndex + stackWordSize

			instructions := []string{}
			instructions = append(instructions, secondArgumentInstructions...)
			instructions = append(instructions, saveRaxOnStackInstruction)
			instructions = append(instructions, firstArgumentInstructions...)
			instructions = append(instructions, addStackInstructionsToRax)

			return instructions

		// case "-":
		// case "*":
		// case "=":
		// case "char=?":
		// (let ((a 1) (b 2)) body-expr body-expr)
		case "let":
			return compileLet(list, environment)

			// (if test-value truebranch falsebranch)
		case "if":
			endLabel := compilerutils.GetUniqueLabel()
			falseLabel := compilerutils.GetUniqueLabel()

			instructionsForTestValue := compileAst(listValues[1], environment)

			compareWithFalseValue := "cmpq $" + getBooleanImmediateRepresentation(false) + ", %rax"
			jumpToFalseBranchIfEqual := "je " + falseLabel

			trueBranchInstructions := compileAst(listValues[2], environment)
			jmpToEndLabelInstruction := "jmp " + endLabel

			emitFalseLabelInstruction := falseLabel + ":"
			falseBranchInstructions := compileAst(listValues[3], environment)

			endOfIfLabelInstruction := endLabel + ":"

			instructions := append(instructionsForTestValue, compareWithFalseValue, jumpToFalseBranchIfEqual)
			instructions = append(instructions, trueBranchInstructions...)
			instructions = append(instructions, jmpToEndLabelInstruction, emitFalseLabelInstruction)
			instructions = append(instructions, falseBranchInstructions...)
			instructions = append(instructions, endOfIfLabelInstruction)
			return instructions

		case "cons":
			instructionsForCarValue := compileAst(listValues[1], environment)
			saveRaxOnStackInstruction := "movq %rax, " + strconv.Itoa(spIndex) + "(%rsp)"
			spIndex = spIndex + stackWordSize

			instructionsForCdrValue := compileAst(listValues[2], environment)
			spIndex := spIndex - stackWordSize

			storeCarValue := []string{"movq " + strconv.Itoa(spIndex) + "(%rsp), %rbx", "movq %rbx, 0(%rsi)"}
			storeCdrValue := "movq %rax, 8(%rsi)"

			makeRaxAPairTag := []string{"movq %rsi, %rax", "orq $1, %rax"}
			bumpEsi := "addq $16, %rsi"

			instructions := instructionsForCarValue
			instructions = append(instructions, saveRaxOnStackInstruction)
			instructions = append(instructions, instructionsForCdrValue...)
			instructions = append(instructions, storeCarValue...)
			instructions = append(instructions, storeCdrValue)
			instructions = append(instructions, makeRaxAPairTag...)
			instructions = append(instructions, bumpEsi)

			return instructions

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

		fetchStackPosToRax := "movq " + strconv.Itoa(stackVariable.spIndex) + "(%rsp), %rax"
		return []string{fetchStackPosToRax}

	default:
		immediateRepresentation := getImmediateValue(ast)
		writeValue := storeImmediateRepresentationInRax(immediateRepresentation)

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

	stdout, err := gccBinaryCmd.Output()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code: %d error message: %s", exitError.ExitCode(), stdout)
			fmt.Println("Compilation failed, binary can be found in: ", assemblyOutputFile)
			panic(string(exitError.Stderr))
		}
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

func main() {
}
