// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"xledger/service/ledger/api/internal/logic/ledger"
	"xledger/service/ledger/api/internal/svc"
)

func CategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := ledger.NewCategoryListLogic(r.Context(), svcCtx)
		resp, err := l.CategoryList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
