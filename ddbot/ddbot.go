package main

import (
	ddCmd "ddbot/dd_cmd"
	"ddbot/models"
	"ddbot/pre_init"
	ddutils "ddbot/utils"
)

func main() {
	// 读取加载程序需要使用的环境变量
	upParams := pre_init.LoadEnv()
	ddutils.ExecUpCommand(upParams)
	engine := pre_init.SetupRouters()
	engine.Run(models.GlobalEnv.TgBotToken,
		models.GlobalEnv.TgUserID,
		ddCmd.DebugMode(true),
		ddCmd.TimeOut(60),
	)
}
