package utils

import (
	"io"
)

func SplitRead(stream io.Reader, expectLen int, packet []byte) (int, error) {
	count := 0
	splitSize := 99
	for count < expectLen {
		receiveSize := splitSize
		if expectLen-count < splitSize {
			receiveSize = expectLen - count
		}
		n, err := stream.Read(packet[count : count+receiveSize])
		if err != nil {
			return count, err
		}
		count += n
	}
	return count, nil
}

// ReadLength []byte length to int length
func ReadLength(header []byte) int {
	length := 0
	if len(header) >= 2 {
		length = ((length & 0x00) | int(header[0])) << 8
		length = length | int(header[1])
	}
	return length
}

func WriteLength(header []byte, length int) {
	if len(header) >= 2 {
		header[0] = byte(length >> 8 & 0xff)
		header[1] = byte(length & 0xff)
	}
}

func Merge(a, b []byte) []byte {
	al := len(a)
	bl := len(b)
	c := make([]byte, len(a)+len(b))
	copy(c[al:al+bl], b)
	copy(c[:al], a)
	return c
}
