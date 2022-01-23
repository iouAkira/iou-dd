package main

import (
	"ddbot/pre_init"
)

func main() {
	// 读取加载程序需要使用的环境变量
	pre_init.LoadEnv()
}
