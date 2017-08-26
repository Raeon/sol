package shift

import "fmt"

/* TODO: Write generator that creates a table */
type Builder struct {
	grammar *Grammar
	table   *Table

	closureMap  map[*Permutation]*Closure
	closureList []*Closure

	stateByID      map[int]*State
	stateByClosure map[*Closure]*State
}

func NewBuilder(grammar *Grammar) *Builder {
	return &Builder{
		grammar:    grammar,
		closureMap: make(map[*Permutation]*Closure),
	}
}

func (b *Builder) Build(rootRule string) *Table {
	b.table = NewTable()

	// Find the root rule to derive the zero-state from
	rule := b.grammar.Get(rootRule)
	if rule == nil {
		panic("Cannot build parsing table from undefined root rule")
	}

	// Fail if body count != 1
	if len(rule.bodies) != 1 {
		panic("Root rule must have exactly 1 body")
	}
	body := rule.bodies[0]

	// Fail if body.perms count != 1
	if len(body.perms) < 1 {
		panic("Root rule must have at least 1 symbol")
	}

	// Get closure from the rule body and map the
	// permutations to the closure to prevent duplicate closures.

	// From this point out we start traversing the state tree.
	// We iterate over all upcoming terminal symbols in the closure.
	// For every terminal symbol we then proceed to grab and register
	// the closure of the state *after* that terminal symbol is parsed.

	// Build closure for the first permutation
	b.buildChildClosures([]*Permutation{body.perms[0]})

	// We map every Closure to a State, and every ID to a State.
	b.stateByClosure = make(map[*Closure]*State)
	b.stateByID = make(map[int]*State)

	// Next, we map every found closure to a corresponding State.
	for id, closure := range b.closureList {
		state := NewState(id)
		state.closure = closure // TODO: remove after debugging

		// Keep references to this state
		b.stateByID[id] = state
		b.stateByClosure[closure] = state

		// Put state in parse table
		b.table.states[id] = state
	}

	// Assign the closure *after* parsing the root node an
	// Accept action to accept the given input on EOF
	successClosure := b.closureMap[body.perms[1]]
	successState := b.stateByClosure[successClosure]
	successState.lookaheads[nil] = NewAcceptAction()

	// Once that's done, we can substitute all Closure
	// references with the newly assigned IDs.
	for closure, state := range b.stateByClosure {

		// Fill the lookaheads and gotos from the Closure
		state.Fill(b, closure)

		// Get a completed rule body
		body := closure.getCompletedRuleBody()
		if body != nil && state != successState {
			// If it exists, map all known token types to reduce action
			action := NewReduceAction(body)

			for _, tt := range b.grammar.tokenizer.types {
				// Log if any given action is already a reduce
				before, ok := state.lookaheads[tt]
				if ok {
					fmt.Printf("overriding action: type=%q, id=%d, name=%s",
						before.actionType, before.stateID, before.ruleName)
				}
				state.lookaheads[tt] = action
			}
			// EOF also!
			state.lookaheads[nil] = action
		}

	}

	// Return the now ready-to-use parsing table!
	return b.table
}

func (b *Builder) buildChildClosures(perms []*Permutation) *Closure {

	// Check if we already have a closure for these perms
	for _, perm := range perms {
		closure, ok := b.closureMap[perm]
		if ok {
			return closure
		}
	}

	// Construct the closure for these items
	closure := NewClosure(perms)

	// Fill the closure with derived items
	for _, perm := range perms {

		// Prevent duplicate closures by mapping a permutation
		// to the current closure
		b.closureMap[perm] = closure

		// Fill the closure using items derived from
		// the current base item
		perm.fillClosure(closure, b)
	}

	// Store the closure in a list for debugging?
	b.closureList = append(b.closureList, closure)

	// Next, we fetch all possible valid symbols that
	// this closure can parse and construct the closure that
	// belong to the states after we parse those symbols.
	nextPermMap := closure.GetNextPermutations()
	for nextSymbol, nextPerms := range nextPermMap {

		// Fetch the child closure for the new base items
		// that match the next parseable symbol
		childClosure := b.buildChildClosures(nextPerms)

		// Add a shift action from the current closure
		// with the given terminal symbol to the child closure.
		if nextSymbol.IsTerminal() {
			tokenType, _ := nextSymbol.(*TokenType)
			closure.lookaheads[tokenType] = childClosure
		} else {
			rule, _ := nextSymbol.(*Rule)
			closure.gotos[rule.name] = childClosure
		}
	}

	return closure
}

func (b *Builder) ToString() string {
	// result := ""
	// for i, closure := range b.closureList {
	// 	if i != 0 {
	// 		result += "\n"
	// 	}
	// 	result += closure.ToString() + "\n"
	// }
	// return result
	return b.table.ToString()
}

/*
	Een bepaalde state betreft de permutaties uit de huidige closure.
	Binnen de closure kan er maximaal 1 permutatie per LHS zijn. Dit is
	omdat wanneer we bij een reduce moeten kijken in de goto kolom
	naar welke state we daarna moeten gaan afhankelijk van de huidige state
	en de LHS.

	In onze implementatie betekent het dat wanneer we een child closure
	maken dat hij LHS goto's krijgt die verwijzen naar de oude state.

	Hieronder een voorbeeld:
	ID	Regel									Lookahead						Goto
	---------------------------------------------------------------------------------------
	8	Value		=	int	.					(...) r5 r5 r5
	9	Value		=	id .					(...) r6 r6 r6
	5	Products	=	Products * . Value		(...)							Value: 6
	6	Products	=	Products * Value .		(...)

	Wat we hier zien is dat wanneer een Value successvol zijn input verwerkt heeft
	hij bij alle terminal tokens die hij NIET kan krijgen zich reduceert naar
	zijn eerstvolgende ouder. Deze ouder was de state waarin hij net een Value
	probeerde te consumeren.
	Vervolgens hebben we dus een Value en de state 5 na een reductie, en we zien hier
	dat in de GOTO table van state 5 staat dat bij een geconsumeerde Value hij van
	state moet veranderen naar state 6, waarin hij successvol een Value geconsumeerd heeft.

	Wat ik echter niet begrijp is dat we hier vervolgens ook nog state 9 hebben.
	Deze state reduceert rechtstreeks naar state 6, welke vervolgens helemaal geen Goto
	tabel heeft klaarstaan.
	We zien eigenlijk dat hij meerdere nummers wil accepteren maar niet meerdere IDs,
	wat op zich vrij raar is, maar dit kan een mankement van het voorbeeld zijn.


*/
