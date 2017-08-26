package shift

import "fmt"

type Closure struct {
	perms []*Permutation

	// Typically these cells determine whether the next parser
	// action is shift or reduce, but at this point they can
	// only be shifts.
	lookaheads map[*TokenType]*Closure

	// These cells show which state to advance to after
	// some reduction's left hand side has created an
	// expected new instance of that symbol.
	gotos map[string]*Closure
}

func NewClosure(p []*Permutation) *Closure {
	c := &Closure{
		perms:      p,
		lookaheads: make(map[*TokenType]*Closure),
		gotos:      make(map[string]*Closure),
	}
	return c
}

func (c *Closure) GetNextPermutations() map[Symbol][]*Permutation {
	// Map next symbol to list of permutations
	perms := make(map[Symbol][]*Permutation)

	// Fill the map
outer:
	for _, perm := range c.perms {
		nextSym := perm.NextSymbol()

		// If nextSym is nil, ignore this item
		if nextSym == nil {
			continue
		}

		// If it exists, then get the permutation for the
		// same rule body AFTER we parsed this symbol.
		nextPerm := perm.body.perms[perm.index+1]

		// Store it in existing list if any
		for k, v := range perms {
			if k.IsEqual(nextSym) {
				perms[k] = append(v, nextPerm)
				continue outer
			}
		}

		// Otherwise, create new list
		perms[nextSym] = []*Permutation{nextPerm}
	}
	return perms
}

func (c *Closure) getCompletedRuleBody() *RuleBody {
	if len(c.perms) != 1 {
		return nil
	}
	perm := c.perms[0]
	if perm.NextSymbol() != nil {
		return nil
	}
	return perm.body
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

	// // Lookaheads
	// result += "\nlookahead ("
	// first := true
	// for tt, _ := range c.lookaheads {
	// 	if !first {
	// 		result += ", "
	// 	}
	// 	first = false
	// 	result += tt.name
	// }
	// result += ")"

	// // Gotos
	// result += "\ngotos ("
	// first = true
	// for name, _ := range c.gotos {
	// 	if !first {
	// 		result += ", "
	// 	}
	// 	first = false
	// 	result += name
	// }
	// result += ")"

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

	// First we find all permutations that are comparable
	// to the current permutation and add them recursively.
	p.fillClosureChild(c, b)

}

func (p *Permutation) fillClosureChild(c *Closure, b *Builder) {

	// Get our own next symbol
	nextSym := p.NextSymbol()
	if nextSym == nil {
		return
	}

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
		producer.fillClosureChild(c, b)
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
		result += sym.ToSymbolString()
	}
	if p.index == len(p.body.symbols) {
		result += " ."
	}

	return fmt.Sprintf("%s = %s", p.body.rule.name, result)
}
