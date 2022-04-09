package web

import (
	"github.com/gowsp/cloud189/pkg/invoker"
)

func (c *api) Login(name, password string) error {
	user := &invoker.User{Name: name, Password: password}
	return c.login(user)
}
