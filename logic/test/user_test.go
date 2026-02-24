package test

import (
	"bluebell/dao"
	"bluebell/logic"
	"bluebell/mocks"
	"bluebell/models"
	"bluebell/settings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignUp(t *testing.T) {
	//todo:mock
	settings.GlobalConfig = &settings.AppConfig{
		StartTime: "2024-01-01 00:00:00",
		MachineID: 1,
	}

	t.Run("exist user", func(t *testing.T) {
		m := &mocks.UserDaoInterface{}
		m.On("CheckUserExist", "test").Return(false, nil)
		userLogic := logic.NewUserLogic(m)
		err := userLogic.SignUp(&models.ParamSignUp{
			Username: "test",
			Password: "123456",
		})
		assert.ErrorIs(t, err, dao.ErrorUserExist, "expected error to be ErrorUserExist")
	})

	t.Run("insert user", func(t *testing.T) {
		m := &mocks.UserDaoInterface{}
		m.On("CheckUserExist", "xiaoming").Return(true, nil)
		m.On("InsertUser", mock.MatchedBy(func(user *models.User) bool {
			return user.Username == "xiaoming" &&
				user.Password == "123456"
		})).Return(nil)
		userLogic := logic.NewUserLogic(m)
		err := userLogic.SignUp(&models.ParamSignUp{
			Username: "xiaoming",
			Password: "123456",
		})
		assert.NoError(t, err, "expected no error when signing up with a new username")
	})

}

func TestSearchUser(t *testing.T) {
	dao := &mocks.UserDaoInterface{}
	dao.On("GetUserByUsername", "test").Return(&models.User{
		UserID:   1,
		Username: "test"}, nil)
	ul := logic.NewUserLogic(dao)
	user, err := ul.SearchUserByName("test")
	assert.NoError(t, err)
	assert.Equal(t, user,
		&models.User{
			UserID:   1,
			Username: "test",
		},
		"expected user to be returned with correct username and userID",
	)

}
