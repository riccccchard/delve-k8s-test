package utils

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"io"
	mrand "math/rand"
	"os"
	"time"
)

func randBytes(x []byte) {
	var seeded bool = false
	length := len(x)
	n, err := rand.Read(x)

	if n != length || err != nil {
		if !seeded {
			mrand.Seed(time.Now().UnixNano())
		}

		for length > 0 {
			length--
			x[length] = byte(mrand.Int31n(256))
		}
	}
}

func UUID() string {
	var x [16]byte
	randBytes(x[:])
	x[6] = (x[6] & 0x0F) | 0x40
	x[8] = (x[8] & 0x3F) | 0x80
	return fmt.Sprintf("%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		x[0], x[1], x[2], x[3], x[4],
		x[5], x[6],
		x[7], x[8],
		x[9], x[10], x[11], x[12], x[13], x[14], x[15])
}

// 计算string的md5值，以32位字符串形式返回
func GetMd5FromString(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}

// 计算[]byte的md5值，以32位字符串形式返回
func GetMd5FromBytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// 计算文件的md5值，以32位字符串形式返回
func GetMd5FromFile(filename string) (string, error) {
	h := md5.New()
	if err := readFile(filename, h); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 辅助函数，扫描文件内容并编码到hash.Hash中
func readFile(filename string, h hash.Hash) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(file)
	for s.Scan() {
		h.Write(s.Bytes())
	}

	return s.Err()
}

// AES是对称加密算法
// Key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
// 下面的AES使用CBC模式和PKCS5Padding填充法
// AES加密，传入plaintext和key，返回ciphertext（plaintext不改变）
func AesEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	padding := aes.BlockSize - len(plaintext)%aes.BlockSize
	ciphertext := make([]byte, aes.BlockSize+len(plaintext)+padding)

	// 初始化向量IV放在ciphertext前面
	iv := ciphertext[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)

	copy(ciphertext[aes.BlockSize:], plaintext)

	// PKCS5Padding填充
	for i := 0; i < padding; i++ {
		ciphertext[aes.BlockSize+len(plaintext)+i] = byte(padding)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return ciphertext, nil
}

// AES解密，传入ciphertext和key，返回plaintext（ciphertext不改变）
func AesDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) <= aes.BlockSize || len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("AesDecrypt error: len(ciphertext) % aes.BlockSize != 0")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	unpadding := int(plaintext[len(plaintext)-1])
	if unpadding <= 0 || unpadding > 16 {
		return nil, errors.New("AesDecrypt error: unpadding <= 0 || unpadding > 16")
	}

	plaintext = plaintext[:len(plaintext)-unpadding]

	return plaintext, nil
}

// RSA加密，传入plaintext和publickey，返回ciphertext（plaintext不改变）
func RsaEncrypt(plaintext, publickey []byte) ([]byte, error) {
	block, _ := pem.Decode(publickey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)

	return rsa.EncryptPKCS1v15(rand.Reader, pub, plaintext)
}

// RSA解密，传入ciphertext和privatekey，返回plaintext（ciphertext不改变）
func RsaDecrypt(ciphertext, privatekey []byte) ([]byte, error) {
	block, _ := pem.Decode(privatekey)
	if block == nil {
		return nil, errors.New("private key error")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
