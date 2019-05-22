package conf

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

// LoadConfig
func (conf *config) LoadConfig(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("config file %s doesn't exists\r\n", filename)
	}

	file, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Printf("read config %s failed by %v\r\n", filename, err)
		msg := fmt.Sprintf("read config failed %s\r\n", err)
		fmt.Println(msg)
		return err
	}
	err = json.Unmarshal(file, &conf)

	if err != nil {
		msg := fmt.Sprintf("json unmarshal config failed %s\r\n", err)
		fmt.Println(msg)
		log.Printf("json unmarshal config failed %s\r\n", err)
		return err
	}
	return nil

}

func (conf *config) Update(in io.Reader) error {
	if err := json.NewDecoder(in).Decode(&conf); err != nil {
		return fmt.Errorf("update config failed %v", err)
	}
	return nil
}

func (conf *config) SetHost(host string) error {

}

// Save the config
func (conf *config) Save(folderPath, name string) error {

	err := fxossUtils.CreateFolder(folderPath)
	if err != nil {
		return err
	}

	fileName := path.Join(folderPath, name)
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

// IsValid checks config is not expired and contains an expected hostname
func (conf *config) IsValid(host string, nowTime time.Time) bool {

	if conf.Host != host {
		// todo need add a logger in here.
		return false
	}
	if isValid, err := checkTimeValid(conf.ExpiredAt, nowTime); !isValid || err != nil {
		return false
	}
	return true
}

// checkTimeValid  convert timeStr to time.Time then checks the time whether after the now
func checkTimeValid(timeStr string, now time.Time) (bool, error) {
	loc, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation(timeLayout, timeStr, loc)
	if err != nil {
		fmt.Errorf("parse time %s failed: %v", timeStr, err)
		return false, err
	}
	return t.After(now), nil
}
