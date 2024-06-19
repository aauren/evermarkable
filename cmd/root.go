package cmd

import (
	goflags "flag"

	"github.com/aauren/evermarkable/cmd/rmcmd"
	"github.com/aauren/evermarkable/pkg/cmdsupport"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var (
	RootCmd = &cobra.Command{
		Use:   "evermarkable",
		Short: "evermarkable is a simple way to sync your documents between Remarkable and Evernote",

	}
)

func Execute() error {
	return RootCmd.Execute()
}

//nolint:gochecknoinits // We don't need to check inits for cmd files
func init() {
	RootCmd.AddCommand(rmcmd.RemarkableCommand)

	RootCmd.PersistentFlags().StringVarP(&cmdsupport.Config.ConfigPath, "config-path", "c", cmdsupport.GetDefaultConfigPath(),
		"sets the config path")

	fs := goflags.NewFlagSet("", goflags.PanicOnError)
	klog.InitFlags(fs)
	RootCmd.PersistentFlags().AddGoFlagSet(fs)

	err := cmdsupport.LoadConfigFile(&cmdsupport.Config)
	if err != nil {
		klog.Fatalf("coud not load config file: %v", err)
	}
}
