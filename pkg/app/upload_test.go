package app

import (
	"log"
	"testing"

	"github.com/gowsp/cloud189/pkg/file"
	"github.com/gowsp/cloud189/pkg/invoker"
)

func init() {
	// os.Setenv("DEBUG", "1")
	// os.Setenv("EXE_MODE", "1")
}
func TestLocal(t *testing.T) {
	api := New(invoker.DefaultPath())
	l := file.NewLocalFile("-11", "D:/tmp/1718847748.mp4")
	api.Uploader().Write(l)
}

func TestFast(t *testing.T) {
	api := New(invoker.DefaultPath())
	l := file.NewFastFile("-11", "fast://6530775E0360627193655900231FC57B:86219024/Test.mp4")
	api.Uploader().Write(l)
}
func TestNet(t *testing.T) {
	file := file.NewURLFile("-11", "http://192.168.5.30:8088/view/%E7%A7%81%E6%9C%89%E4%BA%91/job/build_4u_fmweb_backend_dk/758/artifact/install_bz/nginx.tar")
	api := New(invoker.DefaultPath())
	session := api.conf.Session
	u := Upload{session: session}
	e := u.Write(file)
	log.Println(e)
}
