package drive

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gowsp/cloud189/pkg/util"
)

type User struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}
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

type Config struct {
	path       string
	User       User      `json:"user,omitempty"`
	RSA        RsaConfig `json:"rsa,omitempty"`
	SSON       string    `json:"sson,omitempty"`
	Auth       string    `json:"auth,omitempty"`
	SessionKey string    `json:"session_key,omitempty"`
}

func defaultPath() string {
	dir := mkdir(".config", "cloud189")
	return dir + "/config.json"
}
func mkdir(dirs ...string) string {
	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	for _, dir := range dirs {
		path = path + "/" + dir
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
	return path
}
func OpenConfig(path string) (*Config, error) {
	if path == "" {
		path = defaultPath()
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	config.path = path
	return &config, nil
}
func (config *Config) Save() error {
	f, err := os.OpenFile(config.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(config)
}
