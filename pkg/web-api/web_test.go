package web

import (
	"fmt"
	"testing"

	"github.com/gowsp/cloud189/pkg"
)

var api = NewClient("")

func TestListFile(t *testing.T) {
	f, _ := NewClient("").ListFile("-11")
	fmt.Print(f)
}
func TestGetFile(t *testing.T) {
	f, _ := api.Detail("31442115697812348")
	ext := f.Sys().(pkg.FileExt)
	fmt.Println(ext)
	fmt.Println(ext.DownloadUrl)
	resp, _ := api.invoker.http.Get(ext.DownloadUrl)
	fmt.Println(resp.Request.URL)
	fmt.Println(resp.StatusCode)
}
func TestListFolder(t *testing.T) {
	f, _ := NewClient("").ListDir("-11")
	fmt.Print(f)
}

func TestSearchFolder(t *testing.T) {
	f, err := NewClient("").FindDir("-11", "11")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(f)
}
