package rmcmd

import (
	"fmt"
	"os"
	"path"

	"github.com/aauren/evermarkable/pkg/cmdsupport"
	"github.com/aauren/evermarkable/pkg/remarkable/api"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	downCommand = &cobra.Command{
		Use:   "down",
		Short: "Download a file from remarkable",
		Long:  `Download a file from remarkable in PDF form in the current directory`,
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
		Run:   RemDownRun,
	}
)

//nolint:gochecknoinits // We don't need to check inits for cmd files
func init() {
	RemarkableCommand.AddCommand(downCommand)
}

func RemDownRun(cobraCmd *cobra.Command, args []string) {
	srcPath := args[0]
	var dstPath string
	if len(args) > 1 {
		dstPath = args[1]
	} else {
		dstPath = "."
	}
	ctx := cmdsupport.InitContext(cmdsupport.Config)

	httpClientCtx, err := api.EnsureAuthenticated(ctx)
	if err != nil {
		klog.Errorf("could not ensure authenticated: %v", err)
		os.Exit(1)
	}

	blobStorage := api.NewBlobStorage(httpClientCtx)

	ct, err := api.CreateCacheTree(blobStorage)
	if err != nil {
		klog.Errorf("could not create cache tree: %v", err)
		os.Exit(2)
	}

	node, err := api.GetNodeByPathFromCache(ct, srcPath)
	if err != nil {
		klog.Errorf("could not get node by path: %v", err)
		os.Exit(3)
	}

	klog.V(1).Info("Node by path found")

	if node.IsDirectory() {
		fmt.Printf("node is a directory and cannot be downloaded: %s", node.Name())
	}

	fmt.Printf("Downloading file: %s... ", node.Name())

	fetcher := api.NewDocumentFetcher(blobStorage, ct)

	err = fetcher.Fetch(node.Document.ID, path.Join(dstPath, fmt.Sprintf("%s.zip", node.Name())))
	if err != nil {
		fmt.Println()
		klog.Errorf("could not fetch document: %v", err)
		os.Exit(4)
	}

	fmt.Println("Complete!")
}
