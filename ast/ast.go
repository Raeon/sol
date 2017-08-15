package ast

type Node interface {
	ToString() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
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
