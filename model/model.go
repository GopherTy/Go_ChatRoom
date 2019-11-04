package model

import (
	"chatroom/logger"

	"github.com/go-xorm/xorm"
	"go.uber.org/zap"
)

// User ...
type User struct {
	Name   string
	Passwd string
}

// user  user table
type user struct {
	ID     int64  `xorm:"pk autoincr notnull unique 'id' int "`
	Name   string `xorm:"notnull"`
	Passwd string
}

//
var _Engine *xorm.Engine

func init() {
	engine, err := xorm.NewEngine("mysql", "root:Tyx123456.@/test?charset=utf8")
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create engine fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	_Engine = engine
	u := User{}
	err = CreateTable(u)
	if err != nil {
		return
	}
}

// CreateTable 创建数据库表
func CreateTable(table interface{}) (err error) {
	ok, err := _Engine.IsTableExist(table)
	if err != nil {
		return
	}
	if ok {
		err = _Engine.Sync2(table)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "create uniques fail"); ce != nil {
				ce.Write(zap.Error(err))
			}
			return
		}
		err = _Engine.CreateUniques(table)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "create uniques fail"); ce != nil {
				ce.Write(zap.Error(err))
			}
			return
		}
		err = _Engine.CreateIndexes(table)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, "create indexes fail"); ce != nil {
				ce.Write(zap.Error(err))
			}
			return
		}
		return
	}
	err = _Engine.CreateTables(table)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create table fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	err = _Engine.CreateUniques(table)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create uniques fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	err = _Engine.CreateIndexes(table)
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create indexes fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	return
}

// SignUp ...
func (u User) SignUp(name string) (err error) {
	_, err = _Engine.Insert(&user{
		Name: name,
	})
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "insert data fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	return
}

// Check 用户验证
func (u User) Check() (ok bool, err error) {
	ok, err = _Engine.Where("name=? And passwd = ?", u.Name, u.Passwd).Get(&user{})
	if err != nil {
		return
	}
	return
}
