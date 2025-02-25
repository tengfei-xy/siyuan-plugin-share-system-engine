package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tengfei-xy/go-log"
)

func v2PostHomePageRequest(c *gin.Context) {

	type request struct {
		Appid         string `json:"appid"`
		Docid         string `json:"docid"`
		PluginVersion string `json:"plugin_version"`
	}
	type response struct {
		Link string `json:"link"`
	}
	var err error
	var req request

	log.Info("-----------------")
	log.Info("创建首页")
	if err = c.ShouldBindBodyWithJSON(&req); err != nil {
		log.Error(err)
		badRequest(c)
		return
	}

	if req.Appid == "" || req.Docid == "" {
		badRequest(c)
		return
	}

	log.Infof("appid: %s", req.Appid)
	log.Infof("docid: %s", req.Docid)
	log.Infof("plugin_version: %s", req.PluginVersion)

	// 设置为首页
	_, err = app.db.Exec("update share set home_page=0 where home_page=1")
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	_, err = app.db.Exec("update share set home_page=1 where appid=? and docid=?", req.Appid, req.Docid)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}
	row := app.db.QueryRow(`select link from share where home_page=1`)
	var link string
	if err := row.Scan(&link); err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	okData(c, response{
		Link: app.Basic.ShareBaseLink,
	})

}
func v2DeleteHomePageRequest(c *gin.Context) {
	var err error
	type resquest struct {
		Link string `json:"link"`
	}
	log.Info("-----------------")
	log.Info("删除首页")

	row := app.db.QueryRow(`select link from share where home_page=1`)
	var link string
	if err := row.Scan(&link); err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	_, err = app.db.Exec("update share set home_page=0 where home_page=1")
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	okData(c, resquest{
		Link: fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link),
	})
}
