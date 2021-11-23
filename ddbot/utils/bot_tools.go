package utils

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	models "ddbot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DELETE string = "delete"
	RETURN string = "return"
)

// MakeKeyboardMarkup
// @description   生成回复的InlineKeyboard
// @auth      iouAkira
// @param     ikType string 生成InlineKeyboard的类型
// @param     rowBtns int 每行button的数量
// @param     filePath 生成button按钮内容对应文件的所属目录
// @param     suffix 生成button按钮操作的文件后缀名
// @return    keyboardMarkup tgbotapi.InlineKeyboardMarkup 返回构建好的InlineKeyboardMarkup
func MakeKeyboardMarkup(ikTypePrefix string, rowBtns int, filePath string, suffix string) tgbotapi.InlineKeyboardMarkup {
	var keyboardMarkup tgbotapi.InlineKeyboardMarkup
	ikType := ikTypePrefix[1:]
	switch ikType {
	case "rdc":
		log.Printf("生成rdc类型的InlineKeyboard")
	case "ddnode":
		dirList, fileList := ListFileName(string(os.PathSeparator), filePath, "js")

		var btns []*string
		excludeDir := []string{"node_modules", ".git", "utils", "backUp", "logs", "archive", "docker", "icon"}
		for di, dv := range dirList {
			if IsContain(excludeDir, dv) {
				continue
			}
			dir := dv
			btns = append(btns, &dir)
			if len(btns) == rowBtns || di == len(dirList)-1 {
				var row []tgbotapi.InlineKeyboardButton
				for _, dk := range btns {
					row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("🗂%v", *dk), fmt.Sprintf("%v %v %v/%v", ikTypePrefix, "dir", filePath, *dk)))
				}
				keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, row)
				btns = btns[0:0]
			}
		}
		for i, v := range fileList {
			if len(strings.Split(v, "_")) < 2 {
				continue
			}
			if strings.HasPrefix(v, "JD") || strings.HasPrefix(v, "USER") || strings.HasPrefix(v, "JS") {
				continue
			}
			file := v
			btns = append(btns, &file)
			if len(btns) == rowBtns || i == len(fileList)-1 {
				var row []tgbotapi.InlineKeyboardButton
				for _, k := range btns {
					row = append(row, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("📜%v", strings.Split(*k, "_")[1:]), fmt.Sprintf("%v %v/%v", ikTypePrefix, filePath, *k)))
				}
				keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, row)
				btns = btns[0:0]
			}
		}
	case "logs":
		_, fileList := ListFileName(string(os.PathSeparator), filePath, "log")
		var btns []*string
		for i, v := range fileList {
			if len(strings.Split(v, "_")) < 2 {
				continue
			}
			file := v
			btns = append(btns, &file)
			if len(btns) == rowBtns || i == len(fileList)-1 {
				var row []tgbotapi.InlineKeyboardButton
				for _, k := range btns {
					row = append(row, tgbotapi.NewInlineKeyboardButtonData(*k, fmt.Sprintf("%v %v/%v.log", ikTypePrefix, filePath, *k)))
				}
				keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, row)
				btns = btns[0:0]
			}
		}
	}

	keyboardMarkup.InlineKeyboard = append(keyboardMarkup.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("取消", "/cancel"),
	))

	return keyboardMarkup
}

// MakeReplyKeyboard 构建快捷回复按钮
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

// LoadReplyKeyboardMap 更新快捷回复按钮全局配置
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

// RenewCookie 根据传入的wskey更新对应cookie
func RenewCookie(wskey string) (string, error) {
	renewCK := ""
	//默认签名UA配置
	body := "body=%7B%22to%22%3A%22https%3A%5C%2F%5C%2Fplogin.m.jd.com%5C%2Fcgi-bin%5C%2Fm%5C%2Fthirdapp_auth_page%3Ftoken%3DAAEAIEbEUWDGA_SGHg4sHM5fwfnpt-kFtkZ_boToZQULiH0O%26client_type%3Dapple%26appid%3D1125%26appup_type%3D1%22%2C%22action%22%3A%22to%22%7D"
	sign := "sign=71d364d8bfb90d9d8d2e68385da86671&st=1630423869687&sv=122"
	renewUA := "JD4iPhone/167761 (iPhone; iOS 15.0; Scale/3.00)"
	//读取环境变量签名UA配置
	envBody := GetEnvFromEnvFile(models.GlobalEnv.EnvFilePath, "RENEW_BODY")
	envSign := GetEnvFromEnvFile(models.GlobalEnv.EnvFilePath, "RENEW_SIGN")
	envFullBodyAndSign := GetEnvFromEnvFile(models.GlobalEnv.EnvFilePath, "RENEW_FULL_BODY_SIGN")
	envUA := GetEnvFromEnvFile(models.GlobalEnv.EnvFilePath, "RENEW_UA")

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

// writeCookiesFile 将传入的 cookie 写入文件
func writeCookiesFile(newCookie string) error {
	if CheckDirOrFileIsExist(models.GlobalEnv.CookiesListFilePath) {
		isReplace := false
		cookiesFile, err := ioutil.ReadFile(models.GlobalEnv.CookiesListFilePath)
		if err != nil {
			log.Printf("读取cookies.list文件出错。。%s", err)
			return fmt.Errorf("读取cookies.list文件出错❌\n%v", err)
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

		err = ioutil.WriteFile(models.GlobalEnv.CookiesListFilePath, []byte(output), 0644)
		if err != nil {
			log.Printf("写入cookies.list文件出错 %s", err)
			return fmt.Errorf("写入cookies.list文件出错❌\n%v", err)
		}
		return nil
	} else {
		return fmt.Errorf("%v文件不存在⚠️", models.GlobalEnv.CookiesListFilePath)
	}

}

func CleanCommand(cmd string, offset int) []string {
	cmdMsgSplit := strings.Split(cmd[offset:], " ")
	var arr []string
	for _, v := range cmdMsgSplit {
		if v == "" {
			continue
		}
		arr = append(arr, v)
	}
	return arr
}
