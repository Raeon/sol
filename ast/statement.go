package ast

import (
	"fmt"
	"sol/runtime"
)

type DeclarationStatement struct {
	Identifier string
	Expression Expression
}

func (ds *DeclarationStatement) ToString() string {
	return fmt.Sprintf(
		"let %s = %s\n",
		ds.Identifier,
		ds.Expression.ToString(),
	)
}

func (ds *DeclarationStatement) EvaluateStmt(scope *runtime.Scope) {
	// TODO: Evaluate
}

type ReturnStatement struct {
	Expression Expression
}

func (rs *ReturnStatement) ToString() string {
	return fmt.Sprintf(
		"return %s\n",
		rs.Expression.ToString(),
	)
}

func (rs *ReturnStatement) EvaluateStmt(scope *runtime.Scope) {
	// TODO: Evaluate
}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) ToString() string {
	return fmt.Sprintf(
		"%s\n",
		es.Expression.ToString(),
	)
}

func (es *ExpressionStatement) EvaluateStmt(scope *runtime.Scope) {
	es.Expression.EvaluateExpr(scope)
}

type BlockStatement struct {
	Statements []Statement
}

func (s *BlockStatement) ToString() string {
	str := "{\n"
	for _, stmt := range s.Statements {
		str += stmt.ToString() + "\n"
	}
	return str + "}"
}

func (s *BlockStatement) EvaluateStmt(scope *runtime.Scope) {
	// TODO: Evaluate
}
