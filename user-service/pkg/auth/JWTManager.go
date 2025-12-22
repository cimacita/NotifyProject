package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTData struct {
	UserID uuid.UUID `json:"user_id"`
}

type JWTManager struct {
	secret string
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: secret}
}

func (j *JWTManager) Create(data JWTData) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": data.UserID.String(),
	})
	signed, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (j *JWTManager) Parse(tokenStr string) (*JWTData, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSignMethod
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, ErrMissingUserID
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return &JWTData{UserID: userID}, nil

}
