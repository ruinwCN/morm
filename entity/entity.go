package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	timeFormat  = "2006-01-02 15:04:05"
	timeZone, _ = time.LoadLocation("Asia/Shanghai")
)

const (
	TagTrue    = "true"
	TagFieldId = "Id"
	TagField   = "morm"
	TagTable   = "morm_table"  // table name
	TagIgnore  = "morm_ignore" // ignore
	TagDate    = "morm_date"   // date int64 or time.time
)

type BaseEntity struct {
}

// GetTableName
// 优先级从高到低
// 1. id 的 morm_table 标签
// 2. Struct func TableName
// 3. Entity low string
func GetTableName(obj interface{}) (string, error) {
	tableName := ""
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return "", errors.New("obj kind must ptr")
	}
	tof := reflect.TypeOf(obj).Elem()
	for i := 0; i < tof.NumField(); i++ {
		field := tof.Field(i)
		if field.Name == TagFieldId {
			tableName = field.Tag.Get(TagTable)
			break
		}
	}

	if tableName == "" {
		funcA := reflect.ValueOf(obj).MethodByName("TableName").Call([]reflect.Value{})
		if len(funcA) == 1 {
			tableName = funcA[0].String()
		}
	}

	if tableName == "" {
		tableName = strings.ToLower(tof.Name())
	}
	return tableName, nil
}

func GetMapWithMORM(obj interface{}) (map[string]interface{}, error) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return nil, errors.New("obj kind must ptr")
	}
	obcKeyValueMap := make(map[string]interface{}, 0)
	typeOf := reflect.TypeOf(obj).Elem()
	valueOf := reflect.ValueOf(obj).Elem()
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		if ValidTagFieldIgnore(field.Tag.Get(TagIgnore)) {
			continue
		}
		fieldTagName := field.Tag.Get(TagField)
		if fieldTagName == "" {
			continue
		}
		kind := valueOf.Field(i).Kind()
		if reflect.String == kind {
			obcKeyValueMap[fieldTagName] = valueOf.Field(i).String()
		} else if reflect.Bool == kind {
			obcKeyValueMap[fieldTagName] = valueOf.Field(i).Bool()
		} else if reflect.Int == kind || reflect.Int8 == kind || reflect.Int16 == kind {
			obcKeyValueMap[fieldTagName] = valueOf.Field(i).Int()
		} else if reflect.Int32 == kind || reflect.Int64 == kind {
			if ValidTagFieldDate(field.Tag.Get(TagDate)) {
				obcKeyValueMap[fieldTagName] = Timestamp2Str(valueOf.Field(i).Int())
			} else {
				obcKeyValueMap[fieldTagName] = valueOf.Field(i).Int()
			}
		} else if reflect.Uint == kind || reflect.Uint8 == kind || reflect.Uint16 == kind || reflect.Uint32 == kind || reflect.Uint64 == kind {
			obcKeyValueMap[fieldTagName] = valueOf.Field(i).Uint()
		} else if reflect.Float32 == kind || reflect.Float64 == kind {
			obcKeyValueMap[fieldTagName] = valueOf.Field(i).Float()
		} else if reflect.Struct == kind {
			if ValidTagFieldDate(field.Tag.Get(TagDate)) {
				date, ok := valueOf.Field(i).Interface().(time.Time)
				if !ok {
					return nil, errors.New(fmt.Sprintf("this %s transfer date faile", field.Name))
				}
				obcKeyValueMap[fieldTagName] = Time2Str(date)
			} else {
				valueOfMarshalled, err := json.Marshal(valueOf.Field(i).Interface())
				if err != nil {
					return nil, errors.New(fmt.Sprintf("this %s transfer faile", field.Name))
				} else {
					obcKeyValueMap[fieldTagName] = string(valueOfMarshalled)
				}
			}
		} else if !valueOf.Field(i).IsNil() && (reflect.Slice == kind || reflect.Map == kind || reflect.Interface == kind) {
			// The argument must be a chan, func, interface, map, pointer, or slice value;
			valueOfMarshalled, err := json.Marshal(valueOf.Field(i).Interface())
			if err != nil {
				return nil, errors.New(fmt.Sprintf("this %s transfer faile", field.Name))
			} else {
				obcKeyValueMap[fieldTagName] = string(valueOfMarshalled)
			}
		}
	}
	return obcKeyValueMap, nil
}

func ValidTagFieldIgnore(tagString string) bool {
	if tagString == TagTrue {
		return true
	}
	return false
}

func ValidTagFieldDate(tagString string) bool {
	if tagString == TagTrue {
		return true
	}
	return false
}

func Timestamp2Str(t int64) string {
	return time.Unix(t/1000, 0).In(timeZone).Format(timeFormat)
}

func Time2Str(t time.Time) string {
	return t.Format(timeFormat)
}

func TimeParse(t string) (time.Time, error) {
	str, err := time.Parse(timeFormat, t)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}
	return str, nil
}
