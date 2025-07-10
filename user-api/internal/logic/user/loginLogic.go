package user

import (
	"context"
	"time"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	userMap := map[string]string{
		"john":  "123",
		"stone": "456",
		"hyx":   "789",
	}

	number := "unknown"
	if n, ok := userMap[req.Username]; ok {
		number = n
	}

	return &types.LoginResp{
		Id:       1,
		Name:     number,
		Token:    "token",
		ExpireAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
	}, nil
}
