package api

import (
	"fmt"

	"github.com/juruen/rmapi/api/sync15"
	"github.com/juruen/rmapi/model"
	"k8s.io/klog/v2"
)

func GetNodeByPath(blobStorage BlobEMConfigHolder, path string) (*model.Node, error) {
	ct, err := CreateCacheTree(blobStorage)
	if err != nil {
		return nil, fmt.Errorf("could not create cache tree: %v", err)
	}

	return GetNodeByPathFromCache(ct, path)
}

func GetNodeByPathFromCache(ct *sync15.HashTree, path string) (*model.Node, error) {
	ft := sync15.DocumentsFileTree(ct)
	klog.V(1).Info("File Tree obtained")

	klog.V(1).Infof("looking at path: %s", path)

	rootNode := model.CreateNode(model.Document{
		ID:           "",
		Type:         "CollectionType",
		VissibleName: "/",
	})

	klog.V(1).Info("Getting node by path")

	node, err := ft.NodeByPath(path, &rootNode)
	if err != nil {
		return nil, fmt.Errorf("could not find node: %v", err)
	}

	return node, nil
}
