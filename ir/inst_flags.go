package ir

// InstFlags is a set of instruction bit fields.
type InstFlags uint64

const (
	// flagBlockLead bit is set when instructions is a block leader.
	// A block leader is a first instruction in a basic block.
	flagBlockLead uint64 = 1 << iota

	// flagJumpTarget bit is set when instruction is a jump
	// target of some other instruction. First instruction
	// in a function body is also marked as such.
	// Mostly needed for alignment: if instruction is a
	// jump target, it can be beneficial to align it.
	flagJumpTarget
)

func (f InstFlags) IsBlockLead() bool    { return f.getFlag(flagBlockLead) }
func (f *InstFlags) SetBlockLead(v bool) { f.setFlag(v, flagBlockLead) }

func (f InstFlags) IsJumpTarget() bool    { return f.getFlag(flagJumpTarget) }
func (f *InstFlags) SetJumpTarget(v bool) { f.setFlag(v, flagJumpTarget) }

func (f *InstFlags) setFlag(v bool, mask uint64) {
	if v {
		*f |= InstFlags(mask)
	} else {
		*f &^= InstFlags(mask)
	}
}

func (f *InstFlags) getFlag(mask uint64) bool {
	return *f&InstFlags(mask) != 0
}
