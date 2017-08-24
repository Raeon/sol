package shift

import "fmt"

type ActionType int

const (
	Shift ActionType = iota
	Reduce
	Accept
)

/* TODO: Write table that holds instructions for parser */
type Action struct {
	actionType ActionType
	stateID    int
}

func NewShiftAction(state *State) *Action {
	return &Action{
		actionType: Shift,
		stateID:    state.id,
	}
}

func NewReduceAction(state *State) *Action {
	return &Action{
		actionType: Reduce,
		stateID:    state.id,
	}
}

func NewAcceptAction() *Action {
	return &Action{
		actionType: Accept,
	}
}

type State struct {
	id      int
	closure *Closure

	// These cells determine whether the next parser action is
	// shift or reduce.
	lookaheads map[*TokenType]*Action

	// These cells show which state to advance to after
	// some reduction's left hand side has created an
	// expected new instance of that symbol.
	// Here, the key is the rule name and the value is
	// the stateID to move to after detecting said rule.
	gotos map[string]int
}

func NewState(id int) *State {
	return &State{
		id:         id,
		lookaheads: make(map[*TokenType]*Action),
		gotos:      make(map[string]int),
	}
}

func (s *State) Fill(b *Builder, c *Closure) {
	// Lookaheads are, at this point, only shift actions.
	// Here we convert all closures to their corresponding states.
	for tt, closure := range c.lookaheads {
		s.lookaheads[tt] = NewShiftAction(b.stateByClosure[closure])
	}

	// Gotos goes from map[string]*Closure to map[string]int,
	// where the string is the name of a given rule.
	for name, closure := range c.gotos {
		s.gotos[name] = b.stateByClosure[closure].id
	}
}

func (s *State) ToString() string {
	result := fmt.Sprintf("id = %d\n", s.id)

	result += "lookaheads:\n"
	for tt, action := range s.lookaheads {
		name := ""
		if tt != nil {
			name = tt.name
		} else {
			name = "eof"
		}
		result += "  " + name + "="
		if action.actionType == Reduce {
			result += "r"
			result += fmt.Sprintf("%d\n", action.stateID)
		} else if action.actionType == Shift {
			result += "s"
			result += fmt.Sprintf("%d\n", action.stateID)
		} else if action.actionType == Accept {
			result += "acc\n"
		}
	}

	result += "gotos:\n"
	for name, id := range s.gotos {
		result += fmt.Sprintf("  %s=%d\n", name, id)
	}

	return result
}

type Table struct {
	states map[int]*State
}

func NewTable() *Table {
	return &Table{
		states: make(map[int]*State),
	}
}

func (t *Table) ToString() string {
	result := ""
	for _, state := range t.states {
		result += state.closure.ToString() + "\n"
		result += state.ToString() + "\n"
	}
	return result
}
