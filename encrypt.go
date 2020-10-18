package wechat_work_bot

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"sort"
	"strings"
	"time"
)

func init() {
	mathRand.Seed(time.Now().Unix())
}

type EncryptStrings []string

func (s EncryptStrings) Len() int {
	return len(s)
}

func (s EncryptStrings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s EncryptStrings) Less(i, j int) bool {
	minLen := len(s[i])
	if minLen > len(s[j]) {
		minLen = len(s[j])
	}

	for idx := 0; idx < minLen; idx++ {
		if s[i][idx] > s[j][idx] {
			return false
		} else if s[i][idx] < s[j][idx] {
			return true
		}
	}

	return false
}

func Sha1(token, sEncrypt, nonce string, timestamp int64) string {
	sortList := []string{token, sEncrypt, fmt.Sprint(timestamp), nonce}
	sort.Sort(EncryptStrings(sortList))

	strTarget := strings.Join(sortList, "")

	h := sha1.New()
	h.Write([]byte(strTarget))

	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//使用PKCS7进行填充，IOS也是7
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//填充原文
	blockSize := block.BlockSize()
	rawData = PKCS7Padding(rawData, 32) // 企业微信用的加密是 block size 为32
	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, len(rawData))

	//block大小 16
	iv := cipherText[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, rawData)

	return cipherText, nil
}

func AesCBCDecrypt(encryptData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {

		return nil, fmt.Errorf("ciphertext too short")
	}
	iv := encryptData[:blockSize]

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS7UnPadding(encryptData)
	return encryptData, nil
}

func Encrypt(rawData, key []byte) (string, error) {
	data, err := AesCBCEncrypt(rawData, key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func Decrypt(rawData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return nil, err
	}
	dnData, err := AesCBCDecrypt(data, key)
	if err != nil {
		return nil, err
	}
	return dnData, nil
}

func getRandomStr(n int) string {
	ans := ""
	for i := 0; i < n; i++ {
		ans += fmt.Sprint(mathRand.Intn(9))
	}
	return ans
}
