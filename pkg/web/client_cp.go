package web

import "log"

func (client *Client) Cp(paths ...string) {
	if len(paths) < 2 {
		log.Fatalln("receive at least two parameters")
	}
	client.runTask(COPY, paths...)
}
