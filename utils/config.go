package utils

import (
	"github.com/spf13/viper"
)

type ConfigCloudStore struct {
	ServerDbPath string
	BasePath     string
}

func GetConfiguration() ConfigCloudStore {
	return ConfigCloudStore{
		ServerDbPath: viper.GetString("server.Config.BasePath"),
		BasePath:     viper.GetString("server.Config.BasePath"),
	}
}
