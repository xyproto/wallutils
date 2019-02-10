package crossfade

// Converts a float64 to an uint8 checking clamping the values
// to be compatible with different architectures.
// http://code.google.com/p/go/issues/detail?id=3423
func float64ToUint8(x float64) uint8 {
	if x < 0 {
		return 0
	}
	if x > 255 {
		return 255
	}
	return uint8(int(x + 0.5))
}

// Converts a float64 to an uint8 checking clamping the values
// to be compatible with different architectures.
// http://code.google.com/p/go/issues/detail?id=3423
func float64ToUint16(x float64) uint16 {
	if x < 0 {
		return 0
	}
	if x > 65535 {
		return 65535
	}
	return uint16(int(x + 0.5))
}
