package models

type ShareCode struct {
	ShareCodeInfo []struct {
		ShareCodeName   string `json:"shareCodeName"`
		ShareCodeEnv    string `json:"shareCodeEnv"`
		ShareCodeType   string `json:"shareCodeType"`
		LogFileName     string `json:"logFileName"`
		LogPrefix       string `json:"logPrefix"`
		ScriptFileName string `json:"scriptFileName"`
	} `json:"shareCodeInfo"`
}
