package api

import (
	"fmt"
	"time"

	"github.com/aauren/evermarkable/pkg/model"
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

func getRemarkableConfigFromCtx(httpClientCtx *HTTPClientCtx) (*model.EMRemarkableConfig, error) {
	cfgRaw := httpClientCtx.Context.Value(model.ContextConfigSet)
	if cfgRaw == nil {
		return nil, fmt.Errorf("didn't find config on the HTTPClientCtx context")
	}

	cfg, ok := cfgRaw.(model.EMRootConfig)
	if !ok {
		return nil, fmt.Errorf("config stored in HTTPClientCtx context did not appear to be instance of EMRootConfig")
	}

	return &cfg.Config.Remarkable, nil
}

func getURLProviderFromCtx(httpClientCtx *HTTPClientCtx) (model.EMURLProvider, error) {
	cfgRaw := httpClientCtx.Context.Value(model.ContextConfigSet)
	if cfgRaw == nil {
		return nil, fmt.Errorf("didn't find config on the HTTPClientCtx context")
	}

	cfg, ok := cfgRaw.(model.EMRootConfig)
	if !ok {
		return nil, fmt.Errorf("config stored in HTTPClientCtx context did not appear to be instance of EMRootConfig")
	}

	return cfg.Config.Remarkable.URLs, nil
}
