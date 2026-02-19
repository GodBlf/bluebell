package dao

type UserError struct {
	name string
}

func (e *UserError) Error() string {
	return e.name
}

//var (
//	ErrorUserExist       = errors.New("用户已存在")
//	ErrorUserNotExist    = errors.New("用户不存在")
//	ErrorInvalidPassword = errors.New("用户名或密码错误")
//	ErrorInvalidID       = errors.New("无效的ID")
//)

var (
	ErrorUserExist       = &UserError{name: "用户已存在"}
	ErrorUserNotExist    = &UserError{name: "用户不存在"}
	ErrorInvalidPassword = &UserError{name: "用户名或密码错误"}
	ErrorInvalidID       = &UserError{name: "无效的ID"}
)
