package model

import (
	"chatroom/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"go.uber.org/zap"
)

// User ...
type User struct {
	Name string
}

// user  user table
type user struct {
	ID     int64 `xorm:"pk autoincr notnull unique 'id' int "`
	Name   string
	Passwd string
}

// newEg 创建数据库引擎
func (u User) newEg() (engine *xorm.Engine, err error) {
	engine, err = xorm.NewEngine("mysql", "root:Tyx123456.@/test?charset=utf8")
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, "create engine fail"); ce != nil {
			ce.Write(zap.Error(err))
		}
		return
	}
	ok, err := engine.IsTableExist(&user{})
	if err != nil {
		return
	}
	if ok {
		return
	}
	err = engine.CreateTables(&user{})
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
	engine, err := u.newEg()
	if err != nil {
		return
	}
	_, err = engine.Insert(&user{
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
	eg, err := u.newEg()
	ok, err = eg.Where("name = ?", u.Name).Get(&user{})
	if err != nil {
		return
	}
	return
}
