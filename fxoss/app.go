package fxoss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-delve/delve/pkg/config"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/super1-chen/fxoss/conf"
)

type oss struct {
	User, Password, Host, SshUser, SshPassword string
	Conf                                       *conf.Config
	HttpClient                                 *http.Client
}

var (
	UserKey    = "FXOSS_USER"
	HostKey    = "FXOSS_HOST"
	PWDKey     = "FXOSS_PWD"
	SSHUserKey = "FXOSS_SSH_USER"
	SSHPWDKey  = "FXOSS_SSH_PWD"
	ConfDirKey = "FXOSS_DIR"
)

func New(now time.Time) (oss *oss, err error) {

	user, err := envLookUp(UserKey)
	if err != nil {
		return
	}

	host, err := envLookUp(HostKey)
	if err != nil {
		return
	}

	pwd, err := envLookUp(PWDKey)
	if err != nil {
		return
	}

	sshUser, err := envLookUp(SSHUserKey)
	if err != nil {
		return
	}

	sshPwd, err := envLookUp(SSHPWDKey)
	if err != nil {
		return
	}
	confDir := os.Getenv(ConfDirKey)

	if confDir == "" {
		confDir = "/tmp"
	}

	tokenPath := path.Join(confDir, conf.TokenJson)

	oss.User = user
	oss.Host = host
	oss.Password = pwd
	oss.SshUser = sshUser
	oss.SshPassword = sshPwd
	oss.HttpClient = &http.Client{Timeout: time.Minute} // set default client timeout 60s

	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		if err = oss.UpdateToken(tokenPath); err != nil {
			return nil, err
		}
	}

	conf, err := conf.NewConfig(tokenPath)

	if err != nil {
		return nil, err
	}

	if !conf.IsValid(oss.Host, now) {
		oss.UpdateToken(tokenPath)

	} else {
		oss.Conf = conf
	}

	return

}

func (oss *oss) get(api string) ([]byte, error) {

	return nil, nil
}

func (oss *oss) post(api string, body io.Reader)([]byte, error) {
	return nil, nil

}

// UpdateToken will download a new token from oss and make a new config of itself
func (oss *oss) UpdateToken(folderName string) error {
	c := new(conf.Config)
	b, err := oss.GetNewToken()
	if err != nil {
		return fmt.Errorf("get new token failed %s", err)
	}
	r := bytes.NewReader(b)
	c.Update(r)
	c.Host = oss.Host
	if err := c.Save(folderName, conf.TokenJson); err != nil {
		log.Printf("save token failed %s ", path.Join(folderName, conf.TokenJson))
	}
	oss.Conf = c
	return nil
}

func (oss *oss) GetNewToken() ([]byte, error) {
	api := "/v1/auth/tokens"

	url := fmt.Sprintf("%s%s", oss.Host, api)
	data := map[string]string{"username": oss.User, "password": oss.Password,}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b)

	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, fmt.Errorf("create new request %s failed %v", api, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oss.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("reqest new token %s failed %v", api, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request new token %s status %s", api, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

func envLookUp(key string) (value string, err error) {
	errPrefix := "missing env: %s"
	if value, ok := os.LookupEnv(key); !ok {
		// todo add logger in here
		return value, fmt.Errorf(errPrefix, key)
	}
	return value, nil
}
