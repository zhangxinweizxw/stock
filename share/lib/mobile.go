package lib

import (
    "strings"

	"stock/share/lib/crypto"
	"stock/share/logging"
)

var AesKey []byte = []byte("br3pd91abr3pd91a")
var ManSuffix = "A#$d&(-+}%"

func MobileEncrypt(number string) (string, string, error) {
	maq := number[0:3] + "*" + number[7:]
	scrap := number[3:7] + ManSuffix

	if str, err := crypto.AesEncrypt([]byte(scrap), AesKey); err != nil {
		return "", "", err
	} else {
		base64 := string(crypto.EncodeBase64(str))
		return maq, base64, nil
	}
}

func MobileDecrypt(number, man string) (string, error) {
	var err error
	var base64, data []byte

	if base64, err = crypto.DecodeBase64([]byte(man)); err != nil {
		logging.Error("%v", err)
		return "", err
	}

	if data, err = crypto.AesDecrypt(base64, AesKey); err != nil {
		logging.Error("%v", err)
		return "", err
	}

	str := strings.Replace(string(data), ManSuffix, "", -1)

	return strings.Replace(number, "*", str, -1), nil
}

func MobileSinaEncrypt(number string) (string, error) {
	key := []byte("9f24b80a98680a3599d1932b79e3ef8a")
	if str, err := crypto.AesEncrypt([]byte(number), key); err != nil {
		return "", err
	} else {
		return string(crypto.EncodeBase64(str)), nil
	}
}
