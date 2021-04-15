package lib

import (
	"encoding/xml"
    "fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"stock/share/logging"
)

// returnMessage after send sms
type Send_ReturnMessage struct {
	XMLName       xml.Name `xml:"returnsms"`
	ReturnStatus  string   `xml:"returnstatus"`
	Message       string   `xml:"message"`
	Remainpoint   string   `xml:"remainpoint"`
	TaskID        string   `xml:"taskID"`
	SuccessCounts string   `xml:"successCounts"`
}

type Sms struct {
	Addr     string
	UserId   string
	Account  string
	Password string
}

//
func NewSms(Addr, UserId, Account, Password string) *Sms {
	return &Sms{Addr, UserId, Account, Password}
}

// Send
func (this *Sms) Send(mobile, content string) (error, int) {
	url := fmt.Sprintf("%s?action=send&userid=%s&account=%s&password=%s&mobile=%s&content=%s&sendTime=&extno=",
		this.Addr, this.UserId, this.Account, this.Password, mobile, content)

	logging.Debug("%s", url)

	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return err, 0
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err, 0
	}

	tmp := Send_ReturnMessage{}
	err = xml.Unmarshal(body, &tmp)
	if err != nil {
		return err, 0
	}

	status := strings.TrimSpace(strings.ToLower(tmp.ReturnStatus))
	if status == "success" {
		return nil, 1
	}

	return nil, 1
}
