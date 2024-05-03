package rmcmd

import (
	"github.com/aauren/evermarkable/cmd"
	"github.com/spf13/cobra"
)

var (
	RemarkableCommand = &cobra.Command{
		Use:   "rem",
		Short: "Suite of commands pertaining to remarkable",
	}
)

//nolint:gochecknoinits // We don't need to check inits for cmd files
func init() {
	cmd.RootCmd.AddCommand(RemarkableCommand)
}
