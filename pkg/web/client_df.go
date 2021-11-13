package web

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/gowsp/cloud189-cli/pkg/file"
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
	fmt.Println("Size\t\tUsed\t\tAvail\t\tUse%")
	fmt.Printf("%s\t\t%s\t\t%s\t\t%.2f%%\n",
		file.Readable(space.size()),
		file.Readable(space.used()),
		file.Readable(space.available()),
		space.usedPercent(),
	)
}
