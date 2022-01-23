package pre_init

import (
	ddCmd "ddbot/dd_cmd"
	"ddbot/models"
	ddutils "ddbot/utils"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strconv"
	"strings"
)

//defaultRepoBaseDir := "/iouRepos/dd_scripts"
//defaultDataBaseDir := "/data/dd_data"
const (
	_argEnv            = "env"
	_argUp             = "up"
	_argShareCode      = "commitShareCode"
	_argSyncRepo       = "syncRepo"
	_argRenewCookie    = "renewCookie"
	defaultRepoBaseDir = "/iouRepos/dd_scripts"
	defaultDataBaseDir = "/data/dd_data"
)

var (
	envFilePath           string
	envParams             string
	upParams              string
	repoBaseDir           string
	dataBaseDir           string
	wskeyListFilePath     string
	cookieListFilePath    string
	replyKeyboardFilePath string
	tgUserIDStr           string
	tgBotToken            string
	tgUserID              int64
)

// LoadEnv 使用bot需要的一些配置变量初始化
func LoadEnv() string {
	envFilePath = fmt.Sprintf("%v/env.sh", defaultDataBaseDir)
	app := cli.NewApp()
	app.Usage = "ddBot base on tgAPI"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     _argEnv,
			Value:    "envFilePath",
			Usage:    "[必填项]设置env.sh路径，否则程序将不启动。",
			Required: true,
			//FilePath: envFilePath,
			//DefaultText:
			Destination: &envParams,
			Aliases:     []string{"e"},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:  _argUp,
			Aliases: []string{"u"},
			Usage: "启动bot；[commitShareCode]为提交互助码到助力池；[syncRepo]为同步仓库代码；[renewCookie]为给所有wskey续期",
			Subcommands: []*cli.Command{
				{
					Name:     _argRenewCookie,
					Usage:    "wskey续期cookie",
					Category: "up",
					Action: func(ctx *cli.Context) error {
						fmt.Printf("开始给 %v 里面的全部wskey续期...\n", models.GlobalEnv.CookiesListFilePath)
						ddutils.RenewAllCookie()
						return nil
					},
				},
				{
					Name:     _argShareCode,
					Usage:    "上传互助码",
					Category: "up",
					Action: func(ctx *cli.Context) error {
						fmt.Printf("开始上传互助码...\n")
						ddutils.UploadShareCode(models.GlobalEnv)
						return nil
					},
				},
				{
					Name:     _argSyncRepo,
					Usage:    "同步仓库代码",
					Category: "up",
					Action: func(ctx *cli.Context) error {
						fmt.Printf("开始同步仓库代码。\n")
						if models.GlobalEnv.RepoBaseDir != "" && strings.HasPrefix(models.GlobalEnv.RepoBaseDir, "/") {
							ddutils.SyncRepo(models.GlobalEnv)
						} else {
							fmt.Printf("同步仓库设定的目录[%v]不规范，退出同步。\n", models.GlobalEnv.RepoBaseDir)
						}
						return nil
					},
				},
			},
			Action: func(ctx *cli.Context) error {
				engine := SetupRouters()
				engine.Run(models.GlobalEnv.TgBotToken,
					models.GlobalEnv.TgUserID,
					ddCmd.DebugMode(false),
					ddCmd.TimeOut(60),
				)
				return nil
			},
		},
	}
	app.Before = func(context *cli.Context) error {
		checkAfterInit()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
	return upParams
}

func checkAfterInit() {
	fmt.Println("执行前检查...")
	fmt.Printf("1. env 启动参数值:[%v];\n", envParams)
	if ddutils.CheckDirOrFileIsExist(envParams) {
		envFilePath = envParams
	} else {
		fmt.Printf("[%v] ddbot需要是用相关环境变量配置文件不存在，确认目录文件是否存在\n", envParams)
		os.Exit(0)
	}
	repoBaseDir = ddutils.GetEnvFromEnvFile(envFilePath, "REPO_BASE_DIR")
	if repoBaseDir == "" {
		fmt.Printf("未查找到仓库的基础目录配置信息，停止启动。使用默认仓库路径[%v]\n", defaultRepoBaseDir)
		repoBaseDir = defaultDataBaseDir
	} else {
		fmt.Printf("2. 仓库的基础目录配置信息[%v]\n", repoBaseDir)
	}
	dataBaseDir = ddutils.GetEnvFromEnvFile(envFilePath, "DATA_BASE_DIR")
	if dataBaseDir == "" || !ddutils.CheckDirOrFileIsExist(dataBaseDir) {
		fmt.Printf("未查找到数据存放目录配置信息，停止启动。\n")
		os.Exit(0)
	} else {
		fmt.Printf("3. 数据存放目录配置信息[%v]\n", dataBaseDir)
	}
	wskeyListFilePath = ddutils.GetEnvFromEnvFile(envFilePath, "WSKEY_FILE_PATH")
	if wskeyListFilePath == "" {
		wskeyListFilePath = fmt.Sprintf("%v/cookies_wskey.list", dataBaseDir)
	}

	cookieListFilePath = ddutils.GetEnvFromEnvFile(envFilePath, "DDCK_FILE_PATH")
	if cookieListFilePath == "" {
		cookieListFilePath = fmt.Sprintf("%v/cookies.list", dataBaseDir)
	}

	replyKeyboardFilePath = ddutils.GetEnvFromEnvFile(envFilePath, "REPLY_KEYBOARD_FILE_PATH")
	if replyKeyboardFilePath == "" {
		replyKeyboardFilePath = fmt.Sprintf("%v/reply_keyboard.list", dataBaseDir)
	}

	if ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN_HANDLER") != "" {
		tgBotToken = ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN_HANDLER")
	} else {
		tgBotToken = ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN")
	}
	tgUserIDStr = ddutils.GetEnvFromEnvFile(envFilePath, "TG_USER_ID")
	if tgUserIDStr != "" {
		convTgUserID, err := strconv.ParseInt(tgUserIDStr, 10, 64)
		if err == nil {
			tgUserID = convTgUserID
		}
	}
	if tgBotToken == "" && tgUserID <= 0 {
		fmt.Printf("请检查TG配置信息。\n")
		os.Exit(0)
	}
	replyKeyBoard := map[string]string{
		"选择脚本执行⚡️": "/ddnode",
		"选择日志下载⬇️": "/logs",
		"更新仓库代码🔄":  fmt.Sprintf("/cmd cd %v ; sh iou-entry.sh", repoBaseDir),
		"查看账号🍪":    "/rdc",
		"查看系统进程⛓":  "/cmd ps -ef|grep -v 'grep\\| ts\\|/ts\\| sh'",
		"查看帮助说明📝":  ">help",
	}
	models.GlobalEnv = &models.DDEnv{
		RepoBaseDir:              repoBaseDir,
		DataBaseDir:              dataBaseDir,
		DDnodeBtnFilePath:        repoBaseDir,
		LogsBtnFilePath:          fmt.Sprintf("%v/logs", dataBaseDir),
		CustomFilePath:           fmt.Sprintf("%v/custom_scripts", dataBaseDir),
		CookiesWSKeyListFilePath: wskeyListFilePath,
		CookiesListFilePath:      cookieListFilePath,
		ReplyKeyboardFilePath:    replyKeyboardFilePath,
		EnvFilePath:              envFilePath,
		TgBotToken:               tgBotToken,
		TgUserID:                 tgUserID,
		ReplyKeyBoard:            replyKeyBoard,
	}
}
