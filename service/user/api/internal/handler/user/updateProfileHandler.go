// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"xledger/service/user/api/internal/logic/user"
	"xledger/service/user/api/internal/svc"
	"xledger/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateProfileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateProfileLogic(r.Context(), svcCtx)
		l.UpdateProfile(&req, w)
	}
}
