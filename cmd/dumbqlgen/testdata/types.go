package testdata

import (
	"time"
)

//go:generate go run go.tomakado.io/dumbql/cmd/dumbqlgen -type TestUser -package .

type TestUser struct {
	ID        int64  `dumbql:"id"`
	Name      string `dumbql:"name"`
	Email     string `dumbql:"email"`
	CreatedAt time.Time
	Address   Address `dumbql:"address"`
	Private   bool    `dumbql:"-"`
}

type Address struct {
	Street string `dumbql:"street"`
	City   string `dumbql:"city"`
	State  string `dumbql:"state"`
	Zip    string `dumbql:"zip"`
}

//go:generate go run go.tomakado.io/dumbql/cmd/dumbqlgen -type BenchUser -package .

type BenchUser struct {
	ID        int64  `dumbql:"id"`
	Name      string `dumbql:"name"`
	Email     string `dumbql:"email"`
	Age       int    `dumbql:"age"`
	CreatedAt time.Time
	Active    bool `dumbql:"active"`
}
