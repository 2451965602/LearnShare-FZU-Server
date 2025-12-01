package config

import (
	log "LearnShare/pkg/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	Mysql        *mySQL
	Redis        *redis
	Oss          *oss
	Smtp         *smtp
	Verify       *verify
	Server       *server
	Turnstile    *turnstile
	Logger       *logger
	Cors         *cors
	runtimeViper = viper.New()
)

// Init 目的是初始化配置管理器
func Init() {
	configPath := "./config/config.yaml"

	runtimeViper.SetConfigFile(configPath)
	runtimeViper.SetConfigType("yaml")

	if err := runtimeViper.ReadInConfig(); err != nil {
		log.Fatal("config.Init: 未找到配置文件")
	}

	configMapping()

	// 监听配置热更新
	runtimeViper.OnConfigChange(func(e fsnotify.Event) {
		// 一些场景无法确定在配置变更回调时是否已经初始化，所以此处还要再判断一次
		log.Infof("config: notice config changed: %v\n", e.String())
		configMapping() // 重新映射配置
	})
	runtimeViper.WatchConfig()
}

// configMapping 用于将配置映射到全局变量
func configMapping() {
	c := new(config)
	if err := runtimeViper.Unmarshal(&c); err != nil {
		// 如果配置文件损坏，或者热更新时再次失败，此处需要及时记录日志
		log.Fatalf("config.configMapping: 配置反序列化失败: %v", err)
	}
	Mysql = &c.MySQL
	Redis = &c.Redis
	Oss = &c.OSS
	Smtp = &c.Smtp
	Verify = &c.Verify
	Server = &c.Server
	Turnstile = &c.Turnstile
	Logger = &c.Logger
	Cors = &c.Cors
}
