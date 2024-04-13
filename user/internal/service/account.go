package service

import (
	"common/biz"
	"common/logs"
	"context"
	"core/dao"
	"core/models/entity"
	"core/models/requests"
	"core/repo"
	"framework/msError"
	"time"
	"user/pb"
)

// 创建账号

type AccountService struct {
	accountDao *dao.AccountDao
	redisDao   *dao.RedisDao
	pb.UnimplementedUserServiceServer
}

func NewAccountService(manager *repo.Manager) *AccountService {
	return &AccountService{
		accountDao: dao.NewAccountDao(manager),
		redisDao:   dao.NewRedisDao(manager),
	}
}

func (a *AccountService) Register(ctx context.Context, req *pb.RegisterParams) (*pb.RegisterResponse, error) {
	// 注册逻辑
	logs.Info("register service called")
	if req.LoginPlatform == requests.WeiXin {
		account, err := a.wxRegister(req)
		if err != nil {
			return &pb.RegisterResponse{}, msError.GrpcError(err)
		}

		return &pb.RegisterResponse{
			Uid: account.Uid,
		}, nil
	}
	return &pb.RegisterResponse{}, nil
}

// 1.封装account struct，存入数据库 mongo
// 2.需要生成几个数字做为用户的唯一id  redis自增
func (a *AccountService) wxRegister(req *pb.RegisterParams) (*entity.Account, *msError.Error) {
	account := &entity.Account{
		WxAccount:  req.Account,
		CreateTime: time.Now(),
	}
	uid, err := a.redisDao.NextAccountId()
	if err != nil {
		return account, biz.SqlError
	}
	account.Uid = uid
	err = a.accountDao.SaveAccount(context.TODO(), account)
	if err != nil {
		return account, biz.SqlError
	}

	return account, nil
}
