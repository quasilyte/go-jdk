package x64

func fitsInt8(v int64) bool {
	return v >= -128 && v <= 127
}

func modrm(mod, reg, rm byte) byte {
	return (mod << 6) | (reg << 3) | (rm << 0)
}

func appendInt32(buf []byte, v int32) []byte {
	return append(buf,
		byte(v>>0),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24))
}

func appendInt64(buf []byte, v int64) []byte {
	return append(buf,
		byte(v>>0),
		byte(v>>8),
		byte(v>>16),
		byte(v>>24),
		byte(v>>32),
		byte(v>>40),
		byte(v>>48),
		byte(v>>56))
}
