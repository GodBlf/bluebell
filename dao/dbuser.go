package dao

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const secret = "liwenzhou.com"

type UserDao struct {
	mysqlClient *sqlx.DB
}

func NewDBUser(mysqlClient *sqlx.DB) *UserDao {
	return &UserDao{
		mysqlClient: mysqlClient,
	}
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
