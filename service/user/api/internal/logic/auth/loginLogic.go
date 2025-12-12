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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest, w http.ResponseWriter) {
	// 调用 RPC 服务登录
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		l.Logger.Errorf("Login RPC call failed: %v", err)
		response.Error(w, 500, "Internal server error")
		return
	}

	// 业务错误处理
	if loginResp.Code != 200 {
		response.BusinessError(w, loginResp.Code, loginResp.Message)
		return
	}

	// 构造响应数据
	var loginData *types.LoginData
	if loginResp.Data != nil {
		loginData = &types.LoginData{
			User: types.User{
				Id:        loginResp.Data.User.Id,
				Username:  loginResp.Data.User.Username,
				Email:     loginResp.Data.User.Email,
				Gender:    loginResp.Data.User.Gender,
				Avatar:    loginResp.Data.User.Avatar,
				Status:    loginResp.Data.User.Status,
				CreatedAt: loginResp.Data.User.CreatedAt,
				UpdatedAt: loginResp.Data.User.UpdatedAt,
			},
			AccessToken:  loginResp.Data.AccessToken,
			RefreshToken: loginResp.Data.RefreshToken,
			ExpiresIn:    loginResp.Data.ExpiresIn,
			TokenType:    loginResp.Data.TokenType,
		}
	}

	response.SuccessWithMessage(w, "Login successful", loginData)
}
