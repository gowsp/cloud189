package web

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"sync"

	"github.com/gowsp/cloud189-cli/pkg/util"
)

var config Config
var configPath string
var configSingleton sync.Once
var configPathSingleton sync.Once

func GetConfig() *Config {
	configSingleton.Do(func() {
		f, err := os.OpenFile(getConfigPath(), os.O_CREATE|os.O_RDWR, 0666)
		if os.IsNotExist(err) {
			NewContent().QrLogin()
			return
		}
		err = json.NewDecoder(f).Decode(&config)
		if err != nil {
			NewContent().QrLogin()
		}
	})
	return &config
}

type rsa struct {
	ResCode int32  `json:"res_code,omitempty"`
	Expire  int64  `json:"expire,omitempty"`
	PkId    string `json:"pkId,omitempty"`
	PubKey  string `json:"pubKey,omitempty"`
}

func (r *rsa) encrypt(data string) []byte {
	key := util.Key(r.PubKey)
	d, _ := util.RsaEncrypt(key, []byte(data))
	return d
}

type User struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type Config struct {
	User       User   `json:"user,omitempty"`
	RSA        rsa    `json:"rsa,omitempty"`
	SSON       string `json:"sson,omitempty"`
	Auth       string `json:"auth,omitempty"`
	SessionKey string `json:"session_key,omitempty"`
}

func getConfigPath() string {
	configPathSingleton.Do(func() {
		dir := mkdir(".config", "cloud189")
		configPath = dir + "/config.json"
	})
	return configPath
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
func getConfigFile() *os.File {
	f, err := os.OpenFile(getConfigPath(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	return f
}
func (ctx *Config) storeLoginInfo(url string, sson *http.Cookie) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sson)
	jar, _ := cookiejar.New(nil)
	api := &http.Client{Jar: jar}
	resp, _ := api.Do(req)
	cookie := findCookie(jar.Cookies(resp.Request.URL), "COOKIE_LOGIN_USER")
	config.SSON = sson.Value
	config.Auth = cookie.Value
	config.save()
}
func (config *Config) save() {
	f := getConfigFile()
	defer f.Close()
	err := json.NewEncoder(f).Encode(config)
	if err != nil {
		log.Fatalln(err)
	}
}
