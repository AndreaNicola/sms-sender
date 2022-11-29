package sms_sender

import (
	"errors"
	"github.com/ttacon/libphonenumber"
	"strconv"
)

func ParseAndFormat(number, defaultRegion string) (string, error) {

	if number == "" {
		return "", errors.New("phone number is empty")
	}

	if defaultRegion == "" {
		return "", errors.New("default region is empty")
	}

	parse, err := libphonenumber.Parse(number, defaultRegion)
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

func cleanNumber(formattedPhoneNumber, defaultRegion string) (string, error) {

	num, err := libphonenumber.Parse(formattedPhoneNumber, defaultRegion)
	if err != nil {
		return "", err
	}

	cc := num.CountryCode
	return strconv.FormatInt(int64(*cc), 10) + strconv.FormatUint(*num.NationalNumber, 10), nil

}
