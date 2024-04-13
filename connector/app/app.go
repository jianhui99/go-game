package app

import (
	"common/config"
	"common/logs"
	"connector/route"
	"context"
	"core/repo"
	"framework/connector"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动程序 启动grpc服务，http服务，日志，数据库
func Run(ctx context.Context, serverId string) error {
	// 日志
	logs.InitLog(config.Conf.AppName)
	exit := func() {}
	go func() {
		c := connector.Default()
		exit = c.Close
		manager := repo.NewManager()
		c.RegisterHandler(route.Register(manager))
		c.Run(serverId)
	}()
	stop := func() {
		exit()
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
				logs.Info("quit connector")
				return nil
			case syscall.SIGHUP:
				stop()
				logs.Info("hang up, stop connector")
				return nil
			default:
				return nil
			}
		}
	}
}
