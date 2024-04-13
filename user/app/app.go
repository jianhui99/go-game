package app

import (
	"common/config"
	"common/discovery"
	"common/logs"
	"context"
	"core/repo"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/service"
	"user/pb"
)

// Run 启动程序 启动grpc服务，http服务，日志，数据库
func Run(ctx context.Context) error {
	// 日志
	logs.InitLog(config.Conf.AppName)
	// 启动grpc服务端
	server := grpc.NewServer()
	register := discovery.NewRegister()

	// 注册grpc service 想要数据库 mongo，redis
	// init 数据库管理
	manager := repo.NewManager()

	go func() {
		lis, err := net.Listen("tcp", config.Conf.Grpc.Addr)
		if err != nil {
			logs.Fatal("failed to listen grpc: %v", err)
		}

		err = register.Register(config.Conf.Etcd)
		if err != nil {
			logs.Fatal("failed to register grpc to etcd: %v", err)
		}

		pb.RegisterUserServiceServer(server, service.NewAccountService(manager))

		// 阻塞操作
		err = server.Serve(lis)
		if err != nil {
			logs.Fatal("failed to serve grpc: %v", err)
		}
	}()
	stop := func() {
		server.Stop()
		register.Close()
		manager.Close()

		time.Sleep(3 * time.Second)
		logs.Info("stop app")
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for {
		select {
		case <-ctx.Done():
			stop()
			logs.Info("stop app")
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
				logs.Info("unknow signal")
				return nil
			}
		}
	}
}
