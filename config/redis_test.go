package config

import (
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	s, err := GetRedisConnectionString()
	if err != nil {
		t.Errorf("Get redis connection string error: %s", err.Error())
	}

	fmt.Println(s)
}
