package config

import (
	"fmt"
	"testing"
)

func TestMongo(t *testing.T) {
	s, err := GetMongoURL("test")
	if err != nil {
		t.Errorf("Get mongo url error: %s", err.Error())
	}

	fmt.Println(s)
}
