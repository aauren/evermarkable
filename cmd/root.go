package cmd

import (
	"context"
	goflags "flag"
	"fmt"
	"os/user"
	"path"

	"github.com/aauren/evermarkable/pkg/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	config  = model.EMRootConfig{}
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
	currentUser, err := user.Current()
	if err != nil {
		klog.Fatal("could not get the current user executing script")
	}
	defaultConfigPath := path.Join(currentUser.HomeDir, ".config", "evermarkable", "config.yaml")

	RootCmd.PersistentFlags().StringVarP(&config.ConfigPath, "config-path", "c", defaultConfigPath,
		"sets the config path")

	fs := goflags.NewFlagSet("", goflags.PanicOnError)
	klog.InitFlags(fs)
	RootCmd.PersistentFlags().AddGoFlagSet(fs)

	err = loadConfigFile()
	if err != nil {
		klog.Fatalf("coud not load config file: %v", err)
	}
}

func loadConfigFile() error {
	if config.ConfigPath == "" {
		return fmt.Errorf("config file path should either be the default or set, not blank")
	}

	viper.AddConfigPath(config.ConfigPath)
	viper.SetEnvPrefix(model.EnvPrefix)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			klog.Infof("No config file was found at %s, skipping...", config.ConfigPath)
		} else {
			return fmt.Errorf("there was an error reading the config: %v", err)
		}
	} else {
		err = viper.Unmarshal(&config.Config)
		if err != nil {
			return fmt.Errorf("there was an error unmarshalling the config: %v", err)
		}
	}

	return nil
}

func InitContext() context.Context {
	ctx := context.Background()

	ctx = context.WithValue(ctx, model.ContextConfigSet, config)

	return ctx
}
