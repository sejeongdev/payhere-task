package model

import (
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserAuth ...
type UserAuth struct {
	UID    string `json:"uid" gorm:"primaryKey"`
	Phone  string `json:"phone" gorm:"type:varchar(13)"`
	Secret string `json:"-" gorm:"type:longtext"`

	Token    string `json:"-" gorm:"-"`
	Password string `json:"password" gorm:"-"`
}

var phoneRegex string = "01[016789]{1}[- .]?[0-9]{3,4}[- .]?[0-9]{4}([^0-9]+|$)"

// Validate ...
func (ua UserAuth) Validate() bool {
	return (ua.Phone != "" && ua.validatePhone()) && ua.Password != ""
}

func (ua UserAuth) validatePhone() bool {
	reg := regexp.MustCompile(phoneRegex)
	return reg.MatchString(ua.Phone)
}

// InitRegister ...
func (ua *UserAuth) InitRegister() {
	ua.setSecret()
	ua.UID = uuid.NewString()
}

func (ua *UserAuth) setSecret() (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(ua.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	ua.Secret = string(bytes)
	return nil
}

// LoginValidate ...
func (ua UserAuth) LoginValidate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(ua.Secret), []byte(password))
	return err == nil
}

// SetToken ...
func (ua *UserAuth) SetToken(secret string) {
	claims := &AuthJWTClaims{
		UID:    ua.UID,
		Expire: time.Now().Add(time.Minute * 30).UTC(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return
	}
	ua.Token = t
}

// GetHTTPResponse ...
func (ua UserAuth) GetHTTPResponse() map[string]any {
	return map[string]any{
		"token": ua.Token,
	}
}

// AuthJWTClaims ...
type AuthJWTClaims struct {
	UID    string    `json:"uid"`
	Expire time.Time `json:"expire"`

	jwt.RegisteredClaims
}
