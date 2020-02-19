package symbol

import (
	"math/rand"
	"testing"
)

func TestID(t *testing.T) {
	for i := 0; i < 1000; i++ {
		pkg := uint64(rand.Intn(16777215))
		class := uint64(rand.Intn(16777215))
		member := uint64(rand.Intn(65535))
		id := NewID(pkg, class, member)
		if id.PackageIndex() != uint(pkg) {
			t.Fatalf("id(%x,%x,%x) package mismatch:\nhave: %d\nwant: %d\nbits: %x",
				pkg, class, member, id.PackageIndex(), pkg, id)
		}
		if id.ClassIndex() != uint(class) {
			t.Fatalf("id(%x,%x,%x) class mismatch:\nhave: %d\nwant: %d\nbits: %x",
				pkg, class, member, id.ClassIndex(), class, id)
		}
		if id.MemberIndex() != uint(member) {
			t.Fatalf("id(%x,%x,%x) member mismatch:\nhave: %d\nwant: %d\nbits: %x",
				pkg, class, member, id.MemberIndex(), member, id)
		}
	}
}
