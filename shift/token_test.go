package shift

import "testing"

func TestFetch(t *testing.T) {

	tests := []struct {
		input   string
		pattern string
		match   string
	}{
		{"this is a sentence", "this", "this"},
	}

	for i, tt := range tests {
		tz := NewTokenizer(tt.input)
		tz.Register(NewTokenType("", tt.pattern, 0))
		tz.fetch()
		tz.fetch()

		if tz.current == nil {
			t.Fatalf("TestFetch[%d]: no match, remaining=%s",
				i, tz.Remainder())
		}
		if tz.current.Literal != tt.match {
			t.Fatalf("TestFetch[%d]: expected=%s got=%s",
				i, tt.match, tz.current.Literal)
		}
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		input      string
		firstPatt  string
		firstPrec  int
		secondPatt string
		secondPrec int
		remainder  string
		matches    int
	}{
		{"==", "==", 10, "=", 5, "", 1},
		{"==", "==", 5, "=", 10, "", 2},
	}

	for i, tt := range tests {
		tz := NewTokenizer(tt.input)
		tz.Register(NewTokenType("first", tt.firstPatt, tt.firstPrec))
		tz.Register(NewTokenType("second", tt.secondPatt, tt.secondPrec))
		tz.fetch()
		tz.fetch()

		count := 0
		for tz.current != nil {
			count++
			tz.fetch()
		}

		if count != tt.matches {
			t.Fatalf("TestPrecedence[%d]: matches expected=%d got=%d",
				i, tt.matches, count)
		}
		if tz.Remainder() != tt.remainder {
			t.Fatalf("TestPrecedence[%d]: remainder expected=%s got=%s",
				i, tt.remainder, tz.Remainder())
		}

	}

}

func TestRepetition(t *testing.T) {
	tests := []struct {
		input   string
		pattern string
		output  string
		matches int
	}{
		{"a a a a", "a", "a(a),a(a),a(a),a(a),eof(),", 5},
	}

	for i, tt := range tests {
		tz := NewTokenizer(tt.input)
		tz.Register(NewTokenType("a", tt.pattern, 10))
		tz.Register(NewTokenType("eof", "$", 0))
		tz.fetch()
		tz.fetch()

		output := ""
		for j := 0; j < tt.matches; j++ {
			output += tz.Current().ToString() + ","
			tz.fetch()
		}

		if tt.output != output {
			t.Fatalf("TestRepetition[%d]: expected=%s got=%s",
				i, tt.output, output)
		}

	}

}
