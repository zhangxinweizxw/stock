package lib

import (
    "errors"
    "fmt"
    "net"
    "net/url"
    "os"
    "path/filepath"
    "reflect"
    "regexp"
    "strconv"
    "strings"

    "stock/share/logging"
)

// environment model
// --------------------------------------------------------------------------------

const (
	ENVIRONMENT_MODEL_DEVELOPMENT = "development" // 开发模式
	ENVIRONMENT_MODEL_PRE         = "pre-release" // 预发布模式
	ENVIRONMENT_MODEL_PRODUCTION  = "production"  // 线上模式
	ENVIRONMENT_MODEL_TESTING     = "testing"     // 测试模式
)

func Bool2String(b bool) string {
	if b == true {
		return "true"
	}
	return "false"
}

func GetArrayWithCSV(str [][]string) map[string]string {
	if str == nil || len(str) <= 1 {
		return nil
	}
	result := make(map[string]string)
	columns := strings.Split(str[0][0], "|")
	for i := 1; i < len(str); i++ {
		items := strings.Split(str[i][0], "|")
		if len(items) == len(columns) {
			for n := 0; n < len(items); n++ {
				result[columns[n]] = items[n]
			}
		}
	}
	return result
}

func GetLocalIPAddr() (string, error) {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	return strings.Split(conn.LocalAddr().String(), ":")[0], nil
}

func FiltrationHTML(str string) string {
	if len(str) == 0 {
		return str
	}
	re, _ := regexp.Compile(`style=\\"(.[^"]*)\\"`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile(`style='(.[^"]*)'`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile(`class=\\"(.[^"]*)\\"`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile(`class='(.[^"]*)'`)
	str = re.ReplaceAllString(str, "")

	re, _ = regexp.Compile("\\<br[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "\n")

	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	str = re.ReplaceAllString(str, "")

	return str
}

func TrimIPAddr(addr string) string {
	return strings.Split(addr, ":")[0]
}

// ------------------------------------------------------------------------

func IsEmail(email string) bool {
	if len(email) == 0 {
		return false
	} else {
		reg := regexp.MustCompile(`(?i)[A-Z0-9._%+-]+@(?:[A-Z0-9-]+\.)+[A-Z]{2,6}`)
		return reg.MatchString(email)
	}
}

func IsMobile(number string) bool {
	if len(number) == 0 {
		return false
	} else {
		reg := regexp.MustCompile(`0?1[34578]{1}[0-9]{9}`)
		return reg.MatchString(number)
	}
}

func IsTelephone(number string) bool {
	if len(number) == 0 {
		return false
	} else {
		reg := regexp.MustCompile(`(0\d{2})?\d{8}|(0\d{3})?\d{8}|(0\d{3})?\d{7}`)
		return reg.MatchString(number)
	}
}

// ------------------------------------------------------------------------

func getCurrPath() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

func GetPath(name string) string {
	return getCurrPath() + "/" + name
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

func IsFileExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

// 检查目录是否存在，否则并自动创建
func CheckDir(path string) bool {
	if IsDirExists(path) == false {
		var npath string
		var i int

		paths := strings.Split(path, "/")
		for i = 0; i < len(paths)-2; i++ {
			npath += paths[i] + "/"
		}
		npath += paths[i]
		CheckDir(npath)
		os.Mkdir(path, 0777)
	}
	return true
}

// 检查颜色值是否正确
func CheckColor(color string) bool {
	if len(color) != 7 || color[:1] != "#" {
		return false
	}
	return true
}

func DeletePath(path string) error {
	if IsDirExists(path) == false {
		return nil
	}
	return os.RemoveAll(path)
}

func UnsignedParams(key, timestamp, version string, queryParams url.Values) url.Values {
	params := url.Values{
		"auth_key":       {key},
		"auth_timestamp": {timestamp},
		"auth_version":   {version},
	}

	if queryParams != nil {
		for k, v := range queryParams {
			params[k] = v
		}
	}
	return params
}

func UnescapeUrl(_url url.Values) string {
	unesc, _ := url.QueryUnescape(_url.Encode())
	return unesc
}

func ValidateBinary(limits int, action int) bool {
	binaryList := []int{1, 2, 4, 8, 16, 32}
	if limits == -1 {
		return true
	}
	if action < 1 || action > len(binaryList)-1 {
		return false
	}
	if action == len(binaryList)-1 {
		return limits == binaryList[action]
	}
	return (limits & binaryList[action-1]) >= binaryList[action-1]
}

func InArray(array []string, str string) bool {
	for _, v := range array {
		if strings.EqualFold(v, str) {
			return true
		}
	}
	return false
}

//func Pinyin(str string) string {
//	a := pinyin.NewArgs()
//	a.Separator = ""
//	a.Style = pinyin.FirstLetter
//	a.Fallback = func(r rune, a pinyin.Args) []string {
//		return []string{string([]rune{r})}
//	}
//	first := pinyin.Slug(str, a)
//
//	a.Style = pinyin.Normal
//	secondly := pinyin.Slug(str, a)
//	result := ""
//
//	if len(first) > 0 {
//		result = first + "," + secondly
//	}
//
//	return strings.Replace(result, " ", "", -1)
//}

func InetNtoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func InetAton(ipnr string) int64 {
	bits := strings.Split(strings.Split(ipnr, ":")[0], ".")
	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

//----------------------------config-------------------------------

func LoadGlobalConfig(appConfig interface{}) error {
	globalFile := "conf.d/config.xml"

	var adapter XmlConfig
	if !isExist(globalFile) {

		return errors.New("Configuration File Not Found.")
	}
	adapter.Prase(globalFile, appConfig)
	return nil
}

// Config files by default on : conf/config.xml
func ParseConfig(appConfig interface{}, mode string, configPath string) error {
	var adapter XmlConfig
	var fileName string

	preFile := "conf.d/" + configPath[:strings.LastIndex(configPath, ".")]

	switch mode {
	case ENVIRONMENT_MODEL_DEVELOPMENT:
		fileName = fmt.Sprintf("%v_dev.xml", preFile)
	case ENVIRONMENT_MODEL_TESTING:
		fileName = fmt.Sprintf("%v_test.xml", preFile)
	case ENVIRONMENT_MODEL_PRE:
		fileName = fmt.Sprintf("%v_pre.xml", preFile)
	case ENVIRONMENT_MODEL_PRODUCTION:
		fileName = fmt.Sprintf("%v_pro.xml", preFile)
	default:
		fileName = fmt.Sprintf("%v_def.xml", preFile)
	}

	if isExist(fileName) {
		return adapter.Prase(fileName, appConfig)
	}
	return errors.New(fmt.Sprintf("Configuration File %v Not Found.", fileName))
}

func LoadConfig(appName string, appConfig interface{}) interface{} {
	var configFile string
	LoadGlobalConfig(appConfig)
	app := reflect.ValueOf(appConfig).Elem()
	val := app.FieldByName("Settings").FieldByName("Projects")
	mode := app.FieldByName("Settings").FieldByName("Environment").String()
	for i := 0; i < val.Len(); i++ {
		if (val.Index(i).FieldByName("AppId").String()) == appName {
			configFile = val.Index(i).FieldByName("ConfigFile").String()
			break
		}
	}
	if configFile == "" {
		logging.Fatal("Not fonud config file")
	}
	if err := ParseConfig(&appConfig, mode, configFile); err != nil {
		logging.Fatal(err)
		os.Exit(1)
	}
	return appConfig
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//----------------------------url param-------------------------------
const (
	MAXIMUM_LIMIT = 100
	DEFAULT_LIMIT = 5
)

type Queryer interface {
	Query(param string) string
}

func GetLimit(q Queryer, limitParam string) int {
	limit, _ := strconv.Atoi(q.Query(limitParam))

	if limit < DEFAULT_LIMIT || limit > MAXIMUM_LIMIT {
		return DEFAULT_LIMIT
	}

	return limit
}
