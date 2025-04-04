package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

func (s *StringLiteral) ToSql() (string, []any, error) { //nolint:revive
	return "?", []any{s.StringValue}, nil
}

func (n *NumberLiteral) ToSql() (string, []any, error) { //nolint:revive
	return "?", []any{n.NumberValue}, nil
}

func (b *BoolLiteral) ToSql() (string, []any, error) { //nolint:revive
	return "?", []any{b.BoolValue}, nil
}

func (i Identifier) ToSql() (string, []any, error) { //nolint:revive
	return string(i), nil, nil
}

func (o *OneOfExpr) ToSql() (string, []any, error) { //nolint:revive
	return "?", []any{o.Value()}, nil
}

func (b *BinaryExpr) ToSql() (string, []any, error) { //nolint:revive
	switch b.Op {
	case And:
		return sq.And{b.Left, b.Right}.ToSql()
	case Or:
		return sq.Or{b.Left, b.Right}.ToSql()
	}

	return "", nil, fmt.Errorf("unknown operator %q", b.Op)
}

func (n *NotExpr) ToSql() (string, []any, error) { //nolint:revive
	sql, args, err := n.Expr.ToSql()
	if err != nil {
		return "", nil, err
	}

	return sq.Expr("NOT "+sql, args...).ToSql()
}

func (f *FieldExpr) ToSql() (string, []any, error) { //nolint:revive
	field, value := f.Field.String(), f.Value.Value()

	var sqlizer sq.Sqlizer

	switch f.Op {
	case Equal:
		sqlizer = sq.Eq{field: value}
	case NotEqual:
		sqlizer = sq.NotEq{field: value}
	case GreaterThan:
		sqlizer = sq.Gt{field: value}
	case GreaterThanOrEqual:
		sqlizer = sq.GtOrEq{field: value}
	case LessThan:
		sqlizer = sq.Lt{field: value}
	case LessThanOrEqual:
		sqlizer = sq.LtOrEq{field: value}
	case Like:
		sqlizer = sq.Like{field: value}
	case Exists:
		sqlizer = sq.Expr("? IS NOT NULL", value)
	default:
		return "", nil, fmt.Errorf("unknown operator %q", f.Op)
	}

	return sqlizer.ToSql()
}
