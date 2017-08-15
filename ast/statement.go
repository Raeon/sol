package ast

import (
	"fmt"
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

type ReturnStatement struct {
	Expression Expression
}

func (rs *ReturnStatement) ToString() string {
	return fmt.Sprintf(
		"return %s\n",
		rs.Expression.ToString(),
	)
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
