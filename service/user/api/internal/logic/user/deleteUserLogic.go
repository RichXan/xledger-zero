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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserRequest, w http.ResponseWriter) {
	// 调用 RPC 服务删除用户
	deleteUserResp, err := l.svcCtx.UserRpc.DeleteUser(l.ctx, &user.DeleteUserRequest{
		Id: req.Id,
	})
	if err != nil {
		l.Logger.Errorf("DeleteUser RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if deleteUserResp.Code != 200 {
		response.BusinessError(w, deleteUserResp.Code, deleteUserResp.Message)
		return
	}

	response.SuccessWithMessage(w, "User deleted successfully", nil)
}
