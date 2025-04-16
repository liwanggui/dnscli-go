package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	dCmd = &cobra.Command{
		Use:   "domain",
		Short: "管理域名",
	}

	dAddCmd = &cobra.Command{
		Use:   "add DOMAIN",
		Short: "增加管理域名",
	}

	dDelCmd = &cobra.Command{
		Use:   "del DOMAIN",
		Short: "删除域名",
	}

	dListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "ls"},
		Short:   "查看域名列表",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := createProvider()
			if err != nil {
				cobra.CheckErr(err)
			}
			domains, err := client.ListDomains()
			if err != nil {
				cobra.CheckErr(err)
			}
			for _, domain := range domains {
				fmt.Println(domain)
			}
		},
	}
)

func init() {
	dCmd.AddCommand(dAddCmd)
	dCmd.AddCommand(dDelCmd)
	dCmd.AddCommand(dListCmd)
}
