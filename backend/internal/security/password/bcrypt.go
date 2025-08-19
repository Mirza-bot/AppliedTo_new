package password

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	Hash(pw string) (string, error)
	Verify(hash, pw string) bool
	NeedsRehash(hash string) bool
}

// ---- Options pattern ----

type Option func(*bcryptHasher)

func WithCost(cost int) Option {
	return func(b *bcryptHasher) {
		if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
			cost = bcrypt.DefaultCost
		}
		b.cost = cost
	}
}

// (Optional) If you want to always prehash long passwords:
func WithPrehashLong(enabled bool) Option {
	return func(b *bcryptHasher) { b.prehashLong = enabled }
}

// ---- Implementation ----

type bcryptHasher struct {
	cost        int
	prehashLong bool // default true to handle >72 byte passwords
}

func NewBcrypt(opts ...Option) Hasher {
	h := &bcryptHasher{
		cost:        bcrypt.DefaultCost,
		prehashLong: true,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (b *bcryptHasher) Hash(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}
	pw := []byte(password)
	if b.prehashLong && len(pw) > 72 {
		pw = prehash(pw)
	}
	hash, err := bcrypt.GenerateFromPassword(pw, b.cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (b *bcryptHasher) Verify(hash, password string) bool {
	pw := []byte(password)
	if b.prehashLong && len(pw) > 72 {
		pw = prehash(pw)
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), pw) == nil
}

func (b *bcryptHasher) NeedsRehash(hash string) bool {
	c, err := bcrypt.Cost([]byte(hash))
	return err == nil && c < b.cost
}

// ---- helpers ----

func prehash(b []byte) []byte {
	sum := sha256.Sum256(b)
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum[:])
	return dst
}

