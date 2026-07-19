package services

import (
	"encoding/base32"
	"strings"
	"testing"
)

func TestGenerateTOTPSecret(t *testing.T) {
	secret, err := GenerateTOTPSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(secret) != 32 {
		t.Errorf("expected secret length 32, got %d", len(secret))
	}

	_, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		t.Errorf("secret is not a valid base32 string: %v", err)
	}
	if strings.ToUpper(secret) != secret {
		t.Errorf("expected secret to be uppercase")
	}
}

func TestGenerateRecoveryCodes(t *testing.T) {
	count := 5
	codes, err := GenerateRecoveryCodes(count)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(codes) != count {
		t.Errorf("expected %d codes, got %d", count, len(codes))
	}

	for _, code := range codes {
		if len(code) != 9 {
			t.Errorf("expected code length 9, got %d for code %s", len(code), code)
		}
		if code[4] != '-' {
			t.Errorf("expected code to contain '-' at index 4, got %c", code[4])
		}
		if strings.ToUpper(code) != code {
			t.Errorf("expected code to be uppercase: %s", code)
		}
	}

	if count > 1 && codes[0] == codes[1] {
		t.Errorf("expected codes to be randomly generated, but found duplicates")
	}
}
