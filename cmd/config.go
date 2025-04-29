package cmd

import (
	"fmt"
	"github.com/liwanggui/dnscli-go/config"
	"github.com/spf13/cobra"
	"strings"
)

var (
	providerType string

	cCmd = &cobra.Command{
		Use:   "config",
		Short: "配置管理",
		Long:  `管理 DNS 服务商配置`,
	}

	cAddCmd = &cobra.Command{
		Use:   "add",
		Short: "添加 DNS 服务商配置",
		Long:  `添加 DNS 服务商配置`,
		Run: func(cmd *cobra.Command, args []string) {
			isDefault, _ := cmd.Flags().GetBool("default")
			err := config.AddConfig(configName, providerType, secretID, secretKey, apiToken, apiEmail, apiKey, isDefault)
			if err != nil {
				cobra.CheckErr(err)
			}
		},
	}

	cSetCmd = &cobra.Command{
		Use:          "set-default [config name]",
		Aliases:      []string{"set", "default"},
		Short:        "设置默认 DNS 服务商配置",
		Long:         `设置默认 DNS 服务商配置`,
		Example:      `  dnscli config set-default cf -> 设置 cf 配置为默认配置`,
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		PreRunE:      preRunE,
		Run: func(cmd *cobra.Command, args []string) {
			itemName := args[0]
			if err := config.SetDefaultConfig(itemName); err != nil {
				cobra.CheckErr(err)
			}
		},
	}

	cListCmd = &cobra.Command{
		Use:          "list",
		Aliases:      []string{"ls"},
		Short:        "查看 DNS 服务商配置",
		Long:         `查看 DNS 服务商配置`,
		SilenceUsage: true,
		PreRunE:      preRunE,
		Run: func(cmd *cobra.Command, args []string) {
			defaultName := config.GetDefaultConfigName()
			configNameList := config.GetConfigNames()
			for _, itemName := range configNameList {
				if itemName == defaultName {
					fmt.Printf("* %s (%s)\n", itemName, config.GetConfigType(itemName))
				} else {
					fmt.Printf("  %s (%s)\n", itemName, config.GetConfigType(itemName))
				}
			}
		},
	}
)

func preRunE(cmd *cobra.Command, args []string) error {
	return config.IsConfigFileUsed()
}

func init() {
	cAddCmd.Flags().StringVarP(&providerType, "dnsapi", "p", "",
		fmt.Sprintf("DNS 服务提供商, 取值为 (%s)", strings.Join(config.DnsServiceList, ", ")))

	cAddCmd.Flags().Bool("default", false, "是否设置为默认配置")
	cCmd.AddCommand(cAddCmd)
	cCmd.AddCommand(cListCmd)
	cCmd.AddCommand(cSetCmd)
}
