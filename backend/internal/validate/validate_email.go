package validate

import (
	"errors"
	"net/mail"
	"strings"

	"golang.org/x/net/idna"
)

var (
	ErrInvalid = errors.New("invalid email")
	ErrTooLong = errors.New("email too long")
)

const (
	maxEmailLen = 254
	maxLocalLen = 64
)

func NormalizeAndValidateEmail(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", ErrInvalid
	}

	if addr, err := mail.ParseAddress(s); err == nil && addr.Address != "" {
		s = addr.Address
	}

	at := strings.LastIndexByte(s, '@')
	if at <= 0 || at == len(s)-1 {
		return "", ErrInvalid
	}
	local := s[:at]
	domain := s[at+1:]

	if len(local) > maxLocalLen || len(s) > maxEmailLen {
		return "", ErrTooLong
	}
	if strings.ContainsAny(local, " \t\r\n") || strings.ContainsAny(domain, " \t\r\n") {
		return "", ErrInvalid
	}

	asciiDomain, err := idna.Lookup.ToASCII(domain)
	if err != nil || asciiDomain == "" {
		return "", ErrInvalid
	}

	return strings.ToLower(local + "@" + asciiDomain), nil
}

func IsValidEmail(raw string) bool {
	_, err := NormalizeAndValidateEmail(raw)
	return err == nil
}
