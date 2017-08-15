package main

import (
	"bufio"
	"fmt"
	"os"
	"sol/parser"
	"sol/runtime"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	parser := parser.NewParser()
	env := runtime.NewEnv()

	fmt.Println("sol 0.0.1-alpha")
	for {
		fmt.Printf("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("\nfatal: \"%s\"\n", err.Error())
			return
		}

		prog, err := parser.Parse(input)
		if err != nil {
			fmt.Printf("error: \"%s\"\n", input)
			continue
		}

		result := env.Evaluate(prog)
		fmt.Printf("%s\n", result.ToString())
	}

}
