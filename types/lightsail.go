package lightsail

import (
	"time"

	lgtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

// Lightsail 实例的属性
type LgAttr struct {
	Region     string        // 区域
	Name       string        // 实例名称
	Size       string        // 实例类型
	Image      string        // 镜像
	KeyName    string        // 密钥对名称
	UserData   string        // 用户数据
	Status     string        // 状态
	CreateTime time.Time     // 创建时间
	Tags       []lgtypes.Tag // 标签
}

// 实例捆绑包
type LgBundle struct {
	BundleId     string  // 捆绑包ID
	CpuCount     int32   // CPU数量
	RamSizeInGb  float32 // 内存大小
	DiskSizeInGb int32   // 磁盘大小
}
