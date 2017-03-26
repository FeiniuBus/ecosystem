package config

import (
	"bytes"
	"errors"
	"fmt"
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
			var buffer bytes.Buffer
			var delimiter = ""

			for k, v := range m2 {
				if k == "Host" {
					buffer.WriteString(fmt.Sprintf("%s%s=%s", delimiter, "server", v))
				} else if k == "Port" {
					buffer.WriteString(fmt.Sprintf("%s%s=%v", delimiter, "port", v))
				} else if k == "User" {
					buffer.WriteString(fmt.Sprintf("%s%s=%s", delimiter, "user_id", v))
				} else if k == "Password" {
					buffer.WriteString(fmt.Sprintf("%s%s=%s", delimiter, "password", v))
				} else if k == "Charset" {
					buffer.WriteString(fmt.Sprintf("%s%s=%s", delimiter, "characterset", v))
				}
				delimiter = ";"
			}

			buffer.WriteString(fmt.Sprintf("%s%s=%s", delimiter, "database", database))
			return buffer.String(), nil
		}
	}

	return "", errors.New("读取mysql配置失败")
}
