package main

import (
	"fmt"
	"sol/shift"
)

func main() {

	// reader := bufio.NewReader(os.Stdin)
	// parser := parser.NewParser()
	// env := runtime.NewEnv()

	// fmt.Println("sol 0.0.1-alpha")
	// for {
	// 	fmt.Printf("> ")
	// 	input, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		fmt.Printf("\nfatal: \"%s\"\n", err.Error())
	// 		return
	// 	}

	// 	prog, err := parser.Parse(input)
	// 	if err != nil {
	// 		fmt.Printf("error: \"%s\"\n", input)
	// 		continue
	// 	}

	// 	result := env.Evaluate(prog)
	// 	fmt.Printf("%s\n", result.ToString())
	// }

	ttZero := shift.NewTokenType("0", "0", 0)
	ttOne := shift.NewTokenType("1", "1", 0)
	ttPlus := shift.NewTokenType("+", "\\+", 0)
	ttAsterisk := shift.NewTokenType("*", "\\*", 0)

	g := shift.NewGrammar()

	g.Define("B", shift.NewTokenSymbol(ttZero))
	g.Define("B", shift.NewTokenSymbol(ttOne))
	g.Define("E", shift.NewReferenceSymbol(g.Get("B")))
	g.Define("E",
		shift.NewReferenceSymbol(g.Get("E")),
		shift.NewTokenSymbol(ttPlus),
		shift.NewReferenceSymbol(g.Get("B")),
	)
	g.Define("E",
		shift.NewReferenceSymbol(g.Get("E")),
		shift.NewTokenSymbol(ttAsterisk),
		shift.NewReferenceSymbol(g.Get("B")),
	)
	g.Define("S",
		shift.NewReferenceSymbol(g.Get("E")),
	)

	b := shift.NewBuilder(g)
	b.Build("S")
	fmt.Println(b.ToString())

}
