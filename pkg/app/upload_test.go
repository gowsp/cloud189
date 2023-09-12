package app

import (
	"testing"

	"github.com/gowsp/cloud189/pkg/file"
	"github.com/gowsp/cloud189/pkg/invoker"
)

// func TestUpload(t *testing.T) {
// 	api := NewApi(invoker.DefaultPath())
// 	session := api.conf.Session
// 	u := Upload{session: session}
// 	f, err := os.Open("d:/tmp/test01.png")
// 	if err != nil {
// 		return
// 	}
// 	// api.List()
// 	local := NewLocalFile("-11", f)
// 	u.prepare(local)
// }

//	func TestFS(t *testing.T) {
//		f, _ := os.Open("d:/tmp/test.xlsx")
//		local := NewLocalFile("-11", f)
//		api := NewApi(invoker.DefaultPath())
//		session := api.conf.Session
//		u := Upload{session: session}
//		u.Upload(local)
//	}
func TestNet(t *testing.T) {
	file := file.NewURLFile("-11", "http://192.168.5.30:8088/view/%E7%A7%81%E6%9C%89%E4%BA%91/job/build_4u_fmweb_backend_dk/758/artifact/install_bz/nginx.tar")
	api := New(invoker.DefaultPath())
	session := api.conf.Session
	u := Upload{session: session}
	u.Write(file)
}
