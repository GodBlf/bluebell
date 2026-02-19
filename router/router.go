package router

import (
	"bluebell/controller"
	"bluebell/settings"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewRouter(lc fx.Lifecycle, userController *controller.UserController) *gin.Engine {
	r := gin.Default()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", settings.GlobalConfig.Port),
		Handler: r,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					zap.L().Panic("srv.ListenAndServe() failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				zap.L().Error("srv.Shutdown() failed", zap.Error(err))
			}
			return err

		},
	})

	return r
}
