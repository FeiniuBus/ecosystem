package config

import "bytes"
import "strconv"

const (
	sharedRedisFilename = "redis.json"
	redisSectionName    = "RedisConnectionStrings"
)

type redisConfig struct {
	RedisConnectionStrings *redisNode
}

type redisNode struct {
	Endpoints          []string
	Password           string
	KeepAlive          int
	SyncTimeout        int
	AllowAdmin         bool
	ConnectTimeout     int
	WriteBuffer        int
	AbortOnConnectFail bool
	ConnectRetry       int
}

func newRedisNode() *redisNode {
	return &redisNode{
		KeepAlive:          -1,
		SyncTimeout:        1000,
		AllowAdmin:         false,
		ConnectTimeout:     5000,
		WriteBuffer:        4096,
		AbortOnConnectFail: true,
		ConnectRetry:       3,
	}
}

func (n *redisNode) build() string {
	buffer := bytes.NewBufferString("")
	if len(n.Endpoints) == 0 {
		return buffer.String()
	}

	for _, e := range n.Endpoints {
		if len(e) == 0 {
			continue
		}
		if buffer.Len() > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(e)
	}

	n.append(buffer, "keepAlive", n.KeepAlive)
	n.append(buffer, "syncTimeout", n.SyncTimeout)
	n.append(buffer, "allowAdmin", n.AllowAdmin)
	n.append(buffer, "connectTimeout", n.ConnectTimeout)
	n.append(buffer, "password", n.Password)
	n.append(buffer, "writeBuffer", n.WriteBuffer)
	n.append(buffer, "connectRetry", n.ConnectRetry)
	n.append(buffer, "abortConnect", n.AbortOnConnectFail)

	return buffer.String()
}

func (n *redisNode) append(buffer *bytes.Buffer, prefix string, value interface{}) {
	s := getValue(value)
	if s != "" {
		if buffer.Len() != 0 {
			buffer.WriteString(",")
		}
		if prefix != "" {
			buffer.WriteString(prefix)
			buffer.WriteString("=")
		}
		buffer.WriteString(s)
	}
}

func getValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	}
	return ""
}

// GetRedisConnectionString returns redis connection string
func GetRedisConnectionString() (string, error) {
	n := newRedisNode()
	c := &redisConfig{
		RedisConnectionStrings: n,
	}

	_, err := Load(sharedRedisFilename, &c)
	if err != nil {
		return "", nil
	}

	node := c.RedisConnectionStrings
	return node.build(), nil
}
