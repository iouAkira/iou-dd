package controller

import (
	"log"

	ddCmd "ddbot/dd_cmd"
	models "ddbot/models"
	ddutils "ddbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpHandler 使用说明独立控制器
func HelpHandler(env *models.DDEnv) ddCmd.HandlerFunc {
	return func(ctx *ddCmd.Context) {
		readme := "🙌 <b>使用说明</b> v1.0.0\n" +
			"\n 👉 <b>/spnode</b>  \n        ------  执行JS脚本文件" +
			"\n 👉 <b>/logs</b>    \n        ------  下载日志文件" +
			"\n 👉 <b>/rdc</b>    \n        ------  读取Cookies列表" +
			"\n 👉 <b>/bl</b>    \n        ------  查看cookie收支图表   例：/bl 1 查看第一个cookie" +
			"\n 👉 <b>/env</b>    \n        ------  更新或者替换env.sh内的环境变量 例：/env aaa=\"bbb\"" +
			"\n 👉 <b>/cmd</b>    \n        ------  执行指定命令   例：/cmd ls -l" +
			"\n 👉 <b>/ak</b>    \n        ------  添加/更新快捷回复键盘   例：/ak 键盘显示===/cmd echo 'show reply keyboard'" +
			"\n 👉 <b>/dk</b>    \n        ------  删除快捷回复键盘   例：/dk 键盘显示" +
			"\n 👉 <b>/clk</b>    \n        ------  清空快捷回复键盘   例：/clk" +
			"\n 👉 <b>/dl</b>    \n        ------  通过链接下载文件   例：/dl https://raw.githubusercontent.com/iouAkira/someDockerfile/master/dd_scripts/shell_mod_script.sh" +
			"\n 👉 <b>/renew</b>    \n        ------  通过wskey[cookies_wskey.list]更新cookies.list   例：/renew 1  更行cookies_wskey.list里面的第一个ck"

		//创建信息
		helpMsg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, readme)
		//tgbotapi.ChatRecordAudio
		//修改信息格式
		helpMsg.ParseMode = tgbotapi.ModeHTML
		//创建回复键盘结构体
		tkbs := ddutils.MakeReplyKeyboard(env)
		//赋值给ReplyMarkup[快速回复]
		helpMsg.ReplyMarkup = tkbs
		//发送消息
		if _, err := ctx.Send(helpMsg); err != nil {
			log.Println(err)
		}
	}
}

func CancelController(ctx *ddCmd.Context) {
	if ctx.Update.CallbackQuery != nil {
		c := ctx.Update.CallbackQuery
		edit := tgbotapi.NewEditMessageText(c.Message.Chat.ID, c.Message.MessageID, "操作已经取消")
		_, _ = ctx.Send(edit)
	}
}
