package utility

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

	"github.com/PretendoNetwork/friends-secure/globals"
	"github.com/PretendoNetwork/friends-secure/types"
)

func DecryptToken(encryptedToken []byte) (*types.NEXToken, error) {
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
	decryptedAESKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, globals.RSAPrivateKey, encryptedAESKey, nil)
	if err != nil {
		return nil, err
	}

	// Decrypt the token body
	block, err := aes.NewCipher(decryptedAESKey)
	if err != nil {
		return nil, err
	}

	decryptedBody := make([]byte, len(encryptedBody))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decryptedBody, encryptedBody)

	decryptedBody = decryptedBody[:0x17] // Remove AES padding

	// Verify the token body
	err = verifySignature(decryptedBody, signature, globals.HMACSecret)
	if err != nil {
		return nil, err
	}

	// Unpack the token body to struct
	token := &types.NEXToken{}
	tokenReader := bytes.NewBuffer(decryptedBody)

	err = binary.Read(tokenReader, binary.LittleEndian, token)
	if err != nil {
		return nil, err
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

func ParseRsaPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
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
