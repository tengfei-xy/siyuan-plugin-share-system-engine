package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func init_web() {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()

	// 设置接口密码
	if app.APIEnable {
		if app.SQL.Username != "" {
			log.Infof("启动关闭API的DB接口，启动接口密码保护")
			authorized := g.Group("/api/db", gin.BasicAuth(gin.Accounts{
				app.Username: app.Password,
			}))
			authorized.GET("/table", dbTableGETRequest)
			authorized.DELETE("/table", dbTableDeleteRequest)
			authorized.PUT("/table", dbTablePUTRequest)
			// authorized.GET("/record", dbRecoreGETRequest)
		} else {
			log.Infof("启动API的DB接口，关闭接口密码保护")
			log.Infof("API接口密码保护")
			g.GET("/api/db/table", dbTableGETRequest)
			g.DELETE("/api/db/table", dbTableDeleteRequest)
			g.PUT("/api/db/table", dbTablePUTRequest)
		}
	} else {
		log.Infof("关闭API的DB接口")
	}

	g.POST("/api/upload_args", uploadArgsRequest)
	g.POST("/api/upload_file", uploadFileRequest)
	g.POST("/api/getlink", getLinkRequest)
	g.GET("/api/getlinkall", getLinkAllRequest)
	g.POST("/api/deletelink", deleteLinkRequest)

	g.POST("/api/key", AccessKeyPOSTRequest)
	g.GET("/api/key", AccessKeyGetRequest)
	g.GET("/html/:appid/:docid/*filepath", htmlRequest)
	g.GET("/:id", rootRequest)

	g.OPTIONS("/api/getlink", optionRequest)
	g.OPTIONS("/api/upload_args", optionRequest)
	g.OPTIONS("/api/upload_file", optionRequest)
	g.OPTIONS("/api/deletelink", optionRequest)
	g.OPTIONS("/api/url/:url", optionRequest)
	g.OPTIONS("/api/key", optionRequest)
	g.OPTIONS("/html/:appid/:docid/*filepath", optionRequest)
	g.OPTIONS("/:id", optionRequest)

	g.MaxMultipartMemory = app.Web.FileMaxMB << 20 // 100 MiB
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	if app.Web.SSLEnable {
		log.Infof("启动https服务器，监听地址: %s", app.Basic.ListenPort)
		if err := g.RunTLS(app.Basic.ListenPort, app.Web.SSLCERT, app.Web.SSLKEY); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Infof("启动http服务器，监听地址: %s", app.Basic.ListenPort)
		err := g.Run(app.Basic.ListenPort)
		if err != nil {
			log.Fatal(err)
		}
	}
}
func optionRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("预检")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("原始: %s", c.Request.Header.Get("origin"))
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "content-type, cros-status")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
}
func rootRequest(c *gin.Context) {
	id := c.Params.ByName("id")

	log.Info("-----------------")
	log.Info("请求链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", id)

	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
	var res resStruct
	row := app.db.QueryRow(`select appid,docid,status from share where link=?`, id)
	var appid, docid string
	var status int

	// 扫描数据库
	if err := row.Scan(&appid, &docid, &status); err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, res.setNoShare().toMsg())
			return
		} else {
			c.String(http.StatusOK, res.setErrSystem().toMsg())
			return
		}
	}

	// 如果数据库中设置了禁止访问
	if status == STATUS_LINK_DISABLE {
		c.String(http.StatusNotFound, res.setNoShare().toMsg())
		return
	}

	// 访问 + 1
	_, err := app.db.Exec(`update share set count=count+1 where link=?`, id)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	new_url := fmt.Sprintf("%s/html/%s/%s/index.htm", app.Basic.ShareBaseLink, appid, docid)
	//重定向
	log.Infof("重定向: %s", new_url)
	c.Redirect(http.StatusMovedPermanently, new_url)

}
func htmlRequest(c *gin.Context) {

	appid := c.Params.ByName("appid")
	docid := c.Params.ByName("docid")
	param_filename := c.Params.ByName("filepath")
	filename := filepath.Join(app.SavePath, appid, docid, param_filename)
	if param_filename != "/index.htm" {
		c.File(filename)
		return
	}

	var res resStruct
	var access_key string
	var access_key_enable int

	err := app.db.QueryRow(`select access_key,access_key_enable from share where appid=? and docid=?`, appid, docid).Scan(&access_key, &access_key_enable)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	if access_key_enable == 1 {
		log.Infof("此页面需要访问密码")
		pak := c.Query("access_key")
		if pak != access_key {
			c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(access_key_html))
			return
		}
	}
	c.File(filename)
}
