package models

import (
	"fmt"

	"quant_api/database"

	"github.com/RaveNoX/go-jsonmerge"
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
	ID           int    `db:"id"`
	Scope        string `db:"scope"`
	Name         string `db:"name"`
	Value        []byte `db:"value"`
	ChangedValue []byte `db:"changed_value"`
	CreateTime   string `db:"create_time"`
	UpdateTime   string `db:"update_time"`
	UpdateUser   string `db:"update_user"`
}

func NewConfig(scope, name string, value []byte, updateUser string) *Config {
	return &Config{
		Scope:      scope,
		Name:       name,
		Value:      value,
		UpdateUser: updateUser,
	}
}

func MergeJson(oldValue, newValue []byte) ([]byte, error) {
	merged, info, err := jsonmerge.MergeBytes(oldValue, newValue)
	fmt.Println("merged", string(merged), "info", info, "err", err)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

func (c *Config) MergeValue(value []byte) error {
	if c.ChangedValue == nil {
		c.ChangedValue = []byte("{}")
	}
	if c.Value == nil {
		c.Value = []byte("{}")
	}
	mergedChangeValue, err := MergeJson(c.ChangedValue, value)
	if err != nil {
		return err
	}

	mergedValue, err := MergeJson(c.Value, value)
	if err != nil {
		return err
	}

	c.ChangedValue = mergedChangeValue
	c.Value = mergedValue

	fmt.Println("value", string(value))
	fmt.Println("c.ChangedValue", string(c.ChangedValue))
	fmt.Println("c.Value", string(c.Value))

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

	_, err = db.Exec("INSERT INTO configs(scope, name, value, update_user) VALUES(?, ?, ?, ?)", c.Scope, c.Name, c.Value, c.UpdateUser)

	return err
}

// Save the config to the database
func (c *Config) Save() error {
	db, err := database.GetGlobalDB()
	if err != nil {
		return err
	}

	_, err = db.Exec("UPDATE configs SET changed_value = ?, value = ?, update_user = ? WHERE scope = ? AND name = ?", c.ChangedValue, c.Value, c.UpdateUser, c.Scope, c.Name)

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
