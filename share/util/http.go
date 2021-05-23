package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type HttpUtil struct {
}

func NewHttpUtil() *HttpUtil {
	return &HttpUtil{}
}

// Post接口调用 json格式
func (*HttpUtil) PostJson(param map[string]interface{}, url, cookie string) (error, *http.Response) {

	bytesData, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	reader := bytes.NewReader(bytesData)

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}

	request.Header.Set("Content-Type", "application/json")

	if len(cookie) > 0 {
		request.Header.Set("Cookie", cookie)
	}
	client := http.Client{}
	resp, reqerr := client.Do(request)

	if reqerr != nil {
		fmt.Println(reqerr.Error())
		return err, nil
	}

	return err, resp
}

// Post接口调用 json格式
func (*HttpUtil) Post1Json(param []interface{}, url, cookie string) (error, *http.Response) {

	bytesData, err := json.Marshal(param)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	//fmt.Println("-------------", string(bytesData))
	reader := bytes.NewReader(bytesData)

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}

	request.Header.Set("Content-Type", "application/json")

	if len(cookie) > 0 {
		request.Header.Set("Cookie", cookie)
	}
	//fmt.Println(cookie)
	client := http.Client{}
	resp, reqerr := client.Do(request)

	if reqerr != nil {
		fmt.Println(reqerr.Error())
		return err, nil
	}

	return err, resp
}

// Post接口调用 form格式
func (*HttpUtil) PostForms(param, urls, cookie string) (error, *http.Response) {

	body := strings.NewReader(param)
	request, err := http.NewRequest("POST", urls, body)
	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(cookie) > 0 {
		request.Header.Set("Cookie", cookie)
	}
	client := http.Client{}
	resp, cerr := client.Do(request)

	if cerr != nil {
		fmt.Println(cerr.Error())
		return cerr, nil
	}

	return err, resp
}

// get 方法封装
func (*HttpUtil) GetJson(url string) (error, *http.Response) {
	client := &http.Client{} //生成要访问的url
	//body := strings.NewReader(param)
	reqest, err := http.NewRequest("GET", url, nil) //增加header选项

	if err != nil {
		panic(err)
	} //处理返回结果
	resp, _ := client.Do(reqest)
	return err, resp
}
