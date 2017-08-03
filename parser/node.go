package parser

type ParserNode struct {
	Result   string
	Children []ParserNode
	Error    error
}

type Node interface {
	Name() string
	Children() []*Node
}
