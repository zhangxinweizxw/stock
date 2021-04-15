package sms

import "stock/share/lib/sms/drivers"
const (
	SMS_DRIVER_MST     = "mst"
	SMS_DRIVER_ALIDAYU = "alidayu"
)

type ISms interface {
	GetFee() (float64, error)
	Check(string) (string, error)
	Send(string, string, bool, string) error
	SendWithTemplate(string, string, map[string]string) error
}

type SmsFactory struct {
	sms ISms
}

func NewSmsFactory(kind string) *SmsFactory {
	return &SmsFactory{sms: getSms(kind)}
}

func (this *SmsFactory) GetFee() (float64, error) {
	return this.sms.GetFee()
}

func (this *SmsFactory) Check(content string) (string, error) {
	return this.sms.Check(content)
}

func (this *SmsFactory) Send(mobile string, content string, isLong bool, timing string) error {
	return this.sms.Send(mobile, content, isLong, timing)
}

func (this *SmsFactory) SendWithTemplate(mobile string, templateCode string, params map[string]string) error {
	return this.sms.SendWithTemplate(mobile, templateCode, params)
}

func getSms(kind string) ISms {
	switch kind {
	case SMS_DRIVER_MST:
		return drivers.NewSmsMst()
		break
	case SMS_DRIVER_ALIDAYU:
		return drivers.NewAlidayu()
		break
	default:
		return nil
	}

	return nil
}
