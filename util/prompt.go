package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

type ValidateFunc func(string) error

// String 提示用户输入字符串
func String(message string, valid ValidateFunc) string {
	for {
		fmt.Print(message + ": ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return ""
		}

		// 去除输入末尾的换行符
		input = strings.TrimSpace(input)

		if valid == nil {
			return input
		}

		if err := valid(input); err != nil {
			PrintError(err)
			continue
		}
		return input
	}
}

// Password 提示用户输入密码（输入字符不可见）
func Password(message string, valid ValidateFunc) string {
	for {
		fmt.Printf("%s: ", message)
		password, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return ""
		}
		fmt.Println() // 输入完成后换行
		if valid == nil {
			return string(password)
		}
		if err := valid(string(password)); err != nil {
			PrintError(err)
			continue
		}
		return string(password)
	}
}

// Confirm 提示用户确认操作
func Confirm(message string, defaultValue bool, valid ValidateFunc) bool {
	var prompt string
	if defaultValue {
		prompt = message + " [Y/n]: "
	} else {
		prompt = message + " [y/N]: "
	}

	for {
		fmt.Print(prompt)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return defaultValue
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "" {
			return defaultValue
		}
		if valid == nil {
			return input == "y" || input == "yes"
		}

		if err := valid(input); err != nil {
			PrintError(err)
			continue // 继续循环，等待用户重新输入
		}
	}
}

// Select 提示用户从选项列表中选择
func Select(message string, options []string, defaultIndex int) (string, int) {
	if len(options) == 0 {
		return "", -1
	}

	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}

	for {
		prompt := fmt.Sprintf("请选择 (1-%d)", len(options))
		if defaultIndex >= 0 && defaultIndex < len(options) {
			prompt += fmt.Sprintf(" [%d]", defaultIndex+1)
		}
		prompt += ": "

		fmt.Print(prompt)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil && defaultIndex >= 0 && defaultIndex < len(options) {
			return options[defaultIndex], defaultIndex
		}

		input = strings.TrimSpace(input)
		if input == "" && defaultIndex >= 0 && defaultIndex < len(options) {
			return options[defaultIndex], defaultIndex
		}

		// 尝试将输入转换为数字
		var index int
		_, err = fmt.Sscanf(input, "%d", &index)
		if err != nil || index < 1 || index > len(options) {
			PrintError(fmt.Errorf("%w, 无效的选择，请重新输入", err))
			continue
		}
		return options[index-1], index - 1
	}
}

func PrintError(err error) {
	fmt.Printf("\033[31m%v\033[0m\n", err)
}

func ValidStringNotEmpty(s string) error {
	if s == "" {
		return fmt.Errorf("输入不能为空")
	}
	return nil
}
