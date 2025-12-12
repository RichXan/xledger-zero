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

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileRequest, w http.ResponseWriter) {
	// 从 JWT context 获取用户 ID
	userId := l.ctx.Value("userId").(string)

	// 调用 RPC 服务更新资料
	updateProfileResp, err := l.svcCtx.UserRpc.UpdateProfile(l.ctx, &user.UpdateProfileRequest{
		Id:       userId,
		Username: req.Username,
		Gender:   req.Gender,
		Avatar:   req.Avatar,
	})
	if err != nil {
		l.Logger.Errorf("UpdateProfile RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if updateProfileResp.Code != 200 {
		response.BusinessError(w, updateProfileResp.Code, updateProfileResp.Message)
		return
	}

	// 构造响应数据
	var userData *types.User
	if updateProfileResp.Data != nil {
		userData = &types.User{
			Id:        updateProfileResp.Data.Id,
			Username:  updateProfileResp.Data.Username,
			Email:     updateProfileResp.Data.Email,
			Gender:    updateProfileResp.Data.Gender,
			Avatar:    updateProfileResp.Data.Avatar,
			Status:    updateProfileResp.Data.Status,
			CreatedAt: updateProfileResp.Data.CreatedAt,
			UpdatedAt: updateProfileResp.Data.UpdatedAt,
		}
	}

	response.SuccessWithMessage(w, "Profile updated successfully", userData)
}
