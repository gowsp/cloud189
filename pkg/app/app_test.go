package app

import (
	"fmt"
	"testing"

	"github.com/gowsp/cloud189/pkg/invoker"
)

func TestLogin(t *testing.T) {
	NewApi(invoker.DefaultPath()).PwdLogin("xxxxxxx", "xxxxxxxxxxx")
}
func TestSpace(t *testing.T) {
	space, _ := NewApi(invoker.DefaultPath()).Space()
	fmt.Println(space.Available, space.Capacity)
}
func TestSign(t *testing.T) {
	NewApi(invoker.DefaultPath()).Sign()
}
