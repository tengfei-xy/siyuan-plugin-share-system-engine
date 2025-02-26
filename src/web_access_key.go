package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tengfei-xy/go-log"
)

type accessKeyData struct {
	Appid         string `json:"appid"`
	Docid         string `json:"docid"`
	AccessKey     string `json:"accesskey"`
	PluginVersion string `json:"plugin_version"`
}

func (ak accessKeyData) startup() string {
	log.Infof("启动访问密码")
	log.Infof("Appid:%s", ak.Appid)
	log.Infof("Docid:%s", ak.Docid)
	if ak.AccessKey == "" {
		ak.AccessKey = createRand(4)
		log.Infof("AccessKey:%s (系统指定)", ak.AccessKey)
	} else {
		log.Infof("AccessKey:%s (用户指定)", ak.AccessKey)
	}
	log.Infof("插件版本:%s", ak.PluginVersion)

	var res resStruct

	_, err := app.db.Exec("update share set access_key=?,access_key_enable=1 where appid=? and docid=?", ak.AccessKey, ak.Appid, ak.Docid)
	if err != nil {
		log.Error(err)
		return res.setErrSystem().toString()
	}
	return res.setOK(ak.AccessKey).toString()

}
func (ak accessKeyData) close() string {
	log.Infof("关闭访问密码")
	log.Infof("Appid:%s", ak.Appid)
	log.Infof("Docid:%s", ak.Docid)
	log.Infof("插件版本:%s", ak.PluginVersion)

	var res resStruct

	_, err := app.db.Exec("update share set access_key_enable=0 where appid=? and docid=?", ak.Appid, ak.Docid)
	if err != nil {
		log.Error(err)
		return res.setErrSystem().toString()
	}
	return res.setOK("").toString()
}

func AccessKeyPOSTRequest(c *gin.Context) {

	var res resStruct

	var data accessKeyData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, res.setErrJson())
		return
	}

	var ret string
	act := c.Query("action")
	if act == "enable" {
		ret = data.startup()
	} else if act == "disable" {
		ret = data.close()
	} else {
		ret = res.setErrURL().toString()
	}

	c.String(http.StatusOK, ret)

}
func AccessKeyGetRequest(c *gin.Context) {

	type Data struct {
		AccessKey       string `json:"access_key"`
		AccessKeyEnable bool   `json:"access_key_enable"`
	}

	type getaccesskey struct {
		Err  int    `json:"err"`
		Msg  string `json:"msg"`
		Data Data   `json:"data"`
	}

	res := getaccesskey{
		Err: ERR_CODE_PARAM,
		Msg: ERR_MSG_PARAM,
		Data: Data{
			AccessKey:       "",
			AccessKeyEnable: false,
		},
	}
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	log.Info("-----------------")
	log.Infof("获取访问密码")

	var Appid = c.Query("appid")
	var Docid = c.Query("docid")
	if Appid == "" || Docid == "" {
		c.JSON(http.StatusBadRequest, res)
		return
	}
	log.Infof("Appid:%s", Appid)
	log.Infof("Docid:%s", Docid)

	var AccessKey string
	var AccessKeyEnable int

	if err := app.db.QueryRow(`select access_key,access_key_enable from share where appid=? and docid=? `, Appid, Docid).Scan(&AccessKey, &AccessKeyEnable); err != nil {
		if err == sql.ErrNoRows {
			log.Warn("没有记录")
			res.Err = ERR_CODE_NO_KEY
			res.Msg = ERR_MSG_NO_KEY
			c.JSON(http.StatusOK, res)
			return
		}
		log.Error(err)
		c.JSON(http.StatusBadRequest, res)
		return
	}
	res.Err = ERR_CODE_OK
	res.Msg = ERR_MSG_OK
	res.Data.AccessKey = AccessKey

	if AccessKeyEnable == 0 {
		res.Data.AccessKeyEnable = false
	} else {
		res.Data.AccessKeyEnable = true
	}

	c.JSON(http.StatusOK, res)

}
