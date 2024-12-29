package drive

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/app"
	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/gowsp/cloud189/pkg/util"
)

func init() {
	// os.Setenv("189_MODE", "1")
	os.Setenv("189_MODE", "0")
}

func TestDrive(t *testing.T) {
	api := app.New(invoker.DefaultPath())
	f := New(api)
	info, err := f.Stat("/demo/1")
	if err != nil {
		return
	}
	f.Stat("/demo/2")
	f.Stat("/demo/2")
	fmt.Println(info)
}
func TestWalk(t *testing.T) {
	api := app.New(invoker.DefaultPath())
	f := New(api)
	f.ReadDir("/demo")
	f.ReadDir("/demo")
}
func TestDecode(t *testing.T) {
	data := "F961D48A546BFEFAFC5C17B7D8024A56B3DBC64AF1FA980A0E827D524C0760370F255258EF9F89E524A4BA5274434F46BD6D1E25C47CCF9410CA05C2C10A29B60D0D1B119BF871960A0C78B8177670D6ACEDFE20E9C801201AF66858EBAF910AE00207AFC92897043A82DB19204F0FD3357054406579A88FB4FFCBA51FD1905C503EC7B344864408DCC6BE3593E54E2CB46BADC8757651296D4D8D9B2DC2B9E7093F02E6D8B3C64D7F7097F0FDEBE27FCCFEA190DAB9AFDF3DFF3DB14D89ABED08347ED0310DCF14627641BDA5F0E4CD304C1670D64587F45FC1FF15DDF80FC8"
	v := util.DecryptAES([]byte("C8CAB983D32137EE5F076F204B21BBCD"[0:16]), strings.ToLower(data))
	fmt.Println(v)
}
func TestDelete(t *testing.T) {
	api := app.New(invoker.DefaultPath())
	f := New(api)
	err := f.Delete("/demo1")
	fmt.Println(err)
}
func TestMakeDir(t *testing.T) {
	api := app.New(invoker.DefaultPath())
	f := New(api)
	err := f.Mkdir("/demo1")
	fmt.Println(err)
}
func TestUpload(t *testing.T) {
	api := app.New(invoker.DefaultPath())
	fs := New(api)
	cfg := pkg.UploadConfig{Num: 3}
	// fs.Upload(cfg, "/home", "D:/repo/go/src/github.com/gowsp/cloud189/docs/html")
	// fs.Upload(cfg, "/", "D:/tmp/01.txt")
	fs.Upload(cfg, "/", "../")
	// f, _ := os.Open("drive_test.go")
	// l := file.NewLocalFile("-11", f)
	// api.Uploader().Write(l)
	// api.Uploader().Write(file.NewNetFile("-11", "demo.txt", 111))
	// api.Uploader().Write(file.NewFastFile("-11", "fast://13598043FF94C80EFC76642E10D2121C:1372/demo_file.txt"))
}

func TestDownload(t *testing.T) {
	api := app.New(invoker.DefaultPath())
	f := New(api)
	err := f.Download("D:/", "/demo/page.html")
	if err != nil {
		log.Println(err)
	}
}
