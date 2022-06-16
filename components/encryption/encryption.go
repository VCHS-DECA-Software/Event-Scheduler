package encryption

import (
	"crypto/rand"
	"fmt"
	"main/components/db"
	"math"
	"math/big"
	"time"

	"github.com/asdine/storm"
	"github.com/golang-jwt/jwt"

	"golang.org/x/crypto/bcrypt"
)

var signingKey []byte

type AuthClaim struct {
	Key string `json:"key"`
	jwt.StandardClaims
}

type Key struct {
	Key      string `storm:"id"`
	IssuedAt int64
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IssueToken() (string, string, error) {
	k, err := rand.Int(rand.Reader, big.NewInt(int64(math.Pow(2, 32))))
	if err != nil {
		return "", "", err
	}

	key := Key{
		Key:      k.String(),
		IssuedAt: time.Now().Unix(),
	}

	err = db.Save(&key)
	if err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaim{
		Key: key.Key,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: key.IssuedAt,
		},
	})

	signed, err := token.SignedString(signingKey)
	return signed, key.Key, nil
}

func ParseToken(token string) (AuthClaim, error) {
	validator, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return AuthClaim{}, err
	}
	return validator.Claims.(AuthClaim), nil
}

func CheckToken(token string) (bool, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return false, err
	}

	key, err := db.Get[AuthClaim](claims.Key)
	if err == storm.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if key.IssuedAt != claims.IssuedAt {
		return false, nil
	}

	return true, nil
}
