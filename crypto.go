package main

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
	"encoding/binary"
	"encoding/pem"
	"errors"
)

func decryptToken(encryptedToken []byte) (*nexToken, error) {
	// Split the encoded token into it's parts
	cryptoConfig := encryptedToken[:0x90]
	signature := encryptedToken[0x90:0xA4]
	encryptedBody := encryptedToken[0xA4:]

	encryptedAESKey := cryptoConfig[:128]
	iv := cryptoConfig[128:]

	// Decrypt the AES key
	decryptedAESKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, rsaPrivateKey, encryptedAESKey, nil)
	if err != nil {
		return &nexToken{}, err
	}

	// Decrypt the token body
	block, err := aes.NewCipher(decryptedAESKey)
	if err != nil {
		return &nexToken{}, err
	}

	decryptedBody := make([]byte, len(encryptedBody))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptedBody, encryptedBody)

	decryptedBody = decryptedBody[:0x16] // Remove AES padding

	// Verify the token body
	err = verifySignature(decryptedBody, signature, hmacSecret)
	if err != nil {
		return &nexToken{}, err
	}

	// Unpack the token body to struct
	token := &nexToken{}
	tokenReader := bytes.NewBuffer(decryptedBody)

	err = binary.Read(tokenReader, binary.LittleEndian, token)
	if err != nil {
		return &nexToken{}, err
	}

	return token, nil
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
