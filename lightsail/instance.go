package lightsail

import (
	"fmt"
	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/spf13/cobra"
)

// 操作lightsail实例
// createCmd represents the create command
var InstanceCmd = &cobra.Command{
	Use:   "instanceOP",
	Short: "A brief description of your command",
	Long:  `操作lightsail实例`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取区域
		region, _ := cmd.Flags().GetString("region")
		//获取客户端

		instanceNames, _ := cmd.Flags().GetStringSlice("instanceNames")
		delete, _ := cmd.Flags().GetBool("delete")
		if delete {
			DeleteLg(instanceNames, region)
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

// 删除实例
func DeleteLg(instanceNames []string, region string) {
	if len(instanceNames) == 1 && instanceNames[0] == "all" {

		// 获取区域列表
		regionList := GetRegionList(ctx, region)

		// 获取实例
		for _, region := range regionList {
			instanceList := GetInstanceList(ctx, string(region.Name))
			// 删除实例
			for _, instancename := range instanceList {
				cl := cmd2.GetClient[*lightsail.Client](
					cmd2.WithRegion(string(region.Name)),
					cmd2.WithClientType("lightsail"),
				)
				cl.DeleteInstance(ctx, &lightsail.DeleteInstanceInput{
					InstanceName: &instancename,
				})
				fmt.Printf("正在删除实例: %s区域: %s\n", instancename, region.Name)
			}
		}
		fmt.Println("删除所有实例完成")
	}
}
