package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//生成S3文件路径签名，key不传则采用默认值
func S3Sign(path, timeStamp string, key ...string) string {
	var salt string
	if len(key) == 0 {
		salt = "4bb355a09a551b7dddaef6291fdf6102"
	} else {
		salt = key[0]
	}
	h := md5.New()
	h.Write([]byte(path + timeStamp + salt))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}
