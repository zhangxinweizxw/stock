package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5(src string, short bool) string {
	h := md5.New()
	h.Write([]byte(src))
	data := hex.EncodeToString(h.Sum(nil))
	if short {
		return data[8:24]
	}
	return data
}

func GetMD5Dual(src string, short bool) string {
	h := md5.New()
	h.Write([]byte("#s1;l5" + src + "@.k3"))
	data := hex.EncodeToString(h.Sum(nil))
	if short {
		return data[8:24]
	}
	return data
}

func GetMD5Suffix(src string, short bool, suffix string) string {
	h := md5.New()
	h.Write([]byte(src + suffix))
	data := hex.EncodeToString(h.Sum(nil))
	if short {
		return data[8:24]
	}
	return data
}
