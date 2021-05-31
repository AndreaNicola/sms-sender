package sms_sender

import (
	"errors"
	messagebird "github.com/messagebird/go-rest-api/v6"
	"github.com/messagebird/go-rest-api/v6/sms"
	"github.com/messagebird/go-rest-api/v6/verify"
	"github.com/ttacon/libphonenumber"
	"log"
	"strconv"
)

type SmsSender struct {
	client        *messagebird.Client
	originator    string
	defaultRegion string
}

func (s *SmsSender) cleanNumber(formattedPhoneNumber string) string {

	num, err := libphonenumber.Parse(formattedPhoneNumber, s.defaultRegion)
	if err != nil {
		panic(err)
	}

	cc := num.CountryCode
	return strconv.FormatInt(int64(*cc), 10) + strconv.FormatUint(*num.NationalNumber, 10)

}

func NewSmsSender(apiKey, defaultRegion, originator string) *SmsSender {

	return &SmsSender{
		client:        messagebird.New(apiKey),
		originator:    originator,
		defaultRegion: defaultRegion,
	}

}

func (s *SmsSender) CreateVerifyToken(recipient string) (string, error) {

	cn := s.cleanNumber(recipient)
	v, err := verify.Create(s.client, cn, &verify.Params{
		Originator: s.originator,
		Timeout:    600,
	})

	if err == nil {
		return v.ID, nil
	} else {
		log.Println(err.Error())
		return "", err
	}

}

func (s *SmsSender) VerifyToken(tokenId, token string) (bool, error) {

	v, err := verify.VerifyToken(s.client, tokenId, token)
	if err == nil {
		return v.Status == "verified", nil
	} else {
		return false, nil
	}

}

func (s *SmsSender) SendSms(text string, recipients ...string) error {

	if len(recipients) <= 0 {
		return errors.New("recipients must not be empty")
	}

	cns := make([]string, 0)
	for _, cn := range recipients {
		cns = append(cns, s.cleanNumber(cn))
	}

	_, err := sms.Create(s.client, s.originator, cns, text, nil)
	return err
}
