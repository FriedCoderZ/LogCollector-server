package util

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// RSA解密函数
func DecryptRSA(ciphertext, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("decode private key error")
	}
	parseResult, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	prv := parseResult.(*rsa.PrivateKey)
	return rsa.DecryptPKCS1v15(rand.Reader, prv, ciphertext)
}

func EncryptAES(plaintext string, key []byte) (string, error) {
	// 创建cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, len(plaintext))
	// 加密
	c.Encrypt(ciphertext, []byte(plaintext))
	// Base64编码返回
	ciphertextBase64 := base64.StdEncoding.EncodeToString(ciphertext)
	return ciphertextBase64, nil
}

func DecryptAES(ciphertextBase64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))
	c.Decrypt(plaintext, ciphertext)

	plaintextStr := string(plaintext[:])
	return plaintextStr, nil
}
