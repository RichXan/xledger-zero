// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"net/http"

	"xledger/service/user/api/internal/response"
	"xledger/service/user/api/internal/svc"
	"xledger/service/user/api/internal/types"
	"xledger/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutRequest, w http.ResponseWriter) {
	// 调用 RPC 服务登出
	logoutResp, err := l.svcCtx.UserRpc.Logout(l.ctx, &user.LogoutRequest{
		AccessToken: req.AccessToken,
	})
	if err != nil {
		l.Logger.Errorf("Logout RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if logoutResp.Code != 200 {
		response.BusinessError(w, logoutResp.Code, logoutResp.Message)
		return
	}

	response.SuccessWithMessage(w, "Logout successful", nil)
}
