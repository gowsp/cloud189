package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/file"
	"github.com/gowsp/cloud189/pkg/invoker"
)

func init() {
	os.Setenv("189_MODE", "1")
}
func TestLogin(t *testing.T) {
	New(invoker.DefaultPath()).PwdLogin("xxxxxxx", "xxxxxxxxxxx")
}
func TestQrLogin(t *testing.T) {
	New(invoker.DefaultPath()).QrLogin()
}
func TestSpace(t *testing.T) {
	space, _ := New(invoker.DefaultPath()).Space()
	fmt.Println(space.Available, space.Capacity)
}
func TestSign(t *testing.T) {
	New(invoker.DefaultPath()).Sign()
}
func TestListFile(t *testing.T) {
	f, _ := New(invoker.DefaultPath()).List(file.Root, pkg.FILE)
	fmt.Println(f)
}
func TestListDir(t *testing.T) {
	f, _ := New(invoker.DefaultPath()).List(file.Root, pkg.DIR)
	fmt.Println(f)
}
func TestSearchFile(t *testing.T) {
	f, _ := New(invoker.DefaultPath()).Search(file.Root, pkg.FILE, "1")
	fmt.Println(f)
}
func TestSearchDir(t *testing.T) {
	f, _ := New(invoker.DefaultPath()).Search(file.Root, pkg.DIR, "我")
	fmt.Println(f)
}
func TestMkdir(t *testing.T) {
	New(invoker.DefaultPath()).Mkdir(file.Root, "/demo/1/2/3")
}
func TestDelete(t *testing.T) {
	api := New(invoker.DefaultPath())
	dir, _ := api.Search(file.Root, pkg.DIR, "demo")
	api.Delete(dir...)
}
func TestCopy(t *testing.T) {
	api := New(invoker.DefaultPath())
	f, _ := api.Mkdir(file.Root, "/demo/1/2/3")
	api.Copy(file.Root, f)
}
func TestRename(t *testing.T) {
	api := New(invoker.DefaultPath())
	demo, _ := api.Search(file.Root, pkg.DIR, "demo")
	api.Rename(demo[0], "demo")
}
func TestMove(t *testing.T) {
	api := New(invoker.DefaultPath())
	f, _ := api.Mkdir(file.Root, "/demo/1/2/3")
	api.Move(file.Root, f)
}
func TestGetFolder(t *testing.T) {
	api := New(invoker.DefaultPath())
	f, _ := api.Stat("/demo/1/2/3")
	if f.IsDir() {
		api.DirUsage(f)
	}
}
func TestGetFile(t *testing.T) {
	api := New(invoker.DefaultPath())
	api.Stat("/我的图片")
}
