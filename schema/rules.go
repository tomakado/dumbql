package schema

import "fmt"

func Any(rules ...RuleFunc) RuleFunc {
	return func(field Field, value any) error {
		var err error
		for _, rule := range rules {
			if err = rule(field, value); err == nil {
				return nil
			}
		}
		return err
	}
}

func All(rules ...RuleFunc) RuleFunc {
	return func(field Field, value any) error {
		for _, rule := range rules {
			if err := rule(field, value); err != nil {
				return err
			}
		}
		return nil
	}
}

func InRange[T Numeric](min, max T) RuleFunc { //nolint:revive
	return func(field Field, value any) error {
		if v, ok := value.(T); ok {
			if v < min || v > max {
				return fmt.Errorf("field %q: value must be in range [%v, %v], got %v", field, min, max, v)
			}
			return nil
		}
		return fmt.Errorf("field %q: value must be %T, got %T", field, min, value)
	}
}

func Min[T Numeric](min T) RuleFunc { //nolint:revive
	return func(field Field, value any) error {
		if v, ok := value.(T); ok {
			if v < min {
				return fmt.Errorf("field %q: value must be equal or greater than %v, got %v", field, min, v)
			}
			return nil
		}
		return fmt.Errorf("field %q: value must be %T, got %T", field, min, value)
	}
}

func Max[T Numeric](max T) RuleFunc { //nolint:revive
	return func(field Field, value any) error {
		if v, ok := value.(T); ok {
			if v > max {
				return fmt.Errorf("field %q: value must be equal or less than %v, got %v", field, max, v)
			}
			return nil
		}
		return fmt.Errorf("field %q: value must be %T, got %T", field, max, value)
	}
}

func LenInRange(min, max int) RuleFunc { //nolint:revive
	return func(field Field, value any) error {
		if v, ok := value.(string); ok {
			if len(v) < min || len(v) > max {
				return fmt.Errorf("field %q: len must be in range [%d, %d], got %d", field, min, max, len(v))
			}
			return nil
		}
		return fmt.Errorf("field %q: value must be string, got %v", field, value)
	}
}

func MinLen(min int) RuleFunc { //nolint:revive
	return func(field Field, value any) error {
		if v, ok := value.(string); ok {
			if len(v) < min {
				return fmt.Errorf("field %q: len must be greater than %d, got %d", field, min, len(v))
			}
			return nil
		}
		return fmt.Errorf("field %q: value must be string, got %v", field, value)
	}
}

func MaxLen(max int) RuleFunc { //nolint:revive
	return func(field Field, value any) error {
		if v, ok := value.(string); ok {
			if len(v) > max {
				return fmt.Errorf("field %q: value must be less than %d, got %d", field, max, len(v))
			}
			return nil
		}
		return fmt.Errorf("field %q: value must be string, got %v", field, value)
	}
}

func Is[T ValueType]() RuleFunc {
	return func(field Field, value any) error {
		if v, ok := value.(T); !ok {
			return fmt.Errorf("field %q: value must be %T, got %T", field, v, value)
		}
		return nil
	}
}

func EqualsOneOf(values ...any) RuleFunc {
	return func(field Field, value any) error {
		for _, v := range values {
			if v == value {
				return nil
			}
		}
		return fmt.Errorf("field %q: value must be one of %v, got %v", field, values, value)
	}
}
