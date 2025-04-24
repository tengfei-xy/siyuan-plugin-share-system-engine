package web

import (
	"net/http"
	"sys"

	"github.com/gin-gonic/gin"
	"github.com/tengfei-xy/go-log"
)

func env(app *sys.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("ShareBaseLink", app.Basic.ShareBaseLink)
		c.Set("version", app.Basic.Version)
		c.Set("IsPublicServer", app.Basic.IsPublicServer)
		c.Set("SavePath", app.Basic.SavePath)
		c.Set("FileMaxMB", app.Web.FileMaxMB)
		c.Set("Token", app.Basic.Token)
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
func setToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.MustGet("Token")
		if token == "" {
			c.Next()
			return
		}
		// 解析token
		tk := c.Request.Header.Get("Authorization")
		if tk == "" || tk != token {
			log.Debug("unauthorized")
			unauthorized(c)
			c.Abort()
			return
		}
	}
}
