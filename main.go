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

	g := shift.NewGrammar()

	ttZero := g.Token("0", "0", 0)
	ttOne := g.Token("1", "1", 0)
	ttPlus := g.Token("+", "\\+", 0)
	ttAsterisk := g.Token("*", "\\*", 0)

	rB := g.Rule("B", parseB)
	rB.Body(ttZero)
	rB.Body(ttOne)

	rE := g.Rule("E", parseE)
	rE.Body(rB)
	rE.Body(rE, ttPlus, rB)
	rE.Body(rE, ttAsterisk, rB)

	rS := g.Rule("S", parseS)
	rS.Body(rE)

	b := shift.NewBuilder(g)
	b.Build("S")
	fmt.Printf("%s\n", b.ToString())
}

func parseS(node *shift.Node) shift.Any {
	return nil
}

func parseE(node *shift.Node) shift.Any {
	return nil
}

func parseB(node *shift.Node) shift.Any {
	return nil
}
