// 参考 https://github.com/damonlear/CRC16-CCITT/blob/master/src/main/java/com/damonlear/crc16/CRC16.java
package protocol

import "fmt"

// CRC16-CCITTT校验
func CRC16CCITT(data []byte) (uint16, string) {
	// initial value
	crc := 0x0000
	// poly value reversed 0x1021;
	polynomial := 0x8408

	for i := 0; i < len(data); i++ {
		crc ^= int(data[i]) & 0x000000ff
		for j := 0; j < 8; j++ {
			if (crc & 0x00000001) != 0 {
				crc >>= 1
				crc ^= polynomial
			} else {
				crc >>= 1
			}
		}
	}
	return uint16(crc), fmt.Sprintf("%x", crc)
}
