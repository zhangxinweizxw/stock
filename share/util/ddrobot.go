package util

import (
	"bytes"
	"fmt"
	"net/http"
)

// 钉钉机器人发送消息
type DdRobot struct {
}

func NewDdRobot() *DdRobot {
	return &DdRobot{}
}

// 发送钉钉机器人消息
func (*DdRobot) DdRobotPush(stockInfo string) {
	//logging.Debug("====%v", stockInfo)
	//如果有未发送新闻 请求钉钉webhook
	if stockInfo != "" {

		formt := `
		{
			"msgtype": "text",
			"text": {
	   	"content": "stock|
	                   %v"
			},
			"at": {
	   		"atMobiles": [],
	   	"isAtAll": false
			}
		 }`

		body := fmt.Sprintf(formt, stockInfo)
		jsonValue := []byte(body)
		//发送消息到钉钉群使用webhook

		webHook := "https://oapi.dingtalk.com/robot/send?access_token=ab80dbae81ff47b5b3d60e5c585051ed1266fa01f6207a8d22e1d6c6950ea053"
		_, err := http.Post(webHook, "application/json;charset=utf-8", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("消息发送成功!")
	}

}
