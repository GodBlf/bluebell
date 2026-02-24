package test

import (
	"bluebell/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestSignHandler(t *testing.T) {
	m := &mocks.UserLogicInterface{}
	m.On("SignUp", mock.Anything).Return(nil)
}
