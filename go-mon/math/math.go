package math

func AbsInt32(val int32) int32 {
	if val > 0 {
		return val
	}
	return val * -1
}

func AbsInt64(val int64) int64 {
	if val > 0 {
		return val
	}
	return val * -1
}
