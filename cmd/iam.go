package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// IamCmd 获取IAM相关信息
var IamCmd = &cobra.Command{
	Use:   "sts",
	Short: "A brief description of your command",
	Long:  `获取IAM相关信息 服务相关命令`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sts called")
	},
}

func init() {
	rootCmd.AddCommand(IamCmd)
}
