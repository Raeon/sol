package shift

type Any interface{}

type NodeParser func(node *Node) Any

type Node struct {
	children []*Node
	result   Any
	body     *RuleBody
}

func NewNode(body *RuleBody, children []*Node) *Node {
	return &Node{
		body:     body,
		result:   nil,
		children: children,
	}
}
