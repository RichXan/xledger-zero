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

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest, w http.ResponseWriter) {
	// 从 JWT context 获取用户 ID
	userId := l.ctx.Value("userId").(string)

	// 调用 RPC 服务修改密码
	changePasswordResp, err := l.svcCtx.UserRpc.ChangePassword(l.ctx, &user.ChangePasswordRequest{
		Id:          userId,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		l.Logger.Errorf("ChangePassword RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if changePasswordResp.Code != 200 {
		response.BusinessError(w, changePasswordResp.Code, changePasswordResp.Message)
		return
	}

	response.SuccessWithMessage(w, "Password changed successfully", nil)
}
