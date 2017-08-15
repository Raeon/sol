package parser

import (
	"fmt"
	"strings"
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
		p := NewParser()
		prog, err := p.Parse(tt.input)
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
		{"5 + x", "(5 + x)", false},
		{"{ 5 + 5 6 + 6 }", "{(5 + 5)(6 + 6)}", false},
		{"(5 + 91) - 43", "((5 + 91) - 43)", false},
		{"5 + 5", "(5 + 5)", false},
		{"1 + 2 * 3", "(1 + (2 * 3))", false},
		{"x", "x", false},
		{"5 + 10", "(5 + 10)", false},
		{"let x = 5", "let x = 5", false},
		{"x = 10", "(x = 10)", false},
	}

	for i, tt := range tests {
		p := NewParser()
		prog, err := p.Parse(tt.input)
		if err != nil && !tt.fail {
			t.Fatalf("TestParseExpressionStatement[%d]: %s", i, err.Error())
		}

		progStr := strings.Replace(prog.ToString(), "\n", "", -1)
		if progStr != tt.output {
			t.Fatalf("TestParseExpressionStatement[%d]: expected=\"%s\" got=\"%s\"",
				i, tt.output, progStr)
		}
		fmt.Print(progStr + "\n")
	}
}
