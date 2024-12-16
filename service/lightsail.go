package service

import (
	"context"

	lgtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	ctltypes "github.com/golifez/go-aws-ctl/model"
)

// Lightsail 查询接口
type LgQuerysvc interface {
	// 获取捆绑包列表
	GetBundlesInput(ctx context.Context) (LgBundleList []ctltypes.LgBundle, err error)
	// 获取实例数量
	GetInstancesInput(ctx context.Context)
	// 获取区域实例列表
	GetInstanceListWithRegion(ctx context.Context, region string) (instanceNameList []ctltypes.LgAttr, err error)
	// 获取所有区域实例列表
	GetInstanceList(ctx context.Context) (instanceNameList []ctltypes.LgAttr, err error)
	// 获取区域列表
	GetRegionList(ctx context.Context) (regionList []lgtypes.Region)
	// 获取快照列表
	GetSnapshotList(ctx context.Context) (snapshotList []ctltypes.LgSnapshot, err error)
	// 通过区域获取快照
	GetSnapshotListWithRegion(ctx context.Context, region string) (snapshotList []ctltypes.LgSnapshot, err error)
	// 查询实例防火墙端口
	QueryInstanceFirewallPort(ctx context.Context, instanceName string, region string) error
	// 查询实例状态
	QueryInstanceStatus(ctx context.Context, instanceName string, region string) (string, error)
}

// Lightsail 操作接口
type LgOpsvc interface {
	// 删除全部实例
	DeleteAllLg(ctx context.Context) error
	// 删除实例
	DeleteInstance(instanceName string, region string) error
	// 创建实例
	CreateInstance(lgCreateInstance ctltypes.LgCreateInstance) ([]string, error)
	// 停止实例
	StopInstance(instanceName string, region string) error
	// 启动实例
	StartInstance(instanceName string, region string) error
	// 修改实例标签
	ModifyInstanceTag(instanceName string, region string, tagKey string) error
	// 切换实例公网IP
	ChangeInstancePublicIp(instanceName string, region string) error
	// 打开实例端口
	OpenInstancePort(instanceName string, region string, ports string) error
}
