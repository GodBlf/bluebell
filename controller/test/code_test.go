package test

import (
	"bluebell/controller"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMsg(t *testing.T) {
	var code controller.ResCode = 1006
	assert.Equal(t, "需要登录", code.Msg())
	var code2 controller.ResCode = 1000
	assert.Equal(t, "success", code2.Msg())
}
