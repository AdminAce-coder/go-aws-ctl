/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package backup

import (
	"context"
	"fmt"
	"github.com/golifez/go-aws-ctl/cmd"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/spf13/cobra"
)

// dlRPCmd represents the dlRP command
var dlRPCmd = &cobra.Command{
	Use:   "dlRP",
	Short: "删除Backup恢复点",
	Long:  `删除Backup恢复点`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dlRP called")
	},
}

func init() {
	cmd.BackupCmd.AddCommand(dlRPCmd)
}

func deleteRP(ctx context.Context, bp *backup.Client) {
	// 调用删除函数
	vaultName := "Default"                                                   // 备份库名称
	recoveryPointArn := "arn:aws:ec2:us-east-1::image/ami-0b23b0694fee8a4cc" // 恢复点 ARN
	input := &backup.DeleteRecoveryPointInput{
		BackupVaultName:  aws.String(vaultName),
		RecoveryPointArn: aws.String(recoveryPointArn),
	}
	// 调用删除恢复点的 API
	_, err := bp.DeleteRecoveryPoint(ctx, input)
	if err != nil {
		glog.New().Errorf(ctx, "unable to delete recovery point: %v", err)
	}
	glog.New().Infof(ctx, "删除成功")

}
