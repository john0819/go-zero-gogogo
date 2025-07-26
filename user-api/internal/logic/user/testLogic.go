package user

import (
	"context"

	"github.com/pkg/errors"

	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TestLogic {
	return &TestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TestLogic) Test(req *types.TestReq) (resp *types.TestResp, err error) {
	// if err := l.TestOne(); err != nil {
	// 	logx.Errorf("error : %+v", err)
	// }

	logc.Infof(l.ctx, "enter test logic")
	logc.Errorf(l.ctx, "enter test logic error")

	return &types.TestResp{
		Success: true,
	}, nil
}

func (l *TestLogic) TestOne() error {
	return l.TestTwo()
}

func (l *TestLogic) TestTwo() error {
	return l.TestThree()
}

func (l *TestLogic) TestThree() error {
	return errors.Wrap(errors.New("test three error"), "test three error")
}
