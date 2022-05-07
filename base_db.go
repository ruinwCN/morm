package morm

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ruinwCN/morm/conf"
	"sync"
	"time"
)

var DefaultDBTag = "_defaultDB"

var (
	engines  *Engines
	printLog = false
)

type Engine struct {
	db     *sql.DB
	option conf.DBConfig
}

type Engines struct {
	sync.RWMutex
	engineMap map[string]*Engine
}

func (s *Engines) Add(name string, e *sql.DB, option conf.DBConfig) error {
	s.Lock()
	defer s.Unlock()
	if s.engineMap == nil {
		return errors.New("uninitialized")
	}
	if s.engineMap[name] != nil {
		return errors.New(name + " already exists")
	}
	optionCopy := option
	s.engineMap[name] = &Engine{
		e, optionCopy,
	}
	return nil
}

func (s *Engines) Get(name string) (*Engine, error) {
	s.RLock()
	defer s.RUnlock()
	if s.engineMap == nil {
		return nil, errors.New("uninitialized")
	}
	if name == "" {
		name = DefaultDBTag
	}
	if s.engineMap[name] == nil {
		return nil, errors.New("not find " + name)
	}
	return s.engineMap[name], nil
}

func NewEngines(mysqlConfigs []conf.DBConfig) error {
	once := sync.Once{}
	once.Do(func() {
		engines = &Engines{
			engineMap: make(map[string]*Engine),
		}
	})
	for _, configEntity := range mysqlConfigs {
		if configEntity.Address == "" {
			return fmt.Errorf("error, config address is empty\n")
		}
		if configEntity.Database == "" {
			return fmt.Errorf("error, config database name is empty\n")
		}
		if configEntity.Tag == "" {
			configEntity.Tag = DefaultDBTag
		}
		dsn := configEntity.DSN()
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		db.SetConnMaxLifetime(time.Second * time.Duration(configEntity.MaxLifeTime))
		db.SetMaxOpenConns(configEntity.MaxOpenConns)
		db.SetMaxIdleConns(configEntity.MaxIdleConns)
		if pingErr := db.Ping(); pingErr != nil {
			return fmt.Errorf("ping db error :%w", pingErr)
		}
		err = engines.Add(configEntity.Tag, db, configEntity)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetDB(name string) (*Engine, error) {
	if engines == nil {
		return nil, errors.New("not fund this db")
	}
	return engines.Get(name)
}

func PrintLog(b bool) {
	printLog = b
}
