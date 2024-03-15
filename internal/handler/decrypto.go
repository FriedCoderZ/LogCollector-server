package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// RSA解密函数
// 处理函数
func DecryptedAESKeyHandler(c *gin.Context) ([]byte, error) {

	// 获取经过Base64编码的AES密钥
	encryptedAESKey := c.PostForm("encrypted_aes_key")
	if encryptedAESKey == "" {
		return nil, fmt.Errorf("please provide the encrypted AES key")
	}
	// 解码Base64编码的AES密钥
	encryptedKey, err := base64.StdEncoding.DecodeString(encryptedAESKey)
	if err != nil {
		return nil, nil
	}

	// 读取RSA私钥文件
	privateKeyFile := "privateKey.pem"
	privateKeyData, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, nil
	}

	// 解析RSA私钥
	block, _ := pem.Decode(privateKeyData)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode RSA private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil
	}

	// 解密AES密钥
	decryptedKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedKey)
	if err != nil {
		return nil, nil
	}
	return decryptedKey, nil
}
