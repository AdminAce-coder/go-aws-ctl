/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lgCmd represents the lg command
var LgCmd = &cobra.Command{
	Use:   "lg",
	Short: "A brief description of your command",
	Long:  `AWS Lightsail 服务相关命令`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("lg called")
	},
}

func init() {
	rootCmd.AddCommand(LgCmd)
}
