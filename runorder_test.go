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
		res := removeAtIndex(test.slice, test.index)
		if !reflect.DeepEqual(res, test.exp) {
			t.Fatalf("%s failed, expected %s, got %s", title, test.exp, res)
		}
	}
}

func Test_deleteReference(t *testing.T) {
	tests := map[string]struct {
		inMap map[string][]string
		ref   []string
		exp   map[string][]string
	}{
		"Test deleteReference reference not in map": {
			inMap: map[string][]string{
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
			inMap: map[string][]string{
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
			inMap: map[string][]string{
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
			inMap: map[string][]string{},
			ref:   []string{"d"},
			exp:   map[string][]string{},
		},
	}

	for title, test := range tests {
		deleteReference(test.inMap, test.ref...)
		res := test.inMap
		if !reflect.DeepEqual(res, test.exp) {
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
		inMap map[string][]string
		exp   [][]string
	}{
		"Test calculate A": {
			inMap: map[string][]string{
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
			inMap: map[string][]string{
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
			inMap: map[string][]string{},
			exp:   [][]string{},
		},
	}

	for title, test := range tests {
		res := calculate(test.inMap)
		if !compare(t, res, test.exp) {
			t.Fatalf("%s failed, expected %s, got %s", title, test.exp, res)
		}
	}
}

func Test_cp(t *testing.T) {
	title := "Test cp"
	exp := map[string][]string{
		"a": []string{"b"},
		"b": []string{},
	}
	res := cp(exp)
	if !reflect.DeepEqual(res, exp) || &exp == &res {
		t.Fatalf("%s failed, expected %+v/%+v, got %+v/%+v", title, exp, &exp, res, &res)
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
		res := indexOf(test.slice, test.target)
		if res != test.exp {
			t.Fatalf("%s failed, expected %d, got %d", title, test.exp, res)
		}
	}
}

func Test_checkCircularReference(t *testing.T) {
	tests := map[string]struct {
		inMap     map[string][]string
		shouldErr bool
	}{
		"Test without circular reference": {
			inMap: map[string][]string{
				"a": []string{"b"},
				"c": []string{},
			},
			shouldErr: false,
		},
		"Test with circular reference": {
			inMap: map[string][]string{
				"a": []string{"b"},
				"b": []string{"a"},
			},
			shouldErr: true,
		},
	}

	for title, test := range tests {
		err := checkCircularReference(test.inMap)
		if err != nil && !test.shouldErr {
			t.Fatalf("%s failed, expected error %t, error %s", title, test.shouldErr, err)
		}
	}
}

func Test_New(t *testing.T) {
	tests := map[string]struct {
		inMap     map[string][]string
		copy      bool
		exp       [][]string
		shouldErr bool
	}{
		"Test New": {
			inMap: map[string][]string{
				"a": []string{"b"},
				"b": []string{},
			},
			copy: false,
			exp: [][]string{
				[]string{"b"},
				[]string{"a"},
			},
			shouldErr: false,
		},
		"Test New with copy": {
			inMap: map[string][]string{
				"a": []string{"b"},
				"b": []string{},
			},
			copy: true,
			exp: [][]string{
				[]string{"b"},
				[]string{"a"},
			},
			shouldErr: false,
		},
		"Test New with circular reference": {
			inMap: map[string][]string{
				"a": []string{"b"},
				"b": []string{"a"},
			},
			copy:      false,
			exp:       nil,
			shouldErr: true,
		},
	}

	for title, test := range tests {
		c := cp(test.inMap)
		res, err := New(test.inMap, test.copy)
		if !compare(t, res, test.exp) || (err != nil && !test.shouldErr) || (test.copy && !reflect.DeepEqual(test.inMap, c)) {
			t.Fatalf("%s failed, expected %+v, got %+v, should error %t, error %s, use copy %t, original %+v", title, test.exp, res, test.shouldErr, err, test.copy, test.inMap)
		}
	}
}

func Test_Reverse(t *testing.T) {
	title := "Test Reverse"
	exp := [][]string{
		[]string{"c"},
		[]string{"b"},
		[]string{"a"},
	}
	res := Reverse([][]string{
		[]string{"a"},
		[]string{"b"},
		[]string{"c"},
	})
	if !reflect.DeepEqual(res, exp) {
		t.Fatalf("%s failed, expected %s, got %s", title, exp, res)
	}
}
