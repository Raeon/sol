package ast

import (
	"sol/runtime"
)

type Node interface {
	ToString() string
}

type Statement interface {
	Node
	EvaluateStmt(*runtime.Scope)
}

type Expression interface {
	Node
	EvaluateExpr(*runtime.Scope) *runtime.Object
}

type Program struct {
	Statements []Statement
}

func (p *Program) ToString() string {
	str := ""
	for i, stmt := range p.Statements {
		if i != 0 {
			str += "\n"
		}
		str += stmt.ToString()
	}
	return str
}
