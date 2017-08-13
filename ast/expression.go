package ast

import (
	"fmt"
	"sol/runtime"
)

type IdentifierExpression struct {
	Literal string
}

func (e *IdentifierExpression) ToString() string {
	return string(e.Literal)
}

func (e *IdentifierExpression) EvaluateExpr(*runtime.Scope) *runtime.Object {
	// TODO: get value for this identifier from scope
	return nil
}

type IntegerExpression struct {
	Literal string
}

func (e *IntegerExpression) ToString() string {
	return e.Literal
}

func (e *IntegerExpression) EvaluateExpr(*runtime.Scope) *runtime.Object {
	// TODO: return runtime number
	return nil
}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) ToString() string {
	return fmt.Sprintf(
		"%s %s %s",
		e.Left.ToString(),
		e.Operator,
		e.Right.ToString(),
	)
}

func (e *InfixExpression) EvaluateExpr(*runtime.Scope) *runtime.Object {
	// TODO: Evaluate by getting a function from a map[operator]func
	// for the type of this InfixExpression (infix or prefix)
	return nil
}

type ClosedExpression struct {
	Expression Expression
}

func (e *ClosedExpression) ToString() string {
	return fmt.Sprintf("(%s)", e.Expression.ToString())
}

func (e *ClosedExpression) EvaluateExpr(scope *runtime.Scope) *runtime.Object {
	return e.Expression.EvaluateExpr(scope)
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
func (e *FunctionExpression) EvaluateExpr(scope *runtime.Scope) *runtime.Object {
	// TODO: Evaluate
	return nil
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

func (e *CallExpression) EvaluateExpr(scope *runtime.Scope) *runtime.Object {
	// TODO: evaluate
	return nil
}
