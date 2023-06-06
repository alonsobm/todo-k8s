package user

import (
	"errors"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/chacha20"
	"time"
)

type AuthService struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewAuthService(simmetriKey string) (*AuthService, error) {

	if len(simmetriKey) != chacha20.KeySize {
		return nil, errors.New("invalid key size")
	}

	return &AuthService{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(simmetriKey),
	}, nil

}

func (pm *AuthService) CreateToken(userID int, username string, duration time.Duration) (string, error) {
	payload := NewPayload(userID, username, duration)

	encrypted, err := pm.paseto.Encrypt(pm.symmetricKey, payload, nil)
	if err != nil {
		return "", err
	}

	return encrypted, nil
}

func (pm *AuthService) VerifyToken(token string) (*Payload, error) {

	payload := Payload{}
	err := pm.paseto.Decrypt(token, pm.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	if time.Now().After(payload.ExpiredAt) {
		return nil, errors.New("not valid token")
	}

	return &payload, nil
}

type Payload struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"user_name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userID int, username string, duration time.Duration) *Payload {
	return &Payload{
		Username:  username,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
}

func (pm *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error hashing password")
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not
func (pm *AuthService) CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
