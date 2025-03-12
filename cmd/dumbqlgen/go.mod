module go.tomakado.io/dumbql/cmd/dumbqlgen

go 1.24.0

require golang.org/x/tools v0.31.0

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/Masterminds/squirrel v1.5.4 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/stretchr/testify v1.10.0
	go.tomakado.io/dumbql v0.4.0
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
)

// Use the main module to access the match package
replace go.tomakado.io/dumbql => ../../
