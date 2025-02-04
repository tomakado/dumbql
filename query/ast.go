package query

import (
	"fmt"
	"strconv"
)

// --- AST Definitions ---

// Expr is the interface for all expressions.
type Expr interface {
	fmt.Stringer
}

// BinaryExpr represents a binary operation (AND, OR) between two expressions.
type BinaryExpr struct {
	Left  Expr
	Op    CondOperator // "AND" or "OR"
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
	return fmt.Sprintf("not %s", n.Expr)
}

// FieldExpr represents a field query, e.g. status:200.
type FieldExpr struct {
	Field string
	Op    CompareOperator
	Value Expr
}

func (f *FieldExpr) String() string {
	return fmt.Sprintf("(%s %s %v)", f.Op, f.Field, f.Value)
}

// StringLiteral represents a bare term (a free text search term).
type StringLiteral struct {
	Value string
}

func (t *StringLiteral) String() string {
	return strconv.Quote(t.Value)
}

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%f", n.Value)
}

type IntegerLiteral struct {
	Value int64
}

func (i *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type OneOfExpr struct {
	Values []Expr
}

func (o *OneOfExpr) String() string {
	return fmt.Sprintf("%v", o.Values)
}

type CondOperator uint8

const (
	And CondOperator = iota + 1
	Or
)

func (c CondOperator) String() string {
	switch c {
	case And:
		return "and"
	case Or:
		return "or"
	default:
		return "unknown!"
	}
}

type CompareOperator uint8

const (
	Equal CompareOperator = iota + 1
	NotEqual
	GreaterThan
	GreaterThanOrEqual
	LessThan
	LessThanOrEqual
)

func (c CompareOperator) String() string {
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
	default:
		return "unknown!"
	}
}
