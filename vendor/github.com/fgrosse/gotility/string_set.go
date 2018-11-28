package gotility

import "strings"

// The StringSet is a utility type to manage string sets.
type StringSet map[string]struct{}

func NewStringSet(initialValues ...string) StringSet {
	s := StringSet{}
	for i := range initialValues {
		s[initialValues[i]] = struct{}{}
	}

	return s
}

// Set registers the string in this hash set
func (s StringSet) Set(value string) {
	s[value] = struct{}{}
}

// SetAll registers all values in this hash set.
func (s StringSet) SetAll(value []string) {
	for _, v := range value {
		s.Set(v)
	}
}

// Contains returns true if the given value is contained in this string set.
func (s StringSet) Contains(value string) bool {
	_, exists := s[value]

	return exists
}

// Delete removes the value from this set if it exists.
// The return value indicates if the value existed in the first place.
func (s StringSet) Delete(value string) bool {
	exists := s.Contains(value)
	delete(s, value)

	return exists
}

// All returns all elements of this set
func (s StringSet) All() []string {
	keys := make([]string, len(s))
	i := 0
	for e := range s {
		keys[i] = e
		i++
	}

	return keys
}

// String returns a textual representation of this set in the form [key1, key2, key3]
func (s StringSet) String() string {
	return "[" + strings.Join(s.All(), ", ") + "]"
}
