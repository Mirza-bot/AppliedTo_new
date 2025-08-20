package token

import (
	"errors"
	"maps"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret    []byte
	Issuer    string
	AccessTTL time.Duration
}

func (j *JWT) Sign(claims map[string]any) (string, error) {
	now := time.Now()

	// preallocate + copy incoming claims
	mc := make(jwt.MapClaims, len(claims)+3)
	if claims != nil {
		maps.Copy(mc, claims)
	}

	// set standard fields if missing
	if j.Issuer != "" {
		if _, ok := mc["iss"]; !ok {
			mc["iss"] = j.Issuer
		}
	}
	if _, ok := mc["iat"]; !ok {
		mc["iat"] = now.Unix()
	}
	if _, ok := mc["exp"]; !ok {
		ttl := j.AccessTTL
		if ttl <= 0 {
			ttl = 24 * time.Hour // sensible default fallback
		}
		mc["exp"] = now.Add(ttl).Unix()
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	return t.SignedString(j.Secret)
}

// Verify parses & validates an HS256 token and optionally checks issuer.
// Returns the claims if valid.
func (j *JWT) Verify(tokenStr string) (jwt.MapClaims, error) {
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	if j.Issuer != "" {
		if iss, _ := claims["iss"].(string); iss != j.Issuer {
			return nil, errors.New("invalid issuer")
		}
	}
	return claims, nil
}

