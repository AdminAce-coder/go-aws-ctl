/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// backupCmd represents the backup command
var BackupCmd = &cobra.Command{
	Use:   "backup",
	Short: "A brief description of your command",
	Long:  `backup相关操作`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("backup called")
	},
}

func init() {
	rootCmd.AddCommand(BackupCmd)

}
