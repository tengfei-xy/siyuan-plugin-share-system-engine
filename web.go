package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
	"github.com/tengfei-xy/go-tools"
	tool "github.com/tengfei-xy/go-tools"
)

type uploadArgsReq struct {
	Appid   string `json:"appid"`
	Docid   string `json:"docid"`
	Content string `json:"content"`
	Theme   string `json:"theme"`
	Version string `json:"version"`
	Title   string `json:"title"`
}
type getLinkReq struct {
	Appid string `json:"appid"`
	Docid string `json:"docid"`
}

type deleteLinkReq struct {
	Appid string `json:"appid"`
	Docid string `json:"docid"`
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
	g.MaxMultipartMemory = 100 << 20 // 100 MiB
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	// g.Use(CORSMiddleware())
	log.Infof("服务器启动，监听 %s", app.Basic.Listen)
	g.Run(app.Basic.Listen)
}

// 处理函数，
func shareRequest(c *gin.Context) {
	param := c.Params.ByName("url")

	log.Info("-----------------")
	log.Info("请求链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())
	log.Infof("参数: %s", param)
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:61521")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:61521")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	var res resStruct
	appid := c.Query("appid")
	docid := c.Query("docid")

	if len(appid) == 0 || len(docid) == 0 {
		c.String(http.StatusOK, res.setErrParam().toString())
		return

	}
	log.Infof("appid: %s", appid)
	log.Infof("docid: %s", docid)

	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 设定原始文件的名称
	tmp_zip_file := fmt.Sprintf("%s_%s.zip", appid, docid)

	// 保存原始文件
	err = c.SaveUploadedFile(file, tmp_zip_file)
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

	// 复制原始文件到资源文件夹中
	dst_zip_file := filepath.Join(f, "resources.zip")
	if _, err := tool.FileCopy(dst_zip_file, tmp_zip_file); err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	// 解压文件
	log.Infof("解压文件 文件名:%s", dst_zip_file)
	if err := unzip(f, dst_zip_file); err != nil {
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	// 解压文件 结束
	// 删除原始文件
	if err := tools.FileRemove(tmp_zip_file); err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	update_theme_file(f)
	c.String(http.StatusOK, res.setOK("上传文件成功").toString())
	return
}
func uploadArgsRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Infof("IP: %s", c.ClientIP())
	log.Info("上传参数")
	log.Infof("链接: %s", c.Request.URL.String())

	var res resStruct
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:61521")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
	log.Infof("version: %s", data.Version)
	log.Infof("title: %s", data.Title)

	// 创建资源文件夹，
	// 返回资源文件夹的路径,作为生成的html文件的存放路径
	tmp_html, err := mkdir_all(data.Appid, data.Docid)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}

	content := tempate_html

	content = strings.ReplaceAll(content, "{{ .Theme }}", data.Theme)
	content = strings.ReplaceAll(content, "{{ .Version }}", data.Version)
	content = strings.ReplaceAll(content, "{{ .Title }}", data.Title)
	content = strings.ReplaceAll(content, "{{ .Content }}", data.Content)

	f, err := os.Create(filepath.Join(tmp_html, "index.html"))
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

	link := createRand()
	// 处理请求
	var result sql.Result
	result, err = app.db.Exec("INSERT INTO share(appid,docid,title,link) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE `title` = VALUES(`title`),`link` = VALUES(`link`) ", data.Appid, data.Docid, data.Title, link)
	if err != nil {
		log.Error(err)
		c.String(http.StatusOK, res.setErrSystem().toString())
		return
	}
	n, _ := result.RowsAffected()
	old_link := false
	if n == 0 {
		// 处理请求
		row := app.db.QueryRow(`select link from share where appid=? and docid=? `, data.Appid, data.Docid)

		if err := row.Scan(&link); err != nil {
			log.Error(err)
			c.String(http.StatusOK, res.setErrSystem().toString())
			return
		} else {
			old_link = true
		}

		c.String(http.StatusOK, res.toString())
		return
	}

	full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)

	log.Infof("link: %s", link)
	log.Infof("flink: %s", full_link)
	if old_link {
		log.Infof("上传参数成功")
	} else {
		log.Infof("提取参数成功")
	}
	// 返回数据
	c.String(http.StatusOK, res.setOK(full_link).toString())
}
func getLinkRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("获取链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())

	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:61521")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

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

	c.String(http.StatusOK, res.setOK(fmt.Sprintf("http://124.223.15.220/%s", link)).toString())
}
func deleteLinkRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("删除链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:61521")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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

func mkdir_all(app_id, doc_id string) (string, error) {
	// 创建目录
	f := filepath.Join(app.Basic.SavePath, app_id, doc_id)
	err := os.MkdirAll(f, 0755)
	if err != nil {
		if err != os.ErrExist {
			log.Error(err)
			return f, err
		} else {
			return "", err
		}
	}
	return f, nil
}
func rmdir_all(app_id, doc_id string) error {
	// 删除目录
	f := filepath.Join(app.Basic.SavePath, app_id, doc_id)
	err := os.RemoveAll(f)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
