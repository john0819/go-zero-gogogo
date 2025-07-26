package svc

import (
	"zero-demo/user-api/internal/config"
	"zero-demo/user-api/internal/middleware"
	"zero-demo/user-api/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config         config.Config
	TestMiddleware rest.Middleware
	UserModel      model.UserModel
	Redis          *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:         c,
		TestMiddleware: middleware.NewTestMiddleware().Handle,
		UserModel:      model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		Redis:          redis.MustNewRedis(c.Redis),
	}
}
