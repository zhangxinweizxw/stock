package lib

import (
    "fmt"
	"strconv"
	"strings"
)

var (

	// 加密映射表 (按位映射)
	encryptTable = map[string][]string{
		"0": {"4b", "3d", "c5", "1e", "a7", "68", "be", "13", "b0", "49", "76", "83", "6d", "2a", "07", "0a", "e2", "0c", "2d"},
		"1": {"ee", "03", "e7", "37", "18", "d0", "40", "cd", "01", "da", "ec", "5a", "71", "d3", "ba", "5b", "89", "bd", "44"},
		"2": {"d6", "82", "48", "b3", "1c", "70", "19", "98", "69", "1d", "43", "80", "0d", "6a", "17", "14", "9d", "a5", "8c"},
		"3": {"7b", "57", "77", "de", "61", "3a", "2c", "08", "e6", "a4", "58", "02", "81", "0e", "db", "b6", "28", "47", "d7"},
		"4": {"20", "42", "38", "cc", "35", "ea", "25", "10", "2e", "12", "92", "36", "84", "b7", "53", "c9", "85", "4e", "d4"},
		"5": {"a6", "65", "9e", "aa", "64", "d1", "51", "bc", "52", "4c", "4a", "d9", "05", "09", "95", "4d", "e0", "63", "c0"},
		"6": {"6e", "8e", "60", "5d", "06", "3e", "ad", "46", "d2", "5c", "86", "e3", "d8", "99", "66", "ac", "c3", "24", "a0"},
		"7": {"5e", "7d", "7e", "b5", "67", "9a", "bb", "8d", "41", "6c", "e5", "a3", "72", "c2", "1b", "11", "54", "29", "62"},
		"8": {"0b", "45", "ae", "ce", "8b", "b4", "3b", "c7", "56", "7c", "78", "e8", "b8", "a1", "a9", "d5", "9c", "27", "b2"},
		"9": {"23", "cb", "50", "c1", "c6", "8a", "dd", "55", "30", "dc", "33", "59", "75", "16", "91", "ca", "22", "e9", "15"},
	}

	// 解密映射表
	decryptTable = map[string]string{}

	// 随机填充字符串 (按照最后一位的数字映射对应的)
	randTable = map[int64]string{
		0: "eaa18519ee0848558a3e9025ea5de341",
		1: "a67512311c5643a58de05457647be46c",
		2: "491d03ae833464a080de0cbad36ab607",
		3: "b47812314a274abaa8740a819b2e7737",
		4: "81135e21bbdd474c96e9783ada87e434",
		5: "17a39c1e2dee4cbd816b135d630c6465",
		6: "6ba4a1bcc21549bd82248a83e959ea8d",
		7: "edd78cd1389a4908acbcd47bebedb85b",
		8: "aa6c9e83e4244a69a9ceeb056a958632",
		9: "11e72208721e9871314e37e7c790c95a",
	}

	// 间隔记号
	iDsign = "f5"
)

func init() {

	// 根据加密映射表生成对应的解密映射表
	for i, arr := range encryptTable {
		for _, v := range arr {
			decryptTable[v] = i
		}
	}
}

// ID串加密
func IDEncrypt(i int64) string {
	if i <= 0 {
		return ""
	}

	idString := fmt.Sprintf("%v", i)
	elements := []string{randTable[i%10], iDsign}
	for i, _ := range idString {
		elements = append(elements, encryptTable[string(idString[i])][i])
	}

	result := strings.Join(elements, "")
	length := len(result)

	idLength := len(idString)
	if idLength <= 11 {
		return result[length-24 : length]
	}

	if idLength > 11 && idLength < 16 {
		return result[length-32 : length]
	}

	return result[length-40 : length]
}

// ID串解密
func IDDecrypt(s string) int64 {
	if len(s) == 0 {
		return 0
	}

	values := strings.Split(s, iDsign)
	if len(values) == 1 {
		return 0
	}

	v := values[1]
	idLength := len(v)

	var elements []string
	for i := 0; i < idLength; i++ {
		if i >= idLength-1 {
			return 0
		}
		ch := string([]byte{v[i], v[i+1]})
		elements = append(elements, decryptTable[ch])
		i = i + 1
	}

	result, _ := strconv.ParseInt(strings.Join(elements, ""), 10, 64)
	return result
}
