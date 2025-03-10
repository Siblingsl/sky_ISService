package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptAES AES 加密，支持指定密钥长度（128/192/256位）
func EncryptAES(plainText, key string) (string, error) {
	// 密钥长度为 16 字节（128位）、24 字节（192位）、32 字节（256位）
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("无法创建 AES 加密实例: %v", err)
	}

	// 填充明文至 AES 块大小
	plainBytes := []byte(plainText)
	blockSize := block.BlockSize()
	plainBytes = pkcs7Padding(plainBytes, blockSize)

	// 创建随机初始化向量（IV）
	cipherText := make([]byte, blockSize+len(plainBytes))
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("无法生成随机初始化向量: %v", err)
	}

	// 使用 CBC 模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[blockSize:], plainBytes)

	// 返回加密后的文本，采用 Base64 编码
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptAES AES 解密
func DecryptAES(cipherTextBase64, key string) (string, error) {
	keyBytes := []byte(key)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("无法创建 AES 加密实例: %v", err)
	}

	// 解码 Base64 密文
	cipherText, err := base64.StdEncoding.DecodeString(cipherTextBase64)
	if err != nil {
		return "", fmt.Errorf("无法解码 Base64 密文: %v", err)
	}

	// 检查 AES 块大小
	blockSize := block.BlockSize()
	if len(cipherText) < blockSize {
		return "", fmt.Errorf("密文长度不足")
	}

	// 获取 IV 和加密数据
	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]

	// 使用 CBC 模式解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	// 去除填充
	plainText := pkcs7UnPadding(cipherText)
	return string(plainText), nil
}

// PKCS7 填充（加密时使用）
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, paddingText...)
}

// 去除填充（解密时使用）
func pkcs7UnPadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
