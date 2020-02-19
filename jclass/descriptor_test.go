package jclass

import (
	"testing"
)

func TestMethodDescriptor(t *testing.T) {
	tests := []struct {
		sym    string
		raw    string
		pretty string
	}{
		{"", "()V", "void ()"},
		{"", "(C)C", "char (char)"},
		{"", "(J)S", "short (long)"},
		{"f", "()V", "void f()"},
		{"f", "(I)I", "int f(int)"},
		{"ff", "(B[[[B)F", "float ff(byte, byte[][][])"},
		{"f", "(IZ)Z", "boolean f(int, boolean)"},
		{"f", "(LFoo;)[[I", "int[][] f(Foo)"},
		{"", "(LA;LB;)LC;", "C (A, B)"},
		{"", "(LAa;LBb;ILCc;)LDd;", "Dd (Aa, Bb, int, Cc)"},
		{"", "(DLFoo;I[LBar;)V", "void (double, Foo, int, Bar[])"},
	}

	for _, test := range tests {
		d := MethodDescriptor(test.raw)
		have := d.SignatureString(test.sym)
		want := test.pretty
		if have != want {
			t.Errorf("%q result mismatch:\nhave: %s\nwant: %s", test.raw, have, want)
		}
	}
}
