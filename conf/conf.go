package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	fxossUtils "github.com/super1-chen/fxoss/utils"
)

var timeLayout = "2006-01-02 15:04:05"

type config struct {
	Host      string `json:"host"`
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

// NewConfig create new config is path is
func NewConfig() *config {
	return &config{}
}

// Update config
func (conf *config) Update(data []byte) error {
	if err := json.Unmarshal(data, &conf); err != nil {
		return fmt.Errorf("update config failed %v", err)
	}
	return nil
}

// SetHost set config host
func (conf *config) SetHost(host string) error {
	conf.Host = host
	return nil
}

// GetToken return token of config
func (conf *config) GetToken() string {
	return conf.Token
}

// IsValid checks config is not expired and contains an expected hostname
func (conf *config) IsValid(host string, nowUTCTime time.Time) bool {

	if conf.Host != host {
		return false
	}
	if isValid, err := checkTimeValid(conf.ExpiredAt, nowUTCTime); !isValid || err != nil {
		return false
	}
	return true
}

// Save the config
func (conf *config) Save(fileName string) error {
	dirName := path.Dir(fileName)
	err := fxossUtils.CreateFolder(dirName)
	if err != nil {
		return err
	}

	f, err := os.Create(fileName)

	if err != nil {
		return fmt.Errorf("create file %s failed %v", fileName, err)
	}

	defer f.Close()

	bytesBuffer, err := json.Marshal(&conf)
	if err != nil {
		return fmt.Errorf("json marshal failed %v", err)
	}
	f.Write(bytesBuffer)
	return nil

}

// checkTimeValid  convert timeStr to time.Time then checks the time whether after the now
func checkTimeValid(timeStr string, now time.Time) (bool, error) {
	t, err := time.ParseInLocation(timeLayout, timeStr, time.UTC)
	if err != nil {
		err := fmt.Errorf("parse time %s failed: %v", timeStr, err)
		return false, err
	}
	return t.After(now), nil
}
