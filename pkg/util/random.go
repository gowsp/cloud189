package util

import (
	"math/rand"
	"regexp"
	"strconv"
)

func Random(v string) string {
	reg := regexp.MustCompilePOSIX("[xy]")
	data := reg.ReplaceAllFunc([]byte(v), func(msg []byte) []byte {
		var i int64
		t := int64(16 * rand.Float32())
		if msg[0] == 120 {
			i = t
		} else {
			i = 3&t | 8
		}
		return []byte(strconv.FormatInt(i, 16))
	})
	return string(data)
}
