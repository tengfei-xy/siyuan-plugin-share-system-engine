package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	log "github.com/tengfei-xy/go-log"
)

func init_db(flag flagStruct) {
	DB, err := sql.Open("sqlite3", app.SYSFilename)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("数据库已连接,数据文件路径: %s", app.SYSFilename)
	app.db = DB
	if flag.db_file != "" {
		if err := dbReset(flag.db_file); err != nil {
			log.Fatal(err)
		}
	}
	// 如果不存在表则创建表
	db_check_table()
	// 更新表结构(access_key字段)
	db_check_struct_access_key()
	// 更新表结构(home_page字段)
	db_check_struct_home_page()

}
func dbReset(filename string) error {
	// 删除表
	dbTableDelete()

	// 如果不存在表则创建表
	db_check_table()

	// 根据提供的数据文件路径导入其中的INSERT语句
	err := db_upload(filename)
	if err != nil {
		log.Fatal(err)
	}

	// 更新表结构
	db_check_struct_home_page()
	return err
}
func db_check_table() error {
	var g int

	err := app.db.QueryRow(`SELECT EXISTS (SELECT 1 FROM sqlite_master WHERE type = 'table' AND name = 'share')`).Scan(&g)
	if err != nil {
		log.Fatal(err)
	}
	if g == 1 {
		return nil
	}

	log.Infof("数据库无数据！")
	if _, err := app.db.Exec(`CREATE TABLE share (	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,	appid VARCHAR(40),	docid VARCHAR(30),	title TEXT,	link VARCHAR(30),	update_time DATETIME DEFAULT CURRENT_TIMESTAMP,	expire_time DATETIME,	count INTEGER DEFAULT 0,	status INTEGER DEFAULT 0,	UNIQUE(appid, docid)  )`); err != nil {
		return err
	}
	log.Info("添加数据表完成")
	return nil

}
func db_check_struct_access_key() error {
	var g = 0
	if err := app.db.QueryRow(`SELECT EXISTS( SELECT 1 FROM pragma_table_info('share') WHERE name = 'access_key')`).Scan(&g); err != nil {
		return err
	}
	if g == 1 {
		return nil
	}

	if _, err := app.db.Exec(`alter table share ADD COLUMN access_key VARCHAR(10) DEFAULT ''`); err != nil {
		return err

	}
	if _, err := app.db.Exec(`alter table share ADD COLUMN access_key_enable INTEGER DEFAULT 0 `); err != nil {
		return err
	}
	return nil
}
func db_check_struct_home_page() error {
	var g = 0
	if err := app.db.QueryRow(`SELECT EXISTS( SELECT 1 FROM pragma_table_info('share') WHERE name = 'home_page')`).Scan(&g); err != nil {
		return err
	}
	if g == 1 {
		return nil
	}

	if _, err := app.db.Exec(`alter table share ADD COLUMN home_page INTEGER DEFAULT 0`); err != nil {
		return err

	}
	return nil
}
func db_upload(filename string) error {
	log.Infof("上传文件来自:%s", filename)
	bf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	var f = string(bf)
	var count = 0
	for _, line := range strings.Split(f, "\n") {
		if strings.HasPrefix(line, "INSERT") || strings.HasPrefix(line, "insert") {
			if _, err := app.db.Exec(line); err != nil {
				return err
			}
			count++
			log.Info(line)
		}
	}
	if count == 0 {
		return fmt.Errorf("该文件未找到INSERT语句")
	}
	db_check_struct_access_key()
	db_check_struct_home_page()
	return nil
}
func dbTableGET() (string, error) {
	ret, err := app.db.Query("select id,appid,docid,title,link,update_time,expire_time,count,access_key,access_key_enable from share")
	if err != nil {
		return "", err
	}
	var id, count, access_key_enable int
	var appid, docid, title, link, access_key string
	var update_time, expire_time sql.NullTime
	var buf strings.Builder

	buf.WriteString("id,appid,docid,title,link,update_time,expire_time,count,access_key,access_key_enable\n")

	for ret.Next() {
		err := ret.Scan(&id, &appid, &docid, &title, &link, &update_time, &expire_time, &count, &access_key, &access_key_enable)
		if err != nil {
			return "", err
		}

		buf.WriteString(fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%d,%s,%d\n", id, appid, docid, title, link, update_time.Time.String(), expire_time.Time.String(), count, access_key, access_key_enable))
	}
	return buf.String(), nil
}
func dbTableDelete() error {
	_, err := app.db.Exec("DROP TABLE share")
	if err != nil {
		return err
	}
	return nil
}
