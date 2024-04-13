package api

import (
	"common"
	"common/biz"
	"common/config"
	"common/jwts"
	"common/logs"
	"common/rpc"
	"context"
	"framework/msError"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"user/pb"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) Register(ctx *gin.Context) {
	// 接受参数
	var req pb.RegisterParams
	err2 := ctx.ShouldBind(&req)
	if err2 != nil {
		common.Fail(ctx, biz.RequestDataError)
		return
	}

	response, err := rpc.UserClient.Register(context.TODO(), &req)
	if err != nil {
		common.Fail(ctx, msError.ToError(err))
		return
	}

	uid := response.Uid
	if len(uid) == 0 {
		common.Fail(ctx, biz.SqlError)
		return
	}

	// generate token by uid with jwt
	claims := jwts.CustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token, err := jwts.GenToken(&claims, config.Conf.Jwt.Secret)
	if err != nil {
		logs.Error("register jwt gen token err:%v", err)
		common.Fail(ctx, biz.Fail)
		return
	}

	result := map[string]any{
		"token": token,
		"serverInfo": map[string]any{
			"host": config.Conf.Services["connector"].ClientHost,
			"port": config.Conf.Services["connector"].ClientPort,
		},
	}
	common.Success(ctx, result)
}
