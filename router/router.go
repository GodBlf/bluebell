package router

import (
	"bluebell/controller"
	"bluebell/middleware"
	"bluebell/settings"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func routerGroup(r *gin.Engine, userController *controller.UserController) {
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.Use(middleware.Auth())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	v1 := r.Group("/api/v1")
	// 注册
	v1.POST("/signup", userController.SignUpHandler())
	// 登录
	v1.POST("/login", userController.LoginHandler())

}

func NewRouter(lc fx.Lifecycle, userController *controller.UserController) *gin.Engine {
	r := gin.Default()
	routerGroup(r, userController)
	addr := fmt.Sprintf(":%d", settings.GlobalConfig.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	var ln net.Listener
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				zap.L().Error("failed to bind server address", zap.String("addr", addr), zap.Error(err))
				return err
			}
			ln = listener

			go func() {
				if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
					zap.L().Error("srv.Serve() failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if ln == nil {
				return nil
			}

			err := srv.Shutdown(ctx)
			if err != nil {
				zap.L().Error("srv.Shutdown() failed", zap.Error(err))
			}
			return err

		},
	})

	return r
}
