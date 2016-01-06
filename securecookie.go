package securecookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Config :
var Config = struct {
	SecureKey  string
	CacheDays  int
	Version    int
	KeyVersion int
}{
	SecureKey:  "securecookie",
	CacheDays:  30,
	Version:    2,
	KeyVersion: 0,
}

// FormatField :
func FormatField(value string) string {
	return strconv.Itoa(len(value)) + ":" + value
}

// ConsumeField accept strconv.Itoa(len(value)) + ":" + value and return value
func ConsumeField(field string) (string, error) {
	index := strings.Index(field, ":")
	if index == -1 {
		return "", errors.New("ConsumeField: Can not find sep string \":\". Field: " + field)
	}
	length, err := strconv.Atoi(field[:index])
	if err != nil {
		return "", err
	}
	value := field[index+1:]
	if len(value) != length {
		return "", errors.New("ConsumeField: Unequal lengths. Field: " + field)
	}
	return value, nil
}

// CreateSignature :
func CreateSignature(secure string, value string) string {
	key := []byte(secure)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

// CreateSignedValue :
func CreateSignedValue(secure string, name string, value string) string {
	tosign := strings.Join(
		[]string{
			strconv.Itoa(Config.Version),
			FormatField(strconv.Itoa(Config.KeyVersion)),
			FormatField(strconv.Itoa(int(time.Now().Unix()))),
			FormatField(name),
			FormatField(base64.URLEncoding.EncodeToString([]byte(value))),
			""},
		"|")
	return tosign + CreateSignature(secure, tosign)
}

// SecureCookie :
type SecureCookie struct {
	Version    int
	KeyVersion int
	Timestamp  int
	Key        string
	Value      string
	Signature  string
}

// DecodeFiledsValue :
func DecodeFiledsValue(source string) (*SecureCookie, error) {
	fields := strings.Split(source, "|")
	if len(fields) != 6 {
		return nil, errors.New("DecodeFiledsValue: fields length unequals 6. source: " + source)
	}
	version, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}
	if version != Config.Version {
		return nil, errors.New("DecodeFiledsValue: version unequals " + strconv.Itoa(Config.Version))
	}
	keyversionTemp, err := ConsumeField(fields[1])
	if err != nil {
		return nil, err
	}
	keyversion, err := strconv.Atoi(keyversionTemp)
	if err != nil {
		return nil, err
	}
	timestmpTemp, err := ConsumeField(fields[2])
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.Atoi(timestmpTemp)
	if err != nil {
		return nil, err
	}
	key, err := ConsumeField(fields[3])
	if err != nil {
		return nil, err
	}
	value, err := ConsumeField(fields[4])
	if err != nil {
		return nil, err
	}
	signature := fields[5]
	return &SecureCookie{
		Version:    version,
		KeyVersion: keyversion,
		Timestamp:  timestamp,
		Key:        key,
		Value:      value,
		Signature:  signature,
	}, nil
}

// DecodeSignedValue :
func DecodeSignedValue(secure string, name string, svalue string, maxCacheDay int) (string, error) {
	securecookie, err := DecodeFiledsValue(svalue)
	if err != nil {
		return "", err
	}
	tosign := svalue[:len(svalue)-len(securecookie.Signature)]
	exceptSignature := CreateSignature(secure, tosign)
	if securecookie.Signature != exceptSignature {
		return "", errors.New("DecodeSignedValue: signature consist check failed")
	}
	if securecookie.Key == "" || securecookie.Key != name {
		return "", errors.New("DecodeSignedValue: unknown key")
	}
	if securecookie.Timestamp < (int(time.Now().Unix()) - maxCacheDay*86400) {
		return "", errors.New("DecodeSignedValue: the signature has expired")
	}
	devaluebyte, err := base64.URLEncoding.DecodeString(securecookie.Value)
	if err != nil {
		return "", err
	}
	devalue := string(devaluebyte)
	return devalue, nil
}

// SetSecureCookie :
func SetSecureCookie(writer http.ResponseWriter, key string, value string) (err error) {
	cookie := http.Cookie{
		Name:  key,
		Value: CreateSignedValue(Config.SecureKey, key, value),
	}
	http.SetCookie(writer, &cookie)
	return
}

// GetSecureCookie :
func GetSecureCookie(request *http.Request, key string) (value string, err error) {
	singelCookie, err := request.Cookie(key)
	if err != nil {
		return
	}
	svalue := singelCookie.Value
	value, err = DecodeSignedValue(Config.SecureKey, key, svalue, Config.CacheDays)
	return
}
