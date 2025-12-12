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

type ChangeEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeEmailLogic {
	return &ChangeEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeEmailLogic) ChangeEmail(req *types.ChangeEmailRequest, w http.ResponseWriter) {
	// 从 JWT context 获取用户 ID
	userId := l.ctx.Value("userId").(string)

	// 调用 RPC 服务修改邮箱
	changeEmailResp, err := l.svcCtx.UserRpc.ChangeEmail(l.ctx, &user.ChangeEmailRequest{
		Id:       userId,
		NewEmail: req.NewEmail,
		Password: req.Password,
		Code:     req.Code,
	})
	if err != nil {
		l.Logger.Errorf("ChangeEmail RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if changeEmailResp.Code != 200 {
		response.BusinessError(w, changeEmailResp.Code, changeEmailResp.Message)
		return
	}

	response.SuccessWithMessage(w, "Email changed successfully", nil)
}
