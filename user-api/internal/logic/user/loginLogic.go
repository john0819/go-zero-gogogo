package user

import (
	"context"
	"database/sql"
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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go.opentelemetry.io/otel"
)

var (
	loginCounter = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: "user_api",
		Subsystem: "auth",
		Name:      "login_total",
		Help:      "Total number of login attempts",
		Labels:    []string{"status"},
	})

	tracer = otel.Tracer("user-login")
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

	// 分布式追踪 - 创建span用于业务逻辑追踪
	_, span := tracer.Start(l.ctx, "user.login")
	defer span.End()

	// redis LOCK
	rdsLockKey := fmt.Sprintf("user:%d", req.Id)

	// transaction
	if err := l.svcCtx.UserModel.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		_, err := l.svcCtx.UserModel.TransactInsert(ctx, session, &model.User{
			Id:       uint64(req.Id),
			Name:     sql.NullString{String: req.Username, Valid: true},
			Password: sql.NullString{String: req.Password, Valid: true},
		})
		if err != nil {
			return err
		}
		logc.Infof(ctx, "插入数据user inserted: ID=%d, Name=%s", req.Id, req.Username)
		return nil
	}); err != nil {
		logc.Errorf(l.ctx, "插入数据失败: %v", err)
		return nil, err
	}

	var user *model.User

	_, lockSpan := tracer.Start(l.ctx, "user.lock")
	err = l.WithRedisLock(l.ctx, rdsLockKey, func() error {
		logc.Infof(l.ctx, "获取锁成功")

		user, err = l.svcCtx.UserModel.FindOne(l.ctx, uint64(req.Id))
		if err != nil {
			return err
		}
		return nil
	})
	lockSpan.End()

	if err != nil {
		if err == model.ErrNotFound {
			logc.Infof(l.ctx, "user not found in database, ID: %d", req.Id)
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	logc.Infof(l.ctx, "user found: ID=%d, Name=%s", user.Id, user.Name.String)
	return &types.LoginResp{
		Id:       int64(user.Id),
		Name:     user.Name.String,
		Token:    "token",
		ExpireAt: time.Now().Add(time.Hour * 24).Format(time.RFC3339),
	}, nil
}
