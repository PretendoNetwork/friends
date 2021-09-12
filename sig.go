package main

/*
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type nexToken struct {
	SystemType uint8
	TokenType  uint8
	UserPID    uint32
	TitleID    uint64
	CreatTime  uint64
}

var rsaPrivateKeyBytes []byte
var rsaPrivateKey *rsa.PrivateKey
var hmacSecret []byte

func parseRsaPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("Failed to parse RSA key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func verifySignature(body []byte, expectedSignature []byte, secret []byte) error {

	mac := hmac.New(sha1.New, secret)
	mac.Write(body)

	calculatedSignature := mac.Sum(nil)

	if !bytes.Equal(expectedSignature, calculatedSignature) {
		return errors.New("[ERROR] Calculated signature did not match")
	}

	return nil
}

func init() {
	var err error

	rsaPrivateKeyBytes, err = ioutil.ReadFile("private.pem")
	if err != nil {
		panic(err)
	}

	rsaPrivateKey, err = parseRsaPrivateKey(rsaPrivateKeyBytes)
	if err != nil {
		panic(err)
	}

	hmacSecret, err = ioutil.ReadFile("secret.key")
	if err != nil {
		panic(err)
	}
}

func main() {
	tokenBase64 := "Smbh6g2tL93Ro6RAvCBCGtLzGw.40vhPIjGNyZh9LkineqBFQ-PJGVXZ8Xg.hUZk92KTtrisA0.OnzGi92K-tMud5piu08XWa7hwSaxXNGjRJ.o6HWwi1xapFndyhWOHs5kXtmUiFmyc1AOlcfo8nvSHsPTsEDX2Qc9BeooBJrUeMxDc2LDch0HyiMKGdPnpBg2-Lpekma6fNvCg0ylTqUMqonI.smHCXikr4IhWvvwy9QZDJLM*"
	tokenBase64 = strings.Replace(tokenBase64, ".", "+", -1)
	tokenBase64 = strings.Replace(tokenBase64, "-", "/", -1)
	tokenBase64 = strings.Replace(tokenBase64, "*", "=", -1)

	encryptedToken, _ := base64.StdEncoding.DecodeString(tokenBase64)

	// Split the encoded token into it's parts
	cryptoConfig := encryptedToken[:0x82]
	signature := encryptedToken[0x82:0x96]
	encryptedBody := encryptedToken[0x96:]

	// Parse crypto config into parts
	encryptedAESKey := cryptoConfig[:128]
	point1 := cryptoConfig[128]
	point2 := cryptoConfig[129]

	// Rebuild the IV
	iv := make([]byte, 0)
	iv = append(iv, encryptedAESKey[point1:point1+8]...)
	iv = append(iv, encryptedAESKey[point2:point2+8]...)

	// Decrypt the AES key
	decryptedAESKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedAESKey, nil)
	if err != nil {
		panic(err)
	}

	// Decrypt the token body
	block, err := aes.NewCipher(decryptedAESKey)
	if err != nil {
		panic(err)
	}

	decryptedBody := make([]byte, len(encryptedBody))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptedBody, encryptedBody)

	decryptedBody = decryptedBody[:0x16] // Remove AES padding

	// Verify the token body
	err = verifySignature(decryptedBody, signature, hmacSecret)
	if err != nil {
		panic(err)
	}

	// Unpack the token body to struct
	token := &nexToken{}
	tokenReader := bytes.NewBuffer(decryptedBody)

	err = binary.Read(tokenReader, binary.LittleEndian, token)
	if err != nil {
		panic(err)
	}

	fmt.Println(token)
}
*/
