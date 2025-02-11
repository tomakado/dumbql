package query

import (
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"go.tomakado.io/dumbql/schema"
)

//go:generate go run github.com/mna/pigeon@v1.3.0 -optimize-grammar -optimize-parser -o parser.gen.go grammar.peg

type Expr interface {
	fmt.Stringer
	sq.Sqlizer

	Match(target any, matcher Matcher) bool
	Validate(schema.Schema) (Expr, error)
}

type Valuer interface {
	Value() any
	Match(target any, op FieldOperator) bool
}

// BinaryExpr represents a binary operation (`and`, `or`, `AND`, `OR`) between two expressions.
type BinaryExpr struct {
	Left  Expr
	Op    BooleanOperator // `and` or `or`
	Right Expr
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Op, b.Left, b.Right)
}

// NotExpr represents a NOT expression.
type NotExpr struct {
	Expr Expr
}

func (n *NotExpr) String() string {
	return fmt.Sprintf("(not %s)", n.Expr)
}

// FieldExpr represents a field query, e.g. status:200.
type FieldExpr struct {
	Field Identifier
	Op    FieldOperator
	Value Valuer
}

func (f *FieldExpr) String() string {
	return fmt.Sprintf("(%s %s %v)", f.Op, f.Field, f.Value)
}

// StringLiteral represents a bare term (a free text search term).
type StringLiteral struct {
	StringValue string
}

func (s *StringLiteral) String() string { return strconv.Quote(s.StringValue) }
func (s *StringLiteral) Value() any     { return s.StringValue }

type NumberLiteral struct {
	NumberValue float64
}

func (n *NumberLiteral) String() string { return fmt.Sprintf("%f", n.NumberValue) }
func (n *NumberLiteral) Value() any     { return n.NumberValue }

type IntegerLiteral struct {
	IntegerValue int64
}

func (i *IntegerLiteral) String() string { return strconv.FormatInt(i.IntegerValue, 10) }
func (i *IntegerLiteral) Value() any     { return i.IntegerValue }

type Identifier string

func (i Identifier) Value() any     { return string(i) }
func (i Identifier) String() string { return string(i) }

type OneOfExpr struct {
	Values []Valuer
}

func (o *OneOfExpr) String() string { return fmt.Sprintf("%v", o.Values) }

func (o *OneOfExpr) Value() any {
	vals := make([]any, 0, len(o.Values))

	for _, v := range o.Values {
		vals = append(vals, v.Value())
	}

	return vals
}

type BooleanOperator uint8

const (
	And BooleanOperator = iota + 1
	Or
)

func (c BooleanOperator) String() string {
	switch c {
	case And:
		return "and"
	case Or:
		return "or"
	default:
		return "unknown!"
	}
}

type FieldOperator uint8

const (
	Equal FieldOperator = iota + 1
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
	Like
)

func (c FieldOperator) String() string {
	switch c {
	case Equal:
		return "="
	case NotEqual:
		return "!="
	case GreaterThan:
		return ">"
	case GreaterThanOrEqual:
		return ">="
	case LessThan:
		return "<"
	case LessThanOrEqual:
		return "<="
	case Like:
		return "~"
	default:
		return "unknown!"
	}
}
