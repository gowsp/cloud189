package web

import (
	"github.com/gowsp/cloud189/pkg/drive"
)

func (c *api) Login(name, password string) error {
	user := &drive.User{Name: name, Password: password}
	return c.login(user)
}
