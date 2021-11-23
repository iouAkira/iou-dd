package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	dd_cmd "ddbot/dd_cmd"
	models "ddbot/models"
	ddutils "ddbot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HelpHandler 使用说明独立控制器
func DDNodeHandler(env *models.DDEnv) dd_cmd.HandlerFunc {
	return func(ctx *dd_cmd.Context) {
		cmdMsg := ""
		chatId := int64(0)
		msgId := int(0)
		isCallbackQuery := false
		if ctx.Update.CallbackQuery != nil && ctx.Update.CallbackQuery.Data != "" {
			cmdMsg = ctx.Update.CallbackQuery.Data
			chatId = ctx.Update.CallbackQuery.From.ID
			msgId = ctx.Update.CallbackQuery.Message.MessageID
			isCallbackQuery = true
		} else {
			cmdMsg = ctx.Update.Message.Text
			chatId = ctx.Update.Message.Chat.ID
			msgId = ctx.Update.Message.MessageID
		}
		if cmdMsg == "" {
			return
		}

		bot := ctx.Request
		cmdMsgSplit := ddutils.CleanCommand(cmdMsg, 0)
		log.Println("CleanCommand cmdMsgSplit:", cmdMsgSplit)
		if len(cmdMsgSplit) == 1 || ddutils.IsContain(cmdMsgSplit, "dir") {
			if ddutils.IsContain(cmdMsgSplit, "dir") && cmdMsgSplit[1] == "dir" {
				var editNumericKeyboard = ddutils.MakeKeyboardMarkup(ctx.HandlerPrefixStr, 3, cmdMsgSplit[2], "js")
				respMsg := tgbotapi.NewEditMessageText(chatId, msgId, "请选择要执行的操作⚙️")
				respMsg.ReplyMarkup = &editNumericKeyboard
				bot.Send(respMsg)
			} else {
				var numericKeyboard = ddutils.MakeKeyboardMarkup(ctx.HandlerPrefixStr, 3, env.DDnodeBtnFilePath, "js")
				respMsg := tgbotapi.NewMessage(chatId, "请选择要执行的操作⚙️")
				respMsg.ReplyToMessageID = msgId
				respMsg.ReplyMarkup = numericKeyboard
				bot.Send(respMsg)
			}
		} else {
			var respMsgInfo tgbotapi.Message
			if isCallbackQuery {
				respMsg := tgbotapi.NewEditMessageText(chatId, msgId, fmt.Sprintf("`%v` 正在执行⚡️", strings.Join(cmdMsgSplit, " ")))
				respMsgInfo, _ = bot.Send(respMsg)
			} else {
				respMsg := tgbotapi.NewMessage(chatId, fmt.Sprintf("`%v` 正在执行⚡️", strings.Join(cmdMsgSplit, " ")))
				respMsg.ParseMode = tgbotapi.ModeMarkdown
				respMsg.ReplyToMessageID = msgId
				respMsgInfo, _ = bot.Send(respMsg)
			}
			execResult, isFile, err := ddutils.ExecCommand(ddutils.CleanCommand(cmdMsg[1:], 0), ctx.HandlerPrefixStr[1:], env.LogsBtnFilePath)
			if err != nil {
				log.Println(err)
				if isFile {
					respMsgDel := tgbotapi.NewDeleteMessage(chatId, respMsgInfo.MessageID)
					bot.Send(respMsgDel)
					//需要传入绝对路径
					bytes, _ := ioutil.ReadFile(execResult)
					fileSend := tgbotapi.FileBytes{
						Name:  "bot_exec.log",
						Bytes: bytes,
					}
					respMsgFile := tgbotapi.NewDocument(chatId, fileSend)
					respMsgFile.Caption = fmt.Sprintf("`%v` 执行出错❌", strings.Join(cmdMsgSplit, " "))
					respMsgFile.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(respMsgFile)
					_ = os.Remove(execResult)
				} else {
					respMsgEdit := tgbotapi.NewEditMessageText(chatId, respMsgInfo.MessageID, fmt.Sprintf("`%v` 执行出错❌\n\n```\n%v```", strings.Join(cmdMsgSplit, " "), err))
					respMsgEdit.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(respMsgEdit)
				}
			} else {
				execStatus := "执行成功✅"
				if strings.HasPrefix(execResult, "stderr") {
					execStatus = "执行出错❌"
				}
				//log.Printf(execResult)
				if isFile {
					respMsgDel := tgbotapi.NewDeleteMessage(chatId, respMsgInfo.MessageID)
					bot.Send(respMsgDel)
					//需要传入绝对路径
					bytes, _ := ioutil.ReadFile(execResult)
					fileSend := tgbotapi.FileBytes{
						Name:  "bot_exec.log",
						Bytes: bytes,
					}
					respMsgFile := tgbotapi.NewDocument(chatId, fileSend)
					respMsgFile.Caption = fmt.Sprintf("`%v` %v", strings.Join(cmdMsgSplit, " "), execStatus)
					respMsgFile.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(respMsgFile)
					_ = os.Remove(execResult)
				} else {
					respMsgEdit := tgbotapi.NewEditMessageText(chatId,
						respMsgInfo.MessageID,
						fmt.Sprintf("`%v` %v\n\n```\n%v```", strings.Join(cmdMsgSplit, " "), execStatus, execResult))
					respMsgEdit.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(respMsgEdit)
				}
			}
		}
	}
}
