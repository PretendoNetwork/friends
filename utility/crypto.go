package utility

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	"github.com/CloudnetworkTeam/friends/globals"
	"github.com/CloudnetworkTeam/friends/types"
)

func DecryptToken(encryptedToken []byte) (*types.NEXToken, error) {
	// Decrypt the token body
	block, err := aes.NewCipher(globals.AESKey)
	if err != nil {
		return nil, err
	}

	expectedChecksum := binary.BigEndian.Uint32(encryptedToken[0:4])
	encryptedBody := encryptedToken[4:]

	decrypted := make([]byte, len(encryptedBody))
	iv := make([]byte, 16)
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encryptedBody)

	paddingSize := int(decrypted[len(decrypted)-1])

	if paddingSize < 0 || paddingSize >= len(decrypted) {
		return nil, fmt.Errorf("Invalid padding size %d for token %x", paddingSize, encryptedToken)
	}

	decrypted = decrypted[:len(decrypted)-paddingSize]

	table := crc32.MakeTable(crc32.IEEE)
	calculatedChecksum := crc32.Checksum(decrypted, table)

	if expectedChecksum != calculatedChecksum {
		return nil, errors.New("Checksum did not match. Failed decrypt. Are you using the right key?")
	}

	// Unpack the token body to struct
	token := &types.NEXToken{}
	tokenReader := bytes.NewBuffer(decrypted)

	err = binary.Read(tokenReader, binary.LittleEndian, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}
