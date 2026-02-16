package test

import (
	"bluebell/dao"
	_ "bluebell/initialize"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMySQL(t *testing.T) {
	db := dao.NewDB()
	assert.NoError(t, db.Ping())
}
