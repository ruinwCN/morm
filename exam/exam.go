package exam

import (
	"github.com/ruinwCN/morm"
	"github.com/ruinwCN/morm/conf"
	"github.com/ruinwCN/morm/model_test/model"
	"github.com/ruinwCN/util/uerror"
	"time"
)

func examModelInsert() {
	mysqlConfigs := conf.GetDBConfig("../mysql.yaml")
	err := morm.NewEngines([]conf.DBConfig{mysqlConfigs})
	if err != nil {
		println(uerror.NewError("exam config init error ", err))
		return
	}
	status := model.Status{
		Name:       "test_model",
		CreateTime: time.Now(),
	}
	engine, err := morm.GetDB("")
	if err != nil {
		panic(err)
	} else {
		manager := morm.BaseManager{}
		_, err = manager.Insert(engine, &status)
		if err != nil {
			println(err)
		}
	}
}
