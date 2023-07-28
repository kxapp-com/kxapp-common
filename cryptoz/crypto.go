package cryptoz

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
)

func RC4Crypto(data []byte, key string) []byte {
	cipher, _ := rc4.NewCipher([]byte(key))
	dst := make([]byte, len(data))
	cipher.XORKeyStream(dst, data)
	return dst
}
func EncryptAndEncode(data []byte, key string) string {
	return base64.StdEncoding.EncodeToString(RC4Crypto(data, key))
}
func DecodeAndDecrypt(basedString string, key string) ([]byte, error) {
	data, error := base64.StdEncoding.DecodeString(basedString)
	if error != nil {
		return nil, error
	}
	return RC4Crypto(data, key), nil
}

// GetStringMD5
//
//	 @Description: 获得字符串的MD5加密
//	 @param str
//	 @return string 16进制
//		fmt.Sprintf("%x", md5.Sum(str))
func MD5OfBytes(str []byte) string {
	return fmt.Sprintf("%x", md5.Sum(str))
}

func MD5OfFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err2 := io.Copy(md5hash, f); err2 != nil {
		return "", err2
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}

func NewRSAKeyPairToFile(privateKeyPath, publicKeyPath string) error {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("私钥生成失败: %v", err)
	}

	publicKey := &privateKey.PublicKey

	// 保存私钥到文件
	privateKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return fmt.Errorf("私钥文件创建失败: %v", err)
	}
	defer privateKeyFile.Close()
	err = pem.Encode(privateKeyFile, privateKeyPem)
	if err != nil {
		return fmt.Errorf("私钥文件编码失败: %v", err)
	}

	// 保存公钥到文件
	publicKeyDer, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("公钥序列化失败: %v", err)
	}
	publicKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDer,
	}
	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return fmt.Errorf("公钥文件创建失败: %v", err)
	}
	defer publicKeyFile.Close()
	err = pem.Encode(publicKeyFile, publicKeyPem)
	if err != nil {
		return fmt.Errorf("公钥文件编码失败: %v", err)
	}

	return nil
}

func EncryptRSA(publicKeyPath string, message string) ([]byte, error) {
	// 读取公钥文件
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("公钥文件读取失败: %v", err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyBytes)
	if publicKeyBlock == nil {
		return nil, fmt.Errorf("公钥文件解析失败")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("公钥解析失败: %v", err)
	}
	publicKey, ok := publicKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("无效的公钥")
	}

	// 加密消息
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message))
	if err != nil {
		return nil, fmt.Errorf("加密失败: %v", err)
	}

	return ciphertext, nil
}

func DecryptRSA(privateKeyPath string, ciphertext []byte) (string, error) {
	// 读取私钥文件
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("私钥文件读取失败: %v", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return "", fmt.Errorf("私钥文件解析失败")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("私钥解析失败: %v", err)
	}

	// 解密消息
	plaintext, err := rsa.DecryptPKCS1v15(nil, privateKey, ciphertext)
	if err != nil {
		return "", fmt.Errorf("解密失败: %v", err)
	}

	return string(plaintext), nil
}

func GenerateAESKey32() ([]byte, error) {
	// 生成32字节的随机密钥
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("生成AES密钥失败: %v", err)
	}
	return key, nil
}

func EncryptAES(plainData []byte, key []byte) ([]byte, error) {
	// 使用AES加密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES加密算法失败: %v", err)
	}

	// 使用CTR模式加密
	iv := make([]byte, aes.BlockSize)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, fmt.Errorf("生成随机IV失败: %v", err)
	}

	stream := cipher.NewCTR(block, iv)
	ciphertext := make([]byte, len(plainData))
	stream.XORKeyStream(ciphertext, plainData)

	// 拼接IV和密文
	result := append(iv, ciphertext...)

	return result, nil
}

func DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	// 分离IV和密文
	iv := ciphertext[:aes.BlockSize]
	encryptedData := ciphertext[aes.BlockSize:]

	// 使用AES解密
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES解密算法失败: %v", err)
	}

	// 使用CTR模式解密
	stream := cipher.NewCTR(block, iv)
	plainData := make([]byte, len(encryptedData))
	stream.XORKeyStream(plainData, encryptedData)

	return plainData, nil
}
