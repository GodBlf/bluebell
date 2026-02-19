package initialize

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/settings"
)

func Initialize() {
	settings.InitAppConfig()
	logger.InitLogger(settings.GlobalConfig.LoggerConfig)
	controller.InitTrans("zh")
}
