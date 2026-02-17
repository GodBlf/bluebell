package di

import (
	"bluebell/dao"

	"go.uber.org/fx"
)

var dbModule = fx.Module("dbModule",
	fx.Provide(
		dao.NewDB,
		dao.NewRedis,
	),
)

var controllerModule = fx.Module("controllerModule",
	fx.Provide(),
)

var logicModule = fx.Module("logicModule",
	fx.Provide(),
)
