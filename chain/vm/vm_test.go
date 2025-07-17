package vm

import "testing"

func TestVMAny(t *testing.T) {
	m := New()
	script := []Instruction{{OP_SET, "A"}, {OP_ANY, ""}}
	cases := []struct {
		name  string
		sets  map[string][]string
		votes map[string]bool
		pass  bool
	}{
		{"single yes", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v1": true}, true},
		{"all no", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v1": false, "v2": false}, false},
		{"two yes", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v1": true, "v2": true}, true},
		{"yes outside set", map[string][]string{"A": {"v1"}}, map[string]bool{"v2": true}, false},
		{"bigger set one yes", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v2": true}, true},
		{"bigger set none", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v4": true}, false},
		{"extra set ignored", map[string][]string{"A": {"v1"}, "B": {"v2"}}, map[string]bool{"v1": true}, true},
		{"empty set", map[string][]string{"A": []string{}}, map[string]bool{"v1": true}, false},
		{"repeated yes", map[string][]string{"A": {"v1", "v1"}}, map[string]bool{"v1": true}, true},
		{"partial votes", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v2": false}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.Validate(script, tc.votes, tc.sets)
			if res != tc.pass {
				t.Fatalf("expect %v got %v", tc.pass, res)
			}
		})
	}
}

func TestVMQuorum(t *testing.T) {
	m := New()
	script := []Instruction{{OP_SET, "A"}, {OP_QUORUM, "2"}}
	cases := []struct {
		name  string
		sets  map[string][]string
		votes map[string]bool
		pass  bool
	}{
		{"exact quorum", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v1": true, "v2": true}, true},
		{"above quorum", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v1": true, "v2": true, "v3": true}, true},
		{"below quorum", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v1": true}, false},
		{"no votes", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{}, false},
		{"quorum different set", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v3": true, "v4": true}, false},
		{"large set quorum", map[string][]string{"A": {"v1", "v2", "v3", "v4"}}, map[string]bool{"v1": true, "v4": true}, true},
		{"mixed votes", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v1": true, "v2": false, "v3": true}, true},
		{"duplicate yes", map[string][]string{"A": {"v1", "v1", "v2"}}, map[string]bool{"v1": true}, true},
		{"empty set", map[string][]string{"A": []string{}}, map[string]bool{"v1": true}, false},
		{"quorum not reached", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v3": true}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.Validate(script, tc.votes, tc.sets)
			if res != tc.pass {
				t.Fatalf("expect %v got %v", tc.pass, res)
			}
		})
	}
}

func TestVMAll(t *testing.T) {
	m := New()
	script := []Instruction{{OP_SET, "A"}, {OP_ALL, ""}}
	cases := []struct {
		name  string
		sets  map[string][]string
		votes map[string]bool
		pass  bool
	}{
		{"all yes", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v1": true, "v2": true}, true},
		{"one no", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v1": true, "v2": false}, false},
		{"missing vote", map[string][]string{"A": {"v1", "v2"}}, map[string]bool{"v1": true}, false},
		{"empty set", map[string][]string{"A": []string{}}, map[string]bool{}, true},
		{"large set all yes", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v1": true, "v2": true, "v3": true}, true},
		{"large set one no", map[string][]string{"A": {"v1", "v2", "v3"}}, map[string]bool{"v1": true, "v2": false, "v3": true}, false},
		{"yes outside set", map[string][]string{"A": {"v1"}}, map[string]bool{"v1": true, "v2": true}, true},
		{"duplicate validator", map[string][]string{"A": {"v1", "v1"}}, map[string]bool{"v1": true}, true},
		{"duplicate validator no", map[string][]string{"A": {"v1", "v1"}}, map[string]bool{"v1": false}, false},
		{"no votes", map[string][]string{"A": {"v1"}}, map[string]bool{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.Validate(script, tc.votes, tc.sets)
			if res != tc.pass {
				t.Fatalf("expect %v got %v", tc.pass, res)
			}
		})
	}
}

func TestVMThen(t *testing.T) {
	m := New()
	script := []Instruction{
		{OP_SET, "A"}, {OP_QUORUM, "2"}, {OP_THEN, ""},
		{OP_SET, "B"}, {OP_ALL, ""},
	}
	cases := []struct {
		name  string
		sets  map[string][]string
		votes map[string]bool
		pass  bool
	}{
		{
			"happy path",
			map[string][]string{"A": {"v1", "v2", "v3"}, "B": {"v4", "v5"}},
			map[string]bool{"v1": true, "v2": true, "v4": true, "v5": true},
			true,
		},
		{
			"fail second",
			map[string][]string{"A": {"v1", "v2", "v3"}, "B": {"v4", "v5"}},
			map[string]bool{"v1": true, "v2": true, "v5": true},
			false,
		},
		{
			"fail first",
			map[string][]string{"A": {"v1", "v2", "v3"}, "B": {"v4", "v5"}},
			map[string]bool{"v1": true, "v4": true, "v5": true},
			false,
		},
		{
			"different sets",
			map[string][]string{"A": {"x1", "x2"}, "B": {"y1", "y2"}},
			map[string]bool{"x1": true, "x2": true, "y1": true, "y2": true},
			true,
		},
		{
			"extra votes",
			map[string][]string{"A": {"v1", "v2"}, "B": {"v3"}},
			map[string]bool{"v1": true, "v2": true, "v3": true, "v4": true},
			true,
		},
		{
			"missing all second",
			map[string][]string{"A": {"v1", "v2"}, "B": {"v3"}},
			map[string]bool{"v1": true, "v2": true},
			false,
		},
		{
			"quorum not met",
			map[string][]string{"A": {"v1", "v2"}, "B": {"v3"}},
			map[string]bool{"v1": true, "v3": true},
			false,
		},
		{
			"duplicate validators",
			map[string][]string{"A": {"v1", "v1"}, "B": {"v2", "v2"}},
			map[string]bool{"v1": true, "v2": true},
			true,
		},
		{
			"empty sets",
			map[string][]string{"A": []string{}, "B": []string{}},
			map[string]bool{},
			false,
		},
		{
			"extra then ignored",
			map[string][]string{"A": {"v1"}, "B": {"v2"}},
			map[string]bool{"v1": true, "v2": true},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := m.Validate(script, tc.votes, tc.sets)
			if res != tc.pass {
				t.Fatalf("expect %v got %v", tc.pass, res)
			}
		})
	}
}
