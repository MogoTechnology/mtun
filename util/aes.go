package util

import (
	"bytes"
)

var DefaultKey = []byte("mogo2023mogo2023")

func Decrypt(decrypted []byte, key []byte) ([]byte, error) {
	return decrypted, nil
}

func Encrypt(text []byte, key []byte) ([]byte, error) {
	return text, nil
}

//func Decrypt(decrypted []byte, key []byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//
//	blockSize := block.BlockSize()
//	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
//	origData := make([]byte, len(decrypted))
//	blockMode.CryptBlocks(origData, decrypted)
//	origData = PKCS5UnPadding(origData)
//	return origData, nil
//}
//
//func Encrypt(text []byte, key []byte) ([]byte, error) {
//	// 创建一个新的 AES 加密块，使用给定的密钥
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//
//	blockSize := block.BlockSize()
//	origData := PKCS5Padding(text, blockSize)
//	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
//	crypted := make([]byte, len(origData))
//	blockMode.CryptBlocks(crypted, origData)
//	return crypted, nil
//}

func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
