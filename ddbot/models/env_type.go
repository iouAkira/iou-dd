package models

//DDEnv ddBot环境变量集合
type DDEnv struct {
	RepoBaseDir              string
	DataBaseDir              string
	SpnodeBtnFilePath        string
	LogsBtnFilePath          string
	CookiesListFilePath      string
	CookiesWSKeyListFilePath string
	EnvFilePath              string
	ReplyKeyboardFilePath    string
	CustomFilePath           string
	TgBotToken               string
	TgUserID                 int64
	ReplyKeyBoard            map[string]string
}
