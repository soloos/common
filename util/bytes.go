package util

func ChangeBytesArraySize(bs *[]byte, count int) {
	if len(*bs) > count {
		*bs = (*bs)[:count]
		return
	}

	if len(*bs) == 0 {
		*bs = make([]byte, count)
		return
	}

	if len(*bs) < count {
		*bs = make([]byte, count)
		return
	}
}
