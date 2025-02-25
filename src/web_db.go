package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func dbTableGETRequest(c *gin.Context) {
	r, err := dbTableGET()
	if err != nil {
		internalSystem(c)
	}
	c.String(http.StatusOK, r)
}

func dbTableDeleteRequest(c *gin.Context) {
	err := dbTableDelete()
	if err != nil {
		internalSystem(c)
		return
	}
	var res resStruct
	c.JSON(http.StatusOK, res.setOK("已删除"))
}

func dbTablePUTRequest(c *gin.Context) {
	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}
	// 保存文件
	err = c.SaveUploadedFile(file, "upload.sql")
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	if err := dbReset("upload.sql"); err != nil {
		internalSystem(c)
		return
	}

	r, err := dbTableGET()
	if err != nil {
		internalSystem(c)
		return
	}
	c.String(http.StatusOK, r)
}
