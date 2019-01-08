package main

import (
	"testing"
)

func TestLetPrimitiveAssignments(t *testing.T) {
	runSuccess(t, "(let ((a 1) (b 4)) (+ a b))", "5")
}

func TestLetComplexAssignment(t *testing.T) {
	runSuccess(t, "(let ((a (+ 1 2)) (b 4)) (+ a b))", "7")
}

func TestLetSimpleVariableReference(t *testing.T) {
	runSuccess(t, "(let ((a 1)) a)", "1")
}

func TestManySimpleReferenceExpressionsInBody(t *testing.T) {
	runSuccess(t, "(let ((a 1) (b 2)) a b)", "2")
}

func TestManyComplexExpressionsInBody(t *testing.T) {
	runSuccess(t, "(let ((a 1) (b 2)) (+ a b) (+ b a))", "3")
}

func TestMalformedLetList(t *testing.T) {
	runFail(t, "(let 1)", "Malformed let expression")
}

func TestMalformedLetAssignment(t *testing.T) {
	runFail(t, "(let (1))", "Let must have lists as arguments")
}

func TestMustBeVariableInLetPosition(t *testing.T) {
	runFail(t, "(let ((1 1)))", "Must be variable in let position")
}

func TestVariableAlreadySeen(t *testing.T) {
	runFail(t, "(let ((a 1)(a 2)))", "Variable already seen")
}
