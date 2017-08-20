package shift

type Grammar struct {
	names map[string]Symbol
	rules map[Symbol]*Rule
}

func NewGrammar() *Grammar {
	return &Grammar{
		names: make(map[string]Symbol),
		rules: make(map[Symbol]*Rule),
	}
}

func (g *Grammar) Declare(name string) *Rule {
	nameSym, ok := g.names[name]
	if ok {
		return g.rules[nameSym]
	}

	rule := NewRule(name)
	nameSym = NewReferenceSymbol(rule)

	g.names[name] = nameSym
	g.rules[nameSym] = rule

	return rule
}

func (g *Grammar) Define(name string, symbols ...Symbol) *Rule {
	rule := g.Declare(name)
	rule.AddBody(symbols...)
	return rule
}

func (g *Grammar) Get(name string) *Rule {
	nameSym, ok := g.names[name]
	if !ok {
		return nil
	}
	return g.rules[nameSym]
}

func (g *Grammar) findPermsRelative(sym Symbol, rindex int) []*Permutation {
	result := []*Permutation{}
	for _, rule := range g.rules {
		result = append(result, rule.findPermsRelative(sym, rindex)...)
	}
	return result
}

func (g *Grammar) findPermsProducing(sym Symbol) []*Permutation {
	// Check if this is a non-terminal, which we need.
	refSym, ok := sym.(*ReferenceSymbol)
	if !ok {
		return []*Permutation{}
	}

	// Return all permutations of this rule!
	return refSym.rule.perms
}
