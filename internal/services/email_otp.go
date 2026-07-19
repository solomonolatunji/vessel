package services

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"
)

type emailOTPData struct {
	OTP       string
	NewEmail  string
	ExpiresAt time.Time
}

var (
	emailOTPs = make(map[string]emailOTPData) // userID -> data
	otpMutex  sync.Mutex
)

func GenerateEmailOTP(userID, newEmail string) (string, error) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	otpStr := ""
	num := n.Int64()
	for i := 0; i < 6; i++ {
		otpStr = string(rune('0'+(num%10))) + otpStr
		num /= 10
	}

	emailOTPs[userID] = emailOTPData{
		OTP:       otpStr,
		NewEmail:  newEmail,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	return otpStr, nil
}

func VerifyEmailOTP(userID, otp string) (string, bool) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	data, exists := emailOTPs[userID]
	if !exists {
		return "", false
	}
	if time.Now().After(data.ExpiresAt) {
		delete(emailOTPs, userID)
		return "", false
	}
	if data.OTP != otp {
		return "", false
	}
	delete(emailOTPs, userID)
	return data.NewEmail, true
}
