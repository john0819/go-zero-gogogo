package logic

import (
	"context"

	"zero-demo/grpc-gateway/server/internal/svc"
	"zero-demo/grpc-gateway/server/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.UserReq) (*user.UserResp, error) {
	userMap := map[string]string{
		"1": "john",
		"2": "wong",
	}

	if _, ok := userMap[in.Uid]; !ok {
		loginCounter.Inc("failed")
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	loginCounter.Inc("success")

	return &user.UserResp{
		Name: userMap[in.Uid],
	}, nil
}
