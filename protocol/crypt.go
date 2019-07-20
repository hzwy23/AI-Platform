package protocol

// 加密算法
func Encrypt(key uint32, data []byte) []byte {
	return crypt(key, data)
}

// 解密算法
func Decrypt(key uint32, data []byte) []byte {
	return crypt(key, data)
}

func crypt(key uint32, data []byte) []byte {
	var size = uint32(len(data))
	var idx uint32 = 0

	var buf = make([]byte, 0)

	if key == 0 {
		key = 1
	}

	for idx < size {
		key = KEY_IAI*(key%KEY_MI) + KEY_ICI
		buf = append(buf, data[idx]^byte((key>>20)&0xff))
		idx += 1
	}
	return buf
}
