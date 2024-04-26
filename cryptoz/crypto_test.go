package cryptoz

import (
	"fmt"
	"os"
	"testing"
)

func TestEncryptLong(t *testing.T) {
	privateKeyBytes, _ := os.ReadFile("private.pem")
	publicKeyBytes, _ := os.ReadFile("public.pem")

	// 假设需要加密的数据是字符串"Hello, World!"
	data := []byte("Hello, World!")
	cipherData, e := EncryptRSALongData(publicKeyBytes, data)
	textData, e2 := DecryptRSALongData(privateKeyBytes, cipherData)
	fmt.Println(e, e2, string(textData))
}
