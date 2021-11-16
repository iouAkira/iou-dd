package utils

import (
	"crypto/tls"
	"ddbot/models"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DELETE string = "delete"
	RETURN string = "return"
)

func MakeReplyKeyboard(config *models.DDEnv) tgbotapi.ReplyKeyboardMarkup {
	if CheckDirOrFileIsExist(config.ReplyKeyboardFilePath) {
		cookiesFile, err := ioutil.ReadFile(config.ReplyKeyboardFilePath)
		if err != nil {
			log.Printf("读取%v快捷回复配置文件出错。。%s", config.ReplyKeyboardFilePath, err)
		}

		lines := strings.Split(string(cookiesFile), "\n")
		//log.Printf("%v", lines)
		for _, line := range lines {
			lineSpt := strings.Split(line, "===")
			if len(lineSpt) > 1 {
				config.ReplyKeyBoard[lineSpt[0]] = lineSpt[1]
			}
		}
	}

	var replyKeyBoardKeys []string
	for k := range config.ReplyKeyBoard {
		replyKeyBoardKeys = append(replyKeyBoardKeys, k)
	}
	sort.Strings(replyKeyBoardKeys)

	var allRow [][]tgbotapi.KeyboardButton
	var keys []string

	for i, k := range replyKeyBoardKeys {
		//log.Printf("%v %v", i, replyKeyBoardKeys)
		keys = append(keys, k)
		if len(keys) == 2 || i == len(replyKeyBoardKeys)-1 {
			var row []tgbotapi.KeyboardButton
			for _, vi := range keys {
				row = append(row, tgbotapi.KeyboardButton{Text: vi})
			}
			allRow = append(allRow, row)
			keys = keys[0:0]
		}
	}
	replyKeyboards := tgbotapi.NewReplyKeyboard(allRow...)
	return replyKeyboards
}

func LoadReplyKeyboardMap(config *models.DDEnv) {
	if CheckDirOrFileIsExist(config.ReplyKeyboardFilePath) {
		cookiesFile, err := ioutil.ReadFile(config.ReplyKeyboardFilePath)
		if err != nil {
			log.Printf("读取%v快捷回复配置文件出错。。%s", config.ReplyKeyboardFilePath, err)
		}
		lines := strings.Split(string(cookiesFile), "\n")
		for _, line := range lines {
			lineSpt := strings.Split(line, "===")
			if len(lineSpt) > 1 {
				config.ReplyKeyBoard[lineSpt[0]] = lineSpt[1]
			}
		}
	}
}

func UploadShareCode(ddConfig *models.DDEnv) {
	confFilePath := fmt.Sprintf("%v/dd_sharecode.json", ddConfig.RepoBaseDir)
	if CheckDirOrFileIsExist(confFilePath) {
		conf, err := ioutil.ReadFile(confFilePath)
		if err != nil {
			log.Printf("读取配置文件异常`%v`  Error: %v", confFilePath, err)
		}

		shareCodeConf := models.ShareCode{}
		//下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
		json.Unmarshal(conf, &shareCodeConf)
		for _, codeInfo := range shareCodeConf.ShareCodeInfo {
			// 针对使用仓库脚本校验,参数为公开的配置
			if SubmitShareCodeCheck(ddConfig.RepoBaseDir, codeInfo) {
				log.Printf("仓库使用检查失败，取消上传到助力池。", codeInfo.ShareCodeType)
				continue
			}
			if os.Getenv(codeInfo.ShareCodeEnv) != "" || GetEnvFromEnvFile(ddConfig.EnvFilePath, codeInfo.ShareCodeEnv) != "" {
				log.Printf("本地配置了[%v]互助码，取消上传到助力池。", codeInfo.ShareCodeType)
				continue
			}
			logFile, lerr := ioutil.ReadFile(fmt.Sprintf("%v/%v", ddConfig.LogsBtnFilePath, codeInfo.LogFileName))
			if lerr != nil {
				log.Printf("读取日志文件异常[%v]  Error: %v", confFilePath, lerr)
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
			log.Println(shareCode)
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
				log.Printf("获取请求结果Body报错。。。%v", shareCode)
				continue
			}
			if resp.StatusCode == 200 {
				rbody, _ := ioutil.ReadAll(resp.Body)
				log.Println(string(rbody))
			}
			resp.Body.Close()

			time.Sleep(20 * time.Second)
		}
	}
}

// RenewAllCookie 命令行调用函数
func RenewAllCookie(ddConfig *models.DDEnv) {
	wskeyFile, err := ioutil.ReadFile(ddConfig.CookiesWSKeyListFilePath)

	renewSleep := "Y"
	envRenewSleep := GetEnvFromEnvFile(ddConfig.EnvFilePath, "RENEW_SLEEP")
	if envRenewSleep != "" {
		renewSleep = envRenewSleep
	}
	if err != nil {
		log.Printf("读取cookies文件出错。。%s", err)
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

			reNewCk, renewErr := RenewCookie(lines[i], ddConfig)
			if renewErr != nil || reNewCk == "" || strings.Contains(reNewCk, "pt_key=fake_") {
				if strings.Contains(reNewCk, "pt_key=fake_") {
					log.Printf("续期【账号%v】Cookie失败❌ =====> wskey 已经失效", pin)
				} else {
					log.Printf("续期【账号%v】Cookie失败❌\n%v", pin, renewErr.Error())
				}
				failedCnt += 1
			} else {
				log.Printf("续期【账号%v】 Cookie成功✅️", pin)
				if errW := WriteCookiesFile(reNewCk, ddConfig.CookiesListFilePath); errW == nil {
					log.Printf("写入 cookies.list 成功✅")
					log.Printf("账户(%v) Cookie续期操作已完成✅", pin)
					succCnt += 1
					// }
				} else {
					log.Printf("（%v）写入 cookies.list 失败❌%v", reNewCk, errW.Error())
					failedCnt += 1
				}
			}
			if renewSleep == "Y" {
				log.Printf("休息20秒。。。")
				time.Sleep(10 * time.Second)
			}
		}
	}
	log.Printf("renewCookie续期任务已完成✅【续期成功：%v 个；续期失败：%v 个】", succCnt, failedCnt)
}

func RenewCookie(wskey string, ddConfig *models.DDEnv) (string, error) {
	renewCK := ""
	//默认签名UA配置
	body := "body=%7B%22to%22%3A%22https%3A%5C%2F%5C%2Fplogin.m.jd.com%5C%2Fcgi-bin%5C%2Fm%5C%2Fthirdapp_auth_page%3Ftoken%3DAAEAIEbEUWDGA_SGHg4sHM5fwfnpt-kFtkZ_boToZQULiH0O%26client_type%3Dapple%26appid%3D1125%26appup_type%3D1%22%2C%22action%22%3A%22to%22%7D"
	sign := "sign=71d364d8bfb90d9d8d2e68385da86671&st=1630423869687&sv=122"
	renewUA := "JD4iPhone/167761 (iPhone; iOS 15.0; Scale/3.00)"
	//读取环境变量签名UA配置
	envBody := GetEnvFromEnvFile(ddConfig.EnvFilePath, "RENEW_BODY")
	envSign := GetEnvFromEnvFile(ddConfig.EnvFilePath, "RENEW_SIGN")
	envFullBodyAndSign := GetEnvFromEnvFile(ddConfig.EnvFilePath, "RENEW_FULL_BODY_SIGN")
	envUA := GetEnvFromEnvFile(ddConfig.EnvFilePath, "RENEW_UA")

	if envBody != "" {
		body = envBody
	}
	if envSign != "" {
		sign = envSign
	}

	fullBodySign := body + "&client=apple&clientVersion=10.0.10&openudid=1af79a528dd60b3dda24c41a3799ad095de547d3&" + sign
	//如果配置完整的请求body参数就使用完整配置
	if envFullBodyAndSign != "" {
		if envUA != "" {
			fullBodySign = envFullBodyAndSign
			renewUA = envUA
		} else {
			log.Printf("配置了完整请求参数环境变量：RENEW_FULL_BODY_SIGN，但是似乎忘记了配置需要搭配使用的环境变量：RENEW_UA，故继续使用默认签名请求配置。")
		}
	}

	LofDevLog("自定义签名使用Body==>%v", fullBodySign)
	LofDevLog("自定义签名使用UA==>%v", renewUA)
	genTokenurl := "https://api.m.jd.com/client.action?functionId=genToken"
	payload := strings.NewReader(fullBodySign)
	req, _ := http.NewRequest("POST", genTokenurl, payload)
	LofDevLog("wskey==>%v", wskey)
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "zh-Hans-FR;q=1")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", wskey)
	req.Header.Add("host", "api.m.jd.com")
	req.Header.Add("user-agent", renewUA)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer res.Body.Close()
	respBody, _ := ioutil.ReadAll(res.Body)

	var genTokenResult models.GenTokenResult
	var backErrMsg error
	LofDevLog("genToken resp==> %v", res)
	LofDevLog("genToken resp==> %v", string(respBody))
	if erri := json.Unmarshal([]byte(respBody), &genTokenResult); erri != nil {
		backErrMsg = errors.New("genToken 请求返回异常。")
	}
	if genTokenResult.Code == "0" {
		jmpUrl := fmt.Sprintf("https://un.m.jd.com/cgi-bin/app/appjmp?tokenKey=%v&to=https://plogin.m.jd.com/jd-mlogin/static/html/appjmp_blank.html", genTokenResult.TokenKey)

		//跳过证书验证
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{
			Transport: tr,
			//禁止重定向
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		jmpReq, _ := http.NewRequest("GET", jmpUrl, nil)
		jmpReq.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		jmpReq.Header.Add("user-agent", "jdapp;10.0.10;"+RandomString(20)+";network/wifi;model/iPhone13,2;addressid/0;appBuild/167761;jdSupportDarkMode/1;Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148;supportJDSHWK/1")

		jmpResp, jmpErr := client.Do(jmpReq)
		if jmpErr != nil {
			log.Print(jmpErr)
			return "", jmpErr
		}
		defer jmpResp.Body.Close()

		jmpBody, _ := ioutil.ReadAll(jmpResp.Body)
		cookies := jmpResp.Cookies()
		ptKey := ""
		ptPin := ""
		for i := 0; i < len(cookies); i++ {
			var chkCk *http.Cookie = cookies[i]
			if chkCk.Name == "pt_key" && chkCk.Value != "" {
				ptKey = chkCk.Value
				continue
			}
			if chkCk.Name == "pt_pin" && chkCk.Value != "" {
				ptPin = chkCk.Value
				continue
			}
		}
		if ptKey != "" && ptPin != "" {
			renewCK = fmt.Sprintf("pt_key=%v;pt_pin=%v;", ptKey, ptPin)
			LofDevLog("renewCK==>%v", renewCK)
		} else {
			backErrMsg = fmt.Errorf("appJmp请求未返回正确的Cookies信息。%v", string(jmpBody))
		}
	} else {
		backErrMsg = fmt.Errorf("genToken 请求返回结果不正确。%v", string(respBody))
	}
	return renewCK, backErrMsg
}

// WriteCookiesFile
func WriteCookiesFile(newCookie string, cookiesFilePath string) error {
	if CheckDirOrFileIsExist(cookiesFilePath) {
		isReplace := false
		cookiesFile, err := ioutil.ReadFile(cookiesFilePath)
		if err != nil {
			log.Printf("读取cookies.list文件出错。。%s", err)
			return errors.New(fmt.Sprintf("读取cookies.list文件出错❌\n%v", err))
		}
		lines := strings.Split(string(cookiesFile), "\n")
		for i, line := range lines {
			if strings.Contains(line, strings.Split(newCookie, ";")[1]) {
				isReplace = true
				lines[i] = newCookie
			}
		}

		var output string
		if !isReplace {
			lines = append(lines, newCookie)
		}

		lines = RemoveZero(lines)

		output = fmt.Sprintf("%v\n", strings.Join(lines, "\n"))

		err = ioutil.WriteFile(cookiesFilePath, []byte(output), 0644)
		if err != nil {
			log.Printf("写入cookies.list文件出错 %s", err)
			return errors.New(fmt.Sprintf("写入cookies.list文件出错❌\n%v", err))
		}
		return nil
	} else {
		return errors.New(fmt.Sprintf("%v文件不存在⚠️", cookiesFilePath))
	}

}
