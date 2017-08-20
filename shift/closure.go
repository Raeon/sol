package shift

import "fmt"

type Closure struct {
	perms []*Permutation

	// These cells determine whether the next parser action is
	// shift or reduce.
	lookahead map[*TokenType]*Action

	// These cells show which state to advance to after
	// some reduction's left hand side has created an
	// expected new instance of that symbol.
	lhsGoto map[int]int
}

func NewClosure(p *Permutation) *Closure {
	c := &Closure{
		perms:     []*Permutation{p},
		lookahead: make(map[*TokenType]*Action),
		lhsGoto:   make(map[int]int),
	}
	return c
}

func (c *Closure) GetNextPermutations() []*Permutation {
	perms := []*Permutation{}
	for _, perm := range c.perms {
		nextSym := perm.NextSymbol()

		// If this is nil, ignore it
		if nextSym == nil {
			continue
		}

		// If it *IS*, then get the permutation for the
		// same rule body AFTER we parsed this symbol.
		nextPerm := perm.body.perms[perm.index+1]

		// And we store it in the list.
		perms = append(perms, nextPerm)
	}
	return perms
}

func (c *Closure) Contains(p *Permutation) bool {
	for _, perm := range c.perms {
		if perm == p {
			return true
		}
	}
	return false
}

func (c *Closure) IsEqual(other *Closure) bool {
	return c == other
}

func (c *Closure) ToString() string {
	result := ""

	// Permutations
	for i, perm := range c.perms {
		if i != 0 {
			result += "\n"
		}
		result += perm.ToString()
	}

	// Lookaheads
	result += " ("
	first := true
	for tt, act := range c.lookahead {
		if !first {
			result += ", "
		}
		first = false
		result += fmt.Sprintf("%q=%v", tt.name, act.reduce)
	}
	result += ")"

	return result
}

type Permutation struct {
	body  *RuleBody
	index int
}

func NewPermutation(body *RuleBody, index int) *Permutation {
	return &Permutation{
		body:  body,
		index: index,
	}
}

func (p *Permutation) PreviousSymbol() Symbol {
	if p.index < 1 {
		return nil
	}
	return p.body.symbols[p.index-1]
}

func (p *Permutation) NextSymbol() Symbol {
	if p.index >= len(p.body.symbols) {
		return nil
	}
	return p.body.symbols[p.index]
}

func (p *Permutation) hasSymbolRelative(sym Symbol, rindex int) bool {
	nindex := p.index + rindex
	if nindex < 0 || nindex >= len(p.body.symbols) {
		return false
	}
	return p.body.symbols[nindex].IsEqual(sym)
}

func (p *Permutation) fillClosure(c *Closure, b *Builder) {
	stack := []*Permutation{}

	// First we find all permutations that are comparable
	// to the current permutation and add them recursively.
	p.fillClosureChild(c, b, &stack)

	// Then we get our own next symbol
	nextSym := p.NextSymbol()

	if nextSym == nil {

		// 	// If we don't have a next symbol then we must find
		// 	// all permutations that have consumed the same symbol
		// 	// as we have.
		// 	prevSym := p.PreviousSymbol()

		// 	// Find all permutations that parsed this symbol
		// 	prevPerms := b.grammar.findPermsRelative(prevSym, -1)
		// outer1:
		// 	for _, prevPerm := range prevPerms {

		// 		// Check if we already have it
		// 		for _, perm := range c.perms {
		// 			if perm == prevPerm {
		// 				continue outer1
		// 			}
		// 		}

		// 		// Add it
		// 		c.perms = append(c.perms, prevPerm)
		// 	}

		return
	}

	// We want to add all permutations that produce the same
	// symbol next as this current permutation.
	nextPerms := b.grammar.findPermsRelative(nextSym, 0)

	// And add them to the closure provided they aren't in it yet
outer2:
	for _, nextPerm := range nextPerms {

		// Don't add if already in it
		for _, perm := range c.perms {
			if perm == nextPerm {
				continue outer2
			}
		}

		// Add to closure
		c.perms = append(c.perms, nextPerm)

		// And apply the recursiveness again
		nextPerm.fillClosureChild(c, b, &stack)
	}

}

func (p *Permutation) fillClosureChild(c *Closure, b *Builder, s *[]*Permutation) {

	// Check if we are looping
	for _, perm := range *s {
		if perm == p {
			return
		}
	}

	// Get our own next symbol
	nextSym := p.NextSymbol()
	if nextSym == nil {
		return
	}

	// Add ourselves to the stack
	*s = append(*s, p)

	// Find all perms that produce our next symbol and are
	// at rindex=0
	producers := b.grammar.findPermsProducing(nextSym)

	// Add them to this closure
outer:
	for _, producer := range producers {

		// Has to be rindex=0!
		if producer.index != 0 {
			continue
		}

		// Only add if not already known
		for _, perm := range c.perms {
			if perm == producer {
				continue outer
			}
		}
		c.perms = append(c.perms, producer)

		// Recurse
		producer.fillClosureChild(c, b, s)
	}

}

/*
	Als ik *voor* alle symbols zit (index=0) dan bevat de closure
	naast de huidige permutatie alle permutaties die mijn volgende
	symbol produceren.

	Als ik *voor* een non-terminal symbol zit dan bevat mijn closure:
		- Alle permutaties die de huidige non-terminal produceren
		- Als de huidige non-terminal geproduceerd kan worden door
		  andere non-terminals, dan zitten de permutaties die DIE
		  non-terminals produceren OOK in de closure! Uiteraard zitten
	      zij dan allemaal op index=0.
*/

func (p *Permutation) IsEqual(other *Permutation) bool {
	// Theoretically speaking, this function is pointless, since all
	// possible permutations of a rule body are created only once.
	// But then again, when has theory actually worked as expected
	// in practice on the first try?
	if p == other {
		return true
	}
	return p.body == other.body &&
		p.index == other.index
}

func (p *Permutation) ToString() string {
	result := ""

	// Symbols and demarkator
	for i, sym := range p.body.symbols {
		if i == p.index {
			if i != 0 {
				result += " "
			}
			result += ". "
		}
		if i != 0 && i != p.index {
			result += " "
		}
		result += sym.ToString()
	}
	if p.index == len(p.body.symbols) {
		result += " ."
	}

	return fmt.Sprintf("%s = %s (%d)", p.body.rule.name, result, p.index)
}
