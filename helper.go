package golib

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/jsonapi"
)

const (
	// ErrorDataNotFound error message when data doesn't exist
	ErrorDataNotFound = "data tidak ditemukan"
	// CHARS for setting short random string
	CHARS = "abcdefghijklmnopqrstuvwxyz0123456789"
	// NUMBERS for setting short random number
	NUMBERS = "0123456789"

	// PayloadInvalid constanta
	PayloadInvalid = "payload %s is invalid"

	// this block is for validating URL format
	email        string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	ip           string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	urlSchema    string = `((ftp|sftp|tcp|udp|wss?|https?):\/\/)`
	urlUsername  string = `(\S+(:\S*)?@)`
	urlPath      string = `((\/|\?|#)[^\s]*)`
	urlPort      string = `(:(\d{1,5}))`
	urlIP        string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	urlSubdomain string = `((www\.)|([a-zA-Z0-9]([-\.][-\._a-zA-Z0-9]+)*))`
	urlPattern   string = `^` + urlSchema + `?` + urlUsername + `?` + `((` + urlIP + `|(\[` + ip + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + urlSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + urlPort + `?` + urlPath + `?$`
	area         string = `^\+\d{1,5}$`
	phone        string = `^\d{5,}$`
)

var (
	// ErrBadFormatURL variable for error of url format
	ErrBadFormatURL = errors.New("invalid url format")
	// ErrBadFormatMail variable for error of email format
	ErrBadFormatMail = errors.New("invalid email format")
	// ErrBadFormatPhoneNumber variable for error of email format
	ErrBadFormatPhoneNumber = errors.New("invalid phone format")

	// emailRegexp regex for validate email
	emailRegexp = regexp.MustCompile(email)
	// urlRegexp regex for validate URL
	urlRegexp = regexp.MustCompile(urlPattern)
	// areaRegexp  regex for phone area number using +
	areaRegexp = regexp.MustCompile(area)
	// telpRegexp regex for phone number
	phoneRegexp = regexp.MustCompile(phone)
	// camel regex for string camelcase
	camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")
)

// ValidateEmail function for validating email
func ValidateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrBadFormatMail
	}
	return nil
}

// ValidateURL function for validating url
func ValidateURL(str string) error {
	if !urlRegexp.MatchString(str) {
		return ErrBadFormatURL
	}
	return nil
}

// ValidatePhoneNumber function for validating phone number only
func ValidatePhoneNumber(str string) error {
	if !phoneRegexp.MatchString(str) {
		return ErrBadFormatPhoneNumber
	}
	return nil
}

// ValidatePhoneAreaNumber function for validating area phone number
func ValidatePhoneAreaNumber(str string) error {
	if !areaRegexp.MatchString(str) {
		return ErrBadFormatPhoneNumber
	}
	return nil
}

// StringArrayReplace function for replacing whether string in array
// str string searched string
// list []string array
func StringArrayReplace(str string, listFind, listReplace []string) string {
	for i, v := range listFind {
		if strings.Contains(str, v) {
			str = strings.Replace(str, v, listReplace[i], -1)
		}
	}
	return str
}

// ValidateMaxInput function for validating maximum input
func ValidateMaxInput(input string, limit int) error {
	if len(input) > limit {
		err := errors.New(" value is too long")
		return err
	}

	return nil
}

// ValidateNumeric function for check valid numeric
func ValidateNumeric(str string) bool {
	var num, symbol int
	for _, r := range str {
		if r >= 48 && r <= 57 { //code ascii for [0-9]
			num = +1
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	return num >= 1
}

// ValidateAlphabet function for check alphabet
func ValidateAlphabet(str string) bool {
	var uppercase, lowercase, symbol int
	for _, r := range str {
		if r >= 65 && r <= 90 { //code ascii for [A-Z]
			uppercase = +1
		} else if r >= 97 && r <= 122 { //code ascii for [a-z]
			lowercase = +1
		} else { //except alphabet
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}
	return uppercase >= 1 || lowercase >= 1
}

// ValidateAlphabetWithSpace function for check alphabet with space
func ValidateAlphabetWithSpace(str string) bool {
	var uppercase, lowercase, space, symbol int
	for _, r := range str {
		if r >= 65 && r <= 90 { //code ascii for [A-Z]
			uppercase = +1
		} else if r >= 97 && r <= 122 { //code ascii for [a-z]
			lowercase = +1
		} else if r == 32 { //code ascii for space
			space = +1
		} else { //except alphabet
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}
	return uppercase >= 1 || lowercase >= 1 || space >= 1
}

// ValidateAlphanumeric function for check valid alphanumeric
func ValidateAlphanumeric(str string, must bool) bool {
	var uppercase, lowercase, num, symbol int
	for _, r := range str {
		if r >= 65 && r <= 90 { //code ascii for [A-Z]
			uppercase = +1
		} else if r >= 97 && r <= 122 { //code ascii for [a-z]
			lowercase = +1
		} else if r >= 48 && r <= 57 { //code ascii for [0-9]
			num = +1
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	if must { //must alphanumeric
		return (uppercase >= 1 || lowercase >= 1) && num >= 1
	}

	return uppercase >= 1 || lowercase >= 1 || num >= 1
}

// ValidateAlphanumericWithSpace function for validating string to alpha numeric with space
func ValidateAlphanumericWithSpace(str string, must bool) bool {
	var uppercase, lowercase, num, space, symbol int
	for _, r := range str {
		if r >= 65 && r <= 90 { //code ascii for [A-Z]
			uppercase = +1
		} else if r >= 97 && r <= 122 { //code ascii for [a-z]
			lowercase = +1
		} else if r >= 48 && r <= 57 { //code ascii for [0-9]
			num = +1
		} else if r == 32 { //code ascii for space
			space = +1
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	if must { //must alphanumeric
		return (uppercase >= 1 || lowercase >= 1) && num >= 1 && space >= 1
	}

	return (uppercase >= 1 || lowercase >= 1 || num >= 1) || space >= 1
}

// GenerateRandomID function for generating shipping ID
func GenerateRandomID(length int, prefix ...string) string {
	var strPrefix string

	if len(prefix) > 0 {
		strPrefix = prefix[0]
	}

	yearNow, monthNow, _ := time.Now().Date()
	year := strconv.Itoa(yearNow)[2:len(strconv.Itoa(yearNow))]
	month := int(monthNow)
	RandomString := RandomString(length)

	id := fmt.Sprintf("%s%s%d%s", strPrefix, year, month, RandomString)
	return id
}

// RandomString function for random string
func RandomString(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	charsLength := len(CHARS)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = CHARS[rand.Intn(charsLength)]
	}
	return string(result)
}

// RandomNumber function for random number
func RandomNumber(length int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	charsLength := len(NUMBERS)
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = NUMBERS[rand.Intn(charsLength)]
	}
	return string(result)
}

// StringInSlice function for checking whether string in slice
// str string searched string
// list []string slice
func StringInSlice(str string, list []string, caseSensitive ...bool) bool {
	isCaseSensitive := true
	if len(caseSensitive) > 0 {
		isCaseSensitive = caseSensitive[0]
	}

	if isCaseSensitive {
		for _, v := range list {
			if v == str {
				return true
			}
		}
	} else {
		for _, v := range list {
			if strings.ToLower(v) == strings.ToLower(str) {
				return true
			}
		}
	}

	return false
}

// GetProtocol function for getting http protocol based on TLS
// isTLS bool
func GetProtocol(isTLS bool) string {
	// check tls first to get protocol
	if isTLS {
		return "https://"
	}
	return "http://"
}

// GetHostURL function for getting host of any URL
func GetHostURL(req *http.Request) string {
	return fmt.Sprintf("%s%s", GetProtocol(req.TLS != nil), req.Host)
}

// GetSelfLink function to get self link
func GetSelfLink(req *http.Request) string {
	return fmt.Sprintf("%s%s", GetHostURL(req), req.RequestURI)
}

// MarshalConvertManyPayload function to convert struct response to jsonapi.manypayload so that we can add meta or link data
func MarshalConvertManyPayload(structResponse interface{}) (payload *jsonapi.ManyPayload, err error) {
	// set response marshal jsonapi struct
	p, err := jsonapi.Marshal(structResponse)
	if err != nil {
		return nil, err
	}

	var ok bool
	if payload, ok = p.(*jsonapi.ManyPayload); !ok {
		err = fmt.Errorf(PayloadInvalid, "many payload")
		return nil, err
	}

	return
}

// MarshalConvertOnePayload function to convert struct response to jsonapi.OnePayLoad so that we can add meta or link data
func MarshalConvertOnePayload(structResponse interface{}) (payload *jsonapi.OnePayload, err error) {
	// set response marshal jsonapi struct
	p, err := jsonapi.Marshal(structResponse)
	if err != nil {
		return nil, err
	}

	var ok bool
	if payload, ok = p.(*jsonapi.OnePayload); !ok {
		err = fmt.Errorf(PayloadInvalid, "one payload")
		return nil, err
	}

	return
}

// IdentifyPanic for identify line code in panic recover
func IdentifyPanic(ctx string, rec interface{}) string {
	var name, file string
	var line int
	var pc [16]uintptr

	n := runtime.Callers(3, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line = fn.FileLine(pc)
		name = fn.Name()
		if !strings.HasPrefix(name, "runtime.") {
			break
		}
	}

	var source, githubLink string
	switch {
	case name != "":
		source = fmt.Sprintf("%v:%v", name, line)
	case file != "":
		source = fmt.Sprintf("%v:%v", file, line)
	default:
		source = fmt.Sprintf("pc:%x", pc)
	}

	branch := os.Getenv("SERVER_ENV")
	if branch == "production" {
		branch = "master"
	}

	sign := os.Getenv("PROJECT_NAME")
	i := strings.Index(file, sign)
	if i > 0 {
		githubLink = file[i+len(sign):]
	}

	i = strings.Index(name, sign)
	if i > 0 {
		githubLink = fmt.Sprintf("https://%s/blob/%s%s#L%d", name[:i+len(sign)], branch, githubLink, line)
	}

	if githubLink == "" {
		githubLink = source
	}

	SendNotification("Panic Detected", fmt.Sprintf("*Panic source*: `%s`", githubLink), ctx, fmt.Errorf("%v", rec))
	return fmt.Sprintf("panic: %v", rec)
}

// MaskPassword for mask password string
func MaskPassword(s string) string {
	splitText := strings.Split(s, "&")

	var newText string
	for i, text := range splitText {

		password := strings.Index(text, "password=")
		if password > -1 {
			text = "password=xxxxx"
		}

		newPassword := strings.Index(text, "newPassword=")
		if newPassword > -1 {
			text = "newPassword=xxxxx"
		}

		rePassword := strings.Index(text, "rePassword=")
		if rePassword > -1 {
			text = "rePassword=xxxxx"
		}

		if i < 1 {
			newText = text
		} else {
			newText = newText + "&" + text
		}

	}
	return newText
}

// ValidateLatinOnly func for check valid latin only
func ValidateLatinOnly(str string) bool {
	var uppercase, lowercase, num, allowed, symbol int
	for _, r := range str {
		if r >= 65 && r <= 90 { //code ascii for [A-Z]
			uppercase = +1
		} else if r >= 97 && r <= 122 { //code ascii for [a-z]
			lowercase = +1
		} else if r >= 48 && r <= 57 { //code ascii for [0-9]
			num = +1
		} else if r >= 32 && r <= 47 || r >= 58 && r <= 64 || r >= 91 && r <= 96 || r >= 123 && r <= 126 {
			allowed = +1 //code ascii for [space, coma, ., !, ", #, $, %, &, ', (, ), *, +, -, /, :, ;, <, =, >, ?, @, [, \, ], ^, _, `, {, |, }, ~]
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	return uppercase >= 1 || lowercase >= 1 || num >= 1 || allowed >= 0
}

// CamelToLowerCase func for replace camel to lower case
func CamelToLowerCase(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, " "))
}
