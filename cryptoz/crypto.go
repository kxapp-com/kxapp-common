package cryptoz

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/base64"
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

/*
	func MD5OfFile(path string) (string, error) {
		fd, err := os.Open(path)
		if err != nil {
			return "", err
		}

		defer fd.Close()
		isBytes, err2 := ioutil.ReadAll(fd)
		return fmt.Sprintf("%x", md5.Sum(isBytes)), err2
	}
*/
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
