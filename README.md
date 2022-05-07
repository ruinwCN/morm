## morm

早期熟练反射的一个简单orm-lib。同时也是go-gen中支持gorm和morm的model自动创建基础。

v0.1 

### yaml
```go
# mysql
address: "127.0.0.1"
port : 3306
user_name : "root"
password : ""
database : "event_project"
tag : ""

timeout : 30
max_life_time : 300
max_open_conns: 1024
max_idle_conns : 1024

```

### model
```go
type EventT struct {
	Id         int64     `morm:"id" json:"id" morm_table:"event_t"`
	Ndecimal   float32   `morm:"ndecimal" json:"ndecimal"`
	Nbigint    int64     `morm:"nbigint" json:"nbigint"`
	Nchar      byte      `morm:"nchar" json:"nchar"`
	Ntext      string    `morm:"ntext" json:"ntext"`
	Nfloat     float32   `morm:"nfloat" json:"nfloat"`
	CreateTime time.Time `morm:"create_time" morm_date:"true" json:"create_time"`
}
```
###exam

#### insert
- BaseManager 也可以在dao-model中实现

```go
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
```
