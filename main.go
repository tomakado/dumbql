package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/defer-panic/dumbql/query"
)

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing readline: %v\n", err)
		os.Exit(1)
	}
	defer rl.Close()

	fmt.Println("DumbQL (c) defer panic Inc. 2077")
	fmt.Println("Enter your queries below (type 'exit' or 'quit' to leave).")

	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			}
			continue
		} else if err == io.EOF {
			break
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.EqualFold(trimmed, "exit") || strings.EqualFold(trimmed, "quit") {
			break
		}

		ast, err := query.Parse("stdin", []byte(trimmed))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println(ast.(query.Expr))
	}

	fmt.Println("Goodbye!")
}
