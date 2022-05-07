package conf

import (
	"fmt"
	"testing"
)

func TestGetDBConfig(t *testing.T) {
	c := GetDBConfig("../mysql.yaml")
	fmt.Println(c.DSN())
}
