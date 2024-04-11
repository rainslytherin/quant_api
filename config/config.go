package config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/tidwall/pretty"
)

var globalConfig *Config = nil

func GetGlobalConfig() *Config {
	return globalConfig
}

func InitGlobalConfig(configFile string) {
	cfg, err := readConfigFile(configFile)
	if err != nil {
		fmt.Println("read config file error,", err)
		os.Exit(1)
	}
	globalConfig = cfg
	fmt.Println("Global Config:", globalConfig.String())
}

// read config from config_private.json and config.json
// and return config object

type Config struct {
	Log struct {
		Level string `json:"level"`
	}
	Http struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Database string `json:"database"`
	}
	Backend map[string]struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
}

func (c *Config) ToMap() map[string]interface{} {
	return StructToMap(c)
}

func (c *Config) String() string {
	// use json.Marshal to convert struct to string
	b, _ := json.Marshal(c)
	return string(pretty.Pretty(b))
}

// transform Config object to map[string]string
func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// 获取结构体字段数量
	numFields := value.NumField()

	// 遍历结构体字段
	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		fieldName := value.Type().Field(i).Name

		// 如果字段是嵌套结构体，递归调用structToMap
		if field.Kind() == reflect.Struct {
			result[fieldName] = StructToMap(field.Interface())
		} else {
			result[fieldName] = field.Interface()
		}
	}

	return result
}

func readConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
