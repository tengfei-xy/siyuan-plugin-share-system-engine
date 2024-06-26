package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
	"github.com/tengfei-xy/go-tools"
)

type uploadArgsReq struct {
	Appid           string `json:"appid"`
	Docid           string `json:"docid"`
	Content         string `json:"content"`
	Theme           string `json:"theme"`
	SiyuanVersion   string `json:"version"`
	Title           string `json:"title"`
	HideSYVersion   bool   `json:"hide_version"`
	PluginVersion   string `json:"plugin_version"`
	ExistLinkCreate bool   `json:"exist_link_create"`
	PageWide        string `json:"page_wide"`
}

type getLinkReq struct {
	Appid         string `json:"appid"`
	Docid         string `json:"docid"`
	PluginVersion string `json:"plugin_version"`
}

type deleteLinkReq struct {
	Appid         string `json:"appid"`
	Docid         string `json:"docid"`
	PluginVersion string `json:"plugin_version"`
}
type resStruct struct {
	Err  int    `json:"err"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}
type htmlTemplate struct {
	Content  string
	Theme    string
	Version  string
	Resource string
	Title    string
}

func (r *resStruct) setErrJson() *resStruct {
	r.Err = 1
	r.Msg = "json解析错误"
	r.Data = ""
	return r

}
func (r *resStruct) setErrSystem() *resStruct {
	r.Err = 2
	r.Msg = "系统错误"
	r.Data = ""
	return r

}
func (r *resStruct) setNoShare() *resStruct {
	r.Err = 3
	r.Msg = "此页面没有共享"
	r.Data = ""
	return r

}

func (r *resStruct) setOK(str string) *resStruct {
	r.Err = 0
	r.Msg = "处理完成"
	r.Data = str
	return r
}
func (r *resStruct) setErrParam() *resStruct {
	r.Err = 4
	r.Msg = "参数错误"
	r.Data = ""
	return r
}

// json to string
func (r *resStruct) toString() string {
	v, err := json.Marshal(r)
	if err != nil {
		r.setErrJson()
	}
	return string(v)
}

// json to msg
func (r *resStruct) toMsg() string {
	return r.Msg
}
func init_web() {
	g := gin.Default()

	g.POST("/api/upload_args", uploadArgsRequest)
	g.POST("/api/upload_file", uploadFileRequest)
	g.POST("/api/getlink", getLinkRequest)
	g.POST("/api/deletelink", deleteLinkRequest)
	g.GET("/api/url/:url", shareRequest)

	g.OPTIONS("/api/getlink", optionRequest)
	g.OPTIONS("/api/upload_args", optionRequest)
	g.OPTIONS("/api/upload_file", optionRequest)
	g.OPTIONS("/api/deletelink", optionRequest)
	g.OPTIONS("/api/url/:url", optionRequest)

	g.MaxMultipartMemory = 100 << 20 // 100 MiB
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	log.Infof("服务器启动，监听地址 %s", app.Basic.ListenPort)
	g.Run(app.Basic.ListenPort)
}

// 处理函数，
func shareRequest(c *gin.Context) {
	param := c.Params.ByName("url")

	log.Info("-----------------")
	log.Info("请求链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())
	log.Infof("参数: %s", param)
	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
	var res resStruct
	row := app.db.QueryRow(`select appid,docid,status from share where link=?`, param)
	var appid, docid string
	var status int
	if err := row.Scan(&appid, &docid, &status); err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, res.setNoShare().toMsg())
			return
		} else {
			c.String(http.StatusOK, res.setErrSystem().toMsg())
			return
		}
	}
	if status == STATUS_LINK_DISABLE {
		c.String(http.StatusNotFound, res.setNoShare().toMsg())
		return
	}

	// 访问 + 1
	_, err := app.db.Exec(`update share set count=count+1 where link=?`, param)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	//
	new_url := fmt.Sprintf("%s/%s/%s/index.html", app.Basic.ShareBaseLink, appid, docid)
	//重定向
	log.Infof("重定向: %s", new_url)
	c.Redirect(http.StatusMovedPermanently, new_url)
}
func uploadFileRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("上传文件")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())

	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var res resStruct
	appid := c.Query("appid")
	docid := c.Query("docid")
	log.Infof("appid: %s", appid)
	log.Infof("docid: %s", docid)

	if len(appid) == 0 || len(docid) == 0 {
		c.String(http.StatusOK, res.setErrParam().toString())
		return

	}

	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 创建资源文件夹
	f, err := mkdir_all(appid, docid)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	zip_file := filepath.Join(f, "resources.zip")
	// 保存文件
	err = c.SaveUploadedFile(file, zip_file)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	log.Infof("资源文件: %s", zip_file)

	// 如果压缩包不存在
	if !tools.FileExist(zip_file) {
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 解压文件
	if err := unzip(f, zip_file); err != nil {
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	check_theme_file(f)
	c.String(http.StatusOK, res.setOK("上传文件成功").toString())
	return
}
func uploadArgsRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("上传参数")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())

	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var res resStruct
	// 获取请求数据
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 解析请求数据
	var data uploadArgsReq
	if err := json.Unmarshal(b, &data); err != nil {
		c.String(http.StatusOK, res.setErrJson().toString())
		return
	}

	log.Infof("appid: %s", data.Appid)
	log.Infof("docid: %s", data.Docid)
	log.Infof("theme: %s", data.Theme)
	log.Infof("title: %s", data.Title)
	// 不输出文档的内容，格式为html
	// log.Infof("content: %s", data.Content)
	log.Infof("思源版本: %s", data.SiyuanVersion)
	log.Infof("插件版本: %v", data.PluginVersion)
	log.Infof("页面宽度: %s", data.PageWide)
	log.Infof("标题中隐藏思源版本: %v", data.HideSYVersion)
	log.Infof("链接存在时重新创建: %v", data.ExistLinkCreate)
	// 创建资源文件夹，
	// 返回资源文件夹的路径,作为生成的html文件的存放路径
	tmp_html, err := mkdir_all(data.Appid, data.Docid)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	title_version := func() string {
		if data.HideSYVersion {
			return ""
		} else {
			return "   v" + data.SiyuanVersion

		}
	}

	content := tempate_html

	content = strings.ReplaceAll(content, "{{ .Theme }}", data.Theme)
	content = strings.ReplaceAll(content, "{{ .Version }}", data.SiyuanVersion)
	content = strings.ReplaceAll(content, "{{ .Title }}", data.Title)
	content = strings.ReplaceAll(content, "{{ .Content }}", data.Content)
	content = strings.ReplaceAll(content, "{{ .TitleVersion }}", title_version())

	content = strings.ReplaceAll(content, "{{ .PageWide }}", func(wide string) string {
		if strings.HasSuffix(wide, "%") {
			numstr := strings.TrimRight(wide, "%")
			num, err := strconv.Atoi(numstr)
			if (err != nil) || (num < 0) || (num > 100) {
				return "800px"
			}
		} else if strings.HasSuffix(wide, "px") {
			numstr := strings.TrimRight(wide, "px")
			num, err := strconv.Atoi(numstr)
			if (err != nil) || (num < 0) {
				return "800px"
			}
		} else {
			log.Warnf("宽度格式错误，设置为默认宽度，参数值:%s ", wide)
			return "800px"
		}
		return wide

	}(data.PageWide))

	f, err := os.OpenFile(filepath.Join(tmp_html, "index.html"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, get_file_permission())
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 根据appid和docid查看链接是否存在
	row := app.db.QueryRow(`select link from share where appid=? and docid=? `, data.Appid, data.Docid)
	var link string

	if err := row.Scan(&link); err != nil {
		// 如果不存在就插入
		if err == sql.ErrNoRows {
			link := createRand()
			full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)
			_, err := app.db.Exec("INSERT INTO share(appid,docid,title,link) VALUES(?,?,?,?)  ", data.Appid, data.Docid, data.Title, link)
			if err != nil {
				log.Error(err)
				c.String(http.StatusOK, res.setErrSystem().toString())
				return
			}

			log.Infof("创建参数成功")
			log.Infof("创建分享链接: %s", full_link)
			c.String(http.StatusOK, res.setOK(full_link).toString())
			return
		}
		log.Error(row.Err())
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 链接存在时，更新链接
	if data.ExistLinkCreate {
		// 重新更新链接
		link := createRand()
		full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)
		_, err = app.db.Exec("update share set link=? where appid=? and docid=?", link, data.Appid, data.Docid)
		if err != nil {
			log.Error(err)
			c.String(http.StatusOK, res.setErrSystem().toString())
			return
		}
		log.Infof("更新参数成功")
		log.Infof("更新分享链接: %s", full_link)
		c.String(http.StatusOK, res.setOK(full_link).toString())
		return
	}

	full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)

	log.Infof("提取参数成功")
	log.Infof("当前分享链接: %s", full_link)

	// link 是生成的
	c.String(http.StatusOK, res.setOK(full_link).toString())

}

func getLinkRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("获取链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())

	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	var res resStruct

	// 获取post请求的数据
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	var data getLinkReq
	// 解析json数据
	if err := json.Unmarshal(b, &data); err != nil {

		c.String(http.StatusOK, res.setErrJson().toString())
		return
	}
	log.Infof("appid: %s", data.Appid)
	log.Infof("docid: %s", data.Docid)
	log.Infof("插件版本: %v", data.PluginVersion)

	// 处理请求
	row := app.db.QueryRow(`select link from share where appid=? and docid=? `, data.Appid, data.Docid)
	var link string

	if err := row.Scan(&link); err != nil {
		if err == sql.ErrNoRows {
			log.Infof("处理结果: 没有共享")
			c.String(http.StatusOK, res.setNoShare().toString())
			return
		} else {
			log.Errorf("错误: %v", err)
			c.String(http.StatusOK, res.setErrSystem().toString())
		}

		return
	}
	log.Infof("链接: %s", link)

	// 返回数据

	c.String(http.StatusOK, res.setOK(fmt.Sprintf("%s/%s", app.ShareBaseLink, link)).toString())
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

	return
}
func deleteLinkRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("删除链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())

	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
	var res resStruct
	// 从body中读取数据
	// 获取post请求的数据
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Errorf("错误: %v", err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	var data deleteLinkReq
	// 解析json数据
	if err := json.Unmarshal(b, &data); err != nil {
		log.Errorf("错误: %v", err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	log.Infof("appid: %s", data.Appid)
	log.Infof("docid: %s", data.Docid)
	log.Infof("插件版本: %v", data.PluginVersion)

	_, err = app.db.Exec(`delete from share where appid=? and docid=?`, data.Appid, data.Docid)
	if err != nil {
		log.Errorf("错误: %v", err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	if err := rmdir_all(data.Appid, data.Docid); err != nil {
		log.Errorf("错误: %v", err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	log.Info("删除链接成功")
	// 返回数据
	c.String(http.StatusOK, res.setOK("删除链接成功").toString())
}
