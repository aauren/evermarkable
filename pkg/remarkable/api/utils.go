package api

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func (a *authToken) Missing() bool {
	return a.deviceToken == "" || a.userToken == ""
}

func (a *authToken) IsExpired() (bool, error) {
	userJWT := UserTokenJWT{}
	_, _, err := (&jwt.Parser{}).ParseUnverified(a.userToken, &userJWT)

	if err != nil {
		return false, fmt.Errorf("cannot parse token: %v", err)
	}

	if userJWT.VerifyExpiresAt(time.Now().Unix(), false) {
		return false, nil
	}

	return true, nil
}
