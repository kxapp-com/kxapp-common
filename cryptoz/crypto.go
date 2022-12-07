package cryptoz

import (
	"crypto/rc4"
	"encoding/base64"
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
