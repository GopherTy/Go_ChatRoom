package model

import (
	"chatroom/logger"

	_ "github.com/go-sql-driver/mysql"
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
	ID     int64 `xorm:"pk autoincr notnull unique 'id' int "`
	Name   string
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
}

// createTable 创建数据库表
func (u User) createTable() (err error) {
	ok, err := _Engine.IsTableExist(&user{})
	if err != nil {
		return
	}
	if ok {
		return
	}
	err = _Engine.CreateTables(&user{})
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create table fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	return
}

// SignUp ...
func (u User) SignUp(name string) (err error) {
	err = u.createTable()
	if err != nil {
		return
	}
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
	err = u.createTable()
	ok, err = _Engine.Where("name = ?", u.Name).And("passwd = ?", u.Passwd).Get(&user{})
	if err != nil {
		return
	}
	return
}
