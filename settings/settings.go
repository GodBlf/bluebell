package settings

import (
	"fmt"
	"path"
	"runtime"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	RootPath string = getRootPath() + "/../"
)

func getRootPath() string {
	_, file, _, _ := runtime.Caller(0)
	return path.Dir(file)
}

//日志
type LoggerConfig struct {
	Level        string `mapstructure:"level"`
	Filename     string `mapstructure:"filename"`
	ErrorName    string `mapstructure:"error_filename"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Compress     bool   `mapstructure:"compress"`
	LogInConsole bool   `mapstructure:"log_in_console"`
}

//appconfig singleton
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LoggerConfig `mapstructure:"log"`
	*MySQLConfig  `mapstructure:"mysql"`
	*RedisConfig  `mapstructure:"redis"`
}

var GlobalConfig *AppConfig

func InitAppConfig() {
	v := viper.New()
	v.SetConfigFile(RootPath + "conf/dev.yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件发生变化: %s\n", in.Name)
	})
	a := &AppConfig{}
	err = v.Unmarshal(a)
	if err != nil {
		panic(err)
	}
	GlobalConfig = a

}

//mysql
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}
