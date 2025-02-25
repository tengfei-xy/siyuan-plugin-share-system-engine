package main

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var app appConfig

const version string = "2.4.2"

func prompt() {
	fmt.Printf(`
.___________. _______ .__   __.   _______  _______  _______  __           ___   ___ ____    ____ 
|           ||   ____||  \ |  |  /  _____||   ____||   ____||  |          \  \ /  / \   \  /   / 
.---|  |----.|  |__   |   \|  | |  |  __  |  |__   |  |__   |  |  ______   \  V  /   \   \/   /  
    |  |     |   __|  |  . .  | |  | |_ | |   __|  |   __|  |  | |______|   >   <     \_    _/   
    |  |     |  |____ |  |\   | |  |__| | |  |     |  |____ |  |           /  .  \      |  |     
    |__|     |_______||__| \__|  \______| |__|     |_______||__|          /__/ \__\     |__|     
																								   
`)
	fmt.Println("思源笔记-分享笔记插件服务器: https://github.com/tengfei-xy/siyuan-plugin-share-system-engine")
	fmt.Println("思源笔记-分享笔记插件: https://github.com/tengfei-xy/siyuan-plugin-share-system")
	fmt.Printf("当前服务器版本: v%s\n", version)
}
func main() {
	prompt()
	f := init_flag()
	init_config(f)
	init_env()
	init_db(f)
	init_web()
}
