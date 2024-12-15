package lightsail

import (
	"context"
	"fmt"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/duke-git/lancet/v2/random"
	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"
	ctltypes "github.com/golifez/go-aws-ctl/model"
	"github.com/golifez/go-aws-ctl/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
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
