package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLedgerRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateLedgerRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLedgerRecordLogic {
	return &UpdateLedgerRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateLedgerRecordLogic) UpdateLedgerRecord(in *ledger.UpdateLedgerRecordRequest) (*ledger.UpdateLedgerRecordResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.UpdateLedgerRecordResponse{}, nil
}
