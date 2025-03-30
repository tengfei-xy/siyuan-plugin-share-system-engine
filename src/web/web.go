package web

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"sqlite"
	"sys"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
)

func Init(app *sys.Config) {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()

	// 设置接口密码
	if app.APIEnable {
		if app.SQL.Username != "" {
			log.Infof("启动API的DB接口，启动接口密码保护")
			a := g.Group("/api/db", gin.BasicAuth(gin.Accounts{
				app.Username: app.Password,
			}))
			a.GET("/table", dbTableGETRequest)
			a.DELETE("/table", dbTableDeleteRequest)
			a.PUT("/table", dbTablePUTRequest)
			// authorized.GET("/record", dbRecoreGETRequest)
		} else {
			log.Infof("启动API的DB接口，关闭接口密码保护")
			g.GET("/api/db/table", dbTableGETRequest)
			g.DELETE("/api/db/table", dbTableDeleteRequest)
			g.PUT("/api/db/table", dbTablePUTRequest)
		}
	} else {
		log.Infof("关闭API的DB接口")
	}
	g.Use(cors())
	g.Use(env(app))
	g.POST("/api/v2/link", v2PostLinkRequest)
	g.POST("/api/v2/home_page", v2PostHomePageRequest)
	g.DELETE("/api/v2/home_page", v2DeleteHomePageRequest)
	g.POST("/api/upload_args", uploadArgsRequest)
	g.POST("/api/upload_file", uploadFileRequest)
	g.POST("/api/getlink", getLinkRequest)
	g.GET("/api/getlinkall", getLinkAllRequest)
	g.POST("/api/deletelink", deleteLinkRequest)
	g.GET("/api/info", infoRequest)
	g.POST("/api/key", AccessKeyPOSTRequest)
	g.GET("/api/key", AccessKeyGetRequest)
	g.GET("/html/:appid/:docid/*filepath", htmlRequest)
	g.GET("/:id", linkRequest)
	g.GET("/", rootRequest)

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
func env(app *sys.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ShareBaseLink", app.Basic.ShareBaseLink)
		c.Set("version", app.Basic.Version)
		c.Set("IsPublicServer", app.Basic.IsPublicServer)
		c.Set("SavePath", app.Basic.SavePath)
		c.Set("FileMaxMB", app.Web.FileMaxMB)
	}
}
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Cookie, cros-status")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			log.Info("-----------------")
			log.Info("预检")
			log.Infof("原始: %s", c.Request.Header.Get("origin"))
			c.Status(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func linkRequest(c *gin.Context) {
	id := c.Params.ByName("id")

	log.Info("-----------------")
	log.Info("请求链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", id)

	row := sqlite.DB.QueryRow(`select appid,docid,status from share where link=?`, id)
	var appid, docid string
	var status int

	// 扫描数据库
	if err := row.Scan(&appid, &docid, &status); err != nil {
		if err == sql.ErrNoRows {
			noShare(c)
			return
		} else {
			internalSystem(c)
			return
		}
	}

	// 如果数据库中设置了禁止访问
	if status == sys.STATUS_LINK_DISABLE {
		noShare(c)
		return
	}

	// 访问 + 1
	_, err := sqlite.DB.Exec(`update share set count=count+1 where link=?`, id)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}
	sbl := c.MustGet("ShareBaseLink")
	new_url := fmt.Sprintf("%s/html/%s/%s/index.htm", sbl, appid, docid)
	//重定向
	log.Infof("重定向: %s", new_url)
	c.Redirect(http.StatusMovedPermanently, new_url)
}
func rootRequest(c *gin.Context) {
	sbl := c.MustGet("ShareBaseLink")

	// 从数据库中读取home_apge
	var res resStruct

	row := sqlite.DB.QueryRow(`select appid,docid,status,link from share where home_page=1`)
	var appid, docid, link string
	var status int
	if err := row.Scan(&appid, &docid, &status, &link); err != nil {
		if err == sql.ErrNoRows {
			log.Warn("首页未设置")
			c.JSON(http.StatusNotFound, res.setNoPage())
			return
		} else {
			internalSystem(c)
			return
		}
	}

	// 如果数据库中设置了禁止访问
	if status == sys.STATUS_LINK_DISABLE {
		noShare(c)
		return
	}

	// 访问 + 1
	_, err := sqlite.DB.Exec(`update share set count=count+1 where link=?`, link)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	new_url := fmt.Sprintf("%s/html/%s/%s/index.htm", sbl, appid, docid)
	//重定向
	log.Infof("重定向: %s", new_url)
	c.Redirect(http.StatusMovedPermanently, new_url)

}

func infoRequest(c *gin.Context) {
	isps := c.MustGet("IsPublicServer").(bool)
	v := c.MustGet("version").(string)

	type Resquest struct {
		Version        string `json:"version"`
		IsPublicServer bool   `json:"is_public_server"`
	}
	var res Resquest

	res.IsPublicServer = isps
	res.Version = v
	c.JSON(http.StatusOK, msgOK(res))
}
func htmlRequest(c *gin.Context) {
	sp := c.MustGet("SavePath").(string)

	appid := c.Params.ByName("appid")
	docid := c.Params.ByName("docid")
	param_filename := c.Params.ByName("filepath")
	filename := filepath.Join(sp, appid, docid, param_filename)
	if param_filename != "/index.htm" {
		c.File(filename)
		return
	}

	var access_key string
	var access_key_enable int

	err := sqlite.DB.QueryRow(`select access_key,access_key_enable from share where appid=? and docid=?`, appid, docid).Scan(&access_key, &access_key_enable)
	if err != nil {
		log.Error(err)
		internalSystem(c)
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
