package appconf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
	"sync"
)

/**
Global application configuration
*/
type AppConf struct {
	Accounts []Account
}

type Account struct {
	ServerUrl	string
	Login		string
	Password	string
}

const appConfFileName = "conf.json"
var appConf AppConf
var isLoaded bool
var mux sync.Mutex

func GetConfig() *AppConf {
	if !isLoaded {
		mux.Lock()
		defer mux.Unlock()
		if !isLoaded {
			load();
		}		
	}
	return &appConf
}

func StoreConfig() error {
	confContent, err := json.MarshalIndent(appConf, "", " ")
	if err == nil {
		err = ioutil.WriteFile(appConfFileName, confContent, 0644)
	}
 
	return err
}

func load() {
	var _, err = os.Stat(appConfFileName)
	if !os.IsNotExist(err) {
		confContent, err := ioutil.ReadFile(appConfFileName)
		if err == nil {
			err = json.Unmarshal(confContent, &appConf)
			if err != nil {
  				panic(fmt.Sprintf("Couldn't parse main config file. %v", err))
			}
		} else {
			panic(fmt.Sprintf("Couldn't load main config file. %v", err))
		} 
	}

	isLoaded = true
}
