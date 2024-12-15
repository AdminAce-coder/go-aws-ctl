/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package lightsail

import (
	"context"
	"fmt"

	"github.com/duke-git/lancet/v2/datetime"

	"github.com/golifez/go-aws-ctl/cmd"
	cmd2 "github.com/golifez/go-aws-ctl/cmd"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lgtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gogf/gf/v2/os/glog"
	ctltypes "github.com/golifez/go-aws-ctl/model"
	"github.com/golifez/go-aws-ctl/service"
	"github.com/spf13/cobra"
)

type LgQuery struct {
	nextPageToken *string
	lgc           *lightsail.Client
}

func NewLgQuery() service.LgQuerysvc {
	return &LgQuery{
		lgc: cmd2.GetDefaultAwsLgClient(),
	}
}

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long:  `查询相关.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		lc := lgClinet(ctx)
		instanceType, err := cmd.Flags().GetBool("instanceType")
		if err != nil {
			fmt.Println("获取 instanceType 标志失败:", err)
			return
		}
		instanceList, err := cmd.Flags().GetBool("instanceList")
		if err != nil {
			fmt.Println("获取 instanceList 标志失败:", err)
			return
		}
		lg := NewLgQuery()
		if instanceType {
			lg.GetBundlesInput(ctx, lc)
		}
		if instanceList {
			lg.GetInstancesInput(ctx)
		}
	},
}

// var LgBundleList []LgBundle

func init() {
	cmd.LgCmd.AddCommand(queryCmd)
	queryCmd.Flags().BoolP("instanceType", "i", false, "查看可启动实例的类型")
	queryCmd.Flags().BoolP("instanceList", "l", false, "查询区域的实例数量")
}

// 获取捆绑包
func (lg *LgQuery) GetBundlesInput(ctx context.Context, lgc *lightsail.Client) (LgBundleList []ctltypes.LgBundle, err error) {
	output, err := lgc.GetBundles(ctx, &lightsail.GetBundlesInput{})
	if err != nil {
		glog.New().Error(ctx, err)
		return nil, err
	}
	for _, bundle := range output.Bundles {
		LgBundleList = append(LgBundleList, ctltypes.LgBundle{
			BundleId:     aws.ToString(bundle.BundleId),
			CpuCount:     aws.ToInt32(bundle.CpuCount),
			RamSizeInGb:  aws.ToFloat32(bundle.RamSizeInGb),
			DiskSizeInGb: aws.ToInt32(bundle.DiskSizeInGb),
		})
	}
	return LgBundleList, nil
}

// 实例
var instanceListmap = make(map[string]int)

func (lg *LgQuery) GetInstancesInput(ctx context.Context) {
	// 获取实例区域，使用默认区域创建客户端
	lgc := cmd2.GetDefaultAwsLgClient()
	regionList, err := lgc.GetRegions(ctx, &lightsail.GetRegionsInput{
		IncludeAvailabilityZones: aws.Bool(false),
	})
	if err != nil {
		glog.New().Error(ctx, err)
		return
	}

	// 添加表头
	glog.New().Info(ctx, "\n区域实例统计:")
	glog.New().Info(ctx, "+-----------------+--------------+")
	glog.New().Info(ctx, "| 区域            | 实例数量     |")
	glog.New().Info(ctx, "+-----------------+--------------+")

	for _, rgN := range regionList.Regions {
		// 通过区域获取客户端
		lgwithregion := cmd2.GetClient[*lightsail.Client](cmd2.WithRegion(string(rgN.Name)), cmd2.WithClientType("lightsail"))

		// 处理分页
		var allInstances []string
		var nextPageToken *string

		for {
			instances, err := lgwithregion.GetInstances(ctx, &lightsail.GetInstancesInput{
				PageToken: nextPageToken,
			})
			if err != nil {
				glog.New().Error(ctx, err)
				break
			}

			// 添加当前页的实例
			for _, instance := range instances.Instances {
				allInstances = append(allInstances, *instance.Name)
			}

			// 检查是否还有下一页
			if instances.NextPageToken == nil {
				break
			}
			nextPageToken = instances.NextPageToken
		}

		// 格式化输出每个区域的信息
		glog.New().Infof(ctx, "| %-15s | %-12d |", string(rgN.Name), len(allInstances))
		instanceListmap[string(rgN.Name)] = len(allInstances)
	}

	glog.New().Info(ctx, "+-----------------+--------------+")

	// 输出总计
	var total int
	for _, count := range instanceListmap {
		total += count
	}
	glog.New().Infof(ctx, "| %-15s | %-12d |", "总计", total)
	glog.New().Info(ctx, "+-----------------+--------------+")
}

//获取区域列表

func (lg *LgQuery) GetRegionList(ctx context.Context) (regionList []lgtypes.Region) {
	regionListOutput, err := lg.lgc.GetRegions(ctx, &lightsail.GetRegionsInput{
		IncludeAvailabilityZones: aws.Bool(false),
	})
	if err != nil {
		glog.New().Error(ctx, err)
		return nil
	}
	regionList = regionListOutput.Regions
	return regionList
}

// 获取实例
// var instanceNameList []LgAttr

func (lg *LgQuery) GetInstanceListWithRegion(ctx context.Context, region string) (instanceNameList []ctltypes.LgAttr, err error) {
	// 使用指定region的客户端
	lgc := cmd2.GetClient[*lightsail.Client](cmd2.WithRegion(region), cmd2.WithClientType("lightsail"))

	var nextPageToken *string
	for {
		instanceListOutput, err := lgc.GetInstances(ctx, &lightsail.GetInstancesInput{
			PageToken: nextPageToken,
		})
		if err != nil {
			return nil, err
		}
		account, err := lg.GetAccount(ctx)
		if err != nil {
			return nil, err
		}

		for _, instance := range instanceListOutput.Instances {
			instanceNameList = append(instanceNameList, ctltypes.LgAttr{
				Region:       region,
				Zone:         aws.ToString(instance.Location.AvailabilityZone),
				PublicIp:     aws.ToString(instance.PublicIpAddress),
				Status:       aws.ToString(instance.State.Name),
				Tags:         instance.Tags,
				CreateTime:   *instance.CreatedAt,
				InstanceName: *instance.Name,
				InstanceType: *instance.BundleId,
				Account:      account,
			})
		}

		if instanceListOutput.NextPageToken == nil {
			break
		}
		nextPageToken = instanceListOutput.NextPageToken
	}
	return instanceNameList, nil
}

// 获取所有区域实例列表
func (lg *LgQuery) GetInstanceList(ctx context.Context) (instanceNameList []ctltypes.LgAttr, err error) {
	regionList := lg.GetRegionList(ctx)
	for _, region := range regionList {
		instanceList, err := lg.GetInstanceListWithRegion(ctx, string(region.Name))
		if err != nil {
			return nil, err
		}
		instanceNameList = append(instanceNameList, instanceList...)
	}
	return instanceNameList, nil
}

// 获取账户
func (lg *LgQuery) GetAccount(ctx context.Context) (string, error) {
	// 使用 STS 客户端替代 Lightsail 客户端
	stsClient := cmd2.GetDefaultAwsStsClient()
	result, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *result.Account, nil
}

// 获取所有的快照
func (lg *LgQuery) GetSnapshotList(ctx context.Context) (snapshotList []ctltypes.LgSnapshot, err error) {
	// 获取所有区域
	rg := lg.GetRegionList(ctx)
	for _, region := range rg {
		lgcRe := cmd2.GetAwsLgClient(string(region.Name))
		snapshotListOutput, err := lgcRe.GetInstanceSnapshots(ctx, &lightsail.GetInstanceSnapshotsInput{})
		if err != nil {
			continue
		}
		for _, snapshot := range snapshotListOutput.InstanceSnapshots {
			snapshotList = append(snapshotList, ctltypes.LgSnapshot{
				SnapshotName: *snapshot.Name,
				CreatedAt:    datetime.FormatTimeToStr(*snapshot.CreatedAt, "yyyy-mm-dd"),
				InstanceName: *snapshot.FromInstanceName,
				Region:       string(snapshot.Location.RegionName),
			})
		}
	}
	return snapshotList, nil
}
