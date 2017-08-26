package shift

import "fmt"

type ParseError struct {
	message string
}

func (e *ParseError) Error() string {
	return e.message
}

func err(message string) error {
	return &ParseError{message: message}
}

type Parser struct {
	table     *Table
	tokenizer *Tokenizer

	state       *State
	token       *Token
	outputStack []*Node
	stateStack  []int
}

func newParser(table *Table, tokenizer *Tokenizer) *Parser {
	return &Parser{
		table:       table,
		tokenizer:   tokenizer,
		outputStack: []*Node{},
		stateStack:  []int{},
	}
}

func (p *Parser) Parse(input string) (*Node, error) {

	// Load input into tokenizer
	p.tokenizer.Load(input)

	// Reset all variables
	p.stateStack = []int{0}
	p.outputStack = []*Node{}
	p.state = p.table.states[0]

	for {
		// Get current token
		p.token = p.tokenizer.Current()
		var tokenType *TokenType
		if p.token != nil {
			tokenType = p.token.Type
		}

		// Fetch action for current token
		action := p.state.lookaheads[tokenType]

		// Check for failure
		if action == nil {
			return nil, err(fmt.Sprintf("Unexpected token type: %q",
				tokenType))
		}

		// Execute action
		if action.actionType == Shift {
			// Shift stateID and token onto stacks
			p.stateStack = append(p.stateStack, action.stateID)
			p.outputStack = append(p.outputStack, newNodeToken(p.token))

			// Load said state
			p.state = p.table.states[p.stateStack[len(p.stateStack)-1]]

			// Move to next token
			p.tokenizer.fetch()

		} else if action.actionType == Reduce {

			// Take the amount of nodes necessary
			remaining := len(p.outputStack) - action.reduceCount
			children := make([]*Node, action.reduceCount)
			copy(children, p.outputStack[remaining:]) // before
			p.outputStack = p.outputStack[:remaining] // after

			// Build a new group node and apply the reduction
			node := newNodeGroup(children)
			node.result = action.reducer(node)

			// Push the result onto the output stack
			p.outputStack = append(p.outputStack, node)

			// Return to previous state by removing 'reduceCount' states
			p.stateStack = p.stateStack[:len(p.stateStack)-action.reduceCount]
			p.state = p.table.states[p.stateStack[len(p.stateStack)-1]]

			// Find new state from GOTO table
			newStateID := p.state.gotos[action.ruleName]
			p.stateStack = append(p.stateStack, newStateID)
			p.state = p.table.states[newStateID]

		} else if action.actionType == Accept {
			// TODO: Error if len(p.outputStack) != 1
			return p.outputStack[len(p.outputStack)-1], nil
		}

	}

	return nil, nil
}
