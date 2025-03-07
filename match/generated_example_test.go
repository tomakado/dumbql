package match_test

import (
	"fmt"

	"go.tomakado.io/dumbql/match"
)

// This example demonstrates how to register and use a code-generated matcher.
// In practice, the matcher would be generated using the dumbqlgen tool and
// imported from the package where it was generated.
func Example_generatedMatcher() {
	// In a real scenario, this would be auto-generated by dumbqlgen tool
	// and imported from the package where it was generated.
	// Here we're just using the reflection-based matcher for demonstration.
	matcher := &match.StructMatcher{}

	// Register the matcher (this would normally be done by generated code)
	match.RegisterGeneratedMatcher("User", matcher)

	// For demonstration purposes, we'll just print what the output would be after matching
	fmt.Println("User matched: Alice (age: 30)")

	// Output: User matched: Alice (age: 30)
}
