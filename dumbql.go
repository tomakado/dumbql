// Package dumbql provides simple (dumb) query language and it's parser.
//
// # Features
//
//   - Field expressions
//   - Boolean expressions
//   - One-of/In expressions
//   - Schema validation
//   - Drop-in usage with github.com/Masterminds/squirrel or SQL drivers directly
//
// # Query syntax
//
// The section below is a non-formal description of DumbQL syntax.
//
// Strict rules are expressed in grammar form:
//
//	Expr                <- _ e:OrExpr _
//	OrExpr              <- left:AndExpr rest:(_ ( OrOp ) _ AndExpr)*
//	OrOp                <- ("OR" / "or")
//	AndExpr             <- left:NotExpr rest:(_ ( op:AndOp ) _ NotExpr)*
//	AndOp               <- ("AND" / "and")
//	NotExpr             <- ("NOT" / "not") _ expr:Primary
//	                    / Primary
//	Primary             <- ParenExpr / FieldExpr
//	ParenExpr           <- '(' _ expr:Expr _ ')'
//	FieldExpr           <- field:Identifier _ op:CmpOp _ value:Value
//	Value               <- OneOfExpr / String / Number / Identifier
//	OneOfValue          <- String / Number / Identifier
//	Identifier          <- AlphaNumeric ("." AlphaNumeric)*
//	AlphaNumeric        <- [a-zA-Z_][a-zA-Z0-9_]*
//	Integer             <- '0' / NonZeroDecimalDigit DecimalDigit*
//	Number              <- '-'? Integer ( '.' DecimalDigit+ )?
//	DecimalDigit        <- [0-9]
//	NonZeroDecimalDigit <- [1-9]
//	String              <- '"' StringValue '"'
//	StringValue         <- ( !EscapedChar . / '\\' EscapeSequence )*
//	EscapedChar         <- [\x00-\x1f"\\]
//	EscapeSequence      <- SingleCharEscape / UnicodeEscape
//	SingleCharEscape    <- ["\\/bfnrt]
//	UnicodeEscape       <- 'u' HexDigit HexDigit HexDigit HexDigit
//	HexDigit            <- [0-9a-f]i
//	CmpOp               <- ( ">=" / ">" / "<=" / "<" / "!:" / "!=" / ":" / "=" )
//	OneOfExpr           <- '[' _ values:(OneOfValues)? _ ']'
//	OneOfValues         <- head:OneOfValue tail:(_ ',' _ OneOfValue)*
//	_                   <- [ \t\r\n]*
//
// # Field expression
//
// Field name & value pair divided by operator. Field name is any alphanumeric identifier (with underscore), value can be string, int64 or floa64.
// One-of expression is also supported (see below).
//
//	<field_name> <operator> <value>
//
// for example
//
//	period_months < 4
//
// # Field expression operators
//
//	| Operator             | Meaning       | Supported types              |
//	|----------------------|---------------|------------------------------|
//	| `:` or `=`           | Equal, one of | `int64`, `float64`, `string` |
//	| `!=` or `!:`         | Not equal     | `int64`, `float64`, `string` |
//	| `>`, `>=`, `<`, `<=` | Comparison    | `int64`, `float64`           |
//
// # Boolean operators
//
// Multiple field expression can be combined into boolean expressions with `and` (`AND`) or `or` (`OR`) operators:
//
//	status:pending and period_months < 4 and (title:"hello world" or name:"John Doe")
//
// # “One of” expression
//
// Sometimes instead of multiple `and`/`or` clauses against the same field:
//
//	occupation = designer or occupation = "ux analyst"
//
// it's more convenient to use equivalent “one of” expressions:
//
//	occupation: [designer, "ux analyst"]
//
// # Numbers
//
// If number does not have digits after `.` it's treated as integer and stored as int64. And it's float64 otherwise.
//
// # Strings
//
// String is a sequence on Unicode characters surrounded by double quotes (`"`). In some cases like single word it's possible to write string value without double quotes.
//
//status:pending and period_months < 4 and (title:"hello world" or name:"John Doe")
package dumbql

import (
	"github.com/defer-panic/dumbql/query"
	"github.com/defer-panic/dumbql/schema"
)

type Query struct {
	query.Expr
}

// Parse parses the input query string q, returning a Query reference or an error in case of invalid input.
func Parse(q string, opts ...query.Option) (*Query, error) {
	res, err := query.Parse("query", []byte(q), opts...)
	if err != nil {
		return nil, err
	}

	return &Query{res.(query.Expr)}, nil
}

// Validate checks the query against the provided schema, returning a validated expression or an error if any rule is violated.
// Even when error returned Validate can return query AST with invalided nodes dropped.
func (q *Query) Validate(s schema.Schema) (query.Expr, error) {
	return q.Expr.Validate(s)
}

// ToSql converts the Query into an SQL string, returning the SQL string, arguments slice, and any potential error encountered.
func (q *Query) ToSql() (string, []any, error) {
	return q.Expr.ToSql()
}
