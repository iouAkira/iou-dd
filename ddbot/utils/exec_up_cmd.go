package utils

import (
	"flag"
	"fmt"
	"os"

	"ddbot/models"
)

// ExecUpCommand
// @description   执行命令指定了 -up 参数的启动命令
// @auth      iouAkira
// @param     env 全局程序配置信息
func ExecUpCommand() {

	var upParams string
	flag.StringVar(&upParams, "up", "", "默认为空，为启动bot；commitShareCode为提交互助码到助力池；syncRepo为同步仓库代码；")
	flag.Parse()
	// -up 启动参数 不指定默认启动ddbot
	if upParams != "" {
		fmt.Printf("传入 -up参数：%v ", upParams)
		if upParams == "commitShareCode" {
			fmt.Printf("启动程序指定了 -up 参数为 %v 开始上传互助码。", upParams)
			UploadShareCode(models.GlobalEnv)
		} else if upParams == "syncRepo" {
			fmt.Printf("启动程序指定了 -up 参数为 %v 开始同步仓库代码。", upParams)
			if models.GlobalEnv.RepoBaseDir != "" {
				SyncRepo(models.GlobalEnv)
			} else {
				fmt.Printf("同步仓库设定的目录[%v]不规范，退出同步。", models.GlobalEnv.RepoBaseDir)
			}
		} else if upParams == "renewCookie" {
			fmt.Printf("启动程序指定了 -up 参数为 %v 开始给 %v 里面的全部wskey续期。", upParams, models.GlobalEnv.CookiesListFilePath)
			RenewAllCookie(models.GlobalEnv)
		} else {
			fmt.Printf("请传入传入的对应 -up参数：%v ", upParams)
		}
		os.Exit(0)
	}
}
