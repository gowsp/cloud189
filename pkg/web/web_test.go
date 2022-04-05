package web

import (
	"fmt"
	"testing"
)

func TestListFile(t *testing.T) {
	f, _ := NewApi("").ListFile("-11")
	fmt.Print(f)
}
func TestListFolder(t *testing.T) {
	f, _ := NewApi("").ListDir("-11")
	fmt.Print(f)
}

func TestSearchFolder(t *testing.T) {
	f, err := NewApi("").FindDir("-11", "11")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(f)
}
