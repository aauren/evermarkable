package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/aauren/evermarkable/pkg/model"
	"github.com/juruen/rmapi/api/sync15"
	"k8s.io/klog/v2"
)

func getCachedTreePath() (string, error) {
	cachedir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	rmapiFolder := path.Join(cachedir, model.EMAppName)
	err = os.MkdirAll(rmapiFolder, model.EMDefaultCacheDirModeSet)
	if err != nil {
		return "", err
	}
	cacheFile := path.Join(rmapiFolder, model.EMDTCacheFile)
	return cacheFile, nil
}

func loadCacheTree() (*sync15.HashTree, error) {
	cacheFile, err := getCachedTreePath()
	if err != nil {
		return nil, err
	}

	tree := &sync15.HashTree{}
	if _, err := os.Stat(cacheFile); err == nil {
		b, err := os.ReadFile(cacheFile)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, tree)
		if err != nil {
			klog.Warning("cache corrupt, using blank tree")
			return &sync15.HashTree{}, nil
		}
		if tree.CacheVersion != model.EMCacheVersion {
			klog.Warningf("wrong cache file version (found %d, but on %d), resync", tree.CacheVersion, model.EMCacheVersion)
			return &sync15.HashTree{}, nil
		}
	}
	klog.V(1).Infof("cache loaded: %s", cacheFile)

	return tree, nil
}

func SaveTree(tree *sync15.HashTree) error {
	cacheFile, err := getCachedTreePath()
	klog.V(1).Infof("Writing cache: %s", cacheFile)
	if err != nil {
		return err
	}
	tree.CacheVersion = model.EMCacheVersion
	b, err := json.MarshalIndent(tree, "", "")
	if err != nil {
		return err
	}
	err = os.WriteFile(cacheFile, b, model.EMDefaultCacheFileModeSet)
	return err
}

func CreateCacheTree(blobStorage BlobEMConfigHolder) (*sync15.HashTree, error) {
	config, err := blobStorage.GetEMConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get remarkable config: %v", err)
	}

	cacheTree, err := loadCacheTree()
	if err != nil {
		return nil, fmt.Errorf("could not load cache tree: %v", err)
	}

	klog.V(1).Info("Mirroring docs")

	err = cacheTree.Mirror(blobStorage.GetBlobStorage(), config.Remarkable.GetConcurrency())
	if err != nil {
		return nil, fmt.Errorf("could not Mirror tree: %v", err)
	}
	klog.V(1).Info("Tree has been mirrored")

	err = SaveTree(cacheTree)
	if err != nil {
		return nil, fmt.Errorf("could not save tree: %v", err)
	}
	klog.V(1).Info("Tree has been saved")

	return cacheTree, nil
}
