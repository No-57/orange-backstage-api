package model

import (
	"errors"
	"fmt"
	"orange-backstage-api/infra/util/convert"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Email          string    `gorm:"column:email;uniqueIndex;NOT NULL"`
	Name           string    `gorm:"column:name;uniqueIndex;NOT NULL"`
	HashedPassword []byte    `gorm:"column:hashed_password;NOT NULL"`
	CreatedAt      time.Time `gorm:"column:created_at;NOT NULL"`
	UpdatedAt      time.Time `gorm:"column:updated_at;NOT NULL"`
}

type Claims struct {
	jwt.StandardClaims
}

func (a Account) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword(
		a.HashedPassword, convert.StrToBytes(password),
	) == nil
}

func (a Account) GenJWT(secret []byte, expiredTime time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(expiredTime).Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Subject:   strconv.FormatUint(a.ID, 10),
			Id:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

var (
	ErrExpired = errors.New("token is expired")
	ErrParse   = errors.New("token parse error")
)

func ParseTokenWithSecret(secret []byte, token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			switch {
			case ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0:
				return nil, ErrExpired
			default:
				return nil, ErrParse
			}
		}

		return nil, err
	}

	claims, ok := tokenClaims.Claims.(*Claims)
	if ok && tokenClaims.Valid {
		return claims, nil
	}

	return nil, ErrParse
}
