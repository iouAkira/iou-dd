package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"ddbot/models"
)

// ExecUpCommand
// @description   执行命令指定了 -up 参数的启动命令
// @auth      iouAkira
// @param     upParams 全局程序配置信息
func ExecUpCommand(upParams string) {
	// -up 启动参数 不指定默认启动ddbot
	if upParams != "" {
		fmt.Printf("传入 -up参数：%v \n", upParams)
		if upParams == "commitShareCode" {
			fmt.Printf("启动程序指定了 -up 参数为 %v 开始上传互助码。\n", upParams)
			UploadShareCode(models.GlobalEnv)
		} else if upParams == "syncRepo" {
			fmt.Printf("启动程序指定了 -up 参数为 %v 开始同步仓库代码。\n", upParams)
			if models.GlobalEnv.RepoBaseDir != "" {
				SyncRepo(models.GlobalEnv)
			} else {
				fmt.Printf("同步仓库设定的目录[%v]不规范，退出同步。\n", models.GlobalEnv.RepoBaseDir)
			}
		} else if upParams == "renewCookie" {
			fmt.Printf("启动程序指定了 -up 参数为 %v 开始给 %v 里面的全部wskey续期。\n", upParams, models.GlobalEnv.CookiesListFilePath)
			RenewAllCookie()
		} else {
			fmt.Printf("请传入传入的对应 -up参数：%v \n", upParams)
		}
		os.Exit(0)
	}
}

// UploadShareCode
// @description
// @auth       iouAkira
// @param1     config *models.DDEnv
func UploadShareCode(ddConfig *models.DDEnv) {
	confFilePath := fmt.Sprintf("%v/dd_sharecode.json", ddConfig.RepoBaseDir)
	if CheckDirOrFileIsExist(confFilePath) {
		conf, err := ioutil.ReadFile(confFilePath)
		if err != nil {
			fmt.Printf("读取配置文件异常`%v`  Error: %v\n", confFilePath, err)
		}

		shareCodeConf := models.ShareCode{}
		//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
		_ = json.Unmarshal(conf, &shareCodeConf)
		for _, codeInfo := range shareCodeConf.ShareCodeInfo {
			// 针对使用仓库脚本校验,参数为公开的配置
			if SubmitShareCodeCheck(ddConfig.RepoBaseDir, codeInfo) {
				fmt.Printf("仓库使用检查失败，取消上传[%v]到助力池。\n", codeInfo.ShareCodeType)
				continue
			}
			if os.Getenv(codeInfo.ShareCodeEnv) != "" || GetEnvFromEnvFile(ddConfig.EnvFilePath, codeInfo.ShareCodeEnv) != "" {
				fmt.Printf("本地配置了[%v]互助码，取消上传到助力池。\n", codeInfo.ShareCodeType)
				continue
			}
			logFile, lerr := ioutil.ReadFile(fmt.Sprintf("%v/%v", ddConfig.LogsBtnFilePath, codeInfo.LogFileName))
			if lerr != nil {
				fmt.Printf("读取日志文件异常[%v]  Error: %v\n", confFilePath, lerr)
				continue
			}
			logLines := strings.Split(string(logFile), "\n")
			var shareCode []string
			for _, logLine := range logLines {
				if strings.Contains(logLine, codeInfo.LogPrefix) {
					matchShareCode := strings.Split(logLine, codeInfo.LogPrefix)
					if len(matchShareCode) > 1 {
						shareCode = append(shareCode, matchShareCode[len(matchShareCode)-1])
					}
				}
			}
			shareCode = RemoveRepByLoop(shareCode)
			fmt.Printf("%v\n", shareCode)
			// compile_submit_host 编译时传入
			// compile_submit_token 编译时候传入
			commitUrl := fmt.Sprintf("https://compile_submit_host/api/%v/add/%v", codeInfo.ShareCodeType, strings.Join(shareCode, "&"))
			client := http.Client{}
			req, _ := http.NewRequest("POST", commitUrl, nil)
			req.Header.Add("Accept-Encoding", "identity")
			req.Header.Add("Content-Type", "multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW")
			req.Header.Add("Host", "compile_submit_host")
			req.Header.Add("Cookie", "compile_submit_token")

			resp, perr := client.Do(req)
			if perr != nil {
				fmt.Printf("获取请求结果Body报错。。。%v\n", shareCode)
				continue
			}
			if resp.StatusCode == 200 {
				rbody, _ := ioutil.ReadAll(resp.Body)
				fmt.Printf("%v\n", string(rbody))
			}
			_ = resp.Body.Close()

			time.Sleep(20 * time.Second)
		}
	}
}

// RenewAllCookie
// @description   更新所有cookie
// @auth       iouAkira
func RenewAllCookie() {
	wskeyFile, err := ioutil.ReadFile(models.GlobalEnv.CookiesWSKeyListFilePath)

	renewSleep := "Y"
	envRenewSleep := GetEnvFromEnvFile(models.GlobalEnv.EnvFilePath, "RENEW_SLEEP")
	if envRenewSleep != "" {
		renewSleep = envRenewSleep
	}
	if err != nil {
		fmt.Printf("读取cookies文件出错。。%s\n", err)
	}
	lines := strings.Split(string(wskeyFile), "\n")
	succCnt := 0
	failedCnt := 0
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "#") || lines[i] == "" {
			continue
		} else {
			lines[i] = strings.ReplaceAll(lines[i], "\r", "")
			pin := fmt.Sprintf("%v", i)

			r := regexp.MustCompile(`^(pin=)\S.*?;`)
			uaMatch := r.FindString(lines[i])
			LofDevLog("==>>%v", uaMatch)
			if uaMatch != "" {
				pin = strings.ReplaceAll(uaMatch, "pin=", "")
				pin = strings.ReplaceAll(pin, ";", "")
			}

			reNewCk, renewErr := RenewCookie(lines[i])
			if renewErr != nil || reNewCk == "" || strings.Contains(reNewCk, "pt_key=fake_") {
				if strings.Contains(reNewCk, "pt_key=fake_") {
					fmt.Printf("续期【账号%v】Cookie失败❌ =====> wskey 已经失效\n", pin)
				} else {
					fmt.Printf("续期【账号%v】Cookie失败❌\n%v\n", pin, renewErr.Error())
				}
				failedCnt += 1
			} else {
				fmt.Printf("续期【账号%v】 Cookie成功✅️\n", pin)
				if errW := writeCookiesFile(reNewCk); errW == nil {
					fmt.Printf("写入 cookies.list 成功✅\n")
					fmt.Printf("账户(%v) Cookie续期操作已完成✅\n", pin)
					succCnt += 1
					// }
				} else {
					fmt.Printf("（%v）写入 cookies.list 失败❌%v\n", reNewCk, errW.Error())
					failedCnt += 1
				}
			}
			if renewSleep == "Y" {
				fmt.Printf("休息20秒。。。\n")
				time.Sleep(10 * time.Second)
			}
		}
	}
	fmt.Printf("renewCookie续期任务已完成✅【续期成功：%v 个；续期失败：%v 个】\n", succCnt, failedCnt)
}
