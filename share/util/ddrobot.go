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

		webHook := "https://oapi.dingtalk.com/robot/send?access_token=b3286b5a01ab43f2fef9f0479089e543aaa40333d32d7d52e931c9a204143295"
		_, err := http.Post(webHook, "application/json;charset=utf-8", bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println("消息发送成功!")
	}

}
