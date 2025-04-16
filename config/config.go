package config

import (
	"errors"
	"fmt"
	"sort"

	"github.com/liwanggui/dnscli-go/util"
	"github.com/spf13/viper"
)

const DefaultItemName = "default"

// DnsServiceList 支持的 DNS 服务商列表
var DnsServiceList = []string{"aliyun", "tencent", "cloudflare"}

func IsConfigFileUsed() error {
	if viper.ConfigFileUsed() == "" {
		return fmt.Errorf("config file not exist")
	}
	return nil
}

func GetConfigNames() []string {
	allConfig := viper.GetStringMap("configs")
	var configNameList []string
	for config := range allConfig {
		configNameList = append(configNameList, config)
	}
	sort.Strings(configNameList)
	return configNameList
}

func GetConfigType(name string) string {
	return viper.GetString(fmt.Sprintf("configs.%s.type", name))
}

func GetDefaultConfigName() string {
	return viper.GetString(DefaultItemName)
}

func AddConfig(name, pType, secretID, secretKey, apiToken, apiEmail, apiKey string, isDefault bool) error {
	if err := ValidConfigName(name); err != nil {
		name = util.String("请输入配置名", ValidConfigName)
	}

	if err := ValidProviderType(pType); err != nil {
		pType, _ = util.Select("请选择服务商", DnsServiceList, 0)
	}

	viper.Set(fmt.Sprintf("configs.%s.type", name), pType)

	switch pType {
	case "cloudflare":
		if apiToken == "" {
			apiToken = util.Password("请输入 Cloudflare API Token", nil)
		}

		if err := ValidEmail(apiEmail); err != nil {
			apiEmail = util.String("请输入注册 Cloudflare 平台的邮箱地址", ValidEmail)
		}

		if apiToken == "" && apiKey == "" {
			apiKey = util.Password("请输入 Cloudflare API Key", util.ValidStringNotEmpty)
		} else if apiKey == "" {
			apiKey = util.Password("请输入 Cloudflare API Key", nil)
		}

		viper.Set(fmt.Sprintf("configs.%s.credentials.api_token", name), apiToken)
		viper.Set(fmt.Sprintf("configs.%s.credentials.api_key", name), apiKey)
		viper.Set(fmt.Sprintf("configs.%s.credentials.api_email", name), apiEmail)

	default:
		if secretID == "" {
			secretID = util.Password("请输入 DNS 服务商 API 密钥 ID", util.ValidStringNotEmpty)
		}
		if secretKey == "" {
			secretKey = util.Password("请输入 DNS 服务商 API 密钥 Secret", util.ValidStringNotEmpty)
		}
		viper.Set(fmt.Sprintf("configs.%s.credentials.secret_id", name), secretID)
		viper.Set(fmt.Sprintf("configs.%s.credentials.secret_key", name), secretKey)
	}

	if isDefault || !viper.IsSet(DefaultItemName) {
		viper.Set("default", name)
	}

	return WriteConfig()
}

func SetDefaultConfig(name string) error {
	defaultName := GetDefaultConfigName()
	if name == defaultName {
		return nil
	}

	if !viper.IsSet(fmt.Sprintf("configs.%s", name)) {
		return fmt.Errorf("配置名不存在，请检查后重试: %s\n", name)
	}

	viper.Set("default", name)
	return WriteConfig()
}

func WriteConfig() error {
	err := viper.WriteConfig()
	if err == nil {
		return nil
	}
	var configFileNotFoundError viper.ConfigFileNotFoundError
	if !errors.As(err, &configFileNotFoundError) {
		return err
	}
	// 如果文件不存在，创建并写入配置
	if err = viper.SafeWriteConfig(); err != nil {
		return err
	}
	return nil
}
