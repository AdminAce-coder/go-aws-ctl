/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package lightsail

import (
	cmd2 "github.com/golifez/go-aws-ctl/cmd"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long:  `启动lightsail.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取区域
		region, _ := cmd.Flags().GetString("region")
		group, _ := cmd.Flags().GetString("gourp")
		//获取客户端
		cl := cmd2.GetClient[*lightsail.Client](cmd2.WithRegion(region), cmd2.WithClientType("lightsail"))
		StartLightsail(cl, group)
	},
}

func init() {
	cmd2.LgCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("region", "r", "", "区域如:us-east-1")
	startCmd.Flags().StringP("gourp", "g", "", "请输入实例的分组如:gourpA")
}

func StartLightsail(cl *lightsail.Client, group string) {
	// 获取上下文
	ctx := gctx.New()

	// 获取实例列表
	resp, err := cl.GetInstances(ctx, &lightsail.GetInstancesInput{})
	if err != nil {
		g.Log().Error(ctx, "无法获取实例信息: %v", err)
		return
	}

	// 遍历实例
	for _, instance := range resp.Instances {
		// 实例名拆分
		instanceName := aws.ToString(instance.Name)
		nameParts := strings.Split(instanceName, "-")

		// 确保分割后有足够的部分
		if len(nameParts) < 2 {
			g.Log().Warningf(ctx, "实例名格式不正确: %s", instanceName)
			continue
		}

		// 检查实例组是否匹配
		if group == nameParts[0] {
			// 启动实例
			g.Log().Infof(ctx, "正在启动: %s", instanceName)
			_, err := cl.StartInstance(ctx, &lightsail.StartInstanceInput{
				InstanceName: instance.Name,
			})
			if err != nil {
				g.Log().Errorf(ctx, "启动实例 %s 失败: %v", instanceName, err)
			} else {
				g.Log().Infof(ctx, "成功启动实例: %s", instanceName)
			}
		}
	}
}
