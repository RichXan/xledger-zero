package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLedgerRecordDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLedgerRecordDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLedgerRecordDetailLogic {
	return &GetLedgerRecordDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLedgerRecordDetailLogic) GetLedgerRecordDetail(in *ledger.LedgerRecordDetailRequest) (*ledger.CreateLedgerRecordResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.CreateLedgerRecordResponse{}, nil
}
