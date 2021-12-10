package config

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/gowsp/cloud189-cli/pkg/util"
)

var config Config
var configPath string

// init config file
func InitConfigFile(config string) {
	if config == "" {
		configPath = defaultPath()
		return
	}
	_, err := os.Stat(config)
	if err != nil {
		log.Printf("Config file \"%s\" not found - using defaults", config)
		configPath = defaultPath()
		return
	}
	configPath = config
}

// open cloud189 config file
func OpenConfig() (*Config, error) {
	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type User struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
}

type Config struct {
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

func (c *Config) SsonLogin(url string, sson *http.Cookie) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(sson)
	jar, _ := cookiejar.New(nil)
	api := &http.Client{Jar: jar}
	resp, _ := api.Do(req)
	cookie := util.FindCookie(jar.Cookies(resp.Request.URL), "COOKIE_LOGIN_USER")
	c.SSON = sson.Value
	c.Auth = cookie.Value
}

func (config *Config) Save() {
	f, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		log.Fatalln(err)
	}
}
