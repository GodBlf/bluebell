package dao

import (
	"bluebell/settings"
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewDB(lc fx.Lifecycle) *sqlx.DB {
	mySQLConfig := settings.GlobalConfig.MySQLConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", mySQLConfig.User, mySQLConfig.Password, mySQLConfig.Host, mySQLConfig.Port, mySQLConfig.DB)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Panic("连接数据库失败", zap.Error(err))
		panic(err)
	}
	db.SetMaxOpenConns(mySQLConfig.MaxOpenConns)
	db.SetMaxIdleConns(mySQLConfig.MaxIdleConns)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			err = db.PingContext(ctx)
			if err != nil {
				zap.L().Panic("无法连接到数据库", zap.Error(err))
				return err
			}
			zap.L().Debug("成功连接到数据库")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			zap.L().Debug("正在关闭数据库连接")
			err = db.Close()
			if err != nil {
				zap.L().Error("关闭数据库连接失败", zap.Error(err))
				return err
			}
			zap.L().Debug("数据库连接已关闭")
			return nil
		},
	})
	return db
}
