package shift

/* TODO: Write table that holds instructions for parser */
type Action struct {
	reduce  bool
	closure *Closure
}

func NewShiftAction(closure *Closure) *Action {
	return &Action{
		reduce:  false,
		closure: closure,
	}
}

func NewReduceAction(closure *Closure) *Action {
	return &Action{
		reduce:  true,
		closure: closure,
	}
}

type State struct {
	id int

	// These cells determine whether the next parser action is
	// shift or reduce.
	lookahead map[*TokenType]*Action

	// These cells show which state to advance to after
	// some reduction's left hand side has created an
	// expected new instance of that symbol.
	lhsGoto map[int]int
}

func NewState(id int) *State {
	return &State{
		id: id,
	}
}

type Table struct {
	states map[int]*State
}

func NewTable() *Table {
	return &Table{
		states: make(map[int]*State),
	}
}
