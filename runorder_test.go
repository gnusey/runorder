package runorder

import (
	"reflect"
	"sort"
	"testing"
)

func Test_removeAtIndex(t *testing.T) {
	tests := map[string]struct {
		slice []string
		index int
		exp   []string
	}{
		"Test removeAtIndex last index": {
			slice: []string{"a", "b", "c"},
			index: 2,
			exp:   []string{"a", "b"},
		},
		"Test removeAtIndex first index": {
			slice: []string{"a", "b", "c"},
			index: 0,
			exp:   []string{"b", "c"},
		},
	}

	for title, test := range tests {
		if res := removeAtIndex(test.slice, test.index); !reflect.DeepEqual(res, test.exp) {
			t.Fatalf("%s failed, expected %s, got %s", title, test.exp, res)
		}
	}
}

func Test_deleteReference(t *testing.T) {
	tests := map[string]struct {
		m   map[string][]string
		ref []string
		exp map[string][]string
	}{
		"Test deleteReference reference not in map": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c"},
				"c": []string{},
			},
			ref: []string{"d"},
			exp: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c"},
				"c": []string{},
			},
		},
		"Test deleteReference": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c"},
				"c": []string{},
			},
			ref: []string{"c"},
			exp: map[string][]string{
				"a": []string{"b"},
				"b": []string{},
			},
		},
		"Test deleteReference multiple": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c"},
				"c": []string{},
			},
			ref: []string{"b", "c"},
			exp: map[string][]string{
				"a": []string{},
			},
		},
		"Test deleteReference reference empty map": {
			m:   map[string][]string{},
			ref: []string{"d"},
			exp: map[string][]string{},
		},
		"Test deleteReference with multiple references": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c", "c"},
				"c": []string{},
			},
			ref: []string{"c"},
			exp: map[string][]string{
				"a": []string{"b"},
				"b": []string{},
			},
		},
	}

	for title, test := range tests {
		deleteReference(test.m, test.ref...)
		if res := test.m; !reflect.DeepEqual(res, test.exp) {
			t.Fatalf("%s failed, expected %+v, got %+v", title, test.exp, res)
		}
	}
}

func compare(t *testing.T, a, b [][]string) bool {
	t.Helper()

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		sort.Strings(a[i])
		sort.Strings(b[i])
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

func Test_calculate(t *testing.T) {
	tests := map[string]struct {
		m   map[string][]string
		exp [][]string
	}{
		"Test calculate A": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c"},
				"c": []string{},
			},
			exp: [][]string{
				[]string{"c"},
				[]string{"b"},
				[]string{"a"},
			},
		},
		"Test calculate B": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"c", "d", "e"},
				"c": []string{},
				"d": []string{},
				"e": []string{"f"},
				"f": []string{},
			},
			exp: [][]string{
				[]string{"f", "d", "c"},
				[]string{"e"},
				[]string{"b"},
				[]string{"a"},
			},
		},
		"Test calculate C": {
			m:   map[string][]string{},
			exp: [][]string{},
		},
		"Test calculate D": {
			m: map[string][]string{
				"a": []string{"b"},
			},
			exp: [][]string{
				[]string{"b"},
				[]string{"a"},
			},
		},
	}

	for title, test := range tests {
		if res := calculate(test.m); !compare(t, res, test.exp) {
			t.Fatalf("%s failed, expected %s, got %s", title, test.exp, res)
		}
	}
}

func Test_cp(t *testing.T) {
	exp := map[string][]string{
		"a": []string{"b"},
		"b": []string{},
	}
	if res := cp(exp); !reflect.DeepEqual(res, exp) || &exp == &res {
		t.Fatalf("Test cp failed, expected %+v/%+v, got %+v/%+v", exp, &exp, res, &res)
	}
}

func Test_indexOf(t *testing.T) {
	tests := map[string]struct {
		slice  []string
		target string
		exp    int
	}{
		"Test indexOf element exists": {
			slice:  []string{"a", "b", "c"},
			target: "a",
			exp:    0,
		},
		"Test indexOf element not exists": {
			slice:  []string{"a", "b", "c"},
			target: "d",
			exp:    -1,
		},
	}

	for title, test := range tests {
		if res := indexOf(test.slice, test.target); res != test.exp {
			t.Fatalf("%s failed, expected %d, got %d", title, test.exp, res)
		}
	}
}

func Test_checkCircularReference(t *testing.T) {
	tests := map[string]struct {
		m   map[string][]string
		err bool
	}{
		"Test without circular reference": {
			m: map[string][]string{
				"a": []string{"b"},
				"c": []string{},
			},
			err: false,
		},
		"Test with circular reference": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"a"},
			},
			err: true,
		},
	}

	for title, test := range tests {
		if err := checkCircularReference(test.m); (err != nil) != test.err {
			t.Fatalf("%s failed, expected error %t, error %s", title, test.err, err)
		}
	}
}

func Test_Calculate(t *testing.T) {
	tests := map[string]struct {
		m    map[string][]string
		copy bool
		exp  [][]string
		err  bool
	}{
		"Test Calculate": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{},
			},
			copy: false,
			exp: [][]string{
				[]string{"b"},
				[]string{"a"},
			},
			err: false,
		},
		"Test Calculate with copy": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{},
			},
			copy: true,
			exp: [][]string{
				[]string{"b"},
				[]string{"a"},
			},
			err: false,
		},
		"Test Calculate with circular reference": {
			m: map[string][]string{
				"a": []string{"b"},
				"b": []string{"a"},
			},
			copy: false,
			exp:  nil,
			err:  true,
		},
		"Test Calculate with multiple reference": {
			m: map[string][]string{
				"a": []string{"b", "b"},
				"b": []string{},
			},
			copy: false,
			exp: [][]string{
				[]string{"b"},
				[]string{"a"},
			},
			err: false,
		},
	}

	for title, test := range tests {
		c := cp(test.m)
		if res, err := Calculate(test.m, test.copy); !compare(t, res, test.exp) || (err != nil) != test.err || (test.copy && !reflect.DeepEqual(test.m, c)) {
			t.Fatalf("%s failed, expected %+v, got %+v, should error %t, error %s, use copy %t, original %+v", title, test.exp, res, test.err, err, test.copy, test.m)
		}
	}
}

func Test_Reverse(t *testing.T) {
	exp := [][]string{
		[]string{"c"},
		[]string{"b"},
		[]string{"a"},
	}
	if res := Reverse([][]string{
		[]string{"a"},
		[]string{"b"},
		[]string{"c"},
	}); !reflect.DeepEqual(res, exp) {
		t.Fatalf("Test Reverse failed, expected %s, got %s", exp, res)
	}
}
