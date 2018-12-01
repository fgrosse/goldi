package goldi

// A StringSet represents a set of strings.
type StringSet map[string]struct{}

// Set adds a value to the set.
func (s StringSet) Set(value string) {
	s[value] = struct{}{}
}

// Contains returns true if the given value is contained in this string set.
func (s StringSet) Contains(value string) bool {
	_, exists := s[value]

	return exists
}
