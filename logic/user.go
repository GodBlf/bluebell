package logic

import (
	"bluebell/dao"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"bluebell/settings"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type UserLogicInterface interface {
	SignUp(p *models.ParamSignUp) error
	Login(p *models.ParamLogin) (user *models.User, err error)
	SearchUserByName(username string) (user *models.User, err error)
}
type UserLogic struct {
	userDao dao.UserDaoInterface
}

func NewUserLogic(userDao dao.UserDaoInterface) *UserLogic {
	return &UserLogic{
		userDao: userDao,
	}
}

func (u *UserLogic) SignUp(p *models.ParamSignUp) error {
	//1.dao check exist
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
	//2. dao insert user
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

func (u *UserLogic) Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	err = u.userDao.Login(user)
	if err != nil {
		asType, ok := errors.AsType[*dao.UserError](err)
		if ok {
			zap.L().Error("userDao.Login failed", zap.Error(fmt.Errorf("logic.Login failed, err: %v", asType)))
		}
		zap.L().Error("Login failed", zap.String("username", p.Username), zap.Error(err))
		return nil, err
	}
	//jwt
	token, err := jwt.GenJwtToken(user.UserID, user.Username)
	if err != nil {
		zap.L().Error("jwt.GenJwtToken failed", zap.String("username", p.Username), zap.Error(err))
		return nil, err
	}
	user.Token = token
	return
}

func (u *UserLogic) SearchUserByName(username string) (user *models.User, err error) {
	byUsername, err := u.userDao.GetUserByUsername(username)
	if err != nil {
		ok := errors.Is(err, dao.ErrorUserNotExist)
		if ok {
			zap.L().Error("user exist", zap.String("username", username))
			return nil, dao.ErrorUserNotExist
		}
		err = fmt.Errorf("userDao.GetUserByUsername failed, err: %v", err)
		zap.L().Error("userDao.GetUserByUsername failed", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return byUsername, nil
}
