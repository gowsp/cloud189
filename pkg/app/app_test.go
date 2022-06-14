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
func TestListFile(t *testing.T) {
	f, _ := NewApi(invoker.DefaultPath()).ListFile("-11")
	fmt.Println(f)
}
func TestListDir(t *testing.T) {
	f, _ := NewApi(invoker.DefaultPath()).ListDir("-11")
	fmt.Println(f)
}
func TestFind(t *testing.T) {
	f, _ := NewApi(invoker.DefaultPath()).Find("-11", "data")
	fmt.Println(f)
}
func TestFindDir(t *testing.T) {
	f, _ := NewApi(invoker.DefaultPath()).FindDir("-11", "我的文档")
	fmt.Println(f)
}
func TestMkdir(t *testing.T) {
	NewApi(invoker.DefaultPath()).Mkdir("-11", "demo/test", true)
}
func TestMkdirs(t *testing.T) {
	NewApi(invoker.DefaultPath()).Mkdirs("-11", "demo/test", "demo/test2")
}
func TestDelete(t *testing.T) {
	api := NewApi(invoker.DefaultPath())
	demo, _ := api.Find("-11", "demo")
	f, _ := api.ListFile(demo.Id())
	api.Delete(f...)
}
func TestCopy(t *testing.T) {
	api := NewApi(invoker.DefaultPath())
	demo, _ := api.Find("-11", "demo")
	f, _ := api.ListFile(demo.Id())
	api.Copy("-11", f...)
}
func TestRename(t *testing.T) {
	api := NewApi(invoker.DefaultPath())
	demo, _ := api.Find("-11", "demo")
	files, _ := api.ListFile(demo.Id())
	for _, v := range files {
		api.Rename(v, "1"+v.Name())
	}
}
func TestMove(t *testing.T) {
	api := NewApi(invoker.DefaultPath())
	demo, _ := api.Find("-11", "demo")
	files, _ := api.ListFile(demo.Id())
	api.Move("-11", files...)
}
func TestDownload(t *testing.T) {
	api := NewApi(invoker.DefaultPath())
	demo, _ := api.Find("-11", "demo")
	files, _ := api.ListFile(demo.Id())
	resp, err := api.Download(files[0], 0)
	fmt.Println(resp.StatusCode, err)
}
