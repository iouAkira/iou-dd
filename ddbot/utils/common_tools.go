package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	models "ddbot/models"

	"github.com/go-cmd/cmd"
)

// ListFileName
// @description   返回目录下所有指定后缀名的文件
// @auth      iouAkira
// @param     pathSeparator string 系统目录连接符
// @param     fileDir string 需要列出文件的目录
// @return     suffix string 需要列出文件的后缀名
// @return    fileList []string 返回@fileDir下面所有后缀名为@suffix的文件
func ListFileName(pathSeparator string, fileDir string, suffix string) ([]string, []string) {
	files, _ := ioutil.ReadDir(fileDir)
	var dirList []string
	var fileList []string
	for _, oneFile := range files {
		if oneFile.IsDir() {
			dirList = append(dirList, oneFile.Name())
		} else {
			fileName := strings.Split(oneFile.Name(), ".")
			if fileName[len(fileName)-1] == suffix {
				//fmt.Println(fileDir + pathSeparator + oneFile.Name())
				fileList = append(fileList, strings.ReplaceAll(oneFile.Name(), "."+suffix, ""))
			}
		}
	}
	return dirList, fileList
}

// ExecCommand 执行系统命令
// @param botCmd指令数组
// @return 执行结果
func ExecCommand(botCmd []string, cmdType string, logsPath string) (string, bool, error) {
	var execResult string
	timeStamp := time.Now().UnixNano() / 1e6
	var cmdTimeOut time.Duration = 3600
	if cmdType == "cmd" && botCmd[0] != "ddnode" {
		cmdTimeOut = 120
	}
	log.Printf("func-ExecCommand: %v", strings.Join(botCmd, " "))
	execCmd := cmd.NewCmd("sh", "-c", strings.Join(botCmd, " "))
	statusChan := execCmd.Start()

	//限制超时主动结束执行
	go func() {
		<-time.After(cmdTimeOut * time.Second)
		execCmd.Stop()
	}()

	finalStatus := <-statusChan
	//log.Printf("Stderr %v", finalStatus.Stderr)
	//log.Printf("Stdout%v", finalStatus.Stdout)
	if len(finalStatus.Stderr) > 0 {
		execResult = "stderr:\n" + strings.Join(finalStatus.Stderr, "\n")
		//log.Printf(execResult)
	}
	if len(finalStatus.Stdout) > 0 {
		execResult = execResult + "\n" + strings.Join(finalStatus.Stdout, "\n")
		//log.Printf(execResult)
	}
	var isFile = false
	if len(finalStatus.Stdout)+len(finalStatus.Stderr) >= 50 {
		isFile = true
		logFilePath := fmt.Sprintf("%v/ddbot_%v_%v.log", logsPath, botCmd[0], timeStamp)
		logF, err := os.Create(logFilePath)
		if err != nil {
			fmt.Println(err.Error())
			logF.Close()
		} else {
			_, _ = logF.Write([]byte(execResult))
			logF.Close()
			execResult = logFilePath
		}
	}

	return execResult, isFile, finalStatus.Error
}

// ReplyKeyboardFileOpt
// @description   更新/删除全局变量里面的快捷回复配置文件
// @auth      iouAkira
// @param1     filePath string 文件路径
// @param2     rkb string 快捷回复键盘文本
// @param3     optKey string 快捷回复键盘文本对应的命令
// @return    string 返回操作结果
func ReplyKeyboardFileOpt(rkb string, optKey string, opt string) (string, error) {
	optMsg := ""
	isReplace := false
	cookiesFile, err := ioutil.ReadFile(models.GlobalEnv.ReplyKeyboardFilePath)
	if err != nil {
		log.Printf("读取%v快捷回复配置文件出错。。%v", models.GlobalEnv.ReplyKeyboardFilePath, err)
		return "", fmt.Errorf("读取%v快捷回复配置文件出错❌\n%v", models.GlobalEnv.ReplyKeyboardFilePath, err.Error())
	}
	lines := strings.Split(string(cookiesFile), "\n")
	for i, line := range lines {
		if strings.Trim(line, " ") == "" {
			continue
		}
		if !strings.Contains(line, "===") {
			continue
		}
		if strings.Split(line, "===")[0] == optKey {
			if opt == "D" {
				isReplace = true
				lines[i] = ""
				optMsg = "删除"
				delete(models.GlobalEnv.ReplyKeyBoard, optKey)
			} else if opt == "W" {
				isReplace = true
				lines[i] = rkb
				optMsg = "更新"
			} else {
				continue
			}
		}

	}
	if !isReplace && opt != "D" {
		lines = append(lines, rkb)
		optMsg = "添加"
	}
	lines = RemoveZero(lines)

	output := fmt.Sprintf("%v\n", strings.Join(RemoveZero(lines), "\n"))

	err = ioutil.WriteFile(models.GlobalEnv.ReplyKeyboardFilePath, []byte(output), 0644)
	if err != nil {
		log.Printf("写入%v快捷回复配置文件 %v", models.GlobalEnv.ReplyKeyboardFilePath, err)
		return "", fmt.Errorf("写入%v快捷回复配置文件❌\n%v", models.GlobalEnv.ReplyKeyboardFilePath, err.Error())
	}
	return optMsg, nil
}

// GetEnvFromEnvFile
// @description 获取env.sh里面的环境变量的值
// @auth      iouAkira
func GetEnvFromEnvFile(envFilePath string, envName string) string {
	env, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		log.Printf("读取配置env文件异常 %v", err)
		return ""
	}
	lines := strings.Split(string(env), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(" %v=", envName)) && strings.HasPrefix(line, "export") {
			r := regexp.MustCompile(`(` + envName + `=\"|` + envName + ` =\"|` + envName + `=\'|` + envName + ` =\')(.*)(\"|\')`)
			uaMatch := r.FindStringSubmatch(line)
			LofDevLog("%v match ==> %v", envName, uaMatch)
			if len(uaMatch) >= 3 {
				envValue := uaMatch[2]
				return envValue
			}
		}
	}
	return ""
}

// CheckDirOrFileIsExist
// @description   检查目录或者文件是否存在
// @auth      iouAkira
// @param     path string 文件夹或者文件的绝对路径
func CheckDirOrFileIsExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// IsContain
// @description 判断一个字符串是否存在于某个数字中
// @auth      iouAkira
func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// RemoveRepByLoop
// @description 通过两重循环过滤重复元素（切片长度小于1024的时候，循环来过滤）
// @auth      iouAkira
func RemoveRepByLoop(slc []string) []string {
	var result []string // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// RemoveZero
// @description 清除切片里面的零值
// @auth      iouAkira
// @param	  slice []string
func RemoveZero(slice []string) []string {
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if IfZero(v) {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveZero(slice)
		}
	}
	return slice
}

// RemoveRepByMap
// @description 通过map主键唯一的特性过滤重复元素（切片长度小大于1024的时候，通过map来过滤）
// @auth      iouAkira
func RemoveRepByMap(slc []string) []string {
	var result []string
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// IfZero
// @description 判断一个值是否为零值，只支持string,float,int,time 以及其各自的指针，"%"和"%%"也属于零值范畴，场景是like语句
// @auth      iouAkira
// @param 	  arg interface{}
func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, int32, int16, int64:
		if v == 0 {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	default:
		return false
	}
	return false
}

// RandomString
// @description 获取指定长度的随机字符串
// @auth      iouAkira
// @param 	  n int
func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}
	log.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// LofDevLog 打印开发调试日志
func LofDevLog(format string, v ...interface{}) {
	printLog := false
	envPrintLog := LofDevLogGetEnvFromEnvFile("/data/env.sh", "LOF_DEV_LOG")
	// envPrintLog = "true"
	if envPrintLog == "true" {
		printLog = true
	}
	if printLog {
		fmt.Printf(format, v...)
		fmt.Printf("\n")
	}
}
func LofDevLogGetEnvFromEnvFile(envFilePath string, envName string) string {
	env, err := ioutil.ReadFile(envFilePath)
	if err != nil {
		//log.Printf("读取开发模式日志打印配置参数异常")
		return ""
	}
	lines := strings.Split(string(env), "\n")
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf(" %v=", envName)) && strings.HasPrefix(line, "export") {
			r := regexp.MustCompile(`(` + envName + `=\"|` + envName + ` =\"|` + envName + `=\'|` + envName + ` =\')(.*)(\"|\')`)
			uaMatch := r.FindStringSubmatch(line)
			//log.Printf("%v match ==> %v", envName, uaMatch)
			if len(uaMatch) >= 3 {
				envValue := uaMatch[2]
				return envValue
			}
		}
	}
	return ""

}
