package controller

import (
	"fmt"

	dd_cmd "ddbot/dd_cmd"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// UnknownController 对未知指令的处理
func UnknownController(ctx *dd_cmd.Context) {
	unknowMsg := tgbotapi.NewMessage(ctx.Update.Message.From.ID, fmt.Sprintf("`%v` 指令暂未注册或开启", ctx.Update.Message.Text))
	unknowMsg.ParseMode = tgbotapi.ModeMarkdown
	_, _ = ctx.Send(unknowMsg)
}
