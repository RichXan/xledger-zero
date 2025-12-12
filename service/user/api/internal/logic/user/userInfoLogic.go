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

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest, w http.ResponseWriter) {
	// 调用 RPC 服务获取用户信息
	userInfoResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
		Id: req.Id,
	})
	if err != nil {
		l.Logger.Errorf("UserInfo RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if userInfoResp.Code != 200 {
		response.BusinessError(w, userInfoResp.Code, userInfoResp.Message)
		return
	}

	// 构造响应数据
	var userData *types.User
	if userInfoResp.Data != nil {
		userData = &types.User{
			Id:        userInfoResp.Data.Id,
			Username:  userInfoResp.Data.Username,
			Email:     userInfoResp.Data.Email,
			Gender:    userInfoResp.Data.Gender,
			Avatar:    userInfoResp.Data.Avatar,
			Status:    userInfoResp.Data.Status,
			CreatedAt: userInfoResp.Data.CreatedAt,
			UpdatedAt: userInfoResp.Data.UpdatedAt,
		}
	}

	response.Success(w, userData)
}
