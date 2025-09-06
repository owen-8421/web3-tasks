// Package config TODO
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func InitConfig(path string) {
	viper.SetConfigFile(path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Read config file fail: %s", err.Error()))
	}

	// Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到
			log.Fatal("错误：找不到配置文件 'setting.yaml'")
		} else {
			// 读取配置文件时发生其他错误
			log.Fatalf("错误：读取配置文件失败: %v", err)
		}
	}
	log.Println("配置文件加载成功！")
}
