package runorder

import "fmt"

// removeAtIndex removes an element at a specific index from a slice.
func removeAtIndex(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
}

// deleteReference removes all references to value `ref` in `m`.
func deleteReference(m map[string][]string, refs ...string) {
	for k, v := range m {
		t := v
		for _, ref := range refs {
			i := indexOf(v, ref)
			if i > -1 {
				t = removeAtIndex(t, i)
			}
		}
		m[k] = t
	}
	for _, ref := range refs {
		delete(m, ref)
	}
}

// calculate calculates the run order.
func calculate(m map[string][]string) (r [][]string) {
	if len(m) == 0 {
		return r
	}

	var n []string
	for k, v := range m {
		if v == nil || len(v) == 0 {
			n = append(n, k)
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
		c[k] = v
	}
	return c
}

// indexOf returns the index of an element in a slice. If the element does not
// exist the function returns -1.
func indexOf(s []string, t string) int {
	for i, v := range s {
		if v == t {
			return i
		}
	}
	return -1
}

// checkCircularReference checks for any circular references (jobs that mutually depend
// on one another) and returns an error if one is found.
func checkCircularReference(m map[string][]string) error {
	for k, v := range m {
		for _, r := range v {
			if indexOf(m[r], k) > -1 {
				return fmt.Errorf("error: circular reference found between %s and %s", k, r)
			}
		}
	}
	return nil
}

// New returns the run order. Anything that can run concurrently is stored in
// the same slice. The function will mutate the map passed into it. If this is
// not acceptable set `c` to true to create a copy of the map.
func New(m map[string][]string, c bool) ([][]string, error) {
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
