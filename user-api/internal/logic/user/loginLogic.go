package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"
	"zero-demo/user-api/model"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/metric"
	"github.com/zeromicro/go-zero/core/stores/redis"
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

// 实践锁
func (l *LoginLogic) WithRedisLock(ctx context.Context, key string, fn func() error) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)

	acquired, err := lock.AcquireCtx(ctx)
	if err != nil {
		return err
	}
	if !acquired {
		logc.Errorf(ctx, "无法获取锁")
		return errors.New("lock not acquired")
	}

	defer lock.ReleaseCtx(ctx)

	// do something
	return fn()
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
	logc.Infof(l.ctx, "enter login logic")
	logc.Infof(l.ctx, "querying user with ID: %d", req.Id)

	rdsLockKey := fmt.Sprintf("user:%d", req.Id)

	var user *model.User

	err = l.WithRedisLock(l.ctx, rdsLockKey, func() error {
		logc.Infof(l.ctx, "获取锁成功")

		user, err = l.svcCtx.UserModel.FindOne(l.ctx, uint64(req.Id))
		if err != nil {
			return err
		}
		return nil
	})

	logc.Infof(l.ctx, "user found: ID=%d, Name=%s", user.Id, user.Name.String)
	return &types.LoginResp{
		Id:       int64(user.Id),
		Name:     user.Name.String,
		Token:    "token",
		ExpireAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
	}, nil
}
