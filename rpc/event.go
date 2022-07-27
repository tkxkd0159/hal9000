package rpc

func CheckEvt(attrs []string) (any, bool) {
	if len(attrs) == 1 {
		return attrs[0], true
	}
	return attrs, false
}
