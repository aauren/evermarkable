// Package api provides interfaces to interact with external services.
package api

// This package was taken almost verbatim from https://github.com/juruen/rmapi/blob/master/api/sync15/blobstorage.go - A very special thanks
// for all of @juruen's work that he did for years on the rmapi project!

import (
	"io"
	"net/http"

	"github.com/aauren/evermarkable/pkg/model"
	"k8s.io/klog/v2"
)

type BlobStorage struct {
	http *HTTPClientCtx
}

func NewBlobStorage(http *HTTPClientCtx) *BlobStorage {
	return &BlobStorage{
		http: http,
	}
}

func (b *BlobStorage) GetReader(hash string) (io.ReadCloser, error) {
	url, err := b.GetURL(hash)
	if err != nil {
		return nil, err
	}
	klog.V(2).Infof("get url: %s", url)

	blob, _, err := b.http.GetBlobStream(url)
	return blob, err
}

func (b *BlobStorage) GetURL(hash string) (string, error) {
	klog.V(2).Infof("fetching GET blob url for: %s", hash)

	urls, err := getURLProviderFromCtx(b.http)
	if err != nil {
		return "", err
	}

	var req model.BlobStorageRequest
	var res model.BlobStorageResponse
	req.Method = http.MethodGet
	req.RelativePath = hash
	if err := b.http.Post(UserBearer, urls.SyncWithPath(model.DownloadBlobPath), req, &res); err != nil {
		return "", err
	}
	return res.URL, nil
}

func (b *BlobStorage) GetRootIndex() (string, int64, error) {
	url, err := b.GetURL(model.RemRootPathName)
	if err != nil {
		return "", 0, err
	}
	klog.V(1).Infof("got root get url: %s", url)
	blob, gen, err := b.http.GetBlobStream(url)
	if err == ErrNotFound {
		return "", 0, nil

	}
	if err != nil {
		return "", 0, err
	}
	content, err := io.ReadAll(blob)
	if err != nil {
		return "", 0, err
	}
	klog.V(1).Infof("got root gen: %d", gen)
	return string(content), gen, nil

}
