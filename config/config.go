package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"stock/share/lib"
)

var config *AppConfig

type AppSettings struct {
	AllowOrigin string     `xml:"allowOrigin"`
	EncryFactor string     `xml:"encryFactor"`
	Environment string     `xml:"environment"`
	Listen      string     `xml:"listen"`
	Projects    []Projects `xml:"projects"`
}

type Projects struct {
	AppId      string `xml:"appId"`
	ConfigFile string `xml:"configFile"`
}

type AccessKeys struct {
	ID     string `xml:"id"`
	Secret string `xml:"secret"`
	AESKey string `xml:"aesKey"`
}

// Database
type Database struct {
	DriverName string `xml:"driverName"`
	DataSource string `xml:"dataSource"`
}

type SessionSetting struct {
	On           bool   `xml:"on"`
	ProviderName string `xml:"providerName"`
	Config       string `xml:"config"`
}
type MnsSetting struct {
	Url             string    `xml:"url"`
	AccessKeyId     string    `xml:"accessKeyId"`
	AccessKeySecret string    `xml:"accessKeySecret"`
	Queues          QueueName `xml:"queues"`
}
type QueueName struct {
	SmartCall string `xml:"exchangeSmartCall"`
}

type LogServer struct {
	On   bool   `xml:"on"`
	Addr string `xml:"addr"`
	Port string `xml:"port"`
}

type AppConfig struct {
	AccessKeys AccessKeys     `xml:"accessKeys"`
	Cors       CorsSetting    `xml:"cors"`
	Db         Database       `xml:"database"`
	Mns        MnsSetting     `xml:"mns"`
	Serve      ListenAndServe `xml:"listenAndServe"`
	Session    SessionSetting `xml:"session"`
	Settings   AppSettings    `xml:"appSettings"`
	Log        LogServer      `xml:"logServer"`
	Url        Urls           `xml:"url"`
}

type CorsSetting struct {
	AllowOrigin []string `xml:"allowOrigin"`
}

type ListenAndServe struct {
	Port    string `xml:"port"`
	LogPort string `xml:"logport"`
}

// ----------------------------------------------------------------------
// config.xml 中配置的url
type Urls struct {
	SjsMonthList  string `xml:"sjsMonthList"`
	DfcfStockDayK string `xml:"dfcfStockDayK"`
}

// ----------------------------------------------------------------------

//// 传参使用
//type Context struct {
//	LoginToken string
//	Conf       *Configuration
//}
//
//// 接口返回结构
//type ReturnJson struct {
//	Code int32 `json:"code"`
//	//Msg    string      `json:"msg"`
//	Msg    string      `json:"message"`
//	Result interface{} `json:"result"`
//}
//
//// 动态加密码返回结构
//type GetSecretKey struct {
//	Code        int32  `json:"code"`
//	Message     string `json:"message"`
//	Description string `json:"description"`
//	Result      string `json:"result"`
//	Success     bool   `json:"success"`
//}

//type Configuration struct {
//	Urls    Url    `xml:"url"`
//	Users   User   `xml:"user"`
//	Webhook string `xml:"webhook"`
//}
//
//// 接口配置-----------------------------------------start
//type Url struct {
//	Api string `xml:"api"`
//	Add *Add   `xml:"add"`
//}
//type Add struct {
//	LoginToken                         string `xml:"loginToken"`
//	NewLoginToken                      string `xml:"newLoginToken"`
//	GetLayoutAndDataList               string `xml:"getLayoutAndDataList"`
//	SelectTenant                       string `xml:"selectTenant"`
//	OrganizationDetail                 string `xml:"organizationDetail"`
//	OrganizationDetailLatestInvestment string `xml:"organizationDetailLatestInvestment"`
//	GetDataList                        string `xml:"getDataList"`
//	PublisherReport                    string `xml:"publisherReport"`
//	RecommendRead                      string `xml:"recommendRead"`
//	LeftMenu                           string `xml:"leftMenu"`
//	NavigationList                     string `xml:"navigationList"`
//}
//
//// 接口配置-----------------------------------------end
//
//// 用户信息配置-----------------------------------------start
//type User struct {
//	UserId   string `xml:"userId"`
//	Password string `xml:"password"`
//	TenantId string `xml:"tenantId"`
//}

//// 用户信息配置-----------------------------------------end
//func GetConf() *Configuration {
//
//	xmlFile, err := os.Open("./conf.d/config.xml")
//	if err != nil {
//		fmt.Println("Error opening file:", err)
//		return nil
//	}
//
//	defer xmlFile.Close()
//	var conf Configuration
//	if err := xml.NewDecoder(xmlFile).Decode(&conf); err != nil {
//		fmt.Println("Error Decode file:", err)
//		return nil
//	}
//
//	//fmt.Print(conf.Urls.Add)
//	return &conf
//}
////-----------------------------------------------------------------------------
func Default(appID string) *AppConfig {
	if config == nil {
		var cfg AppConfig
		lib.LoadConfig(appID, &cfg)
		config = &cfg
	}
	return config
}

func Reload() {
	config = nil
}

// ------------------------------------------------------------------------

func getCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	return filepath.Dir(file)
}
