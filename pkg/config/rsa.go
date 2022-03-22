package config

import "github.com/gowsp/cloud189/pkg/util"

type RsaConfig struct {
	ResCode int32  `json:"res_code,omitempty"`
	Expire  int64  `json:"expire,omitempty"`
	PkId    string `json:"pkId,omitempty"`
	PubKey  string `json:"pubKey,omitempty"`
}

func (r *RsaConfig) Encrypt(data string) []byte {
	key := util.Key(r.PubKey)
	d, _ := util.RsaEncrypt(key, []byte(data))
	return d
}
