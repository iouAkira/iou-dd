package pre_init

import (
	"log"

	"ddbot/controller"
	ddCmd "ddbot/dd_cmd"
	models "ddbot/models"
)

// SetupRouters 初始化 ddbot 路由配置
func SetupRouters() *ddCmd.Engine {
	engine := ddCmd.New()
	engine.Use(func(context *ddCmd.Context) {
		if context.Update.Message != nil {
			log.Printf("[%s] %s", context.Update.Message.From.UserName, context.Update.Message.Text)
		}
		if context.Update.CallbackQuery != nil {
			log.Printf("[%s] %s", context.Update.CallbackQuery.From.UserName, context.Update.CallbackQuery.Data)
		}
	})
	engine.RegCommand("/", "cmd", controller.SysCmdHandler(models.GlobalEnv))
	engine.RegCommand(">", "help", controller.HelpHandler(models.GlobalEnv))
	engine.RegCommand("👉", "start", controller.HelpHandler(models.GlobalEnv))
	//engine.Cmd("ddnode", controller.ExecDDnodeController(model.Env, ""))
	//engine.Cmd("ak", controller.AkController(model.Env))
	//engine.Cmd("dk", controller.DkController(model.Env))
	//engine.Cmd("clk", controller.ClearReplyKeyboardController(model.Env))
	//engine.Cmd("dl", controller.DownloadFileByUrlController(model.Env))
	//engine.Cmd("logs", controller.LogController(model.Env))
	engine.RegCommand("👉", "cancel", controller.CancelController)

	return engine
}
