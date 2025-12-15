// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LedgerRecordDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLedgerRecordDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LedgerRecordDetailLogic {
	return &LedgerRecordDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LedgerRecordDetailLogic) LedgerRecordDetail(req *types.LedgerRecordDetailRequest) (resp *types.LedgerRecordResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
