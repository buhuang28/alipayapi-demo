package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//主程序的GetSign
func GetSign(params ...string) string {
	content := ""
	for _, v := range params {
		content += v
	}
	return Md5(content)
}
