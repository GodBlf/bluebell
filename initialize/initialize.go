package initialize

import (
	"bluebell/logger"
	"bluebell/settings"
)

func Initialize() {
	settings.InitAppConfig()
	logger.InitLogger(settings.GlobalConfig.LoggerConfig)

}
