package sms_sender

import "testing"

func TestSmsSender_SendSms(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			t.Error(err)
		}
	}()

	dfss := NewDefaultSmsSender()
	err := dfss.SendSms("Is there anybody out there?", "+393495564972")
	if err != nil {
		t.Error(err)
	}

}
