package morm

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ruinwCN/morm/bufferpool"
	"github.com/ruinwCN/morm/entity"
	"reflect"
	"strconv"
	"strings"
)

type CommandType uint

const (
	CommandTypeEq         CommandType = iota + 1 // =
	CommandTypeNotEq                             // <>
	CommandTypeGT                                // >
	CommandTypeLT                                // <
	CommandTypeGTE                               // >=
	CommandTypeLTE                               // <=
	CommandTypeBetween                           // BETWEEN
	CommandTypeNotBetween                        // NotBetween
	CommandTypeLike                              // LIKE
	CommandTypeIn                                // IN
	CommandTypeAnd                               // AND
	CommandTypeOr                                // OR
)

// BaseManager 操作对象
type BaseManager struct {
}

// BaseOption option
type BaseOption struct {
	commands []*Command
	entity   interface{}
}

// Command obj
type Command struct {
	commandType CommandType
	key         string
	value       interface{}
	values      []interface{}
}

func (baseOption *BaseOption) AddCommand(command *Command) *BaseOption {
	baseOption.commands = append(baseOption.commands, command)
	return baseOption
}

func (baseOption *BaseOption) getWereSqlWord() (*string, error) {
	bufferField := bufferpool.Get()
	defer bufferField.Free()
	bufferField.AppendString("WHERE")
	for _, cmd := range baseOption.commands {
		if CommandTypeAnd == cmd.commandType {
			bufferField.AppendString(" AND")
		} else if CommandTypeOr == cmd.commandType {
			bufferField.AppendString(" OR")
		} else if CommandTypeBetween == cmd.commandType {
			betweenSql, err := sqlBetweenOption("BETWEEN", cmd)
			if err != nil {
				return nil, errors.New("get between sql error. ")
			}
			if betweenSql != "" {
				bufferField.AppendString(betweenSql)
			}
		} else if CommandTypeNotBetween == cmd.commandType {
			notBetween, err := sqlBetweenOption("NOT BETWEEN", cmd)
			if err != nil {
				return nil, errors.New("get not between sql error. ")
			}
			if notBetween != "" {
				bufferField.AppendString(notBetween)
			}
		} else if CommandTypeIn == cmd.commandType {
			inSql, err := sqlInOption(cmd)
			if err != nil {
				return nil, errors.New("get in sql error. ")
			}
			if inSql != "" {
				bufferField.AppendString(inSql)
			}
		} else if CommandTypeEq == cmd.commandType {
			unarySql := sqlUnaryOption("=", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		} else if CommandTypeNotEq == cmd.commandType {
			unarySql := sqlUnaryOption("<>", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		} else if CommandTypeLT == cmd.commandType {
			unarySql := sqlUnaryOption("<", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		} else if CommandTypeGT == cmd.commandType {
			unarySql := sqlUnaryOption(">", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		} else if CommandTypeLTE == cmd.commandType {
			unarySql := sqlUnaryOption("<=", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		} else if CommandTypeGTE == cmd.commandType {
			unarySql := sqlUnaryOption(">=", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		} else if CommandTypeLike == cmd.commandType {
			unarySql := sqlUnaryOption("LIKE", cmd)
			if unarySql != "" {
				bufferField.AppendString(unarySql)
			}
		}

	}
	sqlWord := bufferField.String()
	return &sqlWord, nil
}

// Insert engine obj
func (s *BaseManager) Insert(engine *Engine, obj interface{}) (int64, error) {
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return -1, errors.New("obj kind must ptr")
	}
	tableName, err := entity.GetTableName(obj)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("get table name error %s", err.Error())
	}

	mp, err := entity.GetMapWithMORM(obj)
	if err != nil {
		fmt.Println(err)
		return -1, fmt.Errorf("get key/value error %s", err.Error())
	}
	bufferField := bufferpool.Get()
	defer bufferField.Free()
	bufferPlaceholder := bufferpool.Get()
	defer bufferPlaceholder.Free()
	flag := true
	values := make([]interface{}, 0)
	for key, value := range mp {
		kind := reflect.TypeOf(value).Kind()
		if reflect.Int64 == kind {
		}
		if flag {
			bufferField.AppendString(key)
			bufferPlaceholder.AppendString("?")
			values = append(values, value)
			flag = false
			continue
		}
		bufferField.AppendString(",")
		bufferField.AppendString(key)
		bufferPlaceholder.AppendString(",?")
		values = append(values, value)
	}
	bufferStmt := bufferpool.Get()
	defer bufferStmt.Free()
	bufferStmt.AppendString("INSERT INTO ")
	bufferStmt.AppendString(tableName)
	bufferStmt.AppendString("(")
	bufferStmt.AppendString(bufferField.String())
	bufferStmt.AppendString(") ")
	bufferStmt.AppendString("VALUES (")
	bufferStmt.AppendString(bufferPlaceholder.String())
	bufferStmt.AppendString(");")

	sqlString := bufferStmt.String()
	db := engine.db //Get("")
	if engine == nil {
		return -1, fmt.Errorf("get db error")
	}
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return -1, fmt.Errorf("stmt error %s", err.Error())
	}
	result, err := stmt.Exec(values...)

	if printLog {
		fmt.Println("Prepare :", sqlString)
		fmt.Print("Exec ( ")
		flag := true
		for _, value := range values {
			if flag {
				fmt.Print(value)
				flag = false
			} else {
				fmt.Print(",", value)
			}
		}
		fmt.Print(" );\n")
	}

	if err != nil {
		return -1, fmt.Errorf("stmt error %s", err.Error())
	}

	return result.RowsAffected()
}

// FindOne
// command string, key string, object interface{}
func (s *BaseManager) FindOne(engine *Engine, obj interface{}, option *BaseOption, find interface{}) error {
	if engine == nil {
		return fmt.Errorf("get db error")
	}
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj kind must ptr")
	}
	if reflect.TypeOf(find).Kind() != reflect.Ptr {
		return errors.New("data kind must ptr")
	}
	tableName, err := entity.GetTableName(obj)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("get table name error %s", err.Error())
	}
	mp, err := entity.GetMapWithMORM(obj)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("get key/value error %s", err.Error())
	}
	bufferField := bufferpool.Get()
	defer bufferField.Free()
	bufferField.AppendString("SELECT ")
	flag := true
	for key, _ := range mp {
		if flag {
			bufferField.AppendString(key)
			flag = false
			continue
		}
		bufferField.AppendString(",")
		bufferField.AppendString(key)
	}
	bufferField.AppendString(" FROM ")
	bufferField.AppendString(tableName)
	whereSql, err := option.getWereSqlWord()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("get whereSql error %s", err.Error())
	}
	if len(*whereSql) == 0 {
		return fmt.Errorf("get whereSql error %s", "where sql length is null")
	}
	bufferField.AppendString(" ")
	bufferField.AppendString(*whereSql)
	bufferField.AppendString(";")

	db := engine.db
	if printLog {
		fmt.Println("sql: ", bufferField.String())
	}
	rows, err := db.Query(bufferField.String())
	if err != nil {
		return fmt.Errorf("query error %s", err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("rows Columns %s", err.Error())
	}
	if len(columns) != len(mp) {
		return fmt.Errorf("length not matching")
	}
	raws, err := mysqlRowsToRaws(rows)
	if err != nil {
		return err
	}
	if len(raws) <= 0 {
		return nil
	}

	for _, value := range raws {
		valueMap, err := rawsDataToMap(value, columns, obj)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(valueMap)
			jsonbody, err := json.Marshal(valueMap)
			if err != nil {
				fmt.Println(err)
			}
			if err := json.Unmarshal(jsonbody, find); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%#v\n", find)
		}
	}
	return nil
}

// FindAll
// command string, key string, object interface{}
func (s *BaseManager) FindAll(engine *Engine, obj interface{}, option *BaseOption, find interface{}) error {
	if engine == nil {
		return fmt.Errorf("get db error")
	}
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New("obj kind must ptr")
	}
	if reflect.TypeOf(find).Kind() != reflect.Ptr {
		return errors.New("data kind must ptr")
	}
	tableName, err := entity.GetTableName(obj)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("get table name error %s", err.Error())
	}

	mp, err := entity.GetMapWithMORM(obj)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("get key/value error %s", err.Error())
	}
	bufferField := bufferpool.Get()
	defer bufferField.Free()
	bufferField.AppendString("SELECT ")
	flag := true
	for key, _ := range mp {
		if flag {
			bufferField.AppendString(key)
			flag = false
			continue
		}
		bufferField.AppendString(",")
		bufferField.AppendString(key)
	}
	bufferField.AppendString(" FROM ")
	bufferField.AppendString(tableName)
	whereSql, err := option.getWereSqlWord()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("get whereSql error %s", err.Error())
	}
	if len(*whereSql) == 0 {
		return fmt.Errorf("get whereSql error %s", "where sql length is null")
	}
	bufferField.AppendString(" ")
	bufferField.AppendString(*whereSql)
	bufferField.AppendString(";")

	db := engine.db //Get("")
	if printLog {
		fmt.Println("sql: ", bufferField.String())
	}
	rows, err := db.Query(bufferField.String())
	if err != nil {
		return fmt.Errorf("query error %s", err.Error())
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("rows Columns %s", err.Error())
	}
	if len(columns) != len(mp) {
		return fmt.Errorf("length not matching")
	}
	raws, err := mysqlRowsToRaws(rows)
	if err != nil {
		return err
	}
	if len(raws) <= 0 {
		return nil
	}

	for _, value := range raws {
		valueMap, err := rawsDataToMap(value, columns, obj)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(valueMap)
			jsonBody, err := json.Marshal(valueMap)
			if err != nil {
				fmt.Println(err)
			}
			if err := json.Unmarshal(jsonBody, find); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%#v\n", find)
		}
	}
	return nil
}

// update
// delete

// mysql Raws
func mysqlRowsToRaws(rows *sql.Rows) ([][][]byte, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	length := len(columns)
	if length == 0 {
	}
	raws := make([][][]byte, 0)
	for rows.Next() {
		result := make([][]byte, length)
		dest := make([]interface{}, length)
		for i, _ := range result {
			dest[i] = &result[i]
		}
		if err = rows.Scan(dest...); err != nil {
			return nil, errors.New("scan error : " + err.Error())
		}
		raws = append(raws, result)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New("get rows error : " + err.Error())
	}
	return raws, nil
}

// Raws to object
func rawsDataToMap(data [][]byte, columns []string, obj interface{}) (map[string]interface{}, error) {

	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return nil, errors.New("obj kind must ptr")
	}
	typeOf := reflect.TypeOf(obj).Elem()
	valueOf := reflect.ValueOf(obj).Elem()

	obcKeyValueMap := make(map[string]interface{}, 0)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		fieldTagName := field.Tag.Get(entity.TagField)
		if valueOf.Field(i).CanSet() {
			for k, raw := range data {
				if fieldTagName == columns[k] {
					kind := valueOf.Field(i).Kind()
					if reflect.String == kind {
						obcKeyValueMap[fieldTagName] = string(raw)
					} else if reflect.Bool == kind {
						var v bool
						err := UnmarshalJSON(&v, raw)
						if err == nil {
							obcKeyValueMap[fieldTagName] = v
						} else {
							fmt.Println(fmt.Sprintf("resolution %s error %s", fieldTagName, err.Error()))
						}
					} else if reflect.Int == kind || reflect.Int8 == kind || reflect.Int16 == kind {
						var v int64
						err := UnmarshalJSON(&v, raw)
						if err == nil {
							obcKeyValueMap[fieldTagName] = v
						} else {
							fmt.Println(fmt.Sprintf("resolution %s error %s", fieldTagName, err.Error()))
						}
					} else if reflect.Int32 == kind || reflect.Int64 == kind {
						if ValidTagFieldDate(field.Tag.Get(entity.TagDate)) {
							var v int64
							err := UnmarshalJSON(&v, raw)
							if err == nil {
								obcKeyValueMap[fieldTagName] = v
							} else {
								fmt.Println(fmt.Sprintf("resolution %s error %s", fieldTagName, err.Error()))
							}
						} else {
							var v int64
							err := UnmarshalJSON(&v, raw)
							if err == nil {
								obcKeyValueMap[fieldTagName] = v
							} else {
								fmt.Println(fmt.Sprintf("resolution %s error %s", fieldTagName, err.Error()))
							}
						}
					} else if reflect.Float32 == kind || reflect.Float64 == kind {
						var v float64
						err := UnmarshalJSON(&v, raw)
						if err == nil {
							obcKeyValueMap[fieldTagName] = v
						} else {
							fmt.Println(fmt.Sprintf("resolution %s error %s", fieldTagName, err.Error()))
						}
					} else if reflect.Struct == kind {
						if field.Type.String() == "time.Time" {
							// todo time format change
							dateTimeString := string(raw)
							dateTimeString = strings.Replace(dateTimeString, " ", "T", 1)
							dateTimeString += "Z"
							obcKeyValueMap[fieldTagName] = dateTimeString
						}
					}
					//else if !valueOf.Field(i).IsNil() && (reflect.Slice == kind || reflect.Map == kind || reflect.Interface == kind) {
					//}
				}
			}
		}
	}
	return obcKeyValueMap, nil
}

// unary sql format
func sqlUnaryOption(option string, cmd *Command) string {
	key := cmd.key
	// todo interface to type
	value := reflect.ValueOf(cmd.value)
	kind := reflect.TypeOf(cmd.value).Kind()
	switch kind {
	case reflect.String:
		{
			return fmt.Sprintf(" %s %s %s", key, option, value.String())
		}
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		{
			return fmt.Sprintf(" %s %s %s", key, option, strconv.FormatInt(value.Int(), 10))
		}
	case reflect.Bool:
		{
			return fmt.Sprintf(" %s %s %s", key, option, strconv.FormatBool(value.Bool()))
		}
	case reflect.Float64, reflect.Float32:
		{
			return fmt.Sprintf(" %s %s %s", key, option, strconv.FormatFloat(value.Float(), 'g', 30, 32))
		}
	case reflect.Struct, reflect.Slice, reflect.Map, reflect.Interface, reflect.Array:
		{
			return ""
		}
	default:
		{
			return ""
		}
	}
}

// between format
func sqlBetweenOption(option string, cmd *Command) (string, error) {
	key := cmd.key
	values := cmd.values
	if values == nil || len(values) != 2 {
		return "", errors.New("between/notBetween values must not nil and length equal 2. ")
	}
	begin := cmd.values[0]
	end := cmd.values[1]
	if reflect.TypeOf(begin).Kind() != reflect.TypeOf(end).Kind() {
		return "", errors.New("between/notBetween values type must equal. ")
	}
	rBeginV := reflect.ValueOf(begin)
	rEndV := reflect.ValueOf(end)
	kind := reflect.TypeOf(begin).Kind()
	switch kind {
	case reflect.String:
		{
			return fmt.Sprintf(" %s %s %s AND %s", key, option, rBeginV.String(), rEndV.String()), nil
		}
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		{
			return fmt.Sprintf(" %s %s %d AND %d", key, option, rBeginV.Int(), rEndV.Int()), nil
		}
	case reflect.Float64, reflect.Float32:
		{
			return fmt.Sprintf(" %s %s %s AND %s", key, option, strconv.FormatFloat(rBeginV.Float(), 'g', 30, 32), strconv.FormatFloat(rEndV.Float(), 'g', 30, 32)), nil
		}
	default:
		{
			return "", errors.New("between/not between values type not match. ")
		}
	}
}

// in format
func sqlInOption(cmd *Command) (string, error) {
	key := cmd.key
	values := cmd.values
	if values == nil || len(values) < 1 {
		return "", errors.New("in values must not nil and length gte 1. ")
	}
	bufferField := bufferpool.Get()
	defer bufferField.Free()

	bufferField.AppendString(fmt.Sprintf(" %s IN (", key))
	tag := false
	for _, value := range values {
		rValue := reflect.ValueOf(value)
		kind := reflect.TypeOf(value).Kind()
		switch kind {
		case reflect.String:
			{
				if tag {
					bufferField.AppendString(fmt.Sprintf(", "))
				}
				bufferField.AppendString(fmt.Sprintf(" %s", rValue.String()))
			}
		case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
			{
				if tag {
					bufferField.AppendString(fmt.Sprintf(", "))
				}
				bufferField.AppendString(fmt.Sprintf(" %d", rValue.Int()))
			}
		case reflect.Float64, reflect.Float32:
			{
				if tag {
					bufferField.AppendString(fmt.Sprintf(", "))
				}
				bufferField.AppendString(fmt.Sprintf(" %f", rValue.Float()))
			}
		default:
			{

			}
		}
		tag = true
	}
	bufferField.AppendString(fmt.Sprintf(" )"))
	return bufferField.String(), nil
}

// UnmarshalJSON jason
func UnmarshalJSON(value interface{}, data []byte) error {
	err := json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	return nil
}

func ValidTagFieldIgnore(tagString string) bool {
	if tagString == entity.TagTrue {
		return true
	}
	return false
}

func ValidTagFieldDate(tagString string) bool {
	if tagString == entity.TagTrue {
		return true
	}
	return false
}
