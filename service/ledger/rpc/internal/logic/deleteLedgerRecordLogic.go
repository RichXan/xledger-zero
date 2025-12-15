package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLedgerRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLedgerRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLedgerRecordLogic {
	return &DeleteLedgerRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteLedgerRecordLogic) DeleteLedgerRecord(in *ledger.DeleteLedgerRecordRequest) (*ledger.DeleteLedgerRecordResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.DeleteLedgerRecordResponse{}, nil
}
