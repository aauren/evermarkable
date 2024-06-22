package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// This package was taken almost verbatim from:
// - https://github.com/juruen/rmapi/blob/master/model/document.go
// - https://github.com/juruen/rmapi/blob/master/model/node.go
// A very special thanks for all of @juruen's work that he did for years on the rmapi project!

type Node struct {
	Document *Document
	Children map[string]*Node
	Parent   *Node
}

func CreateNode(document Document) Node {
	return Node{&document, make(map[string]*Node, 0), nil}
}

func (node *Node) Name() string {
	return node.Document.VissibleName
}

func (node *Node) ID() string {
	return node.Document.ID
}
func (node *Node) Version() int {
	return node.Document.Version
}

func (node *Node) IsRoot() bool {
	return node.ID() == ""
}

func (node *Node) IsDirectory() bool {
	return node.Document.Type == "CollectionType"
}

func (node *Node) IsFile() bool {
	return !node.IsDirectory()
}

func (node *Node) EntyExists(id string) bool {
	_, ok := node.Children[id]
	return ok
}

func (node *Node) LastModified() (time.Time, error) {
	return time.Parse(time.RFC3339Nano, node.Document.ModifiedClient)
}

func (node *Node) FindByName(name string) (*Node, error) {
	for _, n := range node.Children {
		if n.Name() == name {
			return n, nil
		}
	}
	return nil, errors.New("entry doesn't exist")
}

type Document struct {
	ID                string
	Version           int
	Message           string
	Success           bool
	BlobURLGet        string
	BlobURLGetExpires string
	ModifiedClient    string
	Type              string
	VissibleName      string
	CurrentPage       int
	Bookmarked        bool
	Parent            string
}

type MetadataDocument struct {
	ID             string
	Parent         string
	VissibleName   string
	Type           string
	Version        int
	ModifiedClient string
}

type DeleteDocument struct {
	ID      string
	Version int
}

type UploadDocumentRequest struct {
	ID      string
	Type    string
	Version int
}

type UploadDocumentResponse struct {
	ID                string
	Version           int
	Message           string
	Success           bool
	BlobURLPut        string
	BlobURLPutExpires string
}

type BlobRootStorageRequest struct {
	Method       string `json:"http_method"`
	Initial      bool   `json:"initial_sync,omitempty"`
	RelativePath string `json:"relative_path"`
	RootSchema   string `json:"root_schema,omitempty"`
	Generation   int64  `json:"generation"`
}

// BlobStorageRequest request
type BlobStorageRequest struct {
	Method       string `json:"http_method"`
	Initial      bool   `json:"initial_sync,omitempty"`
	RelativePath string `json:"relative_path"`
	ParentPath   string `json:"parent_path,omitempty"`
}

// BlobStorageResponse response
type BlobStorageResponse struct {
	Expires            string `json:"expires"`
	Method             string `json:"method"`
	RelativePath       string `json:"relative_path"`
	URL                string `json:"url"`
	MaxUploadSizeBytes int64  `json:"maxuploadsize_bytes,omitempty"`
}

// SyncCompleteRequest payload of the sync completion
type SyncCompletedRequest struct {
	Generation int64 `json:"generation"`
}

func CreateDirDocument(parent, name string) MetadataDocument {
	id := uuid.New()

	return MetadataDocument{
		ID:             id.String(),
		Parent:         parent,
		VissibleName:   name,
		Type:           DirectoryType,
		Version:        1,
		ModifiedClient: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func CreateUploadDocumentRequest(id string, entryType string) UploadDocumentRequest {
	if id == "" {
		newID := uuid.New()

		id = newID.String()
	}

	return UploadDocumentRequest{
		id,
		entryType,
		1,
	}
}

func CreateUploadDocumentMeta(id string, entryType, parent, name string) MetadataDocument {

	return MetadataDocument{
		ID:             id,
		Parent:         parent,
		VissibleName:   name,
		Type:           entryType,
		Version:        1,
		ModifiedClient: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func (meta MetadataDocument) ToDocument() Document {
	return Document{
		ID:             meta.ID,
		Parent:         meta.Parent,
		VissibleName:   meta.VissibleName,
		Type:           meta.Type,
		Version:        1,
		ModifiedClient: meta.ModifiedClient,
	}
}

func (doc Document) ToMetaDocument() MetadataDocument {
	return MetadataDocument{
		ID:             doc.ID,
		Parent:         doc.Parent,
		VissibleName:   doc.VissibleName,
		Type:           doc.Type,
		Version:        doc.Version,
		ModifiedClient: time.Now().UTC().Format(time.RFC3339Nano),
	}
}

func (doc Document) ToDeleteDocument() DeleteDocument {
	return DeleteDocument{
		ID:      doc.ID,
		Version: doc.Version,
	}
}
