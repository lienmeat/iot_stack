package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//SQLDatabase SQL database config
type SQLDatabase struct {
	Driver         string `json:"driver"` //mysql, mssql or other remote-capable db driver supported by gorm
	User           string `json:"user"`
	Password       string `json:"password"`
	Host           string `json:"host"`
	Database       string `json:"database"`
	Port           string `json:"port"`
	MaxConnections int    `json:"max_connections"`
}

type LogConf struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

type Conf struct {
	Sqldb    SQLDatabase `json:"sqldb"`
	Log      LogConf     `json:"log"`
	Port     int         `json:"port"`
	FilesDir string      `json:"files_dir"`
}

var conf Conf

func Get() Conf {
	if conf.Port == 0 {
		conf = Conf{}
		ParseConfig(&conf)
	}
	return conf
}

//ParseConfig Parses CONFIG env variable into the passed in config structure
// config should be a pointer to a struct
func ParseConfig(config interface{}) {
	// Grab configuration from the environment variable
	configFile := os.Getenv("CONFIGFILE")

	if configFile == "" {
		panic("Config file not defined. Use CONFIGFILE environment variable. Ex: CONFIGFILE=./lb-config.json")
	}

	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, config)
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %s", err))
	}
}

//EnvironmentVars convienience for getting a map of all Env vars
func EnvironmentVars() map[string]string {
	env := make(map[string]string)
	for _, v := range os.Environ() {
		//v is in the format of "key=value", so we must split by the first "="
		pair := strings.SplitN(v, "=", 2)
		env[pair[0]] = pair[1]
	}
	return env
}
