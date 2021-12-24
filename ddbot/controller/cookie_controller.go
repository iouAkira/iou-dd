package controller

import (
	"ddbot/dd_cmd"
	"ddbot/models"
	"ddbot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

// ReadCookieHandler 读取 cookie列表,主要是简化读取和到处cookie的功能
// 针对WSKEY用户来说可以针对功能来说增加了从cookie列表扩展功能中选择WSKEY读取的选项
// 具体指令如下:
// /rdc 查看用户所有cookie,以按钮形式展示
// /rdc {id} 查看对应序号的索引的cookie值
func ReadCookieHandler(env *models.DDEnv) dd_cmd.HandlerFunc {
	cookies := utils.CookieCfg{DDEnv: env}
	read := func(ctx *dd_cmd.Context) tgbotapi.Chattable {
		message := ctx.Message(ctx)
		path := ctx.Vars()
		cookiesFrom, err := cookies.ReadCookies(false, path)
		if err != nil {
			return nil
		}
		// 读取id对应的cookie字符串
		if len(path) > 0 {
			respMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, "请选择要执行的操作⚙️\n 可以多参数方式 /rdc 1")
			respMsg.Text = strings.Join(cookiesFrom, "\n")
			ReplyMarkup := dd_cmd.MakeKeyboard().WithCancel().Get()
			respMsg.ReplyMarkup = &ReplyMarkup
			return respMsg
			// 读取cookies列表
		} else {
			respMsg := tgbotapi.NewMessage(message.Chat.ID, "请选择要执行的操作⚙️\n 可以多参数方式 /rdc 1")
			markup := dd_cmd.RdcMarkup{Cmd: "rdc", Prefix: "/", Cookies: cookiesFrom, RowBtns: 2, Suffix: "ws"}
			ReplyMarkup := markup.MakeKeyboardMarkup().WithCancel().Get()
			respMsg.ReplyMarkup = &ReplyMarkup
			return respMsg
		}
	}
	//统一处理
	return func(ctx *dd_cmd.Context) {
		msg := read(ctx)
		if msg == nil {
			return
		}
		if _, err := ctx.Send(msg); err != nil {
			log.Println(err)
			return
		}
	}
}
