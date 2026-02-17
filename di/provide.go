package di

import "go.uber.org/fx"

var Provides = fx.Options(
	dbModule,
	controllerModule,
	logicModule,
)
