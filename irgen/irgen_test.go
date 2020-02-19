package irgen

import "testing"

func TestArgCount(t *testing.T) {
	tests := []struct {
		d    string
		want int
	}{
		{"()V", 0},
		{"()I", 0},
		{"(I)I", 1},
		{"(II)I", 2},
		{"(LFoo;)I", 1},
		{"([[I)I", 1},
		{"([I[I)I", 2},
		{"(LFoo;LBar;)I", 2},
		{"(LObject;I)I", 2},
		{"(Z[LObject;Z)V", 3},
	}

	for _, test := range tests {
		have := argsCount(test.d)
		if have != test.want {
			t.Errorf("argsCount(%q):\nhave: %d\nwant: %d", test.d, have, test.want)
		}
	}
}
