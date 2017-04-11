package config

import (
	"bytes"
	"fmt"
	"net/url"
)

// consts
const (
	DefaultMongoConnectTimeout        = 30
	DefaultMongoMaxConnectionIdleTime = 10
	DefaultMongoMaxConnectionLifeTime = 30
	DefaultMongoMaxConnectionPoolSize = 100
	DefaultMongoMinConnectionPoolSize = 0
	DefaultMongoSocketTimeout         = 0

	sharedMongoFilename = "mongo.json"
)

type mongoConfig struct {
	MongoURL *MongoNode `json:"MongoUrl"`
}

// MongoNode is
type MongoNode struct {
	Endpoints               []string
	AuthenticationMechanism string
	Username                string
	Password                string
	DatabaseName            string
	ConnectTimeout          int
	MaxConnectionIdleTime   int
	MaxConnectionLifeTime   int
	MaxConnectionPoolSize   int
	MinConnectionPoolSize   int
	SocketTimeout           int
}

func newMongoNode() *MongoNode {
	return &MongoNode{
		AuthenticationMechanism: "",
		ConnectTimeout:          DefaultMongoConnectTimeout,
		MaxConnectionIdleTime:   DefaultMongoMaxConnectionIdleTime,
		MaxConnectionLifeTime:   DefaultMongoMaxConnectionLifeTime,
		MaxConnectionPoolSize:   DefaultMongoMaxConnectionPoolSize,
		MinConnectionPoolSize:   DefaultMongoMinConnectionPoolSize,
		SocketTimeout:           DefaultMongoSocketTimeout,
		Username:                "",
		Password:                "",
	}
}

func (n *MongoNode) build() string {
	buffer := bytes.NewBufferString("")
	if len(n.Endpoints) == 0 {
		return buffer.String()
	}

	buffer.WriteString("mongodb://")
	if n.Username != "" {
		buffer.WriteString(url.QueryEscape(n.Username))
		if n.Password != "" {
			buffer.WriteString(fmt.Sprintf(":%s", url.QueryEscape(n.Password)))
		}
		buffer.WriteString("@")
	} else if n.Password != "" {
		buffer.WriteString(fmt.Sprintf(":%s@", n.Password))
	}

	firstServer := true
	for _, s := range n.Endpoints {
		if !firstServer {
			buffer.WriteString(",")
		}
		buffer.WriteString(s)
		firstServer = false
	}

	if n.DatabaseName != "" {
		buffer.WriteString("/")
		buffer.WriteString(n.DatabaseName)
	}

	query := bytes.NewBufferString("")
	if n.AuthenticationMechanism != "" {
		query.WriteString(fmt.Sprintf("authMechanism=%s;", n.AuthenticationMechanism))
	}
	if n.ConnectTimeout != DefaultMongoConnectTimeout {
		query.WriteString(fmt.Sprintf("connectTimeout=%d;", n.ConnectTimeout))
	}
	if n.MaxConnectionIdleTime != DefaultMongoMaxConnectionIdleTime {
		query.WriteString(fmt.Sprintf("maxIdleTime=%d;", n.MaxConnectionIdleTime))
	}
	if n.MaxConnectionLifeTime != DefaultMongoMaxConnectionLifeTime {
		query.WriteString(fmt.Sprintf("maxLifeTime=%d;", n.MaxConnectionLifeTime))
	}
	if n.MaxConnectionPoolSize != DefaultMongoMaxConnectionPoolSize {
		query.WriteString(fmt.Sprintf("maxPoolSize=%d;", n.MaxConnectionPoolSize))
	}
	if n.MinConnectionPoolSize != DefaultMongoMinConnectionPoolSize {
		query.WriteString(fmt.Sprintf("minPoolSize=%d;", n.MinConnectionPoolSize))
	}
	if n.SocketTimeout != DefaultMongoSocketTimeout {
		query.WriteString(fmt.Sprintf("socketTimeout=%d;", n.SocketTimeout))
	}

	if query.Len() != 0 {
		query.Truncate(query.Len() - 1)
		if n.DatabaseName == "" {
			buffer.WriteString("/")
		}
		buffer.WriteString("?")
		buffer.WriteString(query.String())
	}

	return buffer.String()
}

// GetMongoDialInfo returns mongodb url
func GetMongoDialInfo(database string) (*MongoNode, error) {
	n := newMongoNode()
	c := &mongoConfig{
		MongoURL: n,
	}

	_, err := Load(sharedMongoFilename, &c)
	if err != nil {
		return nil, err
	}

	c.MongoURL.DatabaseName = database
	return c.MongoURL, nil
}
