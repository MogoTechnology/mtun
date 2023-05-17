package util

import (
	"bytes"
	"crypto/des"
	"errors"
)

var DefaultKey = []byte("12345678")

func Decrypt(decrypted []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(decrypted))
	dst := out
	bs := block.BlockSize()
	if len(decrypted)%bs != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	for len(decrypted) > 0 {
		block.Decrypt(dst, decrypted[:bs])
		decrypted = decrypted[bs:]
		dst = dst[bs:]
	}
	out = PKCS7UnPadding(out)
	return out, nil
}

func Encrypt(text []byte, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	src := PKCS7Padding(text, bs)
	if len(src)%bs != 0 {
		return nil, errors.New("need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length-unpadding < 0 {
		return nil
	}
	return origData[:(length - unpadding)]
}
