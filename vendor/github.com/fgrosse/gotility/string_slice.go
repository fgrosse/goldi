package gotility

import "math/rand"

// StringSlice adds useful functions to a slice of strings.
type StringSlice []string

// RandomStringSlice returns a new StringSlice filled with n elements.
// The elements are chosen uniformly from the given elements using math/rand.
// Remember to seed math/rand in order to get truly random slices.
func RandomStringSlice(n int, elements ...string) StringSlice {
	s := StringSlice{}
	if len(elements) == 0 {
		elements = append(elements, "")
	}

	for i := 0; i < n; i++ {
		s.Add(elements[rand.Intn(len(elements))])
	}

	return s
}

// Contains returns true if this slice contains the given string.
func (t StringSlice) Contains(s string) bool {
	for i := range t {
		if t[i] == s {
			return true
		}
	}

	return false
}

// Add adds the given string to this slice.
func (t *StringSlice) Add(s string) {
	*t = append(*t, s)
}

// AddAll adds all of the given strings to this slice.
// This is separate from Add due to performance reasons.
func (t *StringSlice) AddAll(s ...string) {
	*t = append(*t, s...)
}

// DeleteByIndex effectively removes the element at index i from this slice by
// appending the slice that starts after the index to the slice that goes up to the index.
func (t *StringSlice) DeleteByIndex(i int) bool {
	if i < 0 || i >= len(*t) {
		return false
	}

	*t = append((*t)[:i], (*t)[i+1:]...)
	return true
}

// SearchAndDelete searches s in this slice and deletes it if it is found.
// Note that deleting by index is of course way faster and should be used if possible.
func (t *StringSlice) DeleteByValue(s string) bool {
	for i := range *t {
		if (*t)[i] == s {
			t.DeleteByIndex(i)
			return true
		}
	}

	return false
}

// Reverse reverses the element in this slice in place.
func (t *StringSlice) Reverse() {
	var tmp string
	for i := 0; i < len(*t)/2; i++ {
		j := len(*t) - i - 1
		tmp = (*t)[i]

		// swap the element
		(*t)[i] = (*t)[j]
		(*t)[j] = tmp
	}
}
