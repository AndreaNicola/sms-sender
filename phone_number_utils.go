package sms_sender

import (
	"errors"
	"github.com/nyaruka/phonenumbers"
	"strconv"
)

func ParseAndFormat(number, defaultRegion string) (string, error) {

	if number == "" {
		return "", errors.New("phone number is empty")
	}

	if defaultRegion == "" {
		return "", errors.New("default region is empty")
	}

	parse, err := phonenumbers.Parse(number, defaultRegion)
	if err != nil {
		panic(err.Error())
	}

	numberType := phonenumbers.GetNumberType(parse)
	valid := phonenumbers.IsValidNumber(parse)

	if valid && *parse.CountryCode == int32(39) && (numberType == phonenumbers.MOBILE || numberType == phonenumbers.FIXED_LINE_OR_MOBILE) {
		return phonenumbers.Format(parse, phonenumbers.INTERNATIONAL), nil
	} else {
		return "", errors.New("phone number is not valid")
	}

}

func cleanNumber(formattedPhoneNumber, defaultRegion string) (string, error) {

	num, err := phonenumbers.Parse(formattedPhoneNumber, defaultRegion)
	if err != nil {
		return "", err
	}

	cc := num.CountryCode
	return strconv.FormatInt(int64(*cc), 10) + strconv.FormatUint(*num.NationalNumber, 10), nil

}
