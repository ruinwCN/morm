package morm

import (
	"fmt"
	"github.com/ruinwCN/morm/conf"
	"testing"
)

func TestNewEngines(t *testing.T) {
	dbConfig := conf.GetDBConfig("mysql.yaml")
	dbs := []conf.DBConfig{dbConfig}
	err := NewEngines(dbs)
	if err != nil {
		fmt.Println(err)
	}
}
