package lexer

import (
	"fmt"
	"strings"
)

type Lexer interface {
	Lex(s *Scanner) (*LexNode, *LexError)
	ToString() string
}

type LexNode struct {
	Name       string
	Value      string
	GroupName  string
	Children   []*LexNode
	Extra      interface{}
	LineNumber int
	LineIndex  int
}

func (n *LexNode) Groups(group string) ([]*LexNode, int) {
	// Perform a breadth-first search through child nodes
	// for the nodes with the given name
	layerCurrent := n.Children
	layerNext := []*LexNode{}
	results := []*LexNode{}
	depth := 0

	for len(layerCurrent) > 0 {

		// Iterate over all children in current layer
		for _, child := range layerCurrent {
			if child.GroupName == group {
				results = append(results, child)
			}
			for _, subchild := range child.Children {
				layerNext = append(layerNext, subchild)
			}
		}

		// If we have results, return them
		if len(results) > 0 {
			return results, depth
		}

		// Then, move to next layer
		depth++
		layerCurrent = layerNext
		layerNext = []*LexNode{}
	}

	return results, depth
}

func (n *LexNode) GroupNodes(group string) []*LexNode {
	nodes, _ := n.Groups(group)
	return nodes
}

func (n *LexNode) Group(group string) (*LexNode, int) {
	matches, depth := n.Groups(group)
	if len(matches) == 0 {
		return nil, -1
	}
	return matches[0], depth
}

func (n *LexNode) GroupNode(group string) *LexNode {
	node, _ := n.Group(group)
	return node
}

func (n *LexNode) GroupDepth(group string) int {
	_, depth := n.Group(group)
	return depth
}

func (n *LexNode) GroupExists(group string) bool {
	return n.GroupDepth(group) >= 0
}

func (n *LexNode) String(depth int) string {
	result := strings.Repeat("  ", depth) + "{\n"
	depth++

	value := n.Value
	value = strings.Replace(value, "\n", "\\n", -1)
	value = strings.Replace(value, "\r", "\\r", -1)
	value = strings.Replace(value, "\t", "\\t", -1)

	result += strings.Repeat("  ", depth) + "\"name\": \"" + n.Name + "\",\n"
	result += strings.Repeat("  ", depth) + "\"value\": \"" + value + "\",\n"
	result += strings.Repeat("  ", depth) + "\"group\": \"" + n.GroupName + "\",\n"
	result += strings.Repeat("  ", depth) + "\"children\": ["
	if len(n.Children) == 0 {
		result += "]\n"
		depth--
		result += strings.Repeat("  ", depth) + "}"
		return result
	}
	result += "\n"
	for _, child := range n.Children {
		result += child.String(depth+1) + ",\n"
	}
	result += strings.Repeat("  ", depth) + "]\n"
	depth--
	result += strings.Repeat("  ", depth) + "}"
	return result
}

type LexError struct {
	err   string
	stack []string
}

func err(s *Scanner, err string, lex Lexer) *LexError {
	e := &LexError{
		err:   err,
		stack: []string{},
	}
	e.Trace(s, lex)
	return e
}

func (e *LexError) Error() string {
	str := e.err
	for _, err := range e.stack {
		str += "\n" + err
	}
	return str
}

func (e *LexError) Trace(s *Scanner, lex Lexer) {
	e.stack = append(e.stack, fmt.Sprintf("at %d:%d lexing %s",
		s.lineNumber, s.lineIndex, lex.ToString()))
}

type AtomLexer struct {
	atom string
}

func Atom(atom string) Lexer {
	return &AtomLexer{atom: atom}
}

func (p *AtomLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	if s.ConsumeString(p.atom) != p.atom {
		return nil, err(s, fmt.Sprintf("Expected atom: %s", p.atom), p)
	}

	return &LexNode{
		Name:       "atom",
		Value:      p.atom,
		LineNumber: s.lineNumber,
		LineIndex:  s.lineIndex,
	}, nil
}

func (p *AtomLexer) ToString() string {
	return fmt.Sprintf("ATOM(\"%s\")", p.atom)
}

type RegexLexer struct {
	pattern    string
	allowEmpty bool
}

func Regex(pattern string, allowEmpty bool) Lexer {
	return &RegexLexer{
		pattern:    pattern,
		allowEmpty: allowEmpty,
	}
}

func (p *RegexLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	str := s.ConsumeRegex(p.pattern)
	if !p.allowEmpty && str == "" {
		return nil, err(s, fmt.Sprintf("Expected regex: %s", p.pattern), p)
	}

	return &LexNode{
		Name:       "regex",
		Value:      str,
		LineNumber: s.lineNumber,
		LineIndex:  s.lineIndex,
	}, nil
}

func (p *RegexLexer) ToString() string {
	return fmt.Sprintf("REGEX(/%s/)", p.pattern)
}

type AndLexer struct {
	children []Lexer
}

func And(children ...Lexer) Lexer {
	return &AndLexer{
		children: children,
	}
}

func (p *AndLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	var nodes []*LexNode
	var value string

	for _, child := range p.children {
		node, err := child.Lex(s)
		if err != nil {
			err.Trace(s, p)
			return nil, err
		}
		nodes = append(nodes, node)
		value += node.Value
	}

	return &LexNode{
		Name:       "and",
		Children:   nodes,
		Value:      value,
		LineNumber: s.lineNumber,
		LineIndex:  s.lineIndex,
	}, nil
}

func (p *AndLexer) ToString() string {
	str := "AND("
	for i, child := range p.children {
		if i != 0 {
			str += ", "
		}
		str += child.ToString()
	}
	str += ")"
	return str
}

type OrLexer struct {
	children []Lexer
}

func Or(children ...Lexer) Lexer {
	return &OrLexer{
		children: children,
	}
}

func (p *OrLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	var errors []string
	for _, child := range p.children {
		s.Push()
		node, err := child.Lex(s)
		if err == nil {
			s.Discard()
			return &LexNode{
				Name:       "or",
				Value:      node.Value,
				Children:   []*LexNode{node},
				LineNumber: s.lineNumber,
				LineIndex:  s.lineIndex,
			}, nil
		}
		errors = append(errors, "\t"+err.Error())
		s.Pop()
	}

	return nil, err(s, fmt.Sprintf(
		"Failed to lex OR:\n%s", strings.Join(errors, "\n")), p)
}

func (p *OrLexer) ToString() string {
	str := "OR("
	for i, child := range p.children {
		if i != 0 {
			str += ", "
		}
		str += child.ToString()
	}
	str += ")"
	return str
}

type RepeatLexer struct {
	child Lexer
	min   int
	max   int
}

func Repeat(child Lexer, min int, max int) Lexer {
	return &RepeatLexer{
		child: child,
		min:   min,
		max:   max,
	}
}

func Optional(child Lexer) Lexer {
	return &RepeatLexer{
		child: child,
		min:   0,
		max:   1,
	}
}

func (p *RepeatLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	var count int
	var nodes []*LexNode
	var value string

	// Loop while we haven't reached the maximum yet
	for p.max < 0 || count <= p.max {

		s.Push()
		node, err := p.child.Lex(s)
		if err != nil {
			// If we haven't reached minimum yet, return error
			if p.min >= 0 && count < p.min {
				s.Discard()
				err.Trace(s, p)
				return nil, err
			}

			// Otherwise, revert and return result
			s.Pop()
			break
		}
		nodes = append(nodes, node)
		value += node.Value
		count++
	}

	// Return the result
	return &LexNode{
		Name:       "repeat",
		Children:   nodes,
		Value:      value,
		LineNumber: s.lineNumber,
		LineIndex:  s.lineIndex,
	}, nil
}

func (p *RepeatLexer) ToString() string {
	return fmt.Sprintf("REPEAT(%s, %d, %d)",
		p.child.ToString(), p.min, p.max)
}

type GroupLexer struct {
	groupName string
	child     Lexer
}

func Group(name string, child Lexer) Lexer {
	return &GroupLexer{
		groupName: name,
		child:     child,
	}
}

func (p *GroupLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	node, err := p.child.Lex(s)
	if node != nil {
		node.GroupName = p.groupName
	}
	if err != nil {
		err.Trace(s, p)
	}
	return node, err
}

func (p *GroupLexer) ToString() string {
	return fmt.Sprintf("GROUP(\"%s\", %s)",
		p.groupName, p.child.ToString())
}

type IgnoreLexer struct {
	child Lexer
}

func Ignore(child Lexer) Lexer {
	return &IgnoreLexer{
		child: child,
	}
}

func (p *IgnoreLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	node, err := p.child.Lex(s)
	if node != nil {
		node.Value = ""
	}
	err.Trace(s, p)
	return node, err
}

func (p *IgnoreLexer) ToString() string {
	return fmt.Sprintf("IGNORE(%s)", p.child.ToString())
}

type FutureLexer struct {
	pointer *Lexer
	name    string
}

func Future(pointer *Lexer, name string) *FutureLexer {
	return &FutureLexer{
		pointer: pointer,
		name:    name,
	}
}

func (p *FutureLexer) Lex(s *Scanner) (*LexNode, *LexError) {
	return (*p.pointer).Lex(s)
}

func (p *FutureLexer) ToString() string {
	return fmt.Sprintf("FUTURE(\"%s\")", p.name)
}

type InterlaceLexer struct {
	outer Lexer
	inner Lexer
}

func Interlace(outer Lexer, inner Lexer) Lexer {
	return &InterlaceLexer{
		outer: outer,
		inner: inner,
	}
}

func (p *InterlaceLexer) Lex(s *Scanner) (*LexNode, *LexError) {

	var nodes []*LexNode
	var outer *LexNode
	var inner *LexNode
	var value string
	var err *LexError

	for {
		s.Push()
		if len(nodes) == 0 && outer == nil {
			outer, err = p.outer.Lex(s)
		} else if inner == nil {
			inner, err = p.inner.Lex(s)
		} else {
			outer, err = p.outer.Lex(s)
		}

		if err != nil {

			// We failed just now, so pop at least once.
			// Since the other half of the pair (if any)
			// isn't used, we need to pop for those as well!
			s.Pop()
			if outer != nil || inner != nil {
				s.Pop()
			}

			// If we already parsed at least one 'outer',
			// then we return success.
			if len(nodes) > 0 {
				return &LexNode{
					Name:       "interlace",
					Value:      value,
					Children:   nodes,
					LineNumber: s.lineNumber,
					LineIndex:  s.lineIndex,
				}, nil
			}

			// Otherwise, we return the error.
			err.Trace(s, p)
			return nil, err
		}

		// Push result if necessary
		if len(nodes) == 0 && outer != nil {
			nodes = append(nodes, outer)
			outer = nil
			s.Discard()
		} else if outer != nil && inner != nil {
			value += inner.Value
			value += outer.Value
			nodes = append(nodes, inner)
			nodes = append(nodes, outer)
			inner = nil
			outer = nil
			s.Discard()
			s.Discard()
		}
	}
}

func (p *InterlaceLexer) ToString() string {
	return fmt.Sprintf("INTERLACE(%s, %s)",
		p.outer.ToString(), p.inner.ToString())
}
