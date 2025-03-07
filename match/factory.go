package match

import (
	"fmt"
	"reflect"

	"go.tomakado.io/dumbql/query"
)

// MatcherType defines the type of matcher to use
type MatcherType int

const (
	// MatcherTypeReflection uses the reflection-based StructMatcher
	MatcherTypeReflection MatcherType = iota

	// MatcherTypeGenerated uses a code-generated matcher
	MatcherTypeGenerated
)

// CreateMatcher creates a matcher for the given target type.
// If matcherType is MatcherTypeGenerated, a code-generated matcher will be used if available,
// otherwise it falls back to the reflection-based StructMatcher.
// Note that code-generated matchers must be registered first with RegisterGeneratedMatcher.
func CreateMatcher(target interface{}, matcherType MatcherType) (query.Matcher, error) {
	switch matcherType {
	case MatcherTypeReflection:
		return &StructMatcher{}, nil
	case MatcherTypeGenerated:
		t := reflect.TypeOf(target)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		// Try to find a code-generated matcher for the type
		matcher, ok := generatedMatchers[t.Name()]
		if !ok {
			// Fall back to reflection-based matcher
			return &StructMatcher{}, nil
		}
		return matcher, nil
	default:
		return nil, fmt.Errorf("unknown matcher type: %v", matcherType)
	}
}

// generatedMatchers is a registry of code-generated matchers
var generatedMatchers = map[string]query.Matcher{}

// RegisterGeneratedMatcher registers a code-generated matcher for the given type name
func RegisterGeneratedMatcher(typeName string, matcher query.Matcher) {
	generatedMatchers[typeName] = matcher
}
