package parser

import (
	"fmt"
	"testing"
)

func TestParseProgram(t *testing.T) {}

func TestParseStatement(t *testing.T) {}

func TestParseDeclarationStatement(t *testing.T) {}

func TestParseReturnStatement(t *testing.T) {
	tests := []struct {
		input  string
		output string
		fail   bool
	}{
		{"let a = 5", "let a = 5", false},
		{"return 5", "return 5", false},
	}

	for i, tt := range tests {
		p := NewParser(tt.input)
		prog, err := p.Parse()
		if err != nil {
			t.Fatalf("TestParseReturnStatement[%d]: %s", i, err.Error())
			return
		}
		fmt.Println(prog.ToString())
	}
}

func TestParseExpressionStatement(t *testing.T) {
	tests := []struct {
		input  string
		output string
		fail   bool
	}{
		// {"5 + add(9, 10)", "5 + add(9, 10)", false},
		{"(5 + 91) - 43", "(5 + 91) - 43", false},
		{"5 + 5", "5 + 5", false},
	}

	for i, tt := range tests {
		p := NewParser(tt.input)
		prog, err := p.Parse()
		if err != nil && !tt.fail {
			t.Fatalf("TestParseExpressionStatement[%d]: %s", i, err.Error())
		}
		fmt.Println(prog.ToString())
	}
}
