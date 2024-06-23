package rmcmd

import (
	"fmt"
	"os"

	"github.com/aauren/evermarkable/pkg/cmdsupport"
	"github.com/aauren/evermarkable/pkg/remarkable/api"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	lsCommand = &cobra.Command{
		Use:   "ls",
		Short: "List items from remarkable API",
		Long:  `Lists various items from remarkable API`,
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run:   LSRun,
	}
)

//nolint:gochecknoinits // We don't need to check inits for cmd files
func init() {
	RemarkableCommand.AddCommand(lsCommand)
}

func LSRun(cobraCmd *cobra.Command, args []string) {
	path := args[0]
	ctx := cmdsupport.InitContext(cmdsupport.Config)

	httpClientCtx, err := api.EnsureAuthenticated(ctx)
	if err != nil {
		klog.Errorf("could not ensure authenticated: %v", err)
		os.Exit(1)
	}

	blobStorage := api.NewBlobStorage(httpClientCtx)

	node, err := api.GetNodeByPath(blobStorage, path)
	if err != nil {
		klog.Errorf("could not get node by path: %v", err)
		os.Exit(2)
	}

	klog.V(1).Info("Node by path found")

	if node.IsFile() {
		fmt.Printf("node is file: %s", node.Name())
	}

	for _, e := range node.Children {
		eType := "d"
		if e.IsFile() {
			eType = "f"
		}
		fmt.Printf("[%s]\t%s\n", eType, e.Name())
	}
}
