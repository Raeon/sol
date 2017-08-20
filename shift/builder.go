package shift

import (
	"fmt"
)

/* TODO: Write generator that creates a table */
type Builder struct {
	grammar     *Grammar
	table       *Table
	closureMap  map[*Permutation]*Closure
	closureList []*Closure
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

	// Get closure from the rule body and map the
	// permutations to the closure to prevent duplicate closures

	// From this point out we start traversing the state tree.
	// We iterate over all upcoming terminal symbols in the closure.
	// For every terminal symbol we then proceed to grab and register
	// the closure of the state *after* that terminal symbol is parsed.

	// TODO: Make sure body has at least 1 symbol and therefore
	// also at least one permutation
	b.buildChildClosures(body.perms[0])

	return b.table
}

func (b *Builder) buildChildClosures(perm *Permutation) *Closure {
	str := perm.ToString()
	closure, ok := b.closureMap[perm]
	if ok {
		// We already have this closure! Ignore it.
		str := closure.ToString()
		fmt.Sprintf(str)
		return closure
	}
	fmt.Sprintf(str)

	// Construct the closure for this terminal
	closure = NewClosure(perm)
	perm.fillClosure(closure, b)
	b.closureList = append(b.closureList, closure)

	// Link all the permutations to the closure to prevent
	// building their closures multiple times
	b.closureMap[closure.perms[0]] = closure

	// Next, we fetch all possible valid symbols that
	// this closure can parse and fetch the closures that
	// belong to the states after we parse those symbols.
	nextPerms := closure.GetNextPermutations()
	for _, nextPerm := range nextPerms {

		// // Get the symbol that this permutation has just consumed
		// prevSymbol := nextPerm.PreviousSymbol()
		// prevTerminal, ok := prevSymbol.(*TokenSymbol)

		// // If it wasn't a terminal we ignore it
		// if !ok {
		// 	continue
		// }

		// Get the symbol we just jumped over
		str := nextPerm.ToString()
		fmt.Sprintf(str)
		prevSymbol := nextPerm.PreviousSymbol()
		prevTerminal, isTerminal := prevSymbol.(*TokenSymbol)

		// Build the closure for this new found state
		childClosure := b.buildChildClosures(nextPerm)

		// If the previous symbol was a terminal as far as we know, then
		// we will register an action in the current closure for shifting
		// to the child closure upon parsing said terminal symbol.
		if isTerminal {
			closure.lookahead[prevTerminal.tokenType] = NewShiftAction(childClosure)
		}

	}

	return closure
}

func (b *Builder) ToString() string {
	result := ""
	for i, closure := range b.closureList {
		if i != 0 {
			result += "\n"
		}
		result += closure.ToString() + "\n"
	}
	return result
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
