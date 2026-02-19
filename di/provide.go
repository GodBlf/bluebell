package di

import (
	"bluebell/controller"
	"bluebell/dao"
	"bluebell/logic"

	"go.uber.org/fx"
)

var ControllerModules = fx.Module("controller", fx.Provide(
	controller.NewUserController,
))

var LogicModules = fx.Module("logic", fx.Provide(
	logic.NewUserLogic,
))

var DaoModules = fx.Module("dao", fx.Provide(
	dao.NewDBUser,
	dao.NewRedis,
))

var Provide = fx.Options(
	ControllerModules,
)
