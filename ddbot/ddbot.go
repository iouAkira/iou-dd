package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"ddbot/models"
	"ddbot/pre_init"
	ddutils "ddbot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var wg sync.WaitGroup
var bot *tgbotapi.BotAPI

func main() {

	// 读取加载程序需要使用的环境变量
	pre_init.LoadEnv()
	ddutils.ExecUpCommand()

	// 启动bot
	if models.GlobalEnv.TgBotToken == "" || models.GlobalEnv.TgUserID == 0 {
		fmt.Printf("Telegram Bot相关环境变量配置不完整，故不启动。(botToken=%v;tgUserID=%v)", models.GlobalEnv.TgBotToken, models.GlobalEnv.TgUserID)
		os.Exit(0)
	}

	var startErr error
	bot, startErr = tgbotapi.NewBotAPI(models.GlobalEnv.TgBotToken)
	if startErr != nil {
		log.Panicf("start bot failed with some error %v", startErr)
		// os.Exit(0)
	}
	log.Printf("Telegram bot stared，Bot info ==> %s %s[%s]", bot.Self.FirstName, bot.Self.LastName, bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	ddutils.LoadReplyKeyboardMap(models.GlobalEnv)
	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		// 接收消息处理
		if update.Message != nil {
			if update.Message.From.ID != models.GlobalEnv.TgUserID {
				continue
			}
			// 文件消息处理
			if update.Message.Document != nil {
				go ddutils.HandlerDocumentMsg(update.Message, bot)
				continue
			}
			// 普通文本消息处理
			switch update.Message.Command() {
			case "help", "start":
				go ddutils.Help(update.Message.Chat.ID, bot)
			case "ak":
				go ddutils.AddReplyKeyboard(update.Message, bot)
			case "dk":
				go ddutils.DelReplyKeyboard(update.Message, bot)
			case "clk":
				go ddutils.ClearReplyKeyboard(update.Message, bot)
			//case "dl":
			//	go downloadFileByUrl(update.Message, bot)
			//case "spnode":
			//	//log.Println(update.Message.Text)
			//	go execSpnode(update.Message, bot, "")
			//case "logs":
			//	//log.Println(update.Message.Text)
			//	go execLogs(update.Message, bot, "")
			//case "renew":
			//	go renewCookieByWSKey(update.Message, bot)
			//case "rdc":
			//	go execReadCookies(update.Message, bot)
			//case "bl":
			//	go beanStats(update.Message, bot)
			//case "env":
			//	go setEnvSH(update.Message, bot)
			//case "cmd":
			//	go execOtherCmd(update.Message, bot, "")
			//case "nty":
			//	go iouNotify(update.Message, bot)
			default:
				go ddutils.UnknownsCommand(update.Message, bot)
			}
		}
		// inlinebutton交互点击callback处理
		if update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "cancel" {
				edit := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID,
					"操作已经取消")
				_, _ = bot.Send(edit)
			} else if update.CallbackQuery.Data == ddutils.DELETE {
				go func() {
					respMsg := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
					_, _ = bot.Send(respMsg)
				}()
			} else {
				go ddutils.HandlerCallBackOption(update.CallbackQuery, bot)
			}
			log.Printf("update.CallbackQuery.Data %v", update.CallbackQuery.Data)
		}
	}
	wg.Wait()
}
