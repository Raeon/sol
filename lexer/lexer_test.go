package lexer

import (
	"testing"
)

func TestAtomLexer(t *testing.T) {
	tests := []struct {
		input     string
		atom      string
		remainder string
	}{
		{"test", "test", ""},
		{"testing", "test", "ing"},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		ap := Atom(tt.atom)

		match, err := ap.Lex(s)
		if err != nil {
			t.Fatalf("TestAtomLexer[%d]: error=%s", i, err.Error())
		}
		if tt.atom != match.Value {
			t.Fatalf("TestAtomLexer[%d]: atom=%s match=%s",
				i, tt.atom, match.Value)
		}
		if tt.remainder != s.Remainder() {
			t.Fatalf("TestAtomLexer[%d]: expected=%s got=%s",
				i, tt.remainder, s.Remainder())
		}
	}
}

func TestRegexLexer(t *testing.T) {
	tests := []struct {
		input      string
		pattern    string
		allowEmpty bool
		expected   string
		fail       bool
	}{
		{"input", "in", false, "in", false},
		{"input", "ba", false, "", true},
		{"input", "[ia]nput", false, "input", false},
		{`"this is a \"string\""`, `"((?:\\\\|\\"|[^"])+)"`, false, `"this is a \"string\""`, false},
		{`"\\\\" "string"`, `"((?:\\\\|\\"|[^"])+)"`, false, `"\\\\"`, false},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		rp := Regex(tt.pattern, tt.allowEmpty)

		match, err := rp.Lex(s)
		if err != nil && !tt.fail {
			t.Fatalf("TestRegexLexer[%d]: unexpected error=%s",
				i, err.Error())
		} else if tt.fail {
			if err == nil {
				t.Fatalf("TestRegexLexer[%d]: expected error", i)
			}
			continue
		}
		if tt.expected != match.Value {
			t.Fatalf("TestRegexLexer[%d]: expected=%s got=%s",
				i, tt.expected, match.Value)
		}
	}
}

func TestAndLexer(t *testing.T) {
	tests := []struct {
		input     string
		atom1     string
		atom2     string
		remainder string
		fail      bool
	}{
		{"firstsecond", "first", "second", "", false},
		{"firstandthird", "first", "and", "third", false},
		{"secondandthird", "first", "second", "secondandthird", true},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		a1 := Atom(tt.atom1)
		a2 := Atom(tt.atom2)
		ap := And(a1, a2)

		match, err := ap.Lex(s)
		if err != nil && !tt.fail {
			t.Fatalf("TestAndLexer[%d]: unexpected error=%s",
				i, err.Error())
		} else if tt.fail {
			if err == nil {
				t.Fatalf("TestAndLexer[%d]: expected error", i)
			}
			continue
		}
		if match.Value != (tt.atom1 + tt.atom2) {
			t.Fatalf("TestAndLexer[%d]: expected=%s match=%s", i,
				tt.atom1+tt.atom2, match.Value)
		}
		if tt.remainder != s.Remainder() {
			t.Fatalf("TestAndLexer[%d]: expected=%s got=%s",
				i, tt.remainder, s.Remainder())
		}
	}
}

func TestOrLexer(t *testing.T) {
	tests := []struct {
		input string
		atom1 string
		atom2 string
		match string
		fail  bool
	}{
		{"A", "A", "B", "A", false},
		{"B", "A", "B", "B", false},
		{"C", "A", "B", "", true},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		a1 := Atom(tt.atom1)
		a2 := Atom(tt.atom2)
		op := Or(a1, a2)

		match, err := op.Lex(s)
		if err != nil && !tt.fail {
			t.Fatalf("TestOrLexer[%d]: unexpected error=%s",
				i, err.Error())
		} else if tt.fail {
			if err == nil {
				t.Fatalf("TestOrLexer[%d]: expected error", i)
			}
			continue
		}
		if match.Value != tt.match {
			t.Fatalf("TestOrLexer[%d]: expected=%s match=%s", i,
				tt.atom1+tt.atom2, match.Value)
		}
	}
}

func TestRepeatLexer(t *testing.T) {
	tests := []struct {
		input    string
		atom     string
		min      int
		max      int
		expected string
		fail     bool
	}{
		{"aaa", "a", 3, 3, "aaa", false},
		{"aaa", "a", 0, 5, "aaa", false},
		{"aaa", "a", 4, 4, "", true},
		{"aaa", "b", 0, 5, "", false},
		{"aaa", "b", 1, 5, "", true},
		{"aaa", "a", 0, -1, "aaa", false},
		{"aaa", "b", -1, -1, "", false},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		a := Atom(tt.atom)
		rp := Repeat(a, tt.min, tt.max)

		match, err := rp.Lex(s)
		if err != nil && !tt.fail {
			t.Fatalf("TestRepeatLexer[%d]: unexpected error=%s",
				i, err.Error())
		} else if tt.fail {
			if err == nil {
				t.Fatalf("TestRepeatLexer[%d]: expected error", i)
			}
			continue
		}
		if match.Value != tt.expected {
			t.Fatalf("TestRepeatLexer[%d]: expected=%s got=%s",
				i, tt.expected, match.Value)
		}
	}
}

func TestInterlaceLexer(t *testing.T) {
	tests := []struct {
		input     string
		outer     string
		inner     string
		result    int
		remainder string
		fail      bool
	}{
		{"ababa", "a", "b", 5, "", false},
		{"ababab", "a", "b", 5, "b", false},
		{"cababa", "a", "b", 0, "cababa", true},
	}

	for i, tt := range tests {
		s := NewScanner(tt.input)
		outer := Atom(tt.outer)
		inner := Atom(tt.inner)
		inter := Interlace(outer, inner)

		node, err := inter.Lex(s)
		if err != nil && !tt.fail {
			t.Fatalf("TestInterlaceLexer[%d]: unexpected error=%s",
				i, err.Error())
		} else if tt.fail {
			if err == nil {
				t.Fatalf("TestInterlaceLexer[%d]: expected error", i)
			}
			continue
		}
		if len(node.Children) != tt.result {
			t.Fatalf("TestInterlaceLexer[%d]: result expected=%d got=%d",
				i, tt.result, len(node.Children))
		}
		if s.Remainder() != tt.remainder {
			t.Fatalf("TestInterlaceLexer[%d]: remainder expected=\"%s\""+
				" got=\"%s\"", i, tt.remainder, s.Remainder())
		}
	}

}

func TestToString(t *testing.T) {
	tests := []struct {
		lex      Lexer
		expected string
	}{
		{Atom("string"), "ATOM(\"string\")"},
		{Regex("regexp?", false), "REGEX(/regexp?/)"},
		{And(Atom("a"), Atom("b")), "AND(ATOM(\"a\"), ATOM(\"b\"))"},
		{Or(Atom("a"), Atom("b")), "OR(ATOM(\"a\"), ATOM(\"b\"))"},
		{Repeat(Atom("a"), 1, 2), "REPEAT(ATOM(\"a\"), 1, 2)"},
		{Group("group", Atom("a")), "GROUP(\"group\", ATOM(\"a\"))"},
		{Ignore(Atom("a")), "IGNORE(ATOM(\"a\"))"},
		{Interlace(Atom("a"), Atom("b")), "INTERLACE(ATOM(\"a\"), ATOM(\"b\"))"},
	}

	for i, tt := range tests {
		out := tt.lex.ToString()
		if out != tt.expected {
			t.Fatalf("TestToString[%d]: expected=\"%s\" got=\"%s\"",
				i, tt.expected, out)
		}
	}
}
