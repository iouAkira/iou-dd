package controller

import (
	"ddbot/dd_cmd"
	"ddbot/models"
	"ddbot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

const IKTYPE string = "rdc"
const (
	CK_READ_LIST  = iota // 读取cookie列表
	CK_READ_INDEX        // 读取id对应的cookie字符串
	CK_READ_SUB          // 读取子命令
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
		switch len(path) {
		case CK_READ_INDEX:
			cookiesFrom, err := cookies.ReadCookies(false, path)
			if err != nil {
				return nil
			}
			log.Println(cookiesFrom)
			if ctx.IsCallBack() {
				oldMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
				_, _ = ctx.Send(oldMsg)
			}
			//respMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, "请选择要执行的操作⚙️\n")
			respMsg := tgbotapi.NewMessage(message.Chat.ID, "请选择要执行的操作⚙️\n 可以多参数方式 /rdc")
			respMsg.Text = strings.Join(cookiesFrom, "\n")
			pinID := utils.GetPinFromCookieText(respMsg.Text)
			ReplyMarkup := dd_cmd.MakeKeyboard().WithCommandStr(fmt.Sprintf("/wskey %s",pinID),"查询WSKEY").Get()
			respMsg.ReplyMarkup = &ReplyMarkup
			return respMsg
		case CK_READ_LIST:
			cookiesFrom, err := cookies.ReadCookies(false, path)
			if err != nil {
				return nil
			}
			respMsg := tgbotapi.NewMessage(message.Chat.ID, "请选择要执行的操作⚙️\n")
			markup := dd_cmd.RdcMarkup{Cmd: IKTYPE, Prefix: "/", Cookies: cookiesFrom, RowBtns: 2, Suffix: "/wskey"}
			ReplyMarkup := markup.MakeKeyboardMarkup().WithCancel().Get()
			respMsg.ReplyMarkup = &ReplyMarkup
			return respMsg
		// 直接重定向到子命令
		case CK_READ_SUB:
			if path[0] == "/rdc" {
				return nil
			}
			ctx.RedirectToCmd(path[0], path[1:]...)
			return nil
		default:
			return nil
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

func ReadWSKeyHandler(env *models.DDEnv) dd_cmd.HandlerFunc {
	cookies := utils.CookieCfg{DDEnv: env}
	return func(ctx *dd_cmd.Context) {
		message := ctx.Message(ctx)
		path := ctx.Vars()
		cookiesFrom, err := cookies.ReadCookies(true, path)
		if err != nil {
			return
		}
		//markup := dd_cmd.RdcMarkup{Cmd: "wskey", Prefix: "/", Cookies: cookiesFrom, RowBtns: 2, Suffix: "/wskey"}
		log.Println(cookiesFrom)
		if ctx.IsCallBack() {
			oldMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
			_, _ = ctx.Send(oldMsg)
		}
		//respMsg := tgbotapi.NewEditMessageText(message.Chat.ID, message.MessageID, "请选择要执行的操作⚙️\n")
		respMsg := tgbotapi.NewMessage(message.Chat.ID, "请选择要执行的操作⚙️\n 可以多参数方式 /wskey xxx")
		respMsg.Text = strings.Join(cookiesFrom, "\n")
		respMsg.ReplyMarkup = dd_cmd.MakeKeyboard().WithCancel().Get()
		if _, err := ctx.Send(respMsg); err != nil {
			log.Println(err)
			return
		}
	}
}
