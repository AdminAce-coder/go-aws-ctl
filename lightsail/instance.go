package lightsail

import (
	"context"
	"fmt"

	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"
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
	// 删除实例
	_, err := l.lgc.DeleteInstance(ctx, &lightsail.DeleteInstanceInput{
		InstanceName: aws.String(instanceName),
	})
	return err
}
