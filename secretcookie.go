package secretcookie

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
	SecretKey string
	CacheDays int
}{
	SecretKey: "secretcookie",
	CacheDays: 30,
}

// SetSecretCookie :
func SetSecretCookie(writer http.ResponseWriter, key string, value string) (err error) {
	cookie := http.Cookie{
		Name:  key,
		Value: createSignedValue(Config.SecretKey, key, value),
	}
	http.SetCookie(writer, &cookie)
	return
}

// GetSecretCookie :
func GetSecretCookie(request *http.Request, key string) (value string, err error) {
	singelCookie, err := request.Cookie(key)
	if err != nil {
		return
	}
	svalue := singelCookie.Value
	value, err = decodeSignedValue(Config.SecretKey, key, svalue, Config.CacheDays)
	return
}

func formatField(s string) string {
	return strconv.Itoa(len(s)) + ":" + s
}

func createSignature(secret string, value string) string {
	key := []byte(secret)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))
	return hex.EncodeToString(mac.Sum(nil))
}

func createSignedValue(secret string, name string, value string) string {
	timestamp := strconv.Itoa(int(time.Now().Unix()))

	b64value := base64.URLEncoding.EncodeToString([]byte(value))
	toSign := strings.Join(
		[]string{
			"2",
			formatField("0"),
			formatField(timestamp),
			formatField(name),
			formatField(b64value),
			""},
		"|")
	return toSign + createSignature(secret, toSign)
}

func consumeField(v string) (value string) {
	sli := strings.Split(v, ":")
	length, _ := strconv.Atoi(sli[0])
	value = sli[1][:length]
	return
}

func decodeFiledsValue(v string) (version int, timestamp int, key string, value string, sig string) {
	sli := strings.Split(v, "|")
	version, _ = strconv.Atoi(consumeField(sli[1]))
	timestamp, _ = strconv.Atoi(consumeField(sli[2]))
	key = consumeField(sli[3])
	value = consumeField(sli[4])
	sig = sli[5]
	return
}

func decodeSignedValue(secret string, name string, svalue string, maxCacheDay int) (string, error) {
	_, timestamp, key, value, passedSig := decodeFiledsValue(svalue)
	signedString := svalue[:len(svalue)-len(passedSig)]
	exceptSig := createSignature(secret, signedString)
	if passedSig != exceptSig {
		return "", errors.New("consistency of signal check failed")
	}
	if key == "" || key != name {
		return "", errors.New("unknown key")
	}
	if timestamp < (int(time.Now().Unix()) - maxCacheDay*86400) {
		return "", errors.New("the signature has expired")
	}
	devaluebyte, err := base64.URLEncoding.DecodeString(value)
	devalue := string(devaluebyte)
	if err != nil {
		return "", err
	}
	return devalue, nil
}
