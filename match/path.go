package match

import "iter"

func Path(s string) iter.Seq[string] {
	return func(yield func(string) bool) {
		start := 0
		for i := range len(s) {
			if s[i] == '.' {
				if !yield(s[start:i]) {
					return
				}
				start = i + 1
			}
		}
		yield(s[start:])
	}
}
