package shift

import "fmt"

type Rule struct {
	name   string
	bodies []*RuleBody
	perms  []*Permutation
	parser NodeReducer
}

func NewRule(name string, parser NodeReducer) *Rule {
	return &Rule{
		name:   name,
		bodies: []*RuleBody{},
		parser: parser,
	}
}

func (r *Rule) Body(symbols ...Symbol) *RuleBody {
	body := NewBody(r, symbols...)
	r.bodies = append(r.bodies, body)
	r.perms = append(r.perms, body.perms...)
	return body
}

func (r *Rule) findPermsRelative(sym Symbol, rindex int) []*Permutation {
	results := []*Permutation{}
	for _, body := range r.bodies {
		results = append(results, body.findPermsRelative(sym, rindex)...)
	}
	return results
}

func (r *Rule) IsTerminal() bool {
	return false
}

func (r *Rule) IsEqual(other Symbol) bool {
	o, ok := other.(*Rule)
	if !ok {
		return false
	}
	return o.name == r.name
}

func (r *Rule) ToSymbolString() string {
	return r.name
}

func (r *Rule) ToString() string {
	result := ""
	for i, body := range r.bodies {
		if i != 0 {
			result += "\n"
		}
		result += body.ToString()
	}
	return result
}

type RuleBody struct {
	rule    *Rule
	symbols []Symbol
	perms   []*Permutation
}

func NewBody(rule *Rule, symbols ...Symbol) *RuleBody {
	b := &RuleBody{
		rule:    rule,
		symbols: symbols,
	}

	// Create permutations
	b.perms = []*Permutation{}
	for i := 0; i <= len(b.symbols); i++ {
		b.perms = append(b.perms, NewPermutation(b, i))
	}

	return b
}

func (b *RuleBody) findPermsRelative(sym Symbol, rindex int) []*Permutation {
	results := []*Permutation{}
	for _, perm := range b.perms {
		if perm.hasSymbolRelative(sym, rindex) {
			results = append(results, perm)
		}
	}
	return results
}

func (b *RuleBody) ToString() string {
	result := ""
	for i, sym := range b.symbols {
		if i != 0 {
			result += " "
		}
		result += sym.ToSymbolString()
	}
	return fmt.Sprintf("%s = %s", b.rule.name, result)
}
