package web

import "log"

func (client *Client) Mv(paths ...string) {
	if len(paths) < 2 {
		log.Fatalln("receive at least two parameters")
	}
	client.runTask(MOVE, paths...)
}
