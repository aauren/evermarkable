package api

// This package was taken almost verbatim from https://github.com/juruen/rmapi/blob/master/api/sync15/apictx.go - A very special thanks
// for all of @juruen's work that he did for years on the rmapi project!

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aauren/evermarkable/pkg/util"
	"github.com/juruen/rmapi/api/sync15"
	"k8s.io/klog/v2"
)

type Fetcher interface {
	Fetch(docID, dstPath string) error
}

type DocumentFetcher struct {
	bs *BlobStorage
	ht *sync15.HashTree
}

func NewDocumentFetcher(blobStorage *BlobStorage, hashTree *sync15.HashTree) Fetcher {
	return &DocumentFetcher{
		bs: blobStorage,
		ht: hashTree,
	}
}

func (d *DocumentFetcher) Fetch(docID, dstPath string) error {
	doc, err := d.ht.FindDoc(docID)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "rmapizip")

	if err != nil {
		return fmt.Errorf("failed to create tmpfile for zip dir: %v", err)
	}
	defer tmp.Close()

	w := zip.NewWriter(tmp)
	defer w.Close()
	for _, f := range doc.Files {
		klog.V(2).Infof("fetching document: %s", f.DocumentID)
		blobReader, err := d.bs.GetReader(f.Hash)
		if err != nil {
			return err
		}
		defer blobReader.Close()
		header := zip.FileHeader{}
		header.Name = f.DocumentID
		header.Modified = time.Now()
		zipWriter, err := w.CreateHeader(&header)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipWriter, blobReader)

		if err != nil {
			return err
		}
	}
	w.Close()
	tmpPath := tmp.Name()
	_, err = util.CopyFile(tmpPath, dstPath)

	if err != nil {
		return fmt.Errorf("failed to copy %s to %s, er: %v", tmpPath, dstPath, err)
	}

	defer os.RemoveAll(tmp.Name())

	return nil
}
