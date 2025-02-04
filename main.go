package main

import (
	"fmt"
	"log"

	"github.com/defer-panic/dumbql/query"
)

func main() {
	// Example query string.
	// You can try different inputs, for example:
	// input := `status:200`
	// input := `NOT error`
	// input := `status:200 AND (extension:"jpg" OR extension:"png")`
	input := `status : 200 and eps < 0.003 and (req.fields.ext:["jpg", "png"])`
	fmt.Printf("input: %s\n", input)

	ast, err := query.Parse("query", []byte(input))
	if err != nil {
		log.Fatalf("Error parsing query: %s", err)
	}
	// fmt.Printf("AST:\n%#v\n", ast.(query.Expr).String())
	fmt.Println(ast.(query.Expr))
}
