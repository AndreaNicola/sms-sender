package sms_sender

import (
	"errors"
	"github.com/ttacon/libphonenumber"
)

func (sms2 *SmsSender) ParseAndFormat(s string) (string, error) {

	parse, err := libphonenumber.Parse(s, sms2.defaultRegion)
	if err != nil {
		panic(err.Error())
	}

	numberType := libphonenumber.GetNumberType(parse)
	valid := libphonenumber.IsValidNumber(parse)

	if valid && *parse.CountryCode == int32(39) && (numberType == libphonenumber.MOBILE || numberType == libphonenumber.FIXED_LINE_OR_MOBILE) {
		return libphonenumber.Format(parse, libphonenumber.INTERNATIONAL), nil
	} else {
		return "", errors.New("phone number is not valid")
	}

}
