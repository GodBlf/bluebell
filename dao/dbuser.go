package dao

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const secret = "liwenzhou.com"

type UserDaoInterface interface {
	CheckUserExist(username string) (bool, error)
	InsertUser(user *models.User) error
	Login(user *models.User) (err error)
	GetUserByUsername(username string) (user *models.User, err error)
}

type UserDao struct {
	mysqlClient *sqlx.DB
}

func NewDBUser(mysqlClient *sqlx.DB) *UserDao {
	return &UserDao{
		mysqlClient: mysqlClient,
	}
}

func (db *UserDao) Login(user *models.User) (err error) {
	//1. judge exist

	opwd := user.Password
	err = db.mysqlClient.Get(user, "select user_id,username,password from user where username=?", user.Username)
	if err != nil {
		ok := errors.Is(sql.ErrNoRows, err)
		if ok {
			return ErrorUserNotExist
		}
		return err
	}
	pwd := encryptPassword(opwd)
	if pwd != user.Password {
		zap.L().Error("用户登录密码错误", zap.String("username", user.Username))
		return ErrorInvalidPassword
	}
	return nil

}

func (db *UserDao) CheckUserExist(username string) (bool, error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.mysqlClient.Get(&count, sqlStr, username); err != nil {
		zap.L().Error("查询用户是否存在失败", zap.String("username", username), zap.Error(err))
		return false, err
	}
	if count > 0 {
		zap.L().Info("用户已存在", zap.String("username", username))
		return false, ErrorUserExist

	}
	return true, nil
}

// InsertUser 想数据库中插入一条新的用户记录
func (db *UserDao) InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	// 执行SQL语句入库
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.mysqlClient.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func (db *UserDao) GetUserById(uid int64) (user *models.User, err error) {
	user = &models.User{}
	err = db.mysqlClient.Get(user, "select user_id... from user where user_id=?", uid)
	if err != nil {
		zap.L().Error("根据id查询用户失败", zap.Int64("uid", uid), zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (db *UserDao) GetUserByUsername(username string) (user *models.User, err error) {
	user = &models.User{}
	err = db.mysqlClient.Get(user, "select user_id... from user where username like ?", username)
	if err != nil {
		zap.L().Error("根据用户名查询用户失败", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	if user == nil {
		zap.L().Error("根据用户名查询用户失败，用户不存在", zap.String("username", username))
		return nil, ErrorUserNotExist
	}
	return user, nil

}
