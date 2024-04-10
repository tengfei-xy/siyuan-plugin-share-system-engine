package main

import "database/sql"

type appConfig struct {
	Mysql `yaml:"mysql"`
	Basic `yaml:"basic"`
	db    *sql.DB
}
type Basic struct {
	Listen        string `yaml:"listen"`
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
}
