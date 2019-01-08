package main

import (
	"schemec/fails"
	"schemec/parser"
	"strconv"
)

/* Compiles one pair in a let-list (let ( THIS GUY -> (a 1) <- THIS GUY
 * by generating code for whatever the symbol (a in this case)
 * and getting the computed value (by convention stored in %eax)
 * storing it off onto the stack and saving the stack position.
 * Then binding the stack position to the variable name.
 */
func compileLetPair(variableSymbol parser.Symbol, letExpression interface{}, environment map[string]StackVariable) []string {
	variableName := variableSymbol.Value

	_, variablePresent := environment[variableName]
	if variablePresent {
		variableSeenError := fails.CompileFail(inputProgram, "Variable already seen", variableSymbol.SourceMarker)
		panic(variableSeenError)
	}

	letExpressionInstructions := compileAst(letExpression, environment)

	// !! TODO !! Refactor to generic saveStack function
	saveEaxOnStackInstruction := "movl %eax, " + strconv.Itoa(spIndex) + "(%rsp)"

	environment[variableName] = StackVariable{variableSymbol, spIndex}
	spIndex = spIndex + stackWordSize

	return append(letExpressionInstructions, saveEaxOnStackInstruction)
}

/* Runs through each pair in the let assignments and
 * (let (THIS GUY -> (a 1) (b 1) <- THIS GUY
 * Checking that it's a list for every expression
 * in the in-argument list.
 *
 * If so checks that it's a symbol in the first position in the
 * list. If so hands it off to compile the code for the let pair.
 */
func compileLetAssignmentList(letList parser.List, environment map[string]StackVariable) []string {
	letBindingsInstructions := []string{}

	for _, letPair := range letList.Value {

		if letExpression, letExpressionOk := letPair.(parser.List); letExpressionOk {
			variableSymbol := letExpression.Value[0]

			if variableName, variableNameOk := variableSymbol.(parser.Symbol); variableNameOk {
				letExpression := letExpression.Value[1]
				letPairInstructions := compileLetPair(variableName, letExpression, environment)

				letBindingsInstructions = append(letBindingsInstructions, letPairInstructions...)
			} else {
				unknownSourceMarker := variableSymbol.(parser.HasSourceMarker)

				variableInLetPositionError := fails.CompileFail(inputProgram, "Must be variable in let position", unknownSourceMarker.GetSourceMarker())
				panic(variableInLetPositionError)
			}
		} else {
			unknownSourceMarker := letPair.(parser.HasSourceMarker)

			letMustHaveListsError := fails.CompileFail(inputProgram, "Let must have lists as arguments", unknownSourceMarker.GetSourceMarker())
			panic(letMustHaveListsError)
		}
	}

	return letBindingsInstructions
}

/* Compiles the let assignments and then the body of the let
 * (let THIS GUY -> ((a 1)) (b 1)) (+ a b) <- THIS GUY)
 *
 * This is done by pushing the bound values in the let onto the
 * stack connecting the variable name with the position on the stack
 * so when the stack value is needed it is fetched from the stack.
 *
 * Once the let is done the stack is reset and the let bindings "removed"
 */
func compileLet(ast parser.List, environment map[string]StackVariable) []string {
	letExpressionInstructions := []string{}

	letAssignmentsList := ast.Value[1]

	if letList, ok := letAssignmentsList.(parser.List); ok {
		letListAssignmentInstructions := compileLetAssignmentList(letList, environment)
		letExpressionInstructions = append(letExpressionInstructions, letListAssignmentInstructions...)
	} else {
		malformedLetExpressionError := fails.CompileFail(inputProgram, "Malformed let expression", ast.SourceMarker)
		panic(malformedLetExpressionError)
	}

	for _, letBodyExpression := range ast.Value[2:] {
		letExpressionInstructions = append(letExpressionInstructions, compileAst(letBodyExpression, environment)...)
	}

	return letExpressionInstructions
}
