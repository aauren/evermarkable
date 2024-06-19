package rmcmd

import (
	"os"

	"github.com/aauren/evermarkable/pkg/cmdsupport"
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
	authClear bool
)

//nolint:gochecknoinits // We don't need to check inits for cmd files
func init() {
	RemarkableCommand.AddCommand(authCommand)

	authCommand.PersistentFlags().BoolVar(&authClear, "clear", false, "clear authentication token")
}

func AuthRun(cobraCmd *cobra.Command, args []string) {
	ctx := cmdsupport.InitContext(cmdsupport.Config)

	if authClear {
		err := api.ClearTokens()
		if err != nil {
			klog.Errorf("could not clear tokens from os keyring: %v", err)
		}
		klog.Infof("Successfully cleared tokens from Remarkable")
		os.Exit(0)
	}

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
