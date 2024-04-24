package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/tengfei-xy/go-log"
	"github.com/tengfei-xy/go-tools"
	"gopkg.in/yaml.v3"
)

type appConfig struct {
	Mysql `yaml:"mysql"`
	Basic `yaml:"basic"`
	db    *sql.DB
	// 仅仅用作判断是否属于docker环境
	docker bool
}
type Basic struct {
	ListenPort    string `yaml:"listen"`
	SavePath      string `yaml:"savePath"`
	ShareBaseLink string `yaml:"shareBaseLink"`
}

type Mysql struct {
	Ip       string `yaml:"ip"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}
type flagStruct struct {
	config_file string
	version     bool
}

type Docker struct {
	// Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}

type Services struct {
	Nginx Nginx `yaml:"spss_nginx"`
	Db    Db    `yaml:"spss_mysql"`
	App   App   `yaml:"spss_engine"`
}

type Nginx struct {
	Image         string   `yaml:"image"`
	ContainerName string   `yaml:"container_name"`
	Ports         []string `yaml:"ports"`
	Volumes       []string `yaml:"volumes"`
	Restart       string   `yaml:"restart"`
}

type Db struct {
	Restart       string      `yaml:"restart"`
	Privileged    bool        `yaml:"privileged"`
	Image         string      `yaml:"image"`
	ContainerName string      `yaml:"container_name"`
	Volumes       []string    `yaml:"volumes"`
	Environment   Environment `yaml:"environment"`
	Links         []string    `yaml:"links"`
}

type Environment struct {
	MYSQLROOTPASSWORD string `yaml:"MYSQL_ROOT_PASSWORD"`
	MYSQLUSER         string `yaml:"MYSQL_USER"`
	MYSQLPASS         string `yaml:"MYSQL_PASS"`
	MYSQLDATABASE     string `yaml:"MYSQL_DATABASE"`
}

type App struct {
	ContainerName string         `yaml:"container_name"`
	Build         Build          `yaml:"build"`
	Environment   AppEnvironment `yaml:"environment"`
	Restart       string         `yaml:"restart"`
	Volumes       []string       `yaml:"volumes"`
}

type Build struct {
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile"`
}

type AppEnvironment struct {
	SPSSSTARTUPENV string `yaml:"SPSS_STARTUP_ENV"`
	ShareBaseLink  string `yaml:"SHARE_BASE_LINK"`
	ListenPort     int    `yaml:"LISTEN_PORT"`
}

func init_config(flag flagStruct) {
	if flag.version {
		fmt.Println(version)
		os.Exit(0)
	}
	log.Infof("读取配置文件")
	app.docker = false

	if tools.FileExist(flag.config_file) {
		log.Infof("发现配置文件:%s", flag.config_file)
		yamlFile, err := os.ReadFile(flag.config_file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(yamlFile, &app)
		if err != nil {
			panic(err)
		}
		log.Infof("资源文件保存位置:%s", app.Basic.SavePath)
		return
	}

	s := os.Getenv("SPSS_STARTUP_ENV")
	save_path := os.Getenv("SAVE_PATH")
	if save_path == "" {
		panic("未设置SAVE_PATH环境变量")
	}
	app.Basic.SavePath = save_path
	if s != "docker" {
		panic("不是预设的SPSS_STARTUP_ENV变量")
	} else if s == "" {
		panic("未设置SPSS_STARTUP_ENV环境变量")
	}

	log.Infof("发现docker环境")
	app.docker = true
	l := []string{"./", "./docker"}
	var find bool = false
	for _, v := range l {
		docker_config_file := filepath.Join(app.Basic.SavePath, v, "docker-compose.yml")
		if tools.FileExist(docker_config_file) {
			log.Infof("发现配置文件: %s", docker_config_file)
			find = true
			app.Basic.SavePath = filepath.Join(app.Basic.SavePath, v)
		}
	}

	if !find {
		panic(fmt.Sprintf("未在%s下找到docker-compose.yml文件", app.Basic.SavePath))
	}

	yamlFile, err := os.ReadFile(filepath.Join(app.Basic.SavePath, "docker-compose.yml"))
	if err != nil {
		panic(err)
	}
	var docker_config Docker
	err = yaml.Unmarshal(yamlFile, &docker_config)
	if err != nil {
		panic(err)
	}
	trans_docker_config(docker_config)
	log.Infof("分享地址的基础链接: %s", app.Basic.ShareBaseLink)

}
func trans_docker_config(d Docker) {

	if len(d.Services.Nginx.Ports) == 0 {
		panic("docker-compose的nginx端口映射为空")
	}

	for _, v := range d.Services.Nginx.Volumes {

		// 获取docker-compose.yaml的serices字段中nginx容器的volumes字段
		// 寻找带有html的字符串,将冒号前的部分其作为资源文件保存位置
		if strings.Contains(v, "html") {
			app.Basic.SavePath = filepath.Join(app.Basic.SavePath, strings.Split(v, ":")[0])
			log.Infof("资源文件保存位置: %s", app.Basic.SavePath)
			break
		}
	}
	app.Mysql.Ip = "spss_mysql"

	app.Mysql.Port = "3306"
	app.Mysql.Database = d.Services.Db.Environment.MYSQLDATABASE
	app.Mysql.Username = "root"
	app.Mysql.Password = d.Services.Db.Environment.MYSQLPASS

	app.Basic.ListenPort = fmt.Sprintf(":%d", d.Services.App.Environment.ListenPort)
	app.Basic.ShareBaseLink = d.Services.App.Environment.ShareBaseLink

}
