/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package lightsail

import (
	"context"
	"fmt"
	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/duke-git/lancet/v2/datetime"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
// go run main.go lg create -r ap-northeast-2  2  -s Ubuntu-1-1733389433 -i nano_3_0
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long:  `创建lightsail实例`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取区域
		region, _ := cmd.Flags().GetString("region")
		num, _ := cmd.Flags().GetInt("number")
		Snapshotname, _ := cmd.Flags().GetString("SnapshotName")
		instanceType, _ := cmd.Flags().GetString("instanceType")
		//获取客户端
		cl := cmd2.GetClient[*lightsail.Client](cmd2.WithRegion(region), cmd2.WithClientType("lightsail"))
		CreateInstancesFromSnapshot(cl, num, &Snapshotname, &instanceType)
	},
}

func init() {
	cmd.LgCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("region", "r", "", "区域如:us-east-1")
	createCmd.Flags().IntP("number", "n", 1, "创建实例数量")
	createCmd.Flags().StringP("SnapshotName", "s", "", "指定快照名启动")
	createCmd.Flags().StringP("instanceType", "i", "nano_3_0", "指定实例类型")
}

func CreateInstancesFromSnapshot(cl *lightsail.Client, num int, snapshotName, instanceType *string) {
	// 使用 WaitGroup 管理 Goroutines
	var wg sync.WaitGroup
	ctx := context.TODO()

	// 预先生成所有实例名称
	instanceNames := make([]string, num)
	for i := 0; i < num; i++ {
		instanceNames[i] = InstanceName() // 确保每个实例有唯一名称
	}

	// 遍历每个实例名称并创建实例
	for _, name := range instanceNames {
		wg.Add(1) // 增加 WaitGroup 计数器

		go func(instanceName string) {
			defer wg.Done() // Goroutine 完成时减少计数器

			log.Printf("Creating instance with name: %s\n", instanceName) // 添加日志输出

			_, err := cl.CreateInstancesFromSnapshot(ctx, &lightsail.CreateInstancesFromSnapshotInput{
				InstanceSnapshotName: snapshotName,
				InstanceNames:        []string{instanceName},        // 实例名称
				AvailabilityZone:     aws.String("ap-northeast-2a"), // 可用区
				BundleId:             instanceType,                  // 套餐 ID
			})
			if err != nil {
				log.Printf("Failed to create instance %s: %v\n", instanceName, err)
				return
			}
			log.Printf("Instance %s created successfully.\n", instanceName)
		}(name) // 将实例名称作为参数传递，确保每个 Goroutine 使用独立的值
	}

	// 等待所有 Goroutines 完成
	wg.Wait()
	log.Println("All instances have been created.")
}
func InstanceName() (name string) {
	// 获取当前时间作为1 格式为 2022-01-28 15:59:33
	instanceName := datetime.GetNowDateTime()
	fmt.Println(instanceName)
	s1 := strings.Split(instanceName, " ") //提取2024-11-18 15:46:41

	s2 := strings.Split(s1[1], ":") // 提取[15 46 41]

	s3 := strings.Join(s2, "-") // 15-46-41

	s1 = append(s1[:1], s3)
	fliest := strings.Join(s1, "-")
	fmt.Println(fliest)
	time.Sleep(1 * time.Second) // 间隔一米

	return fliest
}

func ListInstancetypeLg() {

}
