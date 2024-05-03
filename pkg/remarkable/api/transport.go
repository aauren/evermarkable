package api

// This package was taken almost verbatim from https://github.com/juruen/rmapi/blob/master/transport/transport.go - A very special thanks
// for all of @juruen's work that he did for years on the rmapi project!

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"github.com/aauren/evermarkable/pkg/model"
	"github.com/aauren/evermarkable/pkg/util"
	"k8s.io/klog/v2"
)

type AuthType int

type BodyString struct {
	Content string
}

var (
	ErrWrongGeneration = errors.New("wrong generation")
	ErrNotFound        = errors.New("not found")

	ErrUnauthorized = errors.New("401 Unauthorized")
	ErrConflict     = errors.New("409 Conflict")
)

const (
	EmptyBearer AuthType = iota
	DeviceBearer
	UserBearer

	HeaderGeneration = "x-goog-generation"
	EmptyBody        = ""
)

type HTTPClientCtx struct {
	Client  *http.Client
	Token   *authToken
	Context context.Context
}

func CreateHTTPClientCtx(token *authToken, ctx context.Context) (*HTTPClientCtx, error) {
	var httpClient = &http.Client{Timeout: 5 * 60 * time.Second}

	if cfg := ctx.Value(model.ContextConfigSet); cfg == nil {
		return nil, fmt.Errorf("config was not found within the context, this shouldn't happen, breaking early")
	} else if _, ok := cfg.(model.EMRootConfig); !ok {
		return nil, fmt.Errorf("config on context did not appear to be an instance of EMRootConfig, breaking early")
	}

	return &HTTPClientCtx{httpClient, token, ctx}, nil
}

func (ctx HTTPClientCtx) addAuthorization(req *http.Request, authType AuthType) {
	var header string

	switch authType {
	case EmptyBearer:
		header = "Bearer"
	case DeviceBearer:
		header = fmt.Sprintf("Bearer %s", ctx.Token.deviceToken)
	case UserBearer:
		header = fmt.Sprintf("Bearer %s", ctx.Token.userToken)
	}

	req.Header.Add("Authorization", header)
}

func (ctx HTTPClientCtx) Get(authType AuthType, url string, body interface{}, target interface{}) error {
	bodyReader, err := util.ToIOReader(body)

	if err != nil {
		klog.Errorf("failed to serialize body: %v", err)
		return err
	}

	response, err := ctx.Request(authType, http.MethodGet, url, bodyReader)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return err
	}

	return json.NewDecoder(response.Body).Decode(target)
}

func (ctx HTTPClientCtx) GetStream(authType AuthType, url string) (io.ReadCloser, error) {
	response, err := ctx.Request(authType, http.MethodGet, url, strings.NewReader(""))
	defer func() {
		_ = response.Body.Close()
	}()

	var respBody io.ReadCloser
	if response != nil {
		respBody = response.Body
	}

	return respBody, err
}

func (ctx HTTPClientCtx) Post(authType AuthType, url string, reqBody, resp interface{}) error {
	return ctx.httpRawReq(authType, http.MethodPost, url, reqBody, resp)
}

func (ctx HTTPClientCtx) Put(authType AuthType, url string, reqBody, resp interface{}) error {
	return ctx.httpRawReq(authType, http.MethodPut, url, reqBody, resp)
}

func (ctx HTTPClientCtx) PutStream(authType AuthType, url string, reqBody io.Reader) error {
	return ctx.httpRawReq(authType, http.MethodPut, url, reqBody, nil)
}

func (ctx HTTPClientCtx) Delete(authType AuthType, url string, reqBody, resp interface{}) error {
	return ctx.httpRawReq(authType, http.MethodDelete, url, reqBody, resp)
}

func (ctx HTTPClientCtx) httpRawReq(authType AuthType, verb, url string, reqBody, resp interface{}) error {
	var contentBody io.Reader

	switch r := reqBody.(type) {
	case io.Reader:
		contentBody = r
	default:
		c, err := util.ToIOReader(reqBody)

		if err != nil {
			klog.Errorf("failed to serialize body: %v", err)
			return nil
		}

		contentBody = c
	}

	response, err := ctx.Request(authType, verb, url, contentBody)

	if response != nil {
		defer response.Body.Close()
	}

	if err != nil {
		return err
	}

	// We want to ingore the response
	if resp == nil {
		return nil
	}

	switch r := resp.(type) {
	case *BodyString:
		bodyContent, err := io.ReadAll(response.Body)

		if err != nil {
			return err
		}

		r.Content = string(bodyContent)
	default:
		err := json.NewDecoder(response.Body).Decode(resp)

		if err != nil {
			klog.Errorf("failed to deserialize body (%s), due to: %v", response.Body, err)
			return err
		}
	}
	return nil
}

func (ctx HTTPClientCtx) Request(authType AuthType, verb, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(context.Background(), verb, url, body)
	if err != nil {
		return nil, err
	}

	ctx.addAuthorization(request, authType)
	request.Header.Add("User-Agent", model.EMUserAgent)

	if klog.V(2).Enabled() {
		drequest, err := httputil.DumpRequest(request, true)
		klog.V(2).Infof("request: %s %v", drequest, err)
	}

	response, err := ctx.Client.Do(request)

	if err != nil {
		klog.Errorf("http request failed with: %v", err)
		return nil, err
	}

	if klog.V(2).Enabled() {
		defer response.Body.Close()
		dresponse, err := httputil.DumpResponse(response, true)
		klog.V(2).Infof("%s %v", dresponse, err)
	}

	if response.StatusCode != http.StatusOK {
		klog.V(1).Infof("request failed with status %d", response.StatusCode)
	}

	switch response.StatusCode {
	case http.StatusOK:
		return response, nil
	case http.StatusUnauthorized:
		return response, ErrUnauthorized
	case http.StatusConflict:
		return response, ErrConflict
	default:
		return response, fmt.Errorf("request failed with status %d", response.StatusCode)
	}
}

func (ctx HTTPClientCtx) GetBlobStream(url string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return nil, 0, err
	}
	if response.StatusCode == http.StatusNotFound {
		return nil, 0, ErrNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("GetBlobStream, status code not ok %d", response.StatusCode)
	}
	var gen int64
	if response.Header != nil {
		genh := response.Header.Get(HeaderGeneration)
		if genh != "" {
			klog.V(2).Infof("got generation header: %v", genh)
			gen, err = strconv.ParseInt(genh, 10, 64)
		}
	}

	return response.Body, gen, err
}
