package types

import(
	"math/rand"
	"crypto/md5"
	"fmt"
	"encoding/hex"
)

var (
	alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	alen = len(alphanum)
)

func RandomString(n int) string {
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(alen)]
	}
	return string(bytes)
}

func MD5(val string) string{
	h := md5.New()
	h.Write([]byte(val))
	return fmt.Sprintf("%s", hex.EncodeToString(h.Sum(nil)))
}