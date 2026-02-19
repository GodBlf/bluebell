package logic

import (
	"bluebell/dao"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"bluebell/settings"
	"fmt"

	"go.uber.org/zap"
)

type UserLogic struct {
	userDao *dao.UserDao
}

func NewUserLogic(userDao *dao.UserDao) *UserLogic {
	return &UserLogic{
		userDao: userDao,
	}
}

func (u *UserLogic) SignUp(p *models.ParamSignUp) error {
	exist, err := u.userDao.CheckUserExist(p.Username)
	if err != nil {
		err = fmt.Errorf("userDao.CheckUserExist failed, err: %v", err)
		zap.L().Error("userDao.CheckUserExist failed", zap.Error(err))
		return err
	}
	if !exist {
		zap.L().Info("user exist", zap.String("username", p.Username))
		return dao.ErrorUserExist
	}
	newSnowflake, err := snowflake.NewSnowflake(settings.GlobalConfig.StartTime, settings.GlobalConfig.MachineID)
	if err != nil {
		err = fmt.Errorf("snowflake.NewSnowflake failed, err: %v", err)
		zap.L().Error("snowflake.NewSnowflake failed", zap.Error(err))
		return err
	}
	id := newSnowflake.GetID()
	user := &models.User{
		UserID:   id,
		Username: p.Username,
		Password: p.Password,
	}
	err = u.userDao.InsertUser(user)
	if err != nil {
		err = fmt.Errorf("userDao.InsertUser failed, err: %v", err)
		zap.L().Error("userDao.InsertUser failed", zap.Error(err))
		return err
	}
	return nil
}
