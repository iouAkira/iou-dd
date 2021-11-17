package init

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"ddbot/models"
	ddutils "ddbot/utils"
)

// LoadEnv
// @description   使用bot需要的一些配置变量初始化
// @auth      iouAkira
// @param     envFilePath string env.sh环境变量配置文件的绝对路径
func LoadEnv() {
	defaultRepoBaseDir := "/iouRepos/dd_scripts"
	defaultDataBaseDir := "/data/dd_data"
	envFilePath := fmt.Sprintf("%v/env.sh", defaultDataBaseDir)

	var envParams string
	// StringVar用指定的名称、控制台参数项目、默认值、使用信息注册一个string类型flag，并将flag的值保存到p指向的变量
	flag.StringVar(&envParams, "env", envFilePath, fmt.Sprintf("默认为[%v],如果env.sh文件不存在于该默认路径，请使用-env指定，否则程序将不启动。", envFilePath))
	flag.Parse()
	fmt.Printf("-env 启动参数值:[%v];", envParams)
	if ddutils.CheckDirOrFileIsExist(envParams) {
		envFilePath = envParams
	} else {
		fmt.Printf("[%v] ddbot需要是用相关环境变量配置文件不存在，确认目录文件是否存在", envParams)
		os.Exit(0)
	}

	repoBaseDir := ddutils.GetEnvFromEnvFile(envFilePath, "REPO_BASE_DIR")
	if repoBaseDir == "" {
		log.Printf("未查找到仓库的基础目录配置信息，停止启动。使用默认仓库路径[%v]", defaultRepoBaseDir)
		repoBaseDir = defaultDataBaseDir
	} else {
		log.Printf("仓库的基础目录配置信息[%v]", repoBaseDir)
	}

	dataBaseDir := ddutils.GetEnvFromEnvFile(envFilePath, "DATA_BASE_DIR")
	if dataBaseDir == "" || !ddutils.CheckDirOrFileIsExist(dataBaseDir) {
		log.Printf("未查找到数据存放目录配置信息，停止启动。")
		os.Exit(0)
	} else {
		log.Printf("数据存放目录配置信息[%v]", dataBaseDir)
	}

	wskeyListFilePath := ddutils.GetEnvFromEnvFile(envFilePath, "WSKEY_FILE_PATH")
	if wskeyListFilePath == "" {
		wskeyListFilePath = fmt.Sprintf("%v/cookies_wskey.list", dataBaseDir)
	}

	cookieListFilePath := ddutils.GetEnvFromEnvFile(envFilePath, "DDCK_FILE_PATH")
	if cookieListFilePath == "" {
		cookieListFilePath = fmt.Sprintf("%v/cookies.list", dataBaseDir)
	}

	replyKeyboardFilePath := ddutils.GetEnvFromEnvFile(envFilePath, "REPLY_KEYBOARD_FILE_PATH")
	if replyKeyboardFilePath == "" {
		replyKeyboardFilePath = fmt.Sprintf("%v/reply_keyboard.list", dataBaseDir)
	}
	tgBotToken := ""
	tgUserID := int64(0)
	if ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN_HANDLER") != "" {
		tgBotToken = ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN_HANDLER")
	} else {
		ddutils.GetEnvFromEnvFile(envFilePath, "TG_BOT_TOKEN")
	}

	tgUserIDStr := ddutils.GetEnvFromEnvFile(envFilePath, "TG_USER_ID")
	if tgUserIDStr != "" {
		convTgUserID, err := strconv.ParseInt(tgUserIDStr, 10, 64)
		if err == nil {
			tgUserID = convTgUserID
		}
	}
	replyKeyBoard := map[string]string{
		"选择脚本执行⚡️": "/spnode",
		"选择日志下载⬇️": "/logs",
		"更新仓库代码🔄": "/cmd docker_entrypoint.sh",
		"查看账号🍪":   "/rdc",
		"查看系统进程⛓":  "/cmd ps -ef|grep -v 'grep\\| ts\\|/ts\\| sh'",
		"查看帮助说明📝": "/help",
	}
	models.GlobalEnv = &models.DDEnv{
		RepoBaseDir:              repoBaseDir,
		DataBaseDir:              dataBaseDir,
		SpnodeBtnFilePath:        repoBaseDir,
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
