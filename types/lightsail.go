package lightsail

import (
	lgtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

// Lightsail 实例的属性
type LgAttr struct {
	Region       string        // 区域
	Zone         string        // 可用区
	PublicIp     string        // 公网IP
	Status       string        // 状态
	Tags         []lgtypes.Tag // 标签
	CreateTime   string        // 创建时间
	InstanceName string        // 实例名称
	InstanceType string        // 实例类型
}

// 实例捆绑包
type LgBundle struct {
	BundleId     string  // 捆绑包ID
	CpuCount     int32   // CPU数量
	RamSizeInGb  float32 // 内存大小
	DiskSizeInGb int32   // 磁盘大小
}
