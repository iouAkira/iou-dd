package main

import (
	"ddbot/pre_init"
)

//主程序执行入口
func main() {
	// 读取加载程序需要使用的环境变量
	pre_init.LoadEnv()
}
