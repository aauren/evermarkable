package rmcmd

import (
	"os"

	"github.com/aauren/evermarkable/cmd"
	"github.com/aauren/evermarkable/pkg/remarkable/api"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	authCommand = &cobra.Command{
		Use:   "auth",
		Short: "Authenticate to remarkable, will re-authenticate even if token is existing",
		Long:  `Authenticate to remarkable and store the valid tokens in OS's keyring`,
		Run:   AuthRun,
	}
)

//nolint:gochecknoinits // We don't need to check inits for cmd files
func init() {
	RemarkableCommand.AddCommand(authCommand)
}

func AuthRun(cobraCmd *cobra.Command, args []string) {
	ctx := cmd.InitContext()

	tokens, err := api.LoadTokens()
	if err != nil {
		klog.Errorf("could not load tokens from os keyring: %v", err)
	}
	httpClientCtx, err := api.CreateHTTPClientCtx(tokens, ctx)
	if err != nil {
		klog.Errorf("cloud not create HTTP client ctx: %v", err)
		os.Exit(1)
	}
	err = api.AuthenticateHTTP(httpClientCtx, true)
	if err != nil {
		klog.Errorf("could not authenticate HTTP: %s", err)
		os.Exit(2)
	}

	klog.Infof("Successfully authenticated to Remarkable")
}
