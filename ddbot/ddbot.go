package main

import (
	"ddbot/models"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	ddutils "ddbot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var wg sync.WaitGroup
var bot *tgbotapi.BotAPI

var (
	RepoBaseDir              = "/iouRepos/dd_scripts"
	DataBaseDir              = "/Users/akira-work/data/dd_data"
	EnvFilePath              = fmt.Sprintf("%v/env.sh", DataBaseDir)
	SpnodeBtnFilePath        = fmt.Sprintf(RepoBaseDir)
	LogsBtnFilePath          = fmt.Sprintf("%v/logs", DataBaseDir)
	CookiesListFilePath      = fmt.Sprintf("%v/cookies.list", DataBaseDir)
	CookiesWSKeyListFilePath = fmt.Sprintf("%v/cookies_wskey.list", DataBaseDir)
	ReplyKeyboardFilePath    = fmt.Sprintf("%v/ReplyKeyBoard.list", DataBaseDir)
	CustomFilePath           = fmt.Sprintf("%v/custom_scripts", DataBaseDir)
	TgBotToken               = ""
	TgUserID                 = int64(0)
	ReplyKeyBoard            = map[string]string{
		"选择脚本执行⚡️": "/spnode",
		"选择日志下载⬇️": "/logs",
		"更新仓库代码🔄": "/cmd docker_entrypoint.sh",
		"查看账号🍪":   "/rdc",
		"查看系统进程⛓":  "/cmd ps -ef|grep -v 'grep\\| ts\\|/ts\\| sh'",
		"查看帮助说明📝": "/help",
	}
)

// ddConfig 组合常用的参数
var ddConfig = new(models.DDEnv)

func main() {
	//构建Linux amd64(x86_64)   CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ddBot-amd64 ddBot.go
	//构建Linux arm64(aarch64)  CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ddBot-arm64 ddBot.go
	//构建Linux arm64(armv7,v6) CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o ddBot-arm ddBot.go
	//构建Windows下可执行文件 CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build goBot.go
	//构建macOS下可执行文件   CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build goBot.go
	var envParams string
	var upParams string
	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量
	flag.StringVar(&envParams, "env", EnvFilePath, fmt.Sprintf("默认为[%v],如果env.sh文件不在该默认路径，请使用-env指定，否则程序将不启动。", EnvFilePath))
	flag.StringVar(&upParams, "up", "", "默认为空，为启动bot；commitShareCode为提交互助码到助力池；syncRepo为同步仓库代码；")
	flag.Parse()
	log.Printf("-env 启动参数值:[%v]; -up 启动参数值:[%v]", envParams, upParams)
	if ddutils.CheckDirOrFileIsExist(envParams) {
		EnvFilePath = envParams
	} else {
		log.Printf("[%v] ddbot需要是用相关环境变量配置文件不存在，确认目录文件是否存在", envParams)
		os.Exit(0)
	}
	//读取加载程序需要使用的环境变量
	loadEnv(EnvFilePath)

	// -up 启动参数 不指定默认启动ddbot
	if upParams != "" {
		log.Printf("传入 -up参数：%v ", upParams)
		if upParams == "commitShareCode" {
			log.Printf("启动程序指定了 -up 参数为 %v 开始上传互助码。", upParams)
			ddutils.UploadShareCode(ddConfig)
		} else if upParams == "syncRepo" {
			log.Printf("启动程序指定了 -up 参数为 %v 开始同步仓库代码。", upParams)
			if ddConfig.RepoBaseDir != "" {
				ddutils.SyncRepo(ddConfig)
			} else {
				log.Printf("同步仓库设定的目录[%v]不规范，退出同步。", ddConfig.RepoBaseDir)
			}
		} else if upParams == "renewCookie" {
			log.Printf("启动程序指定了 -up 参数为 %v 开始给 %v 里面的全部wskey续期。", upParams, CookiesWSKeyListFilePath)
			ddutils.RenewAllCookie(ddConfig)
		} else {
			log.Printf("请传入传入的对应 -up参数：%v ", upParams)
		}
		os.Exit(0)
	}

	var startErr error
	bot, startErr = tgbotapi.NewBotAPI(TgBotToken)
	if startErr != nil {
		log.Panicf("start bot failed with some error %v", startErr)
		// os.Exit(0)
	}
	log.Printf("Telegram bot stared，Bot info ==> %s %s[%s]", bot.Self.FirstName, bot.Self.LastName, bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	ddutils.LoadReplyKeyboardMap(ddConfig)
	for update := range updates {

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		// 接收消息处理
		if update.Message != nil {
			if update.Message.From.ID != TgUserID {
				continue
			}
			// 文件消息处理
			if update.Message.Document != nil {
				go ddutils.HandlerDocumentMsg(update.Message, bot, ddConfig)
				continue
			}
			// 普通文本消息处理
			switch update.Message.Command() {
			case "help", "start":
				go ddutils.Help(update.Message.Chat.ID, bot, ddConfig)
			case "ak":
				go ddutils.AddReplyKeyboard(update.Message, bot, ddConfig)
			case "dk":
				go ddutils.DelReplyKeyboard(update.Message, bot, ddConfig)
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
				go ddutils.UnknownsCommand(update.Message, bot, ddConfig)
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
					bot.Send(respMsg)
				}()
			} else {
				go ddutils.HandlerCallBackOption(update.CallbackQuery, bot, ddConfig)
			}
			log.Printf("update.CallbackQuery.Data %v", update.CallbackQuery.Data)
		}
	}
	wg.Wait()
}

// loadEnv
// @description   使用bot需要的一些配置变量初始化
// @auth      iouAkira
// @param     envFilePath string env.sh环境变量配置文件的绝对路径
func loadEnv(envFilePath string) {
	RepoBaseDir = ddutils.GetEnvFromEnvFile(envFilePath, "REPO_BASE_DIR")

	if RepoBaseDir == "" || !ddutils.CheckDirOrFileIsExist(RepoBaseDir) {
		log.Printf("未查找到仓库的基础目录配置信息，停止启动。")
		os.Exit(0)
	} else {
		SpnodeBtnFilePath = fmt.Sprintf(RepoBaseDir)
		log.Printf("仓库的基础目录配置信息[%v]", RepoBaseDir)
	}

	DataBaseDir = ddutils.GetEnvFromEnvFile(envFilePath, "DATA_BASE_DIR")
	if DataBaseDir == "" || !ddutils.CheckDirOrFileIsExist(DataBaseDir) {
		log.Printf("未查找到数据存放目录配置信息，停止启动。")
		os.Exit(0)
	} else {
		LogsBtnFilePath = fmt.Sprintf("%v/logs", DataBaseDir)
		CustomFilePath = fmt.Sprintf("%v/custom_scripts", DataBaseDir)
		log.Printf("数据存放目录配置信息[%v]", DataBaseDir)
	}

	CookiesWSKeyListFilePath = ddutils.GetEnvFromEnvFile(envFilePath, "WSKEY_FILE_PATH")
	if CookiesWSKeyListFilePath == "" {
		CookiesWSKeyListFilePath = fmt.Sprintf("%v/cookies_wskey.list", DataBaseDir)
	}

	CookiesListFilePath = ddutils.GetEnvFromEnvFile(envFilePath, "DDCK_FILE_PATH")
	if CookiesListFilePath == "" {
		CookiesListFilePath = fmt.Sprintf("%v/cookies.list", DataBaseDir)
	}

	ReplyKeyboardFilePath = ddutils.GetEnvFromEnvFile(envFilePath, "REPLY_KEYBOARD_FILE_PATH")
	if ReplyKeyboardFilePath == "" {
		ReplyKeyboardFilePath = fmt.Sprintf("%v/reply_keyboard.list", DataBaseDir)
	}

	TgBotTokenHandler := ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN_HANDLER")
	TgBotTokenNotify := ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN")
	if TgBotTokenHandler != "" {
		TgBotToken = TgBotTokenHandler
	} else if TgBotTokenNotify != "" {
		TgBotToken = TgBotTokenNotify
	}
	TgUserIDStr := ddutils.GetEnvFromEnvFile(envFilePath, "TG_USER_ID")
	if TgUserIDStr != "" {
		convTgUserID, err := strconv.ParseInt(TgUserIDStr, 10, 64)
		if err == nil {
			TgUserID = convTgUserID
		}
	}
	if TgBotToken == "" || TgUserID == 0 {
		log.Printf("Telegram Bot相关环境变量配置不完整，故不启动。(botToken=%v;tgUserID=%v)", TgBotToken, TgUserID)
		os.Exit(0)
	}
	ddConfig = &models.DDEnv{
		RepoBaseDir:              RepoBaseDir,
		DataBaseDir:              DataBaseDir,
		SpnodeBtnFilePath:        SpnodeBtnFilePath,
		LogsBtnFilePath:          LogsBtnFilePath,
		CustomFilePath:           CustomFilePath,
		CookiesWSKeyListFilePath: CookiesWSKeyListFilePath,
		CookiesListFilePath:      CookiesListFilePath,
		ReplyKeyboardFilePath:    ReplyKeyboardFilePath,
		EnvFilePath:              EnvFilePath,
		TgBotToken:               TgBotToken,
		TgUserID:                 TgUserID,
		ReplyKeyBoard:            ReplyKeyBoard,
	}
}
