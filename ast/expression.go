package ast

import (
	"fmt"
)

type IdentifierExpression struct {
	Literal string
}

func (e *IdentifierExpression) ToString() string {
	return e.Literal
}

type IntegerExpression struct {
	Literal string
}

func (e *IntegerExpression) ToString() string {
	return e.Literal
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) ToString() string {
	var left, right string
	if e.Left != nil {
		left = e.Left.ToString()
	}
	if e.Right != nil {
		right = e.Right.ToString()
	}

	return fmt.Sprintf(
		"(%s %s %s)",
		left,
		e.Operator,
		right,
	)
}

type ClosedExpression struct {
	Expression Expression
}

func (e *ClosedExpression) ToString() string {
	return fmt.Sprintf("(%s)", e.Expression.ToString())
}

type FunctionExpression struct {
	Parameters []string
	Body       Statement
}

func (e *FunctionExpression) ToString() string {
	str := "fn("
	for i, param := range e.Parameters {
		if i > 0 {
			str += ", "
		}
		str += param
	}
	return str + ")" + e.Body.ToString()
}

type CallExpression struct {
	Identifier string
	Arguments  []Expression
}

func (e *CallExpression) ToString() string {
	str := e.Identifier + "("
	for i, arg := range e.Arguments {
		if i > 0 {
			str += ", "
		}
		str += arg.ToString()
	}
	return str + ")"
}
