package pre_init

import (
	"log"

	ctl "ddbot/controller"
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
	engine.RegCommandByChar("/", "cmd", ctl.SysCmdHandler(models.GlobalEnv))
	engine.RegCommandByChar("/", "help", ctl.HelpHandler(models.GlobalEnv))
	engine.RegCommandByChar("/", "start", ctl.HelpHandler(models.GlobalEnv))
	engine.RegCommandByChar("/", "ddnode", ctl.DDNodeHandler(models.GlobalEnv))
	engine.RegCommandByChar("/", "rdc", ctl.ReadCookieHandler(models.GlobalEnv))
	engine.RegCommandByChar("/", "wskey", ctl.ReadWSKeyHandler(models.GlobalEnv))
	engine.RegCommandByChar("/", "ext", ctl.ExtansionController)
	//engine.Cmd("ak", controller.AkController(model.Env))
	//engine.Cmd("dk", controller.DkController(model.Env))
	//engine.Cmd("clk", controller.ClearReplyKeyboardController(model.Env))
	//engine.Cmd("dl", controller.DownloadFileByUrlController(model.Env))
	//engine.Cmd("logs", controller.LogController(model.Env))
	engine.RegCommandByChar("/", "cancel", ctl.CancelController)
	// 注册一个未知命令响应函数
	engine.RegCommandByChar("/", "unknow", ctl.UnknownController)

	return engine
}
