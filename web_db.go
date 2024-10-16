package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func dbTableGETRequest(c *gin.Context) {
	r, err := dbTableGET()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(http.StatusOK, r)
}

func dbTableDeleteRequest(c *gin.Context) {
	err := dbTableDelete()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var res resStruct
	c.JSON(http.StatusOK, res.setOK("已删除"))
}

func dbTablePUTRequest(c *gin.Context) {
	var res resStruct
	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, res.setErrSystem())
		return
	}
	// 保存文件
	err = c.SaveUploadedFile(file, "upload.sql")
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	if err := dbReset("upload.sql"); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	r, err := dbTableGET()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, r)
}
