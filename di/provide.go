package di

import (
	"bluebell/controller"
	"bluebell/dao"
	"bluebell/logic"
	"bluebell/router"

	"go.uber.org/fx"
)

// graph theory
var ControllerModules = fx.Module("controller", fx.Provide(
	controller.NewUserController,
))

var LogicModules = fx.Module("logic", fx.Provide(
	fx.Annotate(
		logic.NewUserLogic,
		fx.As(new(logic.UserLogicInterface)),
	),
))

var DaoModules = fx.Module("dao", fx.Provide(
	dao.NewDB,
	fx.Annotate(
		dao.NewDBUser,
		fx.As(new(dao.UserDaoInterface)),
	),
	dao.NewRedis,
))

var Provide = fx.Options(
	fx.Provide(router.NewRouter),
	ControllerModules,
	LogicModules,
	DaoModules,
)
