package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"strings"
	"time"
)

func GenerateTOTPSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)), nil
}

func GenerateTOTPQRUri(accountName, secret string) string {
	issuer := "Codedock"
	return fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		url.QueryEscape(issuer),
		url.QueryEscape(accountName),
		url.QueryEscape(secret),
		url.QueryEscape(issuer),
	)
}

func GenerateRecoveryCodes(count int) ([]string, error) {
	var codes []string
	for i := 0; i < count; i++ {
		buf := make([]byte, 4)
		if _, err := rand.Read(buf); err != nil {
			return nil, err
		}
		codes = append(codes, strings.ToUpper(fmt.Sprintf("%04x-%04x", buf[:2], buf[2:4])))
	}
	return codes, nil
}

func ValidateTOTP(secret, passcode string) bool {
	passcode = strings.TrimSpace(passcode)
	if len(passcode) != 6 {
		return false
	}
	secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		secretBytes, err = base32.StdEncoding.DecodeString(strings.ToUpper(secret))
		if err != nil {
			return false
		}
	}
	now := time.Now().Unix() / 30
	for step := -1; step <= 1; step++ {
		if GenerateTOTPCode(secretBytes, now+int64(step)) == passcode {
			return true
		}
	}
	return false
}

func GenerateTOTPCode(secret []byte, timeStep int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(timeStep))
	mac := hmac.New(sha1.New, secret)
	mac.Write(buf)
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0xf
	code := int64(((int(sum[offset])&0x7f)<<24)|
		((int(sum[offset+1])&0xff)<<16)|
		((int(sum[offset+2])&0xff)<<8)|
		(int(sum[offset+3])&0xff)) % 1000000
	return fmt.Sprintf("%06d", code)
}
