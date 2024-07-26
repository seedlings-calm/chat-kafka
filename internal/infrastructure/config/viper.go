package config

import (
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var globalConf Config

type Config struct {
	Config BaseViper
	L      sync.RWMutex
}

// 配置来决定项目运行模式

func Setup(cfgPath string) {

	// 设置配置文件
	v := viper.New()
	v.SetConfigFile(cfgPath) // 配置文件名
	v.SetConfigType("yml")   // 配置文件类型
	v.AddConfigPath(".")     // 查找配置文件的路径，可以是绝对路径或相对路径

	// 检查配置文件是否存在
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("Config file '%s' does not exist", cfgPath)
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	var conf BaseViper
	if err = v.Unmarshal(&conf); err != nil {
		panic(err)
	}
	setGlobalConf(conf)
}

func setGlobalConf(cfg BaseViper) {
	globalConf.L.Lock()
	globalConf.Config = cfg
	globalConf.L.Unlock()
}

func GetGlobalConf() BaseViper {
	globalConf.L.RLock()
	defer globalConf.L.RUnlock()
	return globalConf.Config
}
