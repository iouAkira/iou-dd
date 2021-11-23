package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"ddbot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Help
// @description   bot交互help,start指令响应
// @auth       iouAkira
// @param1     chatID	int64
// @param2     bot		*tgbotapi.BotAPI
func Help(chatID int64, bot *tgbotapi.BotAPI) {
	readme := "使用帮助说明" +
		"\n\n/spnode 选择执行JS脚本文件" +
		"\n/logs 选择下载日志文件" +
		"\n/rdc 读取Cookies列表" +
		"\n/bl 查看对应cookie收支图表 例如：/bl 1 查看第一个cookie" +
		"\n/env 更新或者替换env.sh内的环境变量 例：/env aaa=\"bbb\"" +
		"\n/cmd 执行任何想要执行的命令 例：/cmd ls -l" +
		"\n/ak 添加/更新快捷回复键盘 例：/ak 键盘显示===/cmd echo 'show reply keyboard'" +
		"\n/dk 删除快捷回复键盘 例：/dk 键盘显示" +
		"\n/clk 清空快捷回复键盘 例：/clk" +
		"\n/dl 通过链接下载文件 例：/dl https://raw.githubusercontent.com/iouAkira/someDockerfile/master/dd_scripts/shell_mod_script.sh" +
		"\n/renew 通过cookies_wskey.list的wskey更新cookies.list 例如：/renew 1  更行cookies_wskey.list里面的第一个ck"

	helpMsg := tgbotapi.NewMessage(chatID, readme)
	log.Printf("处理前：%v", models.GlobalEnv.ReplyKeyBoard)
	tkbs := MakeReplyKeyboard(models.GlobalEnv)
	log.Printf("处理后：%v", models.GlobalEnv.ReplyKeyBoard)

	helpMsg.ReplyMarkup = tkbs
	log.Printf("tkbs：%v", tkbs)
	if _, err := bot.Send(helpMsg); err != nil {
		log.Println(err)
	}
}

// AddReplyKeyboard
// @description   增加/更新快捷回复键盘指令
// @auth       iouAkira
// @param1     akMsg	*tgbotapi.Message
// @param2     bot		*tgbotapi.BotAPI
func AddReplyKeyboard(akMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	rkb := strings.TrimLeft(akMsg.Text[3:], " ")
	if len(strings.Split(rkb, "===")) > 1 {
		if !CheckDirOrFileIsExist(models.GlobalEnv.ReplyKeyboardFilePath) {
			rkbFile, _ := os.Create(models.GlobalEnv.ReplyKeyboardFilePath)
			defer rkbFile.Close()
		}
		optMsg, err := ReplyKeyboardFileOpt(rkb, strings.Split(rkb, "===")[0], "W")
		if err != nil {
			akRespMsg := tgbotapi.NewMessage(akMsg.Chat.ID, err.Error())
			akRespMsg.ReplyToMessageID = akMsg.MessageID
			bot.Send(akRespMsg)
		} else {
			akRespMsgText := fmt.Sprintf("`%v` 快捷回复配置`%v`成功✅", rkb, optMsg)
			tkbs := MakeReplyKeyboard(models.GlobalEnv)
			akRespMsg := tgbotapi.NewMessage(akMsg.Chat.ID, akRespMsgText)
			akRespMsg.ReplyToMessageID = akMsg.MessageID
			akRespMsg.ReplyMarkup = tkbs
			akRespMsg.ParseMode = tgbotapi.ModeMarkdown
			bot.Send(akRespMsg)
		}
	} else {
		akRespMsg := tgbotapi.NewMessage(akMsg.Chat.ID, "快捷回复配置添加格式错误❌\n\n示例：\n/ak 键盘显示===/cmd echo 'show reply keyboard' ")
		akRespMsg.ReplyToMessageID = akMsg.MessageID
		bot.Send(akRespMsg)
	}
}

// DelReplyKeyboard
// @description   删除快捷回复键盘指令
// @auth       iouAkira
// @param1     akMsg	*tgbotapi.Message
// @param2     bot		*tgbotapi.BotAPI
func DelReplyKeyboard(dkMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	rkb := strings.TrimLeft(dkMsg.Text[3:], " ")
	if len(strings.Split(rkb, "===")) == 1 {
		if !CheckDirOrFileIsExist(models.GlobalEnv.ReplyKeyboardFilePath) {
			dkRespMsg := tgbotapi.NewMessage(dkMsg.Chat.ID, "不存在快捷回复配置文件，无法删除不存在的东西⚠️")
			dkRespMsg.ReplyToMessageID = dkMsg.MessageID
			bot.Send(dkRespMsg)
		}
		optMsg, err := ReplyKeyboardFileOpt(rkb, rkb, "D")
		if err != nil {
			dkRespMsg := tgbotapi.NewMessage(dkMsg.Chat.ID, err.Error())
			dkRespMsg.ReplyToMessageID = dkMsg.MessageID
			bot.Send(dkRespMsg)
		} else {
			if optMsg == "" {
				dkRespMsgText := fmt.Sprintf(" 不存在需要`%v`的配置`%v`⚠️", optMsg, rkb)
				dkRespMsg := tgbotapi.NewMessage(dkMsg.Chat.ID, dkRespMsgText)
				dkRespMsg.ReplyToMessageID = dkMsg.MessageID
				dkRespMsg.ParseMode = tgbotapi.ModeMarkdown
				bot.Send(dkRespMsg)
			} else {
				dkRespMsgText := fmt.Sprintf("`%v` 快捷回复配置`%v`成功✅", rkb, optMsg)
				tkbs := MakeReplyKeyboard(models.GlobalEnv)
				dkRespMsg := tgbotapi.NewMessage(dkMsg.Chat.ID, dkRespMsgText)
				dkRespMsg.ReplyToMessageID = dkMsg.MessageID
				dkRespMsg.ReplyMarkup = tkbs
				dkRespMsg.ParseMode = tgbotapi.ModeMarkdown
				bot.Send(dkRespMsg)
			}
		}
	} else {
		akRespMsg := tgbotapi.NewMessage(dkMsg.Chat.ID, "快捷回复配置删除格式错误❌\n\n示例：\n/dk 键盘显示 (就是下面见面按钮显示内容)")
		akRespMsg.ReplyToMessageID = dkMsg.MessageID
		bot.Send(akRespMsg)
	}
}

// ClearReplyKeyboard
// @description   清楚所有快捷回复键盘指令
// @auth       iouAkira
// @param1     akMsg	*tgbotapi.Message
// @param2     bot		*tgbotapi.BotAPI
func ClearReplyKeyboard(clkMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	tgbotapi.NewRemoveKeyboard(true)
	clkRespMsg := tgbotapi.NewMessage(clkMsg.Chat.ID, "快捷回复键盘已清除🆑")
	clkRespMsg.ReplyToMessageID = clkMsg.MessageID
	clkRespMsg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{
		RemoveKeyboard: true,
		Selective:      false,
	}
	if _, err := bot.Send(clkRespMsg); err != nil {
		log.Printf("清除快捷回复键盘报错❌\n%v", err)
	}

}

// HandlerDocumentMsg
// @description   响应bot接收到文件类型消息
// @auth       iouAkira
// @param1     akMsg	*tgbotapi.Message
// @param2     bot		*tgbotapi.BotAPI
func HandlerDocumentMsg(docMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if CheckDirOrFileIsExist(models.GlobalEnv.CustomFilePath) {
		os.MkdirAll(models.GlobalEnv.CustomFilePath, os.ModePerm)
	}
	docF := docMsg.Document
	fileSuffix := strings.ReplaceAll(path.Ext(docF.FileName), ".", "")
	var keyboardMarkup tgbotapi.InlineKeyboardMarkup
	if fileSuffix == "js" || fileSuffix == "sh" || fileSuffix == "py" {
		if CheckDirOrFileIsExist(fmt.Sprintf("%v/%v", models.GlobalEnv.CustomFilePath, docF.FileName)) {
			var existsRow []tgbotapi.InlineKeyboardButton
			existsRow = append(existsRow, tgbotapi.NewInlineKeyboardButtonData("覆盖仅保存💾", fmt.Sprintf("%vFileSave replace", fileSuffix)))
			existsRow = append(existsRow, tgbotapi.NewInlineKeyboardButtonData("覆盖保存并执行⚡️", fmt.Sprintf("%vFileSaveRun replace", fileSuffix)))
			keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, existsRow)
			var existsRow2 []tgbotapi.InlineKeyboardButton
			existsRow2 = append(existsRow2, tgbotapi.NewInlineKeyboardButtonData("重命名仅保存💾", fmt.Sprintf("%vFileSave rename", fileSuffix)))
			existsRow2 = append(existsRow2, tgbotapi.NewInlineKeyboardButtonData("重命名保存并执行⚡", fmt.Sprintf("%vFileSaveRun rename", fileSuffix)))
			keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, existsRow2)
		} else {
			var existsRow []tgbotapi.InlineKeyboardButton
			existsRow = append(existsRow, tgbotapi.NewInlineKeyboardButtonData("仅保存💾", fmt.Sprintf("%vFileSave default", fileSuffix)))
			existsRow = append(existsRow, tgbotapi.NewInlineKeyboardButtonData("保存并执行⚡️", fmt.Sprintf("%vFileSaveRun default", fileSuffix)))
			keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, existsRow)
		}
		keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("取消", "cancel")))
		respMsg := tgbotapi.NewMessage(docMsg.Chat.ID, fmt.Sprintf("文件保存路径为`%v`，该路径在容器挂载目录内，方便查看，且同时会在`%v`保存一份方便执行调用。\n\n请选择对`%v`文件的操作️", models.GlobalEnv.CustomFilePath, models.GlobalEnv.DDnodeBtnFilePath, docF.FileName))
		respMsg.ReplyMarkup = keyboardMarkup
		respMsg.ReplyToMessageID = docMsg.MessageID
		respMsg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(respMsg)
	} else {
		respMsg := tgbotapi.NewMessage(docMsg.Chat.ID, "暂时只支持`js文件`、`shell文件`保存执行等操作⚠️")
		respMsg.ReplyToMessageID = docMsg.MessageID
		respMsg.ParseMode = tgbotapi.ModeMarkdown
		bot.Send(respMsg)
	}
}

// UnknownsCommand
// @description   响应未知指令
// @auth       iouAkira
// @param1     akMsg	*tgbotapi.Message
// @param2     bot		*tgbotapi.BotAPI
func UnknownsCommand(unCmdMsg *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if models.GlobalEnv.ReplyKeyBoard[unCmdMsg.Text] != "" {
		mapCmd := models.GlobalEnv.ReplyKeyBoard[unCmdMsg.Text][1:]
		LofDevLog(models.GlobalEnv.ReplyKeyBoard[unCmdMsg.Text])
		switch strings.Split(mapCmd, " ")[0] {
		case "help", "start":
			Help(unCmdMsg.Chat.ID, bot)
		//case "spnode":
		//	execSpnode(unCmdMsg, bot, replyKeyBoard[unCmdMsg.Text])
		//case "logs":
		//	execLogs(unCmdMsg, bot, replyKeyBoard[unCmdMsg.Text])
		//case "genCode":
		//	go genShareCodeMsg(unCmdMsg, bot, replyKeyBoard[unCmdMsg.Text])
		//case "rdc":
		//	execReadCookies(unCmdMsg, bot)
		//case "cmd":
		//	execOtherCmd(unCmdMsg, bot, replyKeyBoard[unCmdMsg.Text])
		default:
			text := "请勿发送错误的指令消息"
			if _, err := bot.Send(tgbotapi.NewMessage(unCmdMsg.Chat.ID, text)); err != nil {
				log.Println(err)
			}
		}
	} else {
		text := "请勿发送错误的指令消息"
		if _, err := bot.Send(tgbotapi.NewMessage(unCmdMsg.Chat.ID, text)); err != nil {
			log.Println(err)
		}
	}
}

// HandlerCallBackOption
// @description   响应聊天信息里的按钮点击事件
// @auth       iouAkira
// @param1     callbackQuery	*tgbotapi.CallbackQuery
// @param2     bot		*tgbotapi.BotAPI
func HandlerCallBackOption(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	fileOptions := []string{"jsFileSave",
		"jsFileSaveRun",
		"shFileSave",
		"shFileSaveRun",
		"pyFileSave",
		"pyFileSaveRun",
		"jsUrlFileSave",
		"jsUrlFileSaveRun",
		"shUrlFileSave",
		"shUrlFileSaveRun",
		"pyUrlFileSave",
		"pyUrlFileSaveRun",
	}
	cbDataSplit := strings.Split(callbackQuery.Data, " ")
	if len(cbDataSplit) == 1 {
		LofDevLog(callbackQuery.Data)
	} else {
		if IsContain(fileOptions, cbDataSplit[0]) {
			//saveAndRunFile(callbackQuery, bot)
			return
		}
		if strings.HasPrefix(callbackQuery.Data, "logs") {
			editOrgMsg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID,
				fmt.Sprintf("正在获取`%v`文件....", cbDataSplit[1]))
			editOrgMsg.ParseMode = tgbotapi.ModeMarkdown
			if _, err := bot.Send(editOrgMsg); err != nil {
				log.Printf("编辑inlineButton消息出错：%v", err)
			}
			bytes, readErr := ioutil.ReadFile(cbDataSplit[1])
			if readErr != nil {
				editMsg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID,
					editOrgMsg.MessageID,
					fmt.Sprintf("获取`%v`文件出错❌", cbDataSplit[1]))
				editMsg.ParseMode = tgbotapi.ModeMarkdown
				if _, err := bot.Send(editMsg); err != nil {
					log.Println(err)
				}
				return
			}
			fileSend := tgbotapi.FileBytes{
				Name:  cbDataSplit[1],
				Bytes: bytes,
			}

			respMsg := tgbotapi.NewDocument(callbackQuery.Message.Chat.ID, fileSend)
			respMsg.Caption = fmt.Sprintf("获取`%v`文件成功✅️", cbDataSplit[1])
			respMsg.ParseMode = tgbotapi.ModeMarkdown
			if _, err := bot.Send(respMsg); err != nil {
				log.Println(err)
			}

			delMsg := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, editOrgMsg.MessageID)
			bot.Send(delMsg)
		} else if strings.HasPrefix(callbackQuery.Data, "renew") {
			// respMsg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, callbackQuery.Data)
			// message, _ := bot.Send(respMsg)
			// renewCookieByWSKey(&message, bot)
		} else if strings.HasPrefix(callbackQuery.Data, "rdc") {
			////追加功能 1.读取cookie from pin 2.续期cookie from wskey
			//id := strings.Split(callbackQuery.Data, " ")[2]
			//if userCookie, err := ddUtil.ReadCookiesByID(CookiesListFilePath, id); err != nil {
			//	log.Printf("读取cookies.list文件出错。。%s", err)
			//} else {
			//	respMsgEdit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, userCookie)
			//	numericKeyboard := ddUtil.MakeKeyboardMarkup("rdc", 2, CookiesWSKeyListFilePath, id)
			//	respMsgEdit.ReplyMarkup = &numericKeyboard
			//	if _, err = bot.Send(respMsgEdit); err != nil {
			//		log.Printf("发送消息时出错❌%v", err)
			//	}
			//}
		} else {
			respMsg := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID, fmt.Sprintf("`/%v` 正在执行⚡️", strings.Join(cbDataSplit, " ")))
			respMsg.ParseMode = tgbotapi.ModeMarkdown
			respMsgInfo, _ := bot.Send(respMsg)
			execResult, isFile, err := ExecCommand(cbDataSplit, cbDataSplit[0], models.GlobalEnv.LogsBtnFilePath)
			if err != nil {
				log.Println(err)
				if isFile {
					respMsgDel := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, respMsgInfo.MessageID)
					bot.Send(respMsgDel)
					//需要传入绝对路径
					bytes, _ := ioutil.ReadFile(execResult)
					fileSend := tgbotapi.FileBytes{
						Name:  "bot_exec.log",
						Bytes: bytes,
					}
					respMsgFile := tgbotapi.NewDocument(callbackQuery.Message.Chat.ID, fileSend)
					respMsgFile.Caption = fmt.Sprintf("`/%v` 执行出错❌", strings.Join(cbDataSplit, " "))
					respMsgFile.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(respMsgFile)
					_ = os.Remove(execResult)
				} else {
					respMsgEdit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID,
						respMsgInfo.MessageID,
						fmt.Sprintf("`/%v` 执行出错❌\n\n```\n%v```", strings.Join(cbDataSplit, " "), err))
					respMsgEdit.ParseMode = tgbotapi.ModeMarkdown
					_, _ = bot.Send(respMsgEdit)
				}
			} else {
				//log.Printf(execResult)
				execStatus := "执行成功✅"
				if strings.HasPrefix(execResult, "stderr") {
					execStatus = "执行出错❌"
				}
				if isFile {
					respMsgDel := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, respMsgInfo.MessageID)
					bot.Send(respMsgDel)
					//需要传入绝对路径
					bytes, _ := ioutil.ReadFile(execResult)
					fileSend := tgbotapi.FileBytes{
						Name:  "bot_exec.log",
						Bytes: bytes,
					}
					respMsgFile := tgbotapi.NewDocument(callbackQuery.Message.Chat.ID, fileSend)
					respMsgFile.Caption = fmt.Sprintf("`/%v` %v", strings.Join(cbDataSplit, " "), execStatus)
					respMsgFile.ParseMode = tgbotapi.ModeMarkdown
					bot.Send(respMsgFile)
					_ = os.Remove(execResult)
				} else {
					respMsgEdit := tgbotapi.NewEditMessageText(callbackQuery.Message.Chat.ID,
						respMsgInfo.MessageID,
						fmt.Sprintf("`/%v` %v\n\n```\n%v```", strings.Join(cbDataSplit, " "), execStatus, execResult))
					respMsgEdit.ParseMode = tgbotapi.ModeMarkdown
					_, _ = bot.Send(respMsgEdit)
				}
			}
		}
	}
}
