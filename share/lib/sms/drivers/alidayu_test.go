package drivers

import (
    "testing"
)

func Test_SendWithTemplate(t *testing.T) {
	sms := NewAlidayu()
	params := map[string]string{
		"code": "825364",
	}
	err := sms.SendWithTemplate("18602990888", "SMS_10230656", params)
	if err != nil {
		t.Error(err)
	}
}
