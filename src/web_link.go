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

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	log "github.com/tengfei-xy/go-log"
	"github.com/tengfei-xy/go-tools"
	"golang.org/x/net/html"
)

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
type uploadArgsReq struct {
	Appid            string `json:"appid"`
	Docid            string `json:"docid"`
	Content          string `json:"content"`
	Theme            string `json:"theme"`
	SiyuanVersion    string `json:"version"`
	Title            string `json:"title"`
	HideSYVersion    bool   `json:"hide_version"`
	PluginVersion    string `json:"plugin_version"`
	ExistLinkCreate  bool   `json:"exist_link_create"`
	PageWide         string `json:"page_wide"`
	AccessKeyEnable  bool   `json:"access_key_enable"`
	AccessKey        string `json:"access_key"`
	MiniMenu         bool   `json:"mini_menu"`
	TitleImageHeight string `json:"title_image_height"`
	CustomCSS        string `json:"custom_css"`
	HomePage         bool   `json:"home_page"`
}

func v2PostLinkRequest(c *gin.Context) {
	log.Info("-----------------")
	log.Info("获取链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("链接: %s", c.Request.URL.String())

	cros_status := c.Request.Header.Get("cros-status")
	if cros_status == "true" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
	type resdata struct {
		Link     string `json:"link"`
		HomePage bool   `json:"home_page"`
	}

	// 获取post请求的数据
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	var data getLinkReq
	// 解析json数据
	if err := json.Unmarshal(b, &data); err != nil {
		badRequest(c)
		return
	}
	log.Infof("appid: %s", data.Appid)
	log.Infof("docid: %s", data.Docid)
	log.Infof("插件版本: %v", data.PluginVersion)

	// 处理请求
	row := app.db.QueryRow(`select link,home_page from share where appid=? and docid=? `, data.Appid, data.Docid)
	var link string
	var home_page int
	if err := row.Scan(&link, &home_page); err != nil {
		if err == sql.ErrNoRows {
			log.Infof("处理结果: 没有共享")
			noShare(c)
		} else {
			log.Errorf("错误: %v", err)
			internalSystem(c)
		}
		return
	}
	log.Infof("首页: %v", home_page == 1)
	log.Infof("链接: %s", link)
	var url string = app.ShareBaseLink
	if home_page != 1 {
		url += "/" + link
	}
	okData(c, resdata{
		Link:     url,
		HomePage: home_page == 1,
	})
}

// 以下是旧版的函数
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
		c.JSON(http.StatusOK, res.setErrParam())
		return

	}

	// 读取文件
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	// 创建资源文件夹
	f, err := mkdir_all(appid, docid)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	zip_file := filepath.Join(f, "resources.zip")
	// 保存文件
	err = c.SaveUploadedFile(file, zip_file)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}
	log.Infof("资源文件: %s", zip_file)

	// 如果压缩包不存在
	if !tools.FileExist(zip_file) {
		internalSystem(c)
		return
	}

	// 解压文件
	if err := unzip(f, zip_file); err != nil {
		internalSystem(c)
		return
	}

	check_theme_file(f)
	c.JSON(http.StatusOK, res.setOK("上传文件成功"))
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
		internalSystem(c)
		return
	}

	// 解析请求数据
	var data uploadArgsReq
	if err := json.Unmarshal(b, &data); err != nil {
		c.JSON(http.StatusOK, res.setErrJson())
		return
	}
	TitleImageHeight, err := strconv.Atoi(data.TitleImageHeight)
	if err != nil {
		log.Warnf("题头图高度参数异常，恢复为默认参数")
		TitleImageHeight = 30
	}

	access_key_enable := 0
	home_page_enable := 0

	log.Infof("appid: %s", data.Appid)
	log.Infof("docid: %s", data.Docid)
	log.Infof("theme: %s", data.Theme)
	log.Infof("title: %s", data.Title)
	// 不输出文档的内容，格式为html
	// log.Infof("content: %s", data.Content)
	log.Infof("思源版本: %s", data.SiyuanVersion)
	log.Infof("插件版本: %v", data.PluginVersion)
	log.Infof("页面宽度: %s", data.PageWide)
	if data.AccessKeyEnable {
		log.Infof("访问密钥: %s", data.AccessKey)
		access_key_enable = 1
	} else {
		log.Info("访问密钥: 关闭")
	}
	if data.HomePage {
		log.Infof("设置为首页")
		home_page_enable = 1
	} else {
		log.Infof("取消设置为首页")
	}
	log.Infof("缩略图导航菜单: %v", data.MiniMenu)
	log.Warnf("题头图高度: %s", data.TitleImageHeight)
	log.Infof("标题中隐藏思源版本: %v", data.HideSYVersion)
	log.Infof("链接存在时重新创建: %v", data.ExistLinkCreate)

	// 创建资源文件夹
	// 返回资源文件夹的路径,作为生成的html文件的存放路径
	tmp_html, err := mkdir_all(data.Appid, data.Docid)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	title_version := func() string {
		if data.HideSYVersion {
			return ""
		} else {
			return "   v" + data.SiyuanVersion

		}
	}
	if TitleImageHeight >= 0 {
		setTitleImageHeigt(&data.Content, TitleImageHeight)
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

	// 是否使用导航图
	if data.MiniMenu {
		content = strings.ReplaceAll(content, "{{ .MiniMenuStyle }}", mini_menu_style)
		content = strings.ReplaceAll(content, "{{ .MiniMenuScript }}", mini_menu_script)
	} else {
		content = strings.ReplaceAll(content, "{{ .MiniMenuStyle }}", "")
		content = strings.ReplaceAll(content, "{{ .MiniMenuScript }}", "")
	}
	if data.CustomCSS != "" {
		content = strings.ReplaceAll(content, "{{ .CustomCSS }}", data.CustomCSS)
	}

	f, err := os.OpenFile(filepath.Join(tmp_html, "index.htm"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, get_file_permission())
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		log.Error(err)
		internalSystem(c)
		return
	}

	if data.HomePage {
		// 设置为首页
		_, err = app.db.Exec("update share set home_page=0 where appid=?", data.Appid)
		if err != nil {
			log.Error(err)
			internalSystem(c)
			return
		}
	}
	// 根据appid和docid查看链接是否存在
	row := app.db.QueryRow(`select link from share where appid=? and docid=? `, data.Appid, data.Docid)
	var link string

	if err := row.Scan(&link); err != nil {
		// 如果不存在就插入
		if err == sql.ErrNoRows {
			link := createRand(RAND_URL_LENGTH)
			full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)
			_, err := app.db.Exec("INSERT INTO share(appid,docid,title,link,access_key,access_key_enable,home_page) VALUES(?,?,?,?,?,?,?)  ", data.Appid, data.Docid, data.Title, link, data.AccessKey, access_key_enable, home_page_enable)
			if err != nil {
				log.Error(err)
				internalSystem(c)
				return
			}

			log.Infof("创建参数成功")
			log.Infof("创建分享链接: %s", full_link)
			c.JSON(http.StatusOK, res.setOK(full_link))
			return
		}
		log.Error(row.Err())
		internalSystem(c)
		return
	}

	// 链接存在时，更新链接
	if data.ExistLinkCreate {
		// 重新更新链接
		link := createRand(RAND_URL_LENGTH)
		full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)
		_, err = app.db.Exec("update share set link=? where appid=? and docid=? and access_key=? and access_key_enable=? ans home_page=?", link, data.Appid, data.Docid, data.AccessKey, access_key_enable, home_page_enable)
		if err != nil {
			log.Error(err)
			internalSystem(c)
			return
		}
		log.Infof("更新参数成功")
		log.Infof("更新分享链接: %s", full_link)
		c.JSON(http.StatusOK, res.setOK(full_link))
		return
	}

	full_link := fmt.Sprintf("%s/%s", app.Basic.ShareBaseLink, link)

	log.Infof("提取参数成功")
	log.Infof("当前分享链接: %s", full_link)

	// link 是生成的
	c.JSON(http.StatusOK, res.setOK(full_link))
}
func setTitleImageHeigt(content *string, height int) {
	style := fmt.Sprintf(`height: %dvh;overflow: hidden;`, height)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*content))
	if err != nil {
		log.Fatal(err)
	}
	selection := doc.Find("div[data-type=NodeParagraph]").First()

	if selection.Find("img[alt=\"image\"]").Length() == 1 {
		log.Infof("发现题头图")
		node := selection.Nodes[0]
		// 添加一个新的属性
		node.Attr = append(node.Attr, html.Attribute{Key: "style", Val: style})
		// 输出修改后的HTML
		html, err := doc.Html()
		if err != nil {
			log.Fatal(err)
		}
		*content = html
	}

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
		internalSystem(c)
		return
	}

	var data deleteLinkReq
	// 解析json数据
	if err := json.Unmarshal(b, &data); err != nil {
		log.Errorf("错误: %v", err)
		internalSystem(c)
		return
	}
	log.Infof("appid: %s", data.Appid)
	log.Infof("docid: %s", data.Docid)
	log.Infof("插件版本: %v", data.PluginVersion)

	_, err = app.db.Exec(`delete from share where appid=? and docid=?`, data.Appid, data.Docid)
	if err != nil {
		log.Errorf("错误: %v", err)
		internalSystem(c)
		return
	}
	if err := rmdir_all(data.Appid, data.Docid); err != nil {
		log.Errorf("错误: %v", err)
		internalSystem(c)
		return
	}
	log.Info("删除链接成功")
	// 返回数据
	c.JSON(http.StatusOK, res.setOK("删除链接成功"))
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
		internalSystem(c)
		return
	}

	var data getLinkReq
	// 解析json数据
	if err := json.Unmarshal(b, &data); err != nil {

		c.JSON(http.StatusOK, res.setErrJson())
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
			noShare(c)
			return
		} else {
			log.Errorf("错误: %v", err)
			internalSystem(c)
		}

		return
	}
	log.Infof("链接: %s", link)

	// 返回数据
	c.JSON(http.StatusOK, res.setOK(fmt.Sprintf("%s/%s", app.ShareBaseLink, link)))
}
func getLinkAllRequest(c *gin.Context) {
	var res resStruct

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	appid := c.Query("appid")
	log.Info("-----------------")
	log.Info("获取所有链接")
	log.Infof("IP: %s", c.ClientIP())
	log.Infof("appid: %s", appid)
	if appid == "" {
		c.JSON(http.StatusOK, res.setErrParam())
		return
	}

	type data struct {
		Link            string `json:"link"`
		Title           string `json:"title"`
		DocID           string `json:"docid"`
		Enable          bool   `json:"enable"`
		AccessKey       string `json:"access_key"`
		AccessKeyEnable int    `json:"access_key_enable"`
		IsHomePage      bool   `json:"is_home_page"`
	}

	var length int
	if err := app.db.QueryRow("select COUNT(*) from share where appid=?", appid).Scan(&length); err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgInternalSystemErr())
		return
	}
	if length == 0 {
		d := make([]data, 0)

		log.Warnf("未查询到记录 appid=%s", appid)
		c.JSON(http.StatusOK, msgOK(d))
		return
	}

	d := make([]data, length)

	ret, err := app.db.Query("select docid,title,link,access_key,access_key_enable,home_page from share where appid=?", appid)
	if err != nil {
		log.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, msgInternalSystemErr())
		return
	}
	log.Infof("查询到%d条记录", length)
	var i int = 0
	for ret.Next() {
		var homepage int
		err := ret.Scan(&d[i].DocID, &d[i].Title, &d[i].Link, &d[i].AccessKey, &d[i].AccessKeyEnable, &homepage)
		if err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, msgInternalSystemErr())
			return
		}
		if homepage == 1 {
			d[i].IsHomePage = true
		}
		log.Debugf("%s", d[i].DocID)

		d[i].Link = fmt.Sprintf("%s/%s", app.ShareBaseLink, d[i].Link)
		i++
	}
	c.JSON(http.StatusOK, msgOK(d))
}
