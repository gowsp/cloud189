package invoker

import (
	"encoding/json"
	"io"
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

type Session struct {
	LoginName    string `json:"loginName,omitempty"`
	Key          string `json:"sessionKey,omitempty"`
	Secret       string `json:"sessionSecret,omitempty"`
	KeepAlive    int    `json:"keepAlive,omitempty"`
	FileDiffSpan int    `json:"getFileDiffSpan,omitempty"`
	UserInfoSpan int    `json:"getUserInfoSpan,omitempty"`
	FamilyKey    string `json:"familySessionKey,omitempty"`
	FamilySecret string `json:"familySessionSecret,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

func (s *Session) Merge(new Session) {
	s.Key = new.Key
	s.Secret = new.Secret
	s.LoginName = new.LoginName
	s.KeepAlive = new.KeepAlive
	s.FileDiffSpan = new.FileDiffSpan
	s.UserInfoSpan = new.UserInfoSpan
	s.FamilyKey = new.FamilyKey
	s.FamilySecret = new.FamilySecret
}

func (s *Session) Login() bool {
	return s != nil && s.AccessToken != ""
}
func (s *Session) Empty() bool {
	return s == nil || s.Key == "" || s.Secret == ""
}

type Config struct {
	path    string
	User    *User     `json:"user,omitempty"`
	RSA     RsaConfig `json:"rsa,omitempty"`
	SSON    string    `json:"sson,omitempty"`
	Auth    string    `json:"auth,omitempty"`
	Session *Session  `json:"session,omitempty"`
}

func DefaultPath() string {
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
	file := path
	if file == "" {
		file = DefaultPath()
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var config Config
	err = json.NewDecoder(f).Decode(&config)
	if err == io.EOF {
		config = Config{path: path}
		return &config, nil
	} else if err != nil {
		return nil, err
	}
	config.path = path
	return &config, nil
}
func (config *Config) Save() error {
	if config.path == "" {
		return nil
	}
	f, err := os.OpenFile(config.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(config)
}
