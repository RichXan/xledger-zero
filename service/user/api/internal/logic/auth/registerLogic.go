// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"net/http"

	"xledger/service/user/api/internal/response"
	"xledger/service/user/api/internal/svc"
	"xledger/service/user/api/internal/types"
	"xledger/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest, w http.ResponseWriter) {
	// 调用 RPC 服务注册用户
	registerResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("Register RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 根据 RPC 响应码判断
	if registerResp.Code != 200 {
		response.BusinessError(w, registerResp.Code, registerResp.Message)
		return
	}

	// 构造响应数据
	var userData *types.User
	if registerResp.Data != nil {
		userData = &types.User{
			Id:        registerResp.Data.Id,
			Username:  registerResp.Data.Username,
			Email:     registerResp.Data.Email,
			Gender:    registerResp.Data.Gender,
			Avatar:    registerResp.Data.Avatar,
			Status:    registerResp.Data.Status,
			CreatedAt: registerResp.Data.CreatedAt,
			UpdatedAt: registerResp.Data.UpdatedAt,
		}
	}

	response.SuccessWithMessage(w, "User registered successfully", userData)
}
