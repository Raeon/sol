package shift

type Any interface{}

type NodeReducer func(node *Node) Any

type Node struct {
	children []*Node
	token    *Token
	result   Any
}

func newNodeGroup(children []*Node) *Node {
	return &Node{
		children: children,
		result:   nil,
	}
}

func newNodeToken(token *Token) *Node {
	return &Node{
		token:  token,
		result: nil,
	}
}

func (n *Node) ToString() string {
	if n.token != nil {
		return n.token.Literal
	}

	result := "("
	for i, child := range n.children {
		if i != 0 {
			result += " "
		}
		result += child.ToString()
	}
	return result + ")"
}
