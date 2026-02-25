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
	controller.NewCommunityController,
))

var LogicModules = fx.Module("logic", fx.Provide(
	fx.Annotate(
		logic.NewUserLogic,
		fx.As(new(logic.UserLogicInterface)),
	),
	fx.Annotate(
		logic.NewCommunityLogic,
		fx.As(new(logic.CommunityLogicInterface)),
	),
))

var DaoModules = fx.Module("dao", fx.Provide(
	dao.NewDB,
	fx.Annotate(
		dao.NewDBUser,
		fx.As(new(dao.UserDaoInterface)),
	),
	dao.NewRedis,
	fx.Annotate(
		dao.NewCommunityDao,
		fx.As(
			new(dao.CommunityDaoInterface),
		),
	),
))

var Provide = fx.Options(
	fx.Provide(router.NewRouter),
	ControllerModules,
	LogicModules,
	DaoModules,
)
