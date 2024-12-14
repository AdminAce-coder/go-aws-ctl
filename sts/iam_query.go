package sts

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	cmd2 "github.com/golifez/go-aws-ctl/cmd"
	"github.com/golifez/go-aws-ctl/model"
	"github.com/golifez/go-aws-ctl/service"
	"github.com/spf13/cobra"
)

type IamQueryCommand struct {
	StsClint  *sts.Client
	IamClient *iam.Client
}

var IamQueryCmd = &cobra.Command{
	Use:   "query",
	Short: "A brief description of your command",
	Long:  `STS相关查询操作`,
	Run: func(cmd *cobra.Command, args []string) {
		//获取客户端
	},
}

func init() {
	cmd2.LgCmd.AddCommand(IamQueryCmd)
}

func NewIamQueryCommand() service.IamQuery {
	return &IamQueryCommand{
		StsClint:  cmd2.GetDefaultAwsStsClient(),
		IamClient: cmd2.GetDefaultAwsIamClient(),
	}
}

// 查询IAM信息
func (iq *IamQueryCommand) GetIamInfo(ctx context.Context) (model.IamInfo, error) {
	// 初始化 IamInfo，特别是要初始化 Policy map
	iamInfo := model.IamInfo{
		Policy: make(map[string]string),
	}

	output, err := iq.StsClint.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return iamInfo, err
	}

	// 解析用户名 arn:aws:iam::905418123359:user/admin
	if strings.Contains(*output.Arn, ":user/") {
		arns := strings.Split(*output.Arn, "/")
		username := arns[len(arns)-1]

		// 使用 IAM 客户端查询策略
		policies, err := iq.IamClient.ListAttachedUserPolicies(ctx, &iam.ListAttachedUserPoliciesInput{
			UserName: &username,
		})
		if err != nil {
			return iamInfo, err
		}

		// 遍历每个策略并获取详细信息
		for _, policy := range policies.AttachedPolicies {
			policyDetails, err := iq.IamClient.GetPolicy(ctx, &iam.GetPolicyInput{
				PolicyArn: policy.PolicyArn,
			})
			if err != nil {
				fmt.Printf("获取策略详情失败: %v\n", err)
				continue
			}
			// 将策略名和描述添加到iamInfo.Policy中
			if policyDetails.Policy.Description != nil {
				iamInfo.Policy[*policy.PolicyName] = *policyDetails.Policy.Description
			} else {
				iamInfo.Policy[*policy.PolicyName] = "" // 如果没有描述，设置为空字符串
			}
		}
		iamInfo.UserName = username
		iamInfo.Account = *output.Account
		return iamInfo, nil
	}

	return iamInfo, nil
}
