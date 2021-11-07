package web

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"testing"

	"github.com/gowsp/cloud189-cli/pkg/util"
)

func TestQrCodeLogin(t *testing.T) {
	NewContent().QrLogin()
}
func TestPwdLogin(t *testing.T) {
	NewContent().PwdLogin("xxxx", "xxxxxxxx")
}
func TestRefresh(t *testing.T) {
	GetClient().refresh()
}
func TestMkdir(t *testing.T) {
	GetClient().Mkdir("/demo", "/demo/1/2", "/demo/1/3")
}
func TestUp(t *testing.T) {
	GetClient().Up("/demo/1/2", "../../README.md", "../../LICENSE")
}
func TestLs(t *testing.T) {
	GetClient().Ls("/demo/1/2")
}
func TestReadir(t *testing.T) {
	client := GetClient()
	data, _ := client.Stat("/demo/1/2")
	client.Readdir(data.Id(), 0)
}
func TestDownFile(t *testing.T) {
	GetClient().Dl("d:/", "/demo/")
}
func TestDownDir(t *testing.T) {
	GetClient().Dl("d:/", "/demo/")
}
func TestCp(t *testing.T) {
	GetClient().Cp("/demo/1/2", "/demo")
}
func TestMv(t *testing.T) {
	GetClient().Mv("/demo/1/3", "/demo")
}
func TestRm(t *testing.T) {
	GetClient().Rm("/demo/1/2", "/demo")
}

func TestSign(t *testing.T) {
	GetClient().Sign()
}

func TestEncrypt(t *testing.T) {
	l := "d76d88e5631e49cbad1d09ec543735a"
	plaintext := "parentFolderId=-11&fileName=source.tar.xz&fileSize=30708604&sliceSize=10485760&lazyCheck=1"
	v := util.AesEncrypt([]byte(plaintext), []byte(l[:16]))
	p := hex.EncodeToString(v)
	log.Println(p)
	sha := "SessionKey=38633e41-267c-4b02-8ce2-6bcc5aff8b97&Operate=GET&RequestURI=/person/initMultiUpload&Date=1635478138485&params=702c8c4e2b1715417e8e7239bac9e74309f4bcd4120848b2000bbf832d23cdaa61da2651169311d28091e9da2eeb2662c5dbedeba6e6f1df7751d7ef653c27b026a15807a7a4db47ce8d9f705adc74cfdfef91a3beb16582aae06954be3afe79"
	v = util.SHA1(sha, l)
	log.Print(hex.EncodeToString(v))
	key := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCaIaXc6/wgzhS/2PCKinROMWCXp7Kiv+TDpb2kN3rpTC0xthk9y2uZoT5Lz78DL7+CDg+cS4G/1yPUbSHjIOFeglIRqe6mF2mY5sqhKLyBLQwzniAz4B8Y74BZ7OFTftRna43njDGRfUNxH1qiLuuKPiPCqzYHTiko4p5wszjF6QIDAQAB"
	pub_key := util.Key(key)
	v, _ = util.RsaEncrypt(pub_key, []byte(l))
	log.Print(base64.StdEncoding.EncodeToString(v))
}
