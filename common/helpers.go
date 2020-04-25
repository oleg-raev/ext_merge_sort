package common

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func UnlinkSlice(src []byte) []byte {
	res := make([]byte, len(src))
	copy(res, src)
	return res
}
