package main

import (
	"bluebell/di"
	_ "bluebell/initialize"

	"go.uber.org/fx"
)

func main() {
	var fxConfig = fx.Options(fx.Invoke(di.RegisterHook))
	app := fx.New(fxConfig)
	app.Run()
}
