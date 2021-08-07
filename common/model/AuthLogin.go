package model

type AuthLogin struct {
	AccountName string `json:"account_name"`
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
}
