package user

import (
	"context"
	"database/sql"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"
	"zero-demo/user-api/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateReq) (resp *types.UpdateResp, err error) {
	// 正确的插入操作
	result, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Name:     sql.NullString{String: req.Username, Valid: true},
		Password: sql.NullString{String: req.Password, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	// 获取插入的用户ID
	userId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &types.UpdateResp{
		Success: userId, // 返回插入的用户ID表示成功
	}, nil
}
