package runorder

import "github.com/pkg/errors"

// removeAtIndex removes an element at a specific index from a slice.
func removeAtIndex(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

// removeAll removes all occurrences of an element from a slice.
func removeAll(s []string, t string) []string {
	for i := indexOf(s, t); i > notFound; i = indexOf(s, t) {
		s = removeAtIndex(s, i)
	}

	return s
}

// deleteReference removes all references to value `refs` in `m`.
func deleteReference(m map[string][]string, refs ...string) {
	for k, v := range m {
		t := v
		for _, ref := range refs {
			t = removeAll(t, ref)
		}
		m[k] = t
	}
	for _, ref := range refs {
		delete(m, ref)
	}
}

// isEmpty is a helper function for checking if a string slice is empty.
func isEmpty(s []string) bool {
	return s == nil || len(s) == 0
}

// calculate calculates the run order.
func calculate(m map[string][]string) (r [][]string) {
	if len(m) == 0 {
		return r
	}

	var n []string
	for k, v := range m {
		if isEmpty(v) {
			n = append(n, k)
		}
	}
	if len(n) == 0 {
		for _, v := range m {
			for _, r := range v {
				if isEmpty(m[r]) {
					n = append(n, r)
				}
			}
		}
	}

	deleteReference(m, n...)
	r = append(r, n)

	return append(r, calculate(m)...)
}

// cp copies the contents of a map to a new map.
func cp(m map[string][]string) map[string][]string {
	c := make(map[string][]string)
	for k, v := range m {
		c[k] = append([]string{}, v...)
	}

	return c
}

// notFound represents the index of an element not present in a slice.
const notFound int = -1

// indexOf returns the index of an element in a slice. If the element does not
// exist the function returns `notFound`.
func indexOf(s []string, t string) int {
	for i, v := range s {
		if v == t {
			return i
		}
	}

	return notFound
}

// includes is a helper function for checking if an element is present in a slice.
func includes(s []string, t string) bool {
	return indexOf(s, t) > notFound
}

// ErrCircularReference represents a circular reference, where two elements in the map
// are mutually dependent on one another.
var ErrCircularReference = errors.New("error: circular reference")

// checkCircularReference checks for any circular references (jobs that mutually depend
// on one another) and returns an error if one is found.
func checkCircularReference(m map[string][]string) error {
	for k, v := range m {
		for _, r := range v {
			if includes(m[r], k) {
				return errors.WithMessagef(ErrCircularReference, "between %s and %s", k, r)
			}
		}
	}

	return nil
}

// Calculate returns the run order. Anything that can run concurrently is stored
// in the same slice. The function will mutate the map passed into it. If this is
// not acceptable set `c` to true to create a copy of the map.
func Calculate(m map[string][]string, c bool) ([][]string, error) {
	err := checkCircularReference(m)
	if err != nil {
		return nil, err
	}

	if c {
		return calculate(cp(m)), nil
	}
	return calculate(m), nil
}

// Reverse returns a copy of the reversed run order.
func Reverse(r [][]string) (reversed [][]string) {
	for i := len(r) - 1; i >= 0; i-- {
		reversed = append(reversed, r[i])
	}

	return
}
