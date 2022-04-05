package web

import (
	"fmt"
	"testing"
)

var api = NewClient("")

func TestListFile(t *testing.T) {
	f, _ := NewClient("").ListFile("-11")
	fmt.Print(f)
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
