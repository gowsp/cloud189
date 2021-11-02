package web

import "log"

func (client *Client) Rm(paths ...string) {
	if len(paths) == 0 {
		log.Fatalln("one argument must be received")
	}
	client.runTask(DELETE, paths...)
}
