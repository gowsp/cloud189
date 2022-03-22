package web

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/gowsp/cloud189/pkg/file"
)

type space struct {
	Available uint64 `xml:"available,omitempty"`
	Capacity  uint64 `xml:"capacity,omitempty"`
}

func (s *space) size() uint64 {
	return s.Capacity
}
func (s *space) available() uint64 {
	return s.Available
}
func (s *space) used() uint64 {
	return s.Capacity - s.Available
}
func (s *space) usedPercent() float64 {
	return float64(s.used()*100) / float64(s.Capacity)
}

func (client *Client) Df() {
	resp, err := client.api.Get("https://cloud.189.cn/api/open/user/getUserInfoForPortal.action")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if client.isInvalidSession(data) {
		client.Df()
		return
	}
	var space space
	xml.Unmarshal(data, &space)
	fmt.Printf("%-12s%-12s%-12s%s\n", "Size", "Used", "Avail", "Use%")
	fmt.Printf("%-12s%-12s%-12s%.2f%%\n",
		file.ReadableSize(space.size()),
		file.ReadableSize(space.used()),
		file.ReadableSize(space.available()),
		space.usedPercent(),
	)
}
