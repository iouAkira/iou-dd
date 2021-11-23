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

// 通过系统命令执行指定的操作
func SysCmdHandler(env *models.DDEnv) dd_cmd.HandlerFunc {
	return func(c *dd_cmd.Context) {
		// "/cmd ls -l"返回 ["ls","-l"]
		cmdMsgSplit := c.Vars()
		// 返回的是tgbotapi.Message对象，包一层为了兼容CallbackQuery类型按钮交互消息
		cmdMsg := c.Message(c)

		if len(cmdMsgSplit) == 0 {
			respMsg := tgbotapi.NewMessage(cmdMsg.Chat.ID, "⚠️请在cmd后面写需要执行的指令。例：/cmd ls -l")
			respMsg.ParseMode = tgbotapi.ModeMarkdown
			respMsg.ReplyToMessageID = cmdMsg.MessageID
			_, _ = c.Send(respMsg)
		} else {
			respMsg := tgbotapi.NewMessage(cmdMsg.Chat.ID, fmt.Sprintf("`%v` 正在执行⚡️", strings.Join(cmdMsgSplit, " ")))
			respMsg.ParseMode = tgbotapi.ModeMarkdown
			respMsg.ReplyToMessageID = cmdMsg.MessageID
			respMsgInfo, _ := c.Send(respMsg)
			// 系统接口默认的命令: 例如：/cmd 得到时候就是 cmd
			cmd := cmdMsg.Command()
			execResult, isFile, err := ddutils.ExecCommand(cmdMsgSplit, cmd, env.LogsBtnFilePath)
			if err != nil {
				log.Println(err)
				if isFile {
					respMsgDel := tgbotapi.NewDeleteMessage(cmdMsg.Chat.ID, respMsgInfo.MessageID)
					c.Send(respMsgDel)
					//需要传入绝对路径
					bytes, _ := ioutil.ReadFile(execResult)
					fileSend := tgbotapi.FileBytes{
						Name:  "bot_exec.log",
						Bytes: bytes,
					}
					respMsgFile := tgbotapi.NewDocument(cmdMsg.Chat.ID, fileSend)
					respMsgFile.Caption = fmt.Sprintf("`%v` 执行出错❌", strings.Join(cmdMsgSplit, " "))
					respMsgFile.ParseMode = tgbotapi.ModeMarkdown
					c.Send(respMsgFile)
					_ = os.Remove(execResult)
				} else {
					respMsgEdit := tgbotapi.NewEditMessageText(cmdMsg.Chat.ID, respMsgInfo.MessageID, fmt.Sprintf("`%v` 执行出错❌\n\n```\n%v```", strings.Join(cmdMsgSplit, " "), err))
					respMsgEdit.ParseMode = tgbotapi.ModeMarkdown
					c.Send(respMsgEdit)
				}
			} else {
				execStatus := "执行成功✅"
				if strings.HasPrefix(execResult, "stderr") {
					execStatus = "执行出错❌"
				}
				//log.Printf(execResult)
				if isFile {
					respMsgDel := tgbotapi.NewDeleteMessage(cmdMsg.Chat.ID, respMsgInfo.MessageID)
					c.Send(respMsgDel)
					//需要传入绝对路径
					bytes, _ := ioutil.ReadFile(execResult)
					fileSend := tgbotapi.FileBytes{
						Name:  "bot_exec.log",
						Bytes: bytes,
					}
					respMsgFile := tgbotapi.NewDocument(cmdMsg.Chat.ID, fileSend)
					respMsgFile.Caption = fmt.Sprintf("`%v` %v", strings.Join(cmdMsgSplit, " "), execStatus)
					respMsgFile.ParseMode = tgbotapi.ModeMarkdown
					c.Send(respMsgFile)
					_ = os.Remove(execResult)
				} else {
					respMsgEdit := tgbotapi.NewEditMessageText(cmdMsg.Chat.ID,
						respMsgInfo.MessageID,
						fmt.Sprintf("`%v` %v\n\n```\n%v```", strings.Join(cmdMsgSplit, " "), execStatus, execResult))
					respMsgEdit.ParseMode = tgbotapi.ModeMarkdown
					c.Send(respMsgEdit)
				}
			}
		}
	}
}
