package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	log "github.com/tengfei-xy/go-log"
)

var app appConfig

const version string = "0.1.0"

func init_flag() flagStruct {
	var f flagStruct
	flag.StringVar(&f.config_file, "c", "config.yaml", "打开配置文件")
	flag.Parse()
	return f
}

func init_mysql() {
	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", app.Mysql.Username, app.Mysql.Password, app.Mysql.Ip, app.Mysql.Port, app.Mysql.Database))
	if err != nil {
		panic(err)
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		panic(fmt.Sprintf("数据库连接失败%v", err))
	}
	log.Info("数据库已连接")
	app.db = DB
}

func main() {

	log.Infof("版本:v%s", version)
	f := init_flag()
	init_config(f)
	init_mysql()
	init_web()
}
