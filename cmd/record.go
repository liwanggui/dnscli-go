package main

import (
	"fmt"
	"github.com/liwanggui/dnscli-go/config"
	"github.com/liwanggui/dnscli-go/dnsapi"
	"github.com/liwanggui/dnscli-go/util"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	rCmd = &cobra.Command{
		Use:   "record",
		Short: "域名解析记录管理",
		Long:  "域名解析记录管理",
	}

	rAddCmd = &cobra.Command{
		Use:          "create DOMAIN RECORD_NAME RECORD_TYPE RECORD_VALUE",
		Aliases:      []string{"a", "add"},
		Short:        "创建解析记录",
		Example:      `  dnscli create example.com www A 1.1.1.1`,
		Args:         cobra.ExactArgs(4),
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := createProvider()
			if err != nil {
				cobra.CheckErr(err)
			}
			param := dnsapi.CreateParameter(args[0])
			param.Name = args[1]
			param.Type = args[2]
			param.Value = args[3]
			if err := dnsapi.ValidRecordType(param.Type); err != nil {
				cobra.CheckErr(err)
			}
			if err := client.AddRecord(param); err != nil {
				cobra.CheckErr(err)
			}
		},
	}

	rDelCmd = &cobra.Command{
		Use:     "delete DOMAIN RECORD_ID...",
		Aliases: []string{"d", "del"},
		Short:   "删除解析记录",
		Example: "  dnscli record delete 1234567890",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client, err := createProvider()
			if err != nil {
				cobra.CheckErr(err)
			}
			param := dnsapi.CreateParameter(args[0])
			for _, rid := range args[1:] {
				param.ID = rid
				if err := client.DeleteRecord(param); err != nil {
					cobra.CheckErr(err)
				}
				fmt.Printf("deleted %s ok\n", rid)
			}
		},
	}

	rUpdateCmd = &cobra.Command{
		Use:     "update DOMAIN RECORD_ID RECORD_NAME RECORD_TYPE RECORD_VALUE",
		Aliases: []string{"u"},
		Short:   "更新解析记录",
		Long:    `查询 DNS 服务商账号下指定域名的解析记录`,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := createProvider()
			if err != nil {
				cobra.CheckErr(err)
			}
			param := dnsapi.CreateParameter(args[0])
			param.ID = args[1]
			param.Name = args[2]
			param.Type = args[3]
			param.Value = args[4]
			param.Proxied, _ = cmd.Flags().GetBool("proxied")
			param.TTL, _ = cmd.Flags().GetInt("ttl")
			param.Line, _ = cmd.Flags().GetString("line")
			
			if err := dnsapi.ValidRecordType(param.Type); err != nil {
				cobra.CheckErr(err)
			}
			if err := client.UpdateRecord(param); err != nil {
				cobra.CheckErr(err)
			}
			fmt.Printf("updated %s ok\n", args[1])
		},
	}

	rListCmd = &cobra.Command{
		Use:          "list DOMAIN",
		Aliases:      []string{"l", "ls"},
		Short:        "查询解析记录",
		Long:         `查询 DNS 服务商账号下指定域名的解析记录`,
		Example:      `  dnscli record list example.com`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			client, err := createProvider()
			if err != nil {
				cobra.CheckErr(err)
			}
			param := dnsapi.CreateParameter(args[0])
			param.Name, _ = cmd.Flags().GetString("name")
			param.Type, _ = cmd.Flags().GetString("type")
			param.Value, _ = cmd.Flags().GetString("value")
			param.Line, _ = cmd.Flags().GetString("line")
			if param.Type != "" {
				if err := dnsapi.ValidRecordType(param.Type); err != nil {
					cobra.CheckErr(err)
				}
			}
			records, err := client.ListRecords(param)
			if err != nil {
				cobra.CheckErr(err)
			}
			exclField := make([]string, 1)
			cName := getCurrentConfigName()
			cType := config.GetConfigType(cName)
			switch cType {
			case "aliyun":
				exclField = append(exclField, "Proxied", "Updated")
			case "tencent":
				exclField = append(exclField, "Proxied")
			case "cloudflare":
				exclField = append(exclField, "Line")
			}

			fmt.Println("记录数:", len(records))
			table := tablewriter.NewWriter(os.Stdout)
			for i, record := range records {
				n, v := util.GetStructFieldNamesAndValues(record, "table", exclField)
				if i == 0 {
					table.SetHeader(n)
				}
				table.Append(v)
			}
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.Render()
		},
	}
)

func init() {
	rCmd.PersistentFlags().Int("ttl", 0, "解析记录 TTL 值，免费 DNS 解析基本都不支持小于 600s")
	rCmd.PersistentFlags().String("line", "", "解析线路名，需要 DNS 服务商提供支持")
	rCmd.PersistentFlags().Bool("proxied", false, "是否启动 CND 加速，仅 cloudflare 使用 (default: false)")

	rListCmd.Flags().StringP("name", "n", "", "解析记录名")
	rListCmd.Flags().StringP("type", "t", "", fmt.Sprintf("解析记录类型, 取值(%s)", strings.Join(dnsapi.RecordTypes, ",")))
	rListCmd.Flags().StringP("value", "v", "", "解析记录值")

	rCmd.AddCommand(rAddCmd)
	rCmd.AddCommand(rDelCmd)
	rCmd.AddCommand(rListCmd)
	rCmd.AddCommand(rUpdateCmd)
}
