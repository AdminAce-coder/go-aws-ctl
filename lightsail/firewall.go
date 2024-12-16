package lightsail

import (
	"context"
	"fmt"

	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/spf13/cobra"
)

// FirewalldCmd represents the create command
// # 打开所有端口
// go run main.go lg  Fw --region ap-northeast-2  --instanceNames all --ports all
//
// # 打开具体实例的 80 和 443 端口
// go run main.go lg  Fw --regionap-northeast-2 --instanceNames instance1 --ports 80,443
//
// # 打开多个实例的 80-100 端口
// go run main.go lg  Fw --region ap-northeast-2   --instanceNames instance1,instance2 --ports 80-100
var FirewalldCmd = &cobra.Command{
	Use:   "Fw",
	Short: "Enable firewall for Lightsail instances",
	Long:  `Enable specific ports or a range of ports for AWS Lightsail instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取命令行参数
		region, _ := cmd.Flags().GetString("region")
		instanceNames, _ := cmd.Flags().GetStringSlice("instanceNames")
		ports, _ := cmd.Flags().GetStringSlice("ports")

		// 获取 Lightsail 客户端
		cl := cmd2.GetClient[*lightsail.Client](cmd2.WithRegion(region), cmd2.WithClientType("lightsail"))

		// 开启防火墙端口
		if err := openPorts(cl, instanceNames, ports); err != nil {
			fmt.Printf("Failed to open ports: %v\n", err)
		} else {
			fmt.Println("Ports opened successfully.")
		}
	},
}

var ctx = context.Background()

func init() {
	cmd.LgCmd.AddCommand(FirewalldCmd)
	FirewalldCmd.Flags().StringP("region", "r", "", "Region (e.g., us-east-1)")
	FirewalldCmd.Flags().StringSlice("instanceNames", []string{"all"}, "Instance names (e.g., name1,name2)")
	FirewalldCmd.Flags().StringSlice("ports", []string{"all"}, "Ports to open (e.g., 80,443 or 0-65535)")
}

// openPorts opens the specified ports for Lightsail instances.
func openPorts(client *lightsail.Client, instanceNames, ports []string) error {
	var instances []string

	// 处理实例名称
	if len(instanceNames) == 1 && instanceNames[0] == "all" {
		output, err := client.GetInstances(ctx, &lightsail.GetInstancesInput{})
		if err != nil {
			return fmt.Errorf("failed to get instances: %w", err)
		}
		for _, instance := range output.Instances {
			instances = append(instances, aws.ToString(instance.Name))
		}
	} else {
		instances = instanceNames
	}

	// 遍历每个实例并打开指定的端口
	for _, instance := range instances {
		for _, port := range ports {
			if port == "all" {
				// 打开所有端口 (0-65535)
				_, err := client.OpenInstancePublicPorts(ctx, &lightsail.OpenInstancePublicPortsInput{
					InstanceName: aws.String(instance),
					PortInfo: &types.PortInfo{
						FromPort: 0,
						ToPort:   65535,
						Protocol: types.NetworkProtocolAll,
					},
				})
				if err != nil {
					return fmt.Errorf("failed to open all ports for instance %s: %w", instance, err)
				}
			} else {
				// 打开具体的端口
				portRange, err := parsePortRange(port)
				if err != nil {
					return fmt.Errorf("invalid port %s for instance %s: %w", port, instance, err)
				}

				_, err = client.OpenInstancePublicPorts(ctx, &lightsail.OpenInstancePublicPortsInput{
					InstanceName: aws.String(instance),
					PortInfo: &types.PortInfo{
						FromPort: portRange[0],
						ToPort:   portRange[1],
						Protocol: types.NetworkProtocolTcp,
					},
				})
				if err != nil {
					return fmt.Errorf("failed to open port %s for instance %s: %w", port, instance, err)
				}
			}
		}
	}
	return nil
}

// 解析端口范围
func parsePortRange(port string) ([2]int32, error) {
	var portRange [2]int32
	_, err := fmt.Sscanf(port, "%d-%d", &portRange[0], &portRange[1])
	if err == nil {
		return portRange, nil
	}
	return [2]int32{}, fmt.Errorf("invalid port format: %s", port)
}

// 打开防火墙端口
func OpenFirewallPort(client *lightsail.Client, instanceName string, portRange []int32) error {
	// 打开端口
	_, err := client.OpenInstancePublicPorts(ctx, &lightsail.OpenInstancePublicPortsInput{
		InstanceName: aws.String(instanceName),
		PortInfo: &types.PortInfo{
			FromPort: portRange[0], // 使用数组的第一个值
			ToPort:   portRange[1], // 使用数组的第二个值
			Protocol: types.NetworkProtocolTcp,
		},
	})
	return err
}

// 查询实例防火墙端口
func QueryInstanceFirewallPort(client *lightsail.Client, instanceName string) error {
	// 查询端口
	out, err := client.GetInstancePortStates(ctx, &lightsail.GetInstancePortStatesInput{
		InstanceName: aws.String(instanceName),
	})
	fmt.Println(out.PortStates)
	return err
}
