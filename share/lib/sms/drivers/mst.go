package drivers

//
//import (
//    "errors"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"net/url"
//	"strconv"
//
//	iconv "github.com/axgle/mahonia"
//)
//
//const (
//	NAME                   = "microlink"
//	PASSWORD               = "3d1415926"
//	WEB_SERVICE_SEND_URI   = "http://www.139000.com/send/gsend.asp"
//	WEB_SERVICE_GETFEE_URI = "http://www.139000.com/send/getfee.asp"
//	WEB_SERVICE_CHECK_URI  = "http://www.139000.com/send/checkcontent.asp"
//)
//
//var (
//	errorInfos = map[string]string{
//		"6001": "该团队用户无效",
//		"6002": "此用户账号已经被停用",
//		"6003": "此用户密码错误",
//		"6004": "目标手机号码在保护名单内",
//		"6005": "发送内容中含非法字符",
//		"6006": "发送通道不能对用户代收费",
//		"6007": "未找到合适通道给用户发短信",
//		"6008": "无效的手机号码",
//		"6009": "手机号码是黑名单",
//		"6010": "团队用户验证失败",
//		"6011": "团队不具备发送此号码的权限",
//		"6012": "该团队用户设置了IP限制",
//		"6013": "该团队用户余额不足",
//		"6014": "发送短信内容不能为空",
//		"6015": "短信内容超过了最大长度限制",
//		"6016": "团队密码必须大于4个字符",
//		"6017": "查询团队用户余额失败",
//		"6018": "用户没有开通SDK功能或测试已过期",
//		"6019": "此接口已经停止使用",
//		"6020": "此接口为VIP客户专用接口",
//		"6021": "扩展号码未备案",
//		"6022": "团队加密密钥不正确",
//		"6023": "短信商服务器故障",
//	}
//)
//
//type SmsMst struct {
//}
//
//func NewSmsMst() *SmsMst {
//	return &SmsMst{}
//}
//
//func (this *SmsMst) GetFee() (float64, error) {
//	urlString := fmt.Sprintf("%s?name=%s&pwd=%s",
//        WEB_SERVICE_GETFEE_URI, NAME, PASSWORD)
//
//	data, err := sendRequest(urlString)
//	if err != nil {
//		return 0, err
//	}
//
//	errInfo := data.Get("err")
//	errCode := data.Get("errid")
//	if len(errInfo) == 0 {
//		return 0, errors.New(fmt.Sprintf("Error code: %s", errCode))
//	}
//
//	if errCode == "0" {
//		fee, _ := strconv.ParseFloat(data.Get("id"), 64)
//		return fee / 10, nil
//	}
//	return 0, errors.New(fmt.Sprintf("Error code: %s", errorInfos[errCode]))
//}
//
//func (this *SmsMst) Check(content string) (string, error) {
//	enc := iconv.NewEncoder("gbk")
//	msg := enc.ConvertString(content)
//	urlString := fmt.Sprintf("%s?name=%s&pwd=%s&content=%s",
//        WEB_SERVICE_CHECK_URI, NAME, PASSWORD, msg)
//
//	data, err := sendRequest(urlString)
//	if err != nil {
//		return "", err
//	}
//	errCode := data.Get("errid")
//	if errCode == "0" {
//		return "没有包含屏蔽词", nil
//	}
//
//	return fmt.Sprintf("包含屏蔽词：%s", data.Get("err")), errors.New(fmt.Sprintf("Error code: %s", errorInfos[errCode]))
//}
//
//func (this *SmsMst) Send(mobile string, content string, isLong bool, timing string) error {
//	enc := iconv.NewEncoder("gbk")
//	msg := enc.ConvertString(content)
//	urlString := fmt.Sprintf("%s?name=%s&pwd=%s&dst=%s&msg=%s",
//        WEB_SERVICE_SEND_URI, NAME, PASSWORD, mobile, msg)
//
//	data, err := sendRequest(urlString)
//	if err != nil {
//		return err
//	}
//	errCode := data.Get("errid")
//	if errCode == "0" {
//		return nil
//	}
//
//	return errors.New(fmt.Sprintf("Error code: %s", errorInfos[errCode]))
//}
//
//func (this *SmsMst) SendWithTemplate(mobile string, templateCode string, params map[string]string) error {
//	return nil
//}
//
//// 发送请求
//func sendRequest(url string) (url.Values, error) {
//	response, err := http.Get(url)
//	defer response.Body.Close()
//	if err != nil {
//		return nil, err
//	}
//
//	result, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	data := stringToMap(string(result))
//
//	return data, err
//}
//
//// 解析请求返回字符串
//func stringToMap(data string) url.Values {
//	enc := iconv.NewEncoder("utf-8")
//	msg := enc.ConvertString(data)
//	result, _ := url.ParseQuery(msg)
//	return result
//}
