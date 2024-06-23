package api

import (
	"fmt"

	"github.com/juruen/rmapi/model"
	"k8s.io/klog/v2"
)

func GetNodeByPath(blobStorage BlobEMConfigHolder, path string) (*model.Node, error) {
	ft, err := CreateCacheTree(blobStorage)
	if err != nil {
		return nil, fmt.Errorf("could not create cache tree: %v", err)
	}

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
