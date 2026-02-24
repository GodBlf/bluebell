package main

import (
	"bluebell/di"
	"bluebell/initialize"

	"go.uber.org/fx"
)

func main() {
	initialize.Initialize()
	app := fx.New(
		di.Provide,
		fx.Invoke(di.RegisterHook),
	)
	app.Run()
}
