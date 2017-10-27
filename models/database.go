package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)


var logger = logs.NewLogger()

func SyncDataBase()  {
	orm.Debug = true
	createDatabase()
	connDatabase()

	// 数据库别名
	name := "default"
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := true
	// 遇到错误立即返回
	err := orm.RunSyncdb(name, force, verbose)
	if err != nil {
		logger.Debug("%v", err)
	}
	logger.Debug("database init is complete.")
}

func connDatabase()  {
	// 配置数据库连接
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	var dns string
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	logger.Debug("db_name= %s", db_name)
	switch db_type {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", db_user, db_pass, db_host, db_port, db_name)
	default:
		logger.Critical("Database driver is not allowed:", db_type)
	}
	logger.Debug("dns= %v", dns)
	orm.RegisterDataBase("default", db_type, dns)
	orm.SetMaxIdleConns("default", 30)
	orm.SetMaxOpenConns("default", 30)
	orm.RegisterModel(new(UserInfo), new(Administrator), new(UserLogin), new(UserAuth))
	orm.RegisterModel(new(SubSystem), new(SubProperties))
	//orm.RegisterModel(new(Session))

}

func createDatabase() {
	db_type := beego.AppConfig.String("db_type")
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	var dns string
	var sqlString string
	switch db_type {
	case "mysql":
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8", db_user, db_pass, db_host, db_port)
		sqlString = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
	default:
		logger.Critical("Database driver is not allowed:", db_type)
	}
	db, err := sql.Open(db_type, dns)
	if err != nil {
		logger.Emergency("can not open mysql connection!")
		panic(err.Error())
	}
	defer db.Close()

	_, err = db.Exec(sqlString)
	if err != nil {
		logger.Error("create db failed!", err)
	} else {
		logger.Debug("Database %s created", db_name)
	}
}

func init()  {
	logger.Async()
	logger.EnableFuncCallDepth(true)
	if beego.BConfig.RunMode == "dev" {
		logger.SetLogger(logs.AdapterConsole)
	} else {
		logger.SetLogger(logs.AdapterFile, `{"filename":"database.log","daily":true,"maxdays":10}`)
	}
}