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
	UID          string `json:"uid" gorm:"primaryKey"`
	Phone        string `json:"phone" gorm:"type:varchar(13);uniqueIndex"`
	SessionState string `json:"-" gorm:"type:longtext"`
	Secret       string `json:"-" gorm:"type:longtext"`

	Password string         `json:"password" gorm:"-"`
	Token    *UserAuthToken `json:"-" gorm:"-"`
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
	session := uuid.NewString()
	claims := &AuthJWTClaims{
		Session: session,
		UID:     ua.UID,
		Exp:     time.Now().Add(time.Minute * 30).UTC().Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := access.SignedString([]byte(secret))
	if err != nil {
		return
	}

	rclaims := &AuthJWTClaims{
		Session: session,
		UID:     ua.UID,
		Exp:     time.Now().Add(time.Hour * 24).UTC().Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, rclaims)
	refreshToken, err := refresh.SignedString([]byte(secret))
	if err != nil {
		return
	}

	ua.SessionState = session
	ua.Token = &UserAuthToken{
		Access:  accessToken,
		Refresh: refreshToken,
	}
}

// GetHTTPResponse ...
func (ua UserAuth) GetHTTPResponse() map[string]any {
	return nil
}

// AuthJWTClaims ...
type AuthJWTClaims struct {
	Session string `json:"session"`
	UID     string `json:"uid"`
	Exp     int64  `json:"exp"`

	jwt.RegisteredClaims
}

// UserAuthToken ...
type UserAuthToken struct {
	Access  string `json:"-"`
	Refresh string `json:"-"`
}

// GetHTTPResponse ...
func (at UserAuthToken) GetHTTPResponse() map[string]any {
	return map[string]any{
		"accessToken":  at.Access,
		"refreshToken": at.Refresh,
	}
}
