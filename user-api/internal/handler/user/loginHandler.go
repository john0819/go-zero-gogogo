package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-demo/user-api/internal/logic/user"
	"zero-demo/user-api/internal/svc"
	"zero-demo/user-api/internal/types"
)

// 获取用户信息
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
