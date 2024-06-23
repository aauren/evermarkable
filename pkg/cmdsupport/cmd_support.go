package cmdsupport

import (
	"context"
	"fmt"
	"os/user"
	"path"

	"github.com/aauren/evermarkable/pkg/model"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

var (
	Config = model.EMRootConfig{}
)

func GetDefaultConfigPath() string {
	currentUser, err := user.Current()
	if err != nil {
		klog.Fatal("could not get the current user executing script")
	}

	return path.Join(currentUser.HomeDir, ".config", "evermarkable", "config.yaml")
}

func InitContext(config model.EMRootConfig) context.Context {
	ctx := context.Background()

	ctx = context.WithValue(ctx, model.ContextConfigSet, config)

	return ctx
}

func LoadConfigFile(config *model.EMRootConfig) error {
	if config.ConfigPath == "" {
		return fmt.Errorf("config file path should either be the default or set, not blank")
	}

	viper.AddConfigPath(config.ConfigPath)
	viper.SetEnvPrefix(model.EnvPrefix)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			klog.V(1).Infof("No config file was found at %s, skipping...", config.ConfigPath)
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
