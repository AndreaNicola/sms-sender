package sms_sender

import (
	"errors"
	"fmt"
	messagebird "github.com/messagebird/go-rest-api/v9"
	"github.com/messagebird/go-rest-api/v9/balance"
	"github.com/messagebird/go-rest-api/v9/sms"
	"github.com/messagebird/go-rest-api/v9/verify"
	"log"
	"os"
)

type SmsSender struct {
	client        *messagebird.DefaultClient
	originator    string
	defaultRegion string
}

func NewDefaultSmsSender() *SmsSender {

	mbApiKey := os.Getenv("MB_API_KEY")
	if mbApiKey == "" {
		panic("messagebird api key is mandatory")
	}
	mbDefaultRegion := os.Getenv("MB_DEFAULT_REGION")
	if mbDefaultRegion == "" {
		panic("messagebird default region is mandatory")
	}
	mbOriginator := os.Getenv("MB_ORIGINATOR")
	if mbOriginator == "" {
		panic("messagebird originator is mandatory")
	}

	return NewSmsSender(mbApiKey, mbDefaultRegion, mbOriginator)

}

func NewSmsSender(apiKey, defaultRegion, originator string) *SmsSender {

	client := messagebird.New(apiKey)

	// Request the balance information, returned as a balance.Balance object.
	bal, err := balance.Read(client)
	if err != nil {
		log.Println("sms-sender/sms.go -> NewSmsSender -> error reading balance", err.Error())
	}

	// Display the results.
	fmt.Println("Payment: ", bal.Payment)
	fmt.Println("Type:", bal.Type)
	fmt.Println("Amount:", bal.Amount)

	return &SmsSender{
		client:        client,
		originator:    originator,
		defaultRegion: defaultRegion,
	}

}

func (sms2 *SmsSender) CreateVerifyToken(recipient string) (string, error) {

	cn, cnErr := cleanNumber(recipient, sms2.defaultRegion)
	if cnErr != nil {
		log.Println("sms-sender/sms.go -> CreateVerifyToken -> error cleaning number", cnErr.Error())
		return "", cnErr
	}

	v, err := verify.Create(sms2.client, cn, &verify.Params{
		Originator: sms2.originator,
		Timeout:    600,
	})

	if err != nil {
		log.Println("sms-sender/sms.go -> CreateVerifyToken -> error creating verify token", err.Error())
		return "", err

	}
	return v.ID, nil

}

func (sms2 *SmsSender) VerifyToken(tokenId, token string) (bool, error) {

	v, err := verify.VerifyToken(sms2.client, tokenId, token)
	if err != nil {
		log.Println("sms-sender/sms.go -> VerifyToken -> error verifying token", err.Error())
		return false, nil
	}

	return v.Status == "verified", nil

}

func (sms2 *SmsSender) SendSms(text string, recipients ...string) error {

	if len(recipients) <= 0 {
		log.Println("sms-sender/sms.go -> SendSms -> recipients must not be empty")
		return errors.New("recipients must not be empty")
	}

	cns := make([]string, 0)
	for _, r := range recipients {

		cn, cnErr := cleanNumber(r, sms2.defaultRegion)
		if cnErr != nil {
			log.Println("sms-sender/sms.go -> SendSms -> error cleaning number", cnErr.Error())
			continue
		}

		cns = append(cns, cn)
	}

	_, err := sms.Create(sms2.client, sms2.originator, cns, text, nil)
	return err
}

func (sms2 *SmsSender) ParseAndFormat(number string) (string, error) {
	return ParseAndFormat(number, sms2.defaultRegion)
}
