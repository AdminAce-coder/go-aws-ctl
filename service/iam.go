package service

import (
	"context"

	"github.com/golifez/go-aws-ctl/model"
)

// IAM相关接口
type IamQuery interface {
	GetIamInfo(ctx context.Context) (model.IamInfo, error)                                     //查询IAM相关信息
	GetUserNameByAccessKeyId(ctx context.Context, accessKeyId string) (string, error)          //根据ACCESS_ID查询用户名
	GetAccountByAccessKeyId(ctx context.Context, accessKeyId string) (string, error)           //根据ACCESS_ID查询Account
	GetPolicyByAccessKeyId(ctx context.Context, accessKeyId string) (map[string]string, error) //通过ACCESS_KEY_ID查询该用户的策略
}
