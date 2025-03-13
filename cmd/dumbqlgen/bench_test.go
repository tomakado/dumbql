package main_test

import (
	"testing"
	"time"

	"go.tomakado.io/dumbql"
	"go.tomakado.io/dumbql/cmd/dumbqlgen/testdata"
	"go.tomakado.io/dumbql/match"
)

func BenchmarkReflectionRouter(b *testing.B) {
	user := &testdata.BenchUser{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		Active:    true,
	}

	expr, err := dumbql.Parse(`name:"John Doe" AND email:"john@example.com" AND age > 25 AND active:true`)
	if err != nil {
		b.Fatalf("Failed to parse query: %v", err)
	}

	matcher := &match.StructMatcher{} // This uses the reflection router by default

	for b.Loop() {
		_ = expr.Match(user, matcher)
	}
}

func BenchmarkGeneratedRouter(b *testing.B) {
	user := &testdata.BenchUser{
		ID:        123,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		CreatedAt: time.Now(),
		Active:    true,
	}

	expr, err := dumbql.Parse(`name:"John Doe" AND email:"john@example.com" AND age > 25 AND active:true`)
	if err != nil {
		b.Fatalf("Failed to parse query: %v", err)
	}

	matcher := testdata.NewBenchUserMatcher()

	for b.Loop() {
		_ = expr.Match(user, matcher)
	}
}
