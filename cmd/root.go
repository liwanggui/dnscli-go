package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/liwanggui/dnscli-go/config"
	"github.com/liwanggui/dnscli-go/dnsapi"
	"github.com/liwanggui/dnscli-go/dnsapi/aliyun"
	"github.com/liwanggui/dnscli-go/dnsapi/cloudflare"
	"github.com/liwanggui/dnscli-go/dnsapi/tencent"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configName string
	secretID   string
	secretKey  string
	apiToken   string
	apiKey     string
	apiEmail   string
	cfgFile    string
	rootCmd    = &cobra.Command{
		Use:   "dnscli",
		Short: "DNS 记录管理工具",
		Long:  `DNS 记录管理工具, 支持多个DNS服务商`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "配置文件路径 (default: $HOME/.dnscli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&configName, "config-name", "N", "", "使用的 DNS 服务商配置名，用于区分多个不同的配置")
	rootCmd.PersistentFlags().StringVarP(&secretID, "secret-id", "", "", "DNS 服务商 API 密钥 ID")
	rootCmd.PersistentFlags().StringVarP(&secretKey, "secret-key", "", "", "DNS 服务商 API 密钥 Secret")
	rootCmd.PersistentFlags().StringVarP(&apiToken, "api-token", "", "", "DNS 服务商 API Token")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "api-key", "", "", "DNS 服务商 API KEY")
	rootCmd.PersistentFlags().StringVarP(&apiEmail, "api-email", "", "", "注册 DNS 服务商平台的邮箱地址")

	rootCmd.AddCommand(cCmd)
	rootCmd.AddCommand(rCmd)
	rootCmd.AddCommand(dCmd)
}

func createProvider() (dnsapi.DNSAPI, error) {
	if err := config.IsConfigFileUsed(); err != nil {
		return nil, err
	}

	configName = getCurrentConfigName()
	providerType := config.GetConfigType(configName)

	if secretID == "" {
		secretID = viper.GetString(fmt.Sprintf("configs.%s.credentials.secret_id", configName))
	}
	if secretKey == "" {
		secretKey = viper.GetString(fmt.Sprintf("configs.%s.credentials.secret_key", configName))
	}
	if apiToken == "" {
		apiToken = viper.GetString(fmt.Sprintf("configs.%s.credentials.api_token", configName))
	}
	if apiKey == "" {
		apiKey = viper.GetString(fmt.Sprintf("configs.%s.credentials.api_key", configName))
	}
	if apiEmail == "" {
		apiEmail = viper.GetString(fmt.Sprintf("configs.%s.credentials.api_email", configName))
	}

	switch providerType {
	case "aliyun":
		return aliyun.NewClient(secretID, secretKey)
	case "tencent":
		return tencent.NewClient(secretID, secretKey)
	case "cloudflare":
		return cloudflare.NewClient(apiToken, apiEmail, apiKey)
	default:
		return nil, fmt.Errorf("不支持的 DNS 服务提供商: %s", providerType)
	}
}

func getCurrentConfigName() string {
	if configName != "" {
		return configName
	}
	return config.GetDefaultConfigName()
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(fmt.Sprintf("%s%s.dnscli", home, string(os.PathSeparator)))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()

	notFound := &viper.ConfigFileNotFoundError{}
	switch {
	case err != nil && !errors.As(err, notFound):
		cobra.CheckErr(err)
	case err != nil && errors.As(err, notFound):
		// The config file is optional, we shouldn't exit when the config is not found
		break
		// default:
		// 	fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func Execute() error {
	return rootCmd.Execute()
}
