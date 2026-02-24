package test

import (
	"testing"

	"go.uber.org/fx"
)

func TestTmp(t *testing.T) {

}

type tree struct {
	c1 *cy1
	c2 *cy2
	c3 *cy3
}

func Newtree(cc1 *cy1, cc2 *cy2, cc3 *cy3) *tree {
	return &tree{
		c1: cc1,
		c2: cc2,
		c3: cc3,
	}
}

var provied = fx.Provide(
	Newtree,
	Newc1y,
)

type cy1 struct {
}

func Newc1y() *cy1 {
	return &cy1{}
}

type cy2 struct {
}

type cy3 struct {
}
