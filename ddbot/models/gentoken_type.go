package models

type GenTokenResult struct {
	Code     string `json:"code"`
	TokenKey string `json:"tokenKey"`
	URL      string `json:"url"`
}
