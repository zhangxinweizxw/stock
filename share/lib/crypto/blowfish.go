package crypto

import (
	"crypto/cipher"
	"golang.org/x/crypto/blowfish"
)

func BlowfishEncrypt(dst, key, vector []byte) ([]byte, error) {
	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	crypted := make([]byte, len(dst))
	mode := cipher.NewCBCEncrypter(block, vector)
	mode.CryptBlocks(crypted, dst)
	base64 := EncodeBase64(crypted)

	return base64, nil
}

func BlowfishDecrypt(dst, key, vector []byte) ([]byte, error) {
	base64, err := DecodeBase64(dst)
	if err != nil {
		return nil, err
	}

	block, err := blowfish.NewCipher(key)
	if err != nil {
		return nil, err
	}

	data := make([]byte, len(base64))
	mode := cipher.NewCBCDecrypter(block, vector)
	mode.CryptBlocks(data, base64)

	return data, nil
}
