package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"quant_api/database"
)

/*
CREATE TABLE `configs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `scope` varchar(50) NOT NULL DEFAULT 'stock',
  `name` varchar(50) NOT NULL DEFAULT 'default_name',
  `value` json,
  `changed_value` json,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `update_user` varchar(50) NOT NULL DEFAULT 'admin',
  PRIMARY KEY (`id`),
  UNIQUE KEY (`scope`, `name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci
*/

type Config struct {
	ID           int        `db:"id" json:"id"`
	Scope        string     `db:"scope" json:"scope"`
	Name         string     `db:"name" json:"name"`
	Value        JsonObject `db:"value" json:"value"`
	ChangedValue JsonObject `db:"changed_value" json:"changed_value"`
	CreateTime   string     `db:"create_time" json:"create_time"`
	UpdateTime   string     `db:"update_time" json:"update_time"`
	UpdateUser   string     `db:"update_user" json:"update_user"`
} // @name Config

type JsonObject map[string]interface{}

func (pc *JsonObject) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &pc)
		return nil
	case string:
		json.Unmarshal([]byte(v), &pc)
		return nil
	case nil:
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
func (pc *JsonObject) Value() (driver.Value, error) {
	return json.Marshal(pc)
}

func NewConfig(scope, name string, value map[string]interface{}, updateUser string) *Config {
	return &Config{
		Scope:      scope,
		Name:       name,
		Value:      value,
		UpdateUser: updateUser,
	}
}

func MergeJson(oldValueByte, newValueByte []byte) ([]byte, error) {
	oldObj := make(map[string]interface{})
	newObj := make(map[string]interface{})
	if err := json.Unmarshal(oldValueByte, &oldObj); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(newValueByte, &newObj); err != nil {
		return nil, err
	}

	for k, v := range newObj {
		oldObj[k] = v
	}

	mergedValue, err := json.Marshal(oldObj)
	if err != nil {
		return nil, err
	}

	return mergedValue, nil
}

func PrintMap(m map[string]interface{}) {
	for k, v := range m {
		fmt.Println(k, v)
	}
}

func MergeObject(oldObj, newObj map[string]interface{}) map[string]interface{} {
	for k, v := range newObj {
		oldObj[k] = v
	}

	return oldObj
}

func (c *Config) MergeValue(value map[string]interface{}) error {
	if c.ChangedValue == nil {
		c.ChangedValue = make(map[string]interface{})
	}
	if c.Value == nil {
		c.Value = make(map[string]interface{})
	}
	c.ChangedValue = MergeObject(c.ChangedValue, value)
	c.Value = MergeObject(c.Value, value)

	return nil
}

func (c *Config) Delete() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	// use delete sql
	_, err = db.Exec("DELETE FROM configs WHERE scope = ? AND name = ?", c.Scope, c.Name)
	return err
}

func (c *Config) Create() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	valueStr, _ := json.Marshal(c.Value)

	_, err = db.Exec("INSERT INTO configs(scope, name, value, changed_value, update_user) VALUES(?, ?, ?, ?, ?)", c.Scope, c.Name, valueStr, valueStr, c.UpdateUser)

	return err
}

// Save the config to the database
func (c *Config) Save() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	valueStr, _ := json.Marshal(c.Value)
	changeValueStr, _ := json.Marshal(c.ChangedValue)

	_, err = db.Exec("UPDATE configs SET changed_value = ?, value = ?, update_user = ? WHERE scope = ? AND name = ?", changeValueStr, valueStr, c.UpdateUser, c.Scope, c.Name)

	return err
}

// Reload the config from the database
func (c *Config) Reload() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	err = db.Get(c, "SELECT * FROM configs WHERE scope = ? AND name = ?", c.Scope, c.Name)

	return err
}

// Get the config from the database
func GetConfig(scope, name string) (*Config, error) {
	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	var config Config
	err = db.Get(&config, "SELECT * FROM configs WHERE scope = ? AND name = ?", scope, name)

	return &config, err
}

func GetConfigs(scope string) ([]*Config, error) {
	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	configs := make([]*Config, 0)
	err = db.Select(&configs, "SELECT * FROM configs WHERE scope = ?", scope)

	return configs, err
}

func GetConfigsAfterTime(timeStamp int64) ([]*Config, error) {
	db, err := database.GetGlobalDB()
	if err != nil {
		return nil, err
	}

	configs := make([]*Config, 0)
	err = db.Select(&configs, "SELECT * FROM configs WHERE update_time > ?", timeStamp)

	return configs, err
}
