package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lgtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	ctltypes "github.com/golifez/go-aws-ctl/model"
)

// Lightsail 查询接口
type LgQuerysvc interface {
	// 获取实例类型
	GetBundlesInput(ctx context.Context, lgc *lightsail.Client) ([]ctltypes.LgBundle, error)
	// 获取实例数量
	GetInstancesInput(ctx context.Context)
	// 获取区域实例列表
	GetInstanceListWithRegion(ctx context.Context, region string) (instanceNameList []ctltypes.LgAttr, err error)
	// 获取所有区域实例列表
	GetInstanceList(ctx context.Context) (instanceNameList []ctltypes.LgAttr, err error)
	// 获取区域列表
	GetRegionList(ctx context.Context, region string) (regionList []lgtypes.Region)
}
