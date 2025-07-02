package logic

import (
	"context"

	"zero-demo/user-rpc/internal/svc"
	"zero-demo/user-rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SayHelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloLogic {
	return &SayHelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 定义一个 SayHello 一元 rpc 方法，请求体和响应体必填。
func (l *SayHelloLogic) SayHello(in *pb.SayHelloReq) (*pb.SayHelloResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SayHelloResp{}, nil
}
