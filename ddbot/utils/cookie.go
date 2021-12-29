package utils

import (
	"ddbot/models"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

//CookieCfg 操作类参数
type CookieCfg struct {
	*models.DDEnv
}

func (cc CookieCfg) ReadCookies(isWsKey bool, args []string) ([]string, error) {
	var filePath = cc.CookiesListFilePath
	if isWsKey {
		filePath = cc.CookiesWSKeyListFilePath
	}
	if len(args) <= 0 {
		return ReadCookiesFrom(filePath)
	}
	_, err := strconv.Atoi(args[0])
	var ck Iwskey
	if err != nil {
		ck = GetCookieByID(args[0], filePath)
		cookieLine, err := ck.GetCookieContextByID()
		if err != nil {
			return nil, err
		}
		return []string{cookieLine}, err

	} else {
		ck = GetCookieByIndex(args[0], filePath)
		cookieLine, err := ck.GetCookieContextByIndex()
		if err != nil {
			return nil, err
		}
		return []string{cookieLine}, err
	}
}

func GetPinFromCookieText(cookie string) string  {
	var pin string
	ks := strings.Split(cookie, ";")
	for _, k := range ks {
		if strings.Contains(k, "pt_pin") || strings.Contains(k, "pin") {
			pin = strings.Split(k, "=")[1]
		}
	}
	return pin
}

//Read cookies as list from cookies.list file
func ReadCookiesFrom(cookiesFilePath string) ([]string, error) {
	cookiesFile, err := ioutil.ReadFile(cookiesFilePath)
	if err != nil {
		log.Printf("读取"+cookiesFilePath+"文件出错。。%s", err)
		return nil, fmt.Errorf("读取"+cookiesFilePath+"文件出错❌\n%v", err)
	}
	lines := strings.Split(string(cookiesFile), "\n")

	if len(lines) <= 0 {
		return nil, errors.New("当前cookies文件未空")
	}
	var cookies []string
	for _, line := range lines {
		//pick the value of pt_pin from the line
		if !strings.Contains(line, ";") {
			continue
		}
		ks := strings.Split(line, ";")
		for _, k := range ks {
			if strings.Contains(k, "pt_pin") || strings.Contains(k, "pin") {
				cookies = append(cookies, strings.Split(k, "=")[1])
			}
		}
	}
	if len(cookies) == 0 {
		return nil, errors.New("cookies字段不完整，请检查")
	}
	return cookies, nil
}

type Iwskey interface {
	GetCookieLine() (string, error)
	GetCookieContextByIndex() (string, error)
	GetCookieContextByID() (string, error)
}

type PinFile struct {
	FilePath string
	Pin      string
	Index    int
}

func (p PinFile) GetCookieContextByIndex() (string, error) {
	cookiesFile, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		log.Printf("读取cookies文件出错。。%s", err)
		return "", fmt.Errorf("读取cookies文件出错❌\n%v", err)
	}
	lines := strings.Split(string(cookiesFile), "\n")
	line := lines[p.Index-1]
	line = strings.ReplaceAll(line, "\r", "")
	if line == "" {
		return "", fmt.Errorf("不存在的内容,请添加后继续。")
	}
	return line, nil
}
func (p PinFile) GetCookieContextByID() (string, error) {
	cookiesFile, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		log.Printf("读取cookies文件出错。。%s", err)
		return "", fmt.Errorf("读取cookies文件出错❌\n%v", err)
	}
	cookie := ""
	lines := strings.Split(string(cookiesFile), "\n")
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, p.Pin) {
			cookie = line
			break
		}
	}
	if cookie == "" {
		return "", fmt.Errorf("不存在WSKEY,请添加后继续。")
	}
	return cookie, nil
}

func (p PinFile) GetCookieLine() (string, error) {
	cookiesFile, err := ioutil.ReadFile(p.FilePath)
	if err != nil {
		log.Printf("读取cookies文件出错。。%s", err)
		return "", fmt.Errorf("读取cookies文件出错❌\n%v", err)
	}
	cookieIndex := ""
	lines := strings.Split(string(cookiesFile), "\n")
	for i, line := range lines {
		line = strings.ReplaceAll(line, "\r", "")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, p.Pin) {
			cookieIndex = fmt.Sprint(i + 1)
			break
		}
	}
	if cookieIndex == "" {
		return "", fmt.Errorf("不存在WSKEY,请添加后继续。")
	}
	return cookieIndex, nil
}

func GetCookieByID(id, cookieFile string) Iwskey {
	var _ Iwskey = new(PinFile)
	wskey := PinFile{FilePath: cookieFile, Pin: id}
	return wskey
}

func GetCookieByIndex(index, cookieFile string) Iwskey {
	atoi, err := strconv.Atoi(index)
	if err != nil {
		return nil
	}
	var _ Iwskey = new(PinFile)
	wskey := PinFile{FilePath: cookieFile, Index: atoi}
	return wskey
}
