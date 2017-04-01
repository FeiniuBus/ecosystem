package config

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

const (
	sharedMysqlFilename = "mysql.json"
	mysqlSectionName    = "MysqlConnectionStrings"
	defaultMysqlProfile = "default"
)

// GetMysqlConnectionString get mysql connection string
func GetMysqlConnectionString(database, profile string) (string, error) {
	var c map[string]interface{}
	_, err := Load(sharedMysqlFilename, &c)
	if err != nil {
		return "", err
	}

	profileName := profile
	if profileName == "" {
		profileName = defaultMysqlProfile
	}

	if m1, ok := c[mysqlSectionName].(map[string]interface{}); ok {
		if m2, o := m1[profileName].(map[string]interface{}); o {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=%s&parseTime=true&loc=%s", m2["User"], m2["Password"], m2["Host"], m2["Port"], database, m2["Charset"], url.QueryEscape(time.Local.String()))
			return dsn, nil
		}
	}
	return "", errors.New("读取mysql配置失败")
}
