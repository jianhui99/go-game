package discovery

import (
	"common/config"
	"common/logs"
	"context"
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// Register grpc服务注册到etcd
// 原理：创建一个租约，grpc服务注册到etcd，绑定租约
// 过了租约时间，etcd就会删除grpc服务信息
// 实现心跳，完成续租，如果etcd没有就新注册
type Register struct {
	etcdClient  *clientv3.Client                        // etcd连接
	leaseId     clientv3.LeaseID                        // 租约id
	DialTimeout int                                     // 超时时间
	ttl         int                                     // 租约时间
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse // 心跳
	info        Server                                  // 注册的server信息
	closeCh     chan struct{}
}

func NewRegister() *Register {
	return &Register{
		DialTimeout: 3,
	}
}

func (r *Register) Close() {
	r.closeCh <- struct{}{}
}

func (r *Register) Register(conf config.EtcdConf) error {
	// 注册信息
	info := Server{
		Name:    conf.Register.Name,
		Addr:    conf.Register.Addr,
		Weight:  conf.Register.Weight,
		Version: conf.Register.Version,
		Ttl:     conf.Register.Ttl,
	}

	// 建立etcd连接
	var err error
	r.etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   conf.Addrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	})
	if err != nil {
		return err
	}
	r.info = info
	if err = r.register(); err != nil {
		return err
	}
	r.closeCh = make(chan struct{})

	// 放入携程中，根据心跳结果做相应的操作
	go r.watcher()
	return nil
}

func (r *Register) register() error {
	// 1. 创建租约
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*time.Duration(r.DialTimeout))
	defer cancelFunc()
	var err error
	if err = r.createLease(ctx, r.info.Ttl); err != nil {
		return err
	}

	// 2. 心跳检测
	if r.keepAliveCh, err = r.keepAlive(); err != nil {
		return err
	}

	// 3. 绑定租约
	data, _ := json.Marshal(r.info)

	return r.binLease(ctx, r.info.BuildRegisterKey(), string(data))
}

// createLease ttl秒
func (r *Register) binLease(ctx context.Context, key, value string) error {
	// put动作
	_, err := r.etcdClient.Put(ctx, key, value, clientv3.WithLease(r.leaseId))
	if err != nil {
		logs.Error("bind lease err:", err)
		return err
	}
	logs.Info("bind lease success, key=%s", key)
	return nil
}

// createLease ttl秒
func (r *Register) createLease(ctx context.Context, ttl int64) error {
	grant, err := r.etcdClient.Grant(ctx, ttl)
	if err != nil {
		logs.Error("create lease err: %v", err)
		return err
	}
	r.leaseId = grant.ID
	return nil
}

// keepAlive 心跳检测
func (r *Register) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	alive, err := r.etcdClient.KeepAlive(context.Background(), r.leaseId)
	if err != nil {
		logs.Error("keep alive err:", err)
		return alive, err
	}
	return alive, nil
}

// watcher 续约 新注册 close 注销
func (r *Register) watcher() {
	//租约到期了 是不是需要去检查是否自动注册
	ticker := time.NewTicker(time.Duration(r.info.Ttl) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				logs.Error("close and unregister failed,err:%v", err)
			}
			//租约撤销
			if _, err := r.etcdClient.Revoke(context.Background(), r.leaseId); err != nil {
				logs.Error("close and Revoke lease failed,err:%v", err)
			}
			if r.etcdClient != nil {
				r.etcdClient.Close()
			}
			logs.Info("unregister etcd...")
		case res := <-r.keepAliveCh:
			logs.Info("续约成功,%v", res)
			if res != nil {
				if err := r.register(); err != nil {
					logs.Error("keepAliveCh register failed,err:%v", err)
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					logs.Error("ticker register failed,err:%v", err)
				}
			}
		}
	}
}

func (r *Register) unregister() error {
	_, err := r.etcdClient.Delete(context.Background(), r.info.BuildRegisterKey())
	return err
}
