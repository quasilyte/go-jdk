package symbol

// ID is a compact and globally unique symbol identifier.
//
// ID layout:
//	- 3 bytes for package index (within VM)
//	- 3 bytes for class index (within package)
//	- 2 bytes for member index (within class)
//
// Limits:
//	- 16_777_215 packages
//	- 16_777_215 classes within package
//	- 65_535 class members
type ID uint64

// NewID constructs a symbol ID out of its components.
func NewID(pkg, class, member uint64) ID {
	return ID((pkg << (8 * 5)) | (class << (8 * 2)) | (member << (8 * 0)))
}

func (id ID) PackageIndex() uint { return uint(id >> (8 * 5)) }

func (id ID) ClassIndex() uint { return uint((id << (8 * 3)) >> (8 * 5)) }

func (id ID) MemberIndex() uint { return uint(uint16(id)) }
