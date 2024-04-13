package app

import (
	"common/config"
	"common/logs"
	"context"
	"fmt"
	"gateway/router"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序 启动grpc服务，http服务，日志，数据库
func Run(ctx context.Context) error {
	// 日志
	logs.InitLog(config.Conf.AppName)

	go func() {
		// 启动gin, 注册路由
		r := router.RegisterRouter()

		// http接口
		if err := r.Run(fmt.Sprintf(":%d", config.Conf.HttpPort)); err != nil {
			logs.Fatal("gin run err: %v", err)
		}
	}()
	stop := func() {
		time.Sleep(3 * time.Second)
		logs.Info("stop app")
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for {
		select {
		case <-ctx.Done():
			stop()

			return nil
		case s := <-c:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				stop()
				logs.Info("quit app")
				return nil
			case syscall.SIGHUP:
				stop()
				logs.Info("hang up, stop app")
				return nil
			default:
				return nil
			}
		}
	}
}
