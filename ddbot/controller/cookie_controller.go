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
func ReadCookieHandler(env *models.DDEnv) dd_cmd.HandlerFunc {
	cookies := utils.CookieCfg{DDEnv: env}
	return func(ctx *dd_cmd.Context) {
		message := ctx.Message(ctx)
		path := ctx.Vars()
		cookiesFrom, err := cookies.ReadCookies(false, path)
		if err != nil {
			return
		}
		markup := dd_cmd.RdcMarkup{Cmd: "rdc", Prefix: "/", Cookies: cookiesFrom, RowBtns: 2, Suffix: "ws"}
		//todo 待完善功能
		respMsg := tgbotapi.NewMessage(message.Chat.ID, "请选择要执行的操作⚙️\n 可以多参数方式 /rdc ws 1或者/rdc ws")
		if len(path) > 0 {
			respMsg.Text = strings.Join(cookiesFrom,"\n")
		}else {

			respMsg.ReplyMarkup = markup.MakeKeyboardMarkup().WithCancel()
			respMsg.ReplyToMessageID = message.MessageID
		}
		if _, err := ctx.Send(respMsg); err != nil {
			log.Println(err)
			return
		}
	}
}
