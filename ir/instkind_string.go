// Code generated by "stringer -type=InstKind -trimprefix=Inst"; DO NOT EDIT.

package ir

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[InstInvalid-0]
	_ = x[InstIload-1]
	_ = x[InstLload-2]
	_ = x[InstAload-3]
	_ = x[InstRet-4]
	_ = x[InstIret-5]
	_ = x[InstLret-6]
	_ = x[InstAret-7]
	_ = x[InstCallStatic-8]
	_ = x[InstCallGo-9]
	_ = x[InstIcmp-10]
	_ = x[InstLcmp-11]
	_ = x[InstJump-12]
	_ = x[InstJumpEqual-13]
	_ = x[InstJumpNotEqual-14]
	_ = x[InstJumpGtEq-15]
	_ = x[InstJumpGt-16]
	_ = x[InstJumpLt-17]
	_ = x[InstImul-18]
	_ = x[InstIdiv-19]
	_ = x[InstIadd-20]
	_ = x[InstLadd-21]
	_ = x[InstFadd-22]
	_ = x[InstIsub-23]
	_ = x[InstIneg-24]
	_ = x[InstLneg-25]
	_ = x[InstDadd-26]
	_ = x[InstConvL2I-27]
	_ = x[InstConvF2I-28]
	_ = x[InstConvD2I-29]
	_ = x[InstConvI2L-30]
	_ = x[InstConvI2B-31]
	_ = x[InstNewBoolArray-32]
	_ = x[InstNewCharArray-33]
	_ = x[InstNewFloatArray-34]
	_ = x[InstNewDoubleArray-35]
	_ = x[InstNewByteArray-36]
	_ = x[InstNewShortArray-37]
	_ = x[InstNewIntArray-38]
	_ = x[InstNewLongArray-39]
	_ = x[InstIntArraySet-40]
	_ = x[InstIntArrayGet-41]
	_ = x[InstArrayLen-42]
}

const _InstKind_name = "InvalidIloadLloadAloadRetIretLretAretCallStaticCallGoIcmpLcmpJumpJumpEqualJumpNotEqualJumpGtEqJumpGtJumpLtImulIdivIaddLaddFaddIsubInegLnegDaddConvL2IConvF2IConvD2IConvI2LConvI2BNewBoolArrayNewCharArrayNewFloatArrayNewDoubleArrayNewByteArrayNewShortArrayNewIntArrayNewLongArrayIntArraySetIntArrayGetArrayLen"

var _InstKind_index = [...]uint16{0, 7, 12, 17, 22, 25, 29, 33, 37, 47, 53, 57, 61, 65, 74, 86, 94, 100, 106, 110, 114, 118, 122, 126, 130, 134, 138, 142, 149, 156, 163, 170, 177, 189, 201, 214, 228, 240, 253, 264, 276, 287, 298, 306}

func (i InstKind) String() string {
	if i < 0 || i >= InstKind(len(_InstKind_index)-1) {
		return "InstKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _InstKind_name[_InstKind_index[i]:_InstKind_index[i+1]]
}
