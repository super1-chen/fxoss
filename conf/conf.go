package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	//"path"

	//fxossUtils "github.com/super1-chen/fxoss/utils"
)

const (
	TokenJson = "token.json"
)

type Config struct {
	Host      string `json:"host"`
	Token     string `json:"token"`
	ExpiredAt string `json:"expired_at"`
}

func NewConfig(path string) (config *Config, err error) {
	config = &Config{}
	err = LoadConfig(path, config)
	return config, err
}

func LoadConfig(filename string, conf *Config) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("config file %s doesn't exists\r\n", filename)
	}

	file, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Printf("read config %s failed by %v\r\n", filename, err)
		msg := fmt.Sprintf("read config failed %s\r\n", err)
		//fxossUtils.ErrorPrintln(msg, false)
		fmt.Println(msg)
		return err
	}
	err = json.Unmarshal(file, &conf)

	if err != nil {
		msg := fmt.Sprintf("json unmarshal config failed %s\r\n", err)
		fmt.Println(msg)
		//fxossUtils.ErrorPrintln(msg, false)
		log.Printf("json unmarshal config failed %s\r\n", err)
		return err
	}

	return nil

}

func (conf *Config)Save(fileDir string){
	if _, err := os.Stat(fileDir); os.IsNotExist(err){
		//TODO NEED ADD LOG PRINT
		// log.info("create new log)
	}
}