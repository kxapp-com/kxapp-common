package utilz

import (
	"golang.org/x/exp/rand"
	"strings"
	"time"
)

func RandomStringCharset(length int, charset string) string {
	len1 := len(charset)
	//fmt.Println(len)
	random := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	var sb strings.Builder
	for i := 0; i < length; i++ {
		number := random.Intn(len1)
		sb.WriteByte(charset[number])
	}
	return sb.String()
}
func RandomString(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandomStringCharset(length, str)
}
