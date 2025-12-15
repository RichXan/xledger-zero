// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"xledger/service/ledger/api/internal/logic/ledger"
	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"
)

func LedgerStatisticsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LedgerStatisticsRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := ledger.NewLedgerStatisticsLogic(r.Context(), svcCtx)
		resp, err := l.LedgerStatistics(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
