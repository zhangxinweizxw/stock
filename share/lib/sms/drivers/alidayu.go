package drivers

const (
	APP_KEY    = "23466956"
	APP_SECRET = "41ad9ec16967ab3b92aaa1eea502466f"
	SIGN_NAME  = "首证投顾"
)

func init() {
	alidayu.InitAlidayu(true, APP_KEY, APP_SECRET)
}

type SmsAlidayu struct {
}

func NewAlidayu() *SmsAlidayu {
	return &SmsAlidayu{}
}

func (this *SmsAlidayu) GetFee() (float64, error) {
	return 0, nil
}

func (this *SmsAlidayu) Check(content string) (string, error) {
	return "", nil
}

func (this *SmsAlidayu) Send(mobile string, content string, isLong bool, timing string) error {
	return nil
}

func (this *SmsAlidayu) SendWithTemplate(mobile string, templateCode string, params map[string]string) error {
	msg := alidayu.NewMessageSms(SIGN_NAME).SetTel(mobile).SetContent(templateCode, params)
	if err := alidayu.SendMessageDirect(msg); err != nil {
		return err
	}
	return nil
}
