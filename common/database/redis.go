package database

import (
	"common/config"
	"common/logs"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisManager struct {
	Client        *redis.Client        // 单机
	ClusterClient *redis.ClusterClient // 集群
}

func NewRedis() *RedisManager {
	var clusterClient *redis.ClusterClient
	var client *redis.Client

	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()

	addrs := config.Conf.Database.RedisConf.ClusterAddrs

	// 非集群，单节点
	if len(addrs) == 0 {
		logs.Info("redis cluster address is empty, use default cluster address")
		client = redis.NewClient(&redis.Options{
			Addr:         config.Conf.Database.RedisConf.Addr,
			PoolSize:     config.Conf.Database.RedisConf.PoolSize,
			MinIdleConns: config.Conf.Database.RedisConf.MinIdleConns,
			Password:     config.Conf.Database.RedisConf.Password,
		})
	} else {
		logs.Info("redis cluster address is %v", addrs)
		clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        config.Conf.Database.RedisConf.ClusterAddrs,
			PoolSize:     config.Conf.Database.RedisConf.PoolSize,
			MinIdleConns: config.Conf.Database.RedisConf.MinIdleConns,
			Password:     config.Conf.Database.RedisConf.Password,
		})
	}

	if clusterClient != nil {
		if err := clusterClient.Ping(ctx).Err(); err != nil {
			logs.Fatal("redis cluster connect err: ", err)
			return nil
		}

		err := clusterClient.Incr(ctx, "test_1").Err()
		if err != nil {
			logs.Fatal("redis cluster connect err: ", err)
		}
	}

	if client != nil {
		if err := client.Ping(ctx).Err(); err != nil {
			logs.Fatal("redis client connect err: ", err)
			return nil
		}

		err := client.Incr(ctx, "test_1").Err()
		if err != nil {
			logs.Fatal("redis client connect err: ", err)
		}
	}

	return &RedisManager{
		Client:        client,
		ClusterClient: clusterClient,
	}
}

func (r *RedisManager) Close() {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			logs.Error("redis close err: ", err)
		}
	}

	if r.ClusterClient != nil {
		if err := r.ClusterClient.Close(); err != nil {
			logs.Error("redis cluster close err: ", err)
		}
	}
}

func (r *RedisManager) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r.Client != nil {
		return r.Client.Set(ctx, key, value, expiration).Err()

	}

	if r.ClusterClient != nil {
		return r.ClusterClient.Set(ctx, key, value, expiration).Err()
	}

	return nil
}
