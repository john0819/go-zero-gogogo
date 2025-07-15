package user

import (
	"context"
	"errors"
	"time"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"
	"zero-demo/user-api/model"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
)

var (
	loginCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "user_api",
		Subsystem: "auth",
		Name:      "login_total",
		Help:      "Total number of login attempts",
		Labels:    []string{"status"},
	})
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
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, uint64(req.Id))
	if err != nil && err != model.ErrNotFound {
		return nil, errors.New("user not found")
	}

	return &types.LoginResp{
		Id:       int64(user.Id),
		Name:     user.Name.String,
		Token:    "token",
		ExpireAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
	}, nil
}
