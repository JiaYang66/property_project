package conf

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"relocate/util/logging"
	"time"
)

type config struct {
	App      app      `yaml:"app"`
	Server   server   `yaml:"server"`
	Database database `yaml:"database"`
	Redis    redis    `yaml:"redis"`
}

type app struct {
	JwtSecret string `yaml:"jwtSecret"`
}

type server struct {
	RunMode      string        `yaml:"runMode"`
	Log          bool          `yaml:"log"`
	HttpPort     int           `yaml:"httpPort"`
	ReadTimeout  time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
}

type database struct {
	Type string `yaml:"type"`
	Url  string `yaml:"url"`
}

type redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}

var Data *config

const defaultConfigFile = "config.yaml"

// 初始化程序配置
func Setup() {
	v := viper.New()
	v.SetConfigFile(defaultConfigFile)
	if err := v.ReadInConfig(); err != nil {
		logging.Fatalf("读取%s配置文件发生错误: %s", defaultConfigFile, err)
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logging.Infof("%s配置文件发生了修改", e.Name)
		loadConfig(v)
	})
	loadConfig(v)
}

func loadConfig(v *viper.Viper) {
	Data = &config{}
	if err := v.Unmarshal(&Data); err != nil {
		logging.Fatalf("解析%s配置发生错误: %s", defaultConfigFile, err)
	}
	Data.Server.ReadTimeout *= time.Second
	Data.Server.WriteTimeout *= time.Second
}
