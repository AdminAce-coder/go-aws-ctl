package service

import (
	"context"
	"github.com/golifez/go-aws-ctl/model"
)

// IAM相关接口
type IamQuery interface {
	GetIamInfo(ctx context.Context) (model.IamInfo, error) //查询IAM相关信息
}
