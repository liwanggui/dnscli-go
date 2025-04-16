package config

import (
	"fmt"
	"regexp"

	"github.com/spf13/viper"
)

// ValidConfigName 判断配置名是否有效
func ValidConfigName(name string) error {
	if name == "" {
		return fmt.Errorf("配置名不允许为空")
	}

	configs := viper.GetStringMap("configs")
	for k := range configs {
		if name == k {
			return fmt.Errorf("配置名已存在：%s", name)
		}
	}
	return nil
}

// ValidProviderType 判断服务商是否支持
func ValidProviderType(providerType string) error {
	for _, p := range DnsServiceList {
		if p == providerType {
			return nil
		}
	}
	return fmt.Errorf("暂不支持此 DNS 服务商：%s", providerType)
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidEmail 判断邮箱地址是否合法
func ValidEmail(email string) error {
	if emailRegex.MatchString(email) {
		return nil
	}
	return fmt.Errorf("无效的邮箱地址：%s", email)
}
