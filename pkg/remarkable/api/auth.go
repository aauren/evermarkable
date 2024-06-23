package api

// This package was taken almost verbatim from https://github.com/juruen/rmapi/blob/master/api/auth.go - A very special thanks
// for all of @juruen's work that he did for years on the rmapi project!

import (
	"context"
	"fmt"

	"github.com/aauren/evermarkable/internal/keyring"
	"github.com/aauren/evermarkable/pkg/cli"
	"github.com/aauren/evermarkable/pkg/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"k8s.io/klog/v2"
)

type authToken struct {
	deviceToken string
	userToken   string
}

type deviceTokenRequest struct {
	Code       string `json:"code"`
	DeviceDesc string `json:"deviceDesc"`
	DeviceID   string `json:"deviceID"`
}

type UserTokenJWT struct {
	Auth0 struct {
		UserID string
		Email  string
	} `json:"auth0-profile"`
	Scopes string
	*jwt.StandardClaims
}

func AuthenticateHTTP(httpClientCtx *HTTPClientCtx, reAuth bool) error {
	token := httpClientCtx.Token

	if token.deviceToken == "" {
		err := refreshDeviceToken(httpClientCtx, token)
		if err != nil {
			return err
		}
	}

	if token.userToken == "" || reAuth {
		err := refreshUserToken(httpClientCtx, token)
		if err != nil {
			return err
		}
	}

	return nil
}

func ClearTokens() error {
	err := keyring.DeleteSecretFromStore(model.DeviceTokenSecName)
	if err != nil {
		return fmt.Errorf("could not delete device token from store: %v", err)
	}

	err = keyring.DeleteSecretFromStore(model.UserTokenSecName)
	if err != nil {
		return fmt.Errorf("could not delete user token from store: %v", err)
	}

	return nil

}

func LoadTokens() (*authToken, error) {
	token := authToken{}

	devTok, err := keyring.GetSecretFromStore(model.DeviceTokenSecName)
	if err != nil {
		if keyring.ErrorIsNotFound(err) {
			klog.Infof("device token not found in keyring")
		} else {
			return &token, fmt.Errorf("could not retrieve device token from store: %v", err)
		}
	}

	userTok, err := keyring.GetSecretFromStore(model.UserTokenSecName)
	if err != nil {
		if keyring.ErrorIsNotFound(err) {
			klog.Infof("user token not found in keyring")
		} else {
			return &token, fmt.Errorf("could not retrieve user token from store: %v", err)
		}
	}

	token.deviceToken = devTok
	token.userToken = userTok
	return &token, nil
}

// EnsureAuthenticated ensures that the tokens are loaded, valid, and not expired.
// If the tokens are missing or expired, it re-authenticates.
func EnsureAuthenticated(ctx context.Context) (*HTTPClientCtx, error) {
	tokens, err := LoadTokens()
	if err != nil {
		klog.Errorf("could not load tokens from os keyring: %v", err)
		return nil, err
	}

	httpClientCtx, err := CreateHTTPClientCtx(tokens, ctx)
	if err != nil {
		klog.Errorf("could not create HTTP client ctx: %v", err)
		return nil, err
	}

	expired, err := tokens.IsExpired()
	if err != nil {
		klog.Errorf("could not check if token is expired: %v", err)
		return nil, err
	}

	if tokens.Missing() || expired {
		klog.Infof("Tokens are missing or expired, re-authenticating")
		err = AuthenticateHTTP(httpClientCtx, true)
		if err != nil {
			klog.Errorf("could not authenticate HTTP: %v", err)
			return nil, err
		}
	}

	return httpClientCtx, nil
}

func newDeviceToken(http *HTTPClientCtx, code string) (string, error) {
	uuid := uuid.New()

	req := deviceTokenRequest{code, model.DefaultDeviceDesc, uuid.String()}

	urlProv, err := http.GetURLProvider()
	if err != nil {
		return "", err
	}

	resp := BodyString{}
	err = http.Post(EmptyBearer, urlProv.AuthWithPath(model.DeviceTokenPath), req, &resp)

	if err != nil {
		return "", fmt.Errorf("unable to create device token: %v", err)
	}

	return resp.Content, nil
}

func newUserToken(http *HTTPClientCtx) (string, error) {
	urlProv, err := http.GetURLProvider()
	if err != nil {
		return "", err
	}

	resp := BodyString{}
	err = http.Post(DeviceBearer, urlProv.AuthWithPath(model.UserDevicePath), nil, &resp)

	if err != nil {
		return "", err
	}

	return resp.Content, nil
}

func refreshDeviceToken(httpClientCtx *HTTPClientCtx, token *authToken) error {
	devToken, err := newDeviceToken(httpClientCtx, cli.ReadCode())

	if err != nil {
		return err
	}

	klog.V(1).Infof("device token obtained: %v", devToken)

	token.deviceToken = devToken

	return keyring.SaveSecretInStore(model.DeviceTokenSecName, devToken)
}

func refreshUserToken(httpClientCtx *HTTPClientCtx, token *authToken) error {
	userToken, err := newUserToken(httpClientCtx)

	if err != nil {
		return err
	}

	klog.V(1).Infof("user token obtained: %v", userToken)

	token.userToken = userToken

	return keyring.SaveSecretInStore(model.UserTokenSecName, userToken)
}
