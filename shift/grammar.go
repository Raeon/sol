package shift

import "fmt"

type Grammar struct {
	tokenizer *Tokenizer
	names     map[string]Symbol
	rules     map[Symbol]*Rule
}

func NewGrammar() *Grammar {
	return &Grammar{
		tokenizer: NewTokenizer(),
		names:     make(map[string]Symbol),
		rules:     make(map[Symbol]*Rule),
	}
}

func (g *Grammar) Rule(name string, parser NodeReducer) *Rule {
	nameSym, ok := g.names[name]
	if ok {
		return g.rules[nameSym]
	}

	rule := NewRule(name, parser)

	g.names[name] = rule
	g.rules[rule] = rule

	return rule
}

func (g *Grammar) Token(name, pattern string, precedence int) *TokenType {
	return g.tokenizer.Type(name, pattern, precedence)
}

func (g *Grammar) Parser(rootRule string) *Parser {
	b := NewBuilder(g)
	tbl := b.Build(rootRule)
	fmt.Println(b.ToString())
	return newParser(tbl, b.grammar.tokenizer)
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
	rule, ok := sym.(*Rule)
	if !ok {
		return []*Permutation{}
	}

	// Return all permutations of this rule!
	return rule.perms
}
