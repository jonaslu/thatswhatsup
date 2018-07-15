package main

import (
	"errors"
	"schemec/parser"
	"strconv"
)

func compileLetPair(variableSymbol parser.Symbol, letExpression interface{}, environment map[string]StackVariable) []string {
	variableValue := variableSymbol.Value

	_, variablePresent := environment[variableValue]
	if variablePresent {
		logAndQuit(errors.New("Variable already seen"))
	}

	letExpressionInstructions := compileAst(letExpression, environment)

	// !! TODO !! Refactor to generic saveStack function
	saveEaxOnStackInstruction := "movl %eax, " + strconv.Itoa(spIndex) + "(%rsp)"

	environment[variableValue] = StackVariable{variableSymbol, spIndex}
	spIndex = spIndex + stackWordSize

	return append(letExpressionInstructions, saveEaxOnStackInstruction)
}

func compileLetList(letList parser.List, environment map[string]StackVariable) []string {
	letInstructions := []string{}

	for _, letPair := range letList.Value {

		if letExpression, letExpressionOk := letPair.(parser.List); letExpressionOk {
			variableSymbol := letExpression.Value[0]

			if variableName, variableNameOk := variableSymbol.(parser.Symbol); variableNameOk {
				letExpression := letExpression.Value[1]
				letInstructions = append(letInstructions, compileLetPair(variableName, letExpression, environment)...)
			} else {
				logAndQuit(errors.New("Must be variable in let position"))
			}
		} else {
			logAndQuit(errors.New("Let must have lists as arguments"))
		}

	}

	return letInstructions
}

func compileLet(ast parser.List, environment map[string]StackVariable) []string {
	letListInstructions := []string{}

	letListDefinition := ast.Value[1]

	if letList, ok := letListDefinition.(parser.List); ok {
		letListInstructions = append(letListInstructions, compileLetList(letList, environment)...)
	} else {
		logAndQuit(errors.New("Malformed let expression"))
	}

	for _, letBodyExpression := range ast.Value[2:] {
		letListInstructions = append(letListInstructions, compileAst(letBodyExpression, environment)...)
	}

	return letListInstructions
}
