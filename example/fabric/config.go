package fabric

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	//监听改变动态跟新配置
	go watchConfig()
	//加载配置
	loadConfig()
}

//监听配置改变
func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		//改变重新加载
		loadConfig()
	})
}

//加载配置
func loadConfig() {
	viper.SetConfigName("feconfig") // name of kubeconfig file
	viper.AddConfigPath(".")        // optionally look for kubeconfig in the working directory
	viper.AddConfigPath("../")      // path to look for the kubeconfig file in
	err := viper.ReadInConfig()     // Find and read the feconfig.yaml file
	if err != nil {                 // Handle errors reading the kubeconfig file
		os.Exit(-1)
	}
	//全局配置
	Config = viper.GetViper()
}
