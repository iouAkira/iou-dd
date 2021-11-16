package utils

func CloneRepoCheck() bool {
	// clone仓库校验使用场景，逻辑不公开，编译时候补充
	return true
}

func SubmitShareCodeCheck(repoDir string, info struct {
	ShareCodeName  string `json:"shareCodeName"`
	ShareCodeEnv   string `json:"shareCodeEnv"`
	ShareCodeType  string `json:"shareCodeType"`
	LogFileName    string `json:"logFileName"`
	LogPrefix      string `json:"logPrefix"`
	ScriptFileName string `json:"scriptFileName"`
}, ) bool {
	// 提交互助码前提条件交易，逻辑不公开，编译时候补充
	return true
}
