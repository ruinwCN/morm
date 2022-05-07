package tmp

import (
	"encoding/json"
	"fmt"
	"github.com/ruinwCN/morm/model_test/model"
	"reflect"
	"time"
)

func structToJson() {
	type SubModel struct {
		SName string `json:"sname" morm:"sname"`
	}
	type Model1 struct {
		Id       int64         `json:"id" morm:"id" morm_table:"model1"`
		Name     string        `json:"name" morm:"name"`
		Ignore   string        `json:"ignore" morm_ignore:"true"`
		Age      int32         `json:"age" morm:"age"`
		Duration time.Duration `json:"duration" morm:"duration"`
		List     []string      `json:"list" morm:"list"`
		Sub      SubModel      `json:"submodel" morm:"submodel"`
	}
	m := &Model1{1, "name1", "ig", 12, 12345, []string{"1", "2"}, SubModel{"123"}}

	model_marshalled, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(model_marshalled))
	}

	tMap := make(map[string]Model1, 0)
	tMap["t"] = *m
	tmap_marshalled, err := json.Marshal(tMap)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(tmap_marshalled))
	}
}

func EntityTestReflect() {
	obj := model.Status{Id: 12, Name: "name", CreateTime: time.Now()}
	getType := reflect.TypeOf(obj)
	fmt.Println("get Type is :", getType.Name())

	getValue := reflect.ValueOf(obj)
	fmt.Println("get all Fields is:", getValue)
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)

		value := getValue.Field(i).Interface()
		subType := reflect.TypeOf(value)

		fmt.Println(subType.Kind().String())
		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		fmt.Println(field.Tag)
		fmt.Println(field.Tag.Get("morm"))
	}
}

type bProtocol interface {
	GetValue() string
}
type BaseEntiry struct {
	value string
}

func (self *BaseEntiry) GetValue() string {
	return self.value
}

type OneEntity struct {
	BaseEntiry
}
type TwoEntity struct {
	BaseEntiry
}

func EntityGetValue(i bProtocol) string {
	return i.GetValue()
}

func TryEntity() {
	one := &OneEntity{BaseEntiry{
		"One",
	}}
	two := &TwoEntity{BaseEntiry{
		"Two",
	}}
	fmt.Println(EntityGetValue(one))
	fmt.Println(EntityGetValue(two))
}

func TryReflect() {

	//static type
	var i int = 233
	p := reflect.ValueOf(&i)
	switch p.Interface().(type) {
	case *int:
		{
			cp := p.Interface().(*int)
			fmt.Println(cp)
		}
	default:
		println("done")
	}

	//concrete type
	m := model1{dataString: "s", dataInt: 1}

	//var mInterface interface{} = m
	//rt := reflect.TypeOf(mInterface)
	rt := reflect.TypeOf(m)
	fmt.Println("rt.Name() :" + rt.Name())

	rv := reflect.ValueOf(m)
	fmt.Println("value : ", rv)

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := rv.Field(i)
		fmt.Printf("field's name: %s, type: %v, value: %v\n", field.Name, field.Type, value)
	}
}

type model1 struct {
	dataString string
	dataInt    int
}

func (m *model1) Echo() {
	fmt.Printf("echo 1")
}

func FindModelReflect() {
}
