package main

import (
	"common/config"
	"common/metrics"
	"connector/app"
	"context"
	"fmt"
	"framework/game"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "connector",
	Short:   "connector 管理连接， session以及路由请求",
	Long:    "connector 管理连接， session以及路由请求",
	Run:     func(cmd *cobra.Command, args []string) {},
	PostRun: func(cmd *cobra.Command, args []string) {},
}

var (
	configFile    string
	gameConfigDir string
	serverId      string
)

func init() {
	rootCmd.Flags().StringVar(&configFile, "config", "application.yml", "app config")
	rootCmd.Flags().StringVar(&gameConfigDir, "gameDir", "../config", "game config directory")
	rootCmd.Flags().StringVar(&serverId, "serverId", "", "server id, required")
	_ = rootCmd.MarkFlagRequired("serverId")
}

func main() {
	// 1. 加载配置
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	config.InitConfig(configFile)
	game.InitConfig(gameConfigDir)

	// 2. 启动监控
	go func() {
		err := metrics.Server(fmt.Sprintf("0.0.0.0:%d", config.Conf.MetricPort))
		if err != nil {
			//panic(err)
			log.Println(err)
		}
	}()

	// 3. 启动grpc服务端
	err := app.Run(context.Background(), serverId)
	if err != nil {
		log.Println(err)
		os.Exit(-1)
	}
}
