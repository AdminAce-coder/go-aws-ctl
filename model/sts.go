package model

type IamInfo struct {
	Account  string            //账户
	UserName string            //用户名
	Policy   map[string]string //策略名和描述
}
