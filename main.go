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

	ttInt := g.Token("int", "[0-9]+", 0)
	ttID := g.Token("id", "[a-zA-Z]+", 0)
	ttPlus := g.Token("+", "\\+", 0)
	ttAsterisk := g.Token("*", "\\*", 0)

	//ttEqual := g.Token("=", "=", 0)

	rValue := g.Rule("value", parseB)
	rValue.Body(ttInt)
	rValue.Body(ttID)

	rProducts := g.Rule("products", parseB)
	rProducts.Body(rProducts, ttAsterisk, rValue)
	rProducts.Body(rValue)

	rSums := g.Rule("sums", parseB)
	rSums.Body(rSums, ttPlus, rProducts)
	rSums.Body(rProducts)

	rGoal := g.Rule("goal", parseB)
	rGoal.Body(rSums)

	// b := shift.NewBuilder(g)
	// tbl := b.Build("S'")
	// fmt.Println(tbl.ToString())

	parser := g.Parser("goal")
	node, err := parser.Parse("5 + 5")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	} else {
		fmt.Printf("success: %s\n", node.ToString())
	}
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
