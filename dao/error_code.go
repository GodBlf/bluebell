package dao

type UserError struct {
	name string
}

func NewUserError(name string) *UserError {
	return &UserError{name: name}
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
	ErrorUserExist       = NewUserError("用户已存在")
	ErrorUserNotExist    = NewUserError("用户不存在")
	ErrorInvalidPassword = NewUserError("用户名或密码错误")
	ErrorInvalidID       = NewUserError("无效的ID")
)
