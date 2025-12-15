// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LedgerRecordListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLedgerRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LedgerRecordListLogic {
	return &LedgerRecordListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LedgerRecordListLogic) LedgerRecordList(req *types.LedgerRecordListRequest) (resp *types.LedgerRecordListResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
