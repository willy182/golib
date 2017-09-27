package golib

import (
	"strings"
	"testing"
)

func TestRandomString(t *testing.T) {
	length := 10
	rs := RandomString(length)

	// set invalid string that
	// should not be contained in random string
	invalidString := `!@#$%^&*()_+`

	if len(rs) != length {
		t.Errorf("length of random string is not equal %d", length)
	}

	if strings.Contains(rs, invalidString) {
		t.Fatal("random string contains symbols")
	}
}

func TestValidateEmail(t *testing.T) {
	var (
		email string
		err   error
	)

	// test valid email should be not error/valid
	email = "julius.bernhard@bhinneka.com"
	if err = ValidateEmail(email); err != nil {
		t.Fatal("testing valid email is not valid")
	}

	// test invalid email should be error
	email = "julius.@bernhard@bhinneka.com"
	if err = ValidateEmail(email); err == nil {
		t.Fatal("testing invalid email is not valid")
	}
}

func TestValidateURL(t *testing.T) {
	var (
		url string
		err error
	)

	url = "http://www.bhinneka.com"
	if err = ValidateURL(url); err != nil {
		t.Fatal("testing 1st valid URL is not valid")
	}

	url = "www.bhinneka.com"
	if err = ValidateURL(url); err != nil {
		t.Fatal("testing 2nd valid URL is not valid")
	}

	url = "ftp://www.bhinneka.com"
	if err = ValidateURL(url); err != nil {
		t.Fatal("testing 3rd valid URL is not valid")
	}

	url = "https:///www.bhinneka.com"
	if err = ValidateURL(url); err == nil {
		t.Fatal("testing invalid URL is not valid")
	}
}

func TestValidatePhoneNumber(t *testing.T) {
	var (
		tel string
		err error
	)

	tel = "08119889788"
	if err = ValidatePhoneNumber(tel); err != nil {
		t.Fatal("testing valid phone number is not valid")
	}

	tel = "081-1988-9788"
	if err = ValidatePhoneNumber(tel); err == nil {
		t.Fatal("testing 1st invalid phone number is not valid")
	}

	tel = "0811"
	if err = ValidatePhoneNumber(tel); err == nil {
		t.Fatal("testing 2nd invalid phone number - not greater than 5 chars is not valid")
	}
}

func TestValidatePhoneAreaNumber(t *testing.T) {
	var (
		area string
		err  error
	)

	area = "+62"
	if err = ValidatePhoneAreaNumber(area); err != nil {
		t.Fatal("testing valid area number is not valid")
	}

	area = "+6 2"
	if err = ValidatePhoneAreaNumber(area); err == nil {
		t.Fatal("testing 1st invalid area number is not valid")
	}

	area = "+"
	if err = ValidatePhoneAreaNumber(area); err == nil {
		t.Fatal("testing 2nd invalid area number is not valid")
	}
}

func TestValidateAlphaNumeric(t *testing.T) {
	var (
		alpha string
	)

	alpha = "Some days are beautiful."
	if !ValidateAlphaNumeric(alpha) {
		t.Fatal("testing valid alpha numeric is not valid")
	}

	alpha = "Some days are beautiful. :) :*"
	if ValidateAlphaNumeric(alpha) {
		t.Fatal("testing 1st invalid alpha numeric is not valid")
	}

	alpha = `<img src="http://example.com/image.jpg" />`
	if ValidateAlphaNumeric(alpha) {
		t.Fatal("testing 2nd invalid alpha numeric is not valid")
	}
}
