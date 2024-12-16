package lightsail

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/random"
	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"
	ctltypes "github.com/golifez/go-aws-ctl/model"
	"github.com/golifez/go-aws-ctl/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/spf13/cobra"
)

type LgInstanceOpCommand struct {
	lgc *lightsail.Client
}

func NewLgInstanceOpCommand() service.LgOpsvc {
	return &LgInstanceOpCommand{
		lgc: cmd2.GetDefaultAwsLgClient(),
	}
}

// 操作lightsail实例
// createCmd represents the create command
var InstanceCmd = &cobra.Command{
	Use:   "instanceOP",
	Short: "A brief description of your command",
	Long:  `操作lightsail实例`,
	Run: func(cmd *cobra.Command, args []string) {
		delete, _ := cmd.Flags().GetBool("delete")
		// 获取区域
		region, _ := cmd.Flags().GetString("region")
		//获取客户端
		lgo := NewLgInstanceOpCommand()
		instanceNames, _ := cmd.Flags().GetStringSlice("instanceNames")

		if delete {
			for _, instanceName := range instanceNames {
				err := lgo.DeleteInstance(instanceName, region)
				if err != nil {
					fmt.Println("删除实例失败:", err)
					return
				}
			}
		}
	},
}

func init() {
	cmd.LgCmd.AddCommand(InstanceCmd)
	InstanceCmd.Flags().StringP("region", "r", "", "区域如:us-east-1")
	InstanceCmd.Flags().StringSlice("instanceNames", []string{"all"}, "实例名列表（例如: name1,name2,name3）")
	// 是否删除
	InstanceCmd.Flags().BoolP("delete", "d", false, "是否删除实例")
}

// 删除所有实例
func (l *LgInstanceOpCommand) DeleteAllLg(ctx context.Context) error {
	// // 先获取所有的区域
	lgq := NewLgQuery()
	// 获取所有的实例
	instanceList, err := lgq.GetInstanceList(ctx)
	if err != nil {
		return err
	}
	// 删除所有的实例
	for _, instance := range instanceList {
		l.DeleteInstance(instance.InstanceName, instance.Region)
	}
	return nil
}

// 删除实例
func (l *LgInstanceOpCommand) DeleteInstance(instanceName string, region string) error {
	// 获取带区域的客户端
	lgcwithRegion := cmd2.GetAwsLgClient(region)
	// 删除实例
	_, err := lgcwithRegion.DeleteInstance(ctx, &lightsail.DeleteInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	return err
}

// 通过快照创建实例
func (l *LgInstanceOpCommand) CreateInstance(lgCreateInstance ctltypes.LgCreateInstance) error {
	// 获取带区域的客户端
	lgcwithRegion := cmd2.GetAwsLgClient(lgCreateInstance.Region)
	// 创建实例
	for i := 0; i < lgCreateInstance.Num; i++ {
		instanceName := l.AutoInstanceName(lgCreateInstance)
		_, err := lgcwithRegion.CreateInstancesFromSnapshot(ctx, &lightsail.CreateInstancesFromSnapshotInput{
			InstanceNames:        []string{instanceName},
			AvailabilityZone:     aws.String(lgCreateInstance.AvailabilityZone),
			BundleId:             aws.String(lgCreateInstance.BundleId),
			InstanceSnapshotName: aws.String(lgCreateInstance.SnapshotName),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 自动生成实例名称
func (l *LgInstanceOpCommand) AutoInstanceName(lgCreateInstance ctltypes.LgCreateInstance) string {
	// 获取当天的日期
	currentDate := datetime.GetNowDate()
	// 生成一个随机字符串
	randomString := random.RandNumeralOrLetter(6)
	// 返回实例名称
	if lgCreateInstance.IsAutoName {
		return fmt.Sprintf("%s-%s", currentDate, randomString)
	}
	return lgCreateInstance.InstanceName
}

// 停止实例
func (l *LgInstanceOpCommand) StopInstance(instanceName string, region string) error {
	// 获取带区域的客户端
	lgcwithRegion := cmd2.GetAwsLgClient(region)
	// 停止实例
	_, err := lgcwithRegion.StopInstance(ctx, &lightsail.StopInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	return err
}

// 启动实例
func (l *LgInstanceOpCommand) StartInstance(instanceName string, region string) error {
	// 获取带区域的客户端
	lgcwithRegion := cmd2.GetAwsLgClient(region)
	// 启动实例
	_, err := lgcwithRegion.StartInstance(ctx, &lightsail.StartInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	return err
}

// 切换lightsail实例公网IP
func (l *LgInstanceOpCommand) ChangeInstancePublicIp(instanceName string, region string) error {
	lgcwithRegion := cmd2.GetAwsLgClient(region)

	// 先分离当前的静态IP(如果有的话)
	// _, err := lgcwithRegion.DetachStaticIp(ctx, &lightsail.DetachStaticIpInput{
	// 	StaticIpName: aws.String(instanceName),
	// })
	// if err != nil {
	// 	return err
	// }

	// 停止实例
	_, err := lgcwithRegion.StopInstance(ctx, &lightsail.StopInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	if err != nil {
		return err
	}

	// 获取实例状态
	for {
		instance, err := lgcwithRegion.GetInstance(ctx, &lightsail.GetInstanceInput{
			InstanceName: aws.String(instanceName),
		})
		if err != nil {
			return err
		}
		fmt.Println(instance.Instance.State)
		if *instance.Instance.State.Name == "stopped" {
			break
		}
		time.Sleep(10 * time.Second)
	}
	// 启动实例 - 启动时会自动分配新的公网IP
	_, err = lgcwithRegion.StartInstance(ctx, &lightsail.StartInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	if err != nil {
		return err
	}
	return nil
}

// 修改实例标签，只添加Key
func (l *LgInstanceOpCommand) ModifyInstanceTag(instanceName string, region string, tagKey string) error {
	lgcwithRegion := cmd2.GetAwsLgClient(region)
	// 修改实例标签
	_, err := lgcwithRegion.TagResource(ctx, &lightsail.TagResourceInput{
		ResourceName: aws.String(instanceName),
		Tags:         []types.Tag{{Key: aws.String(tagKey)}},
	})
	return err
}

// 打开实例端口
func (l *LgInstanceOpCommand) OpenInstancePort(instanceName string, region string, portRange []int32) error {
	lgcwithRegion := cmd2.GetAwsLgClient(region)
	// 打开实例端口
	err := OpenFirewallPort(lgcwithRegion, instanceName, portRange)
	if err != nil {
		return err
	}
	return nil
}
