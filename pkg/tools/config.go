package tools

import (
	"github.com/spf13/viper"
)

func InitConfig(filePath string, bc interface{}) error {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(&bc)
}
