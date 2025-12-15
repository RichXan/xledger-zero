// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLedgerRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLedgerRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLedgerRecordLogic {
	return &UpdateLedgerRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLedgerRecordLogic) UpdateLedgerRecord(req *types.UpdateLedgerRecordRequest) (resp *types.LedgerRecordResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
