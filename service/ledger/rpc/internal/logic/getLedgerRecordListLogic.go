package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLedgerRecordListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLedgerRecordListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLedgerRecordListLogic {
	return &GetLedgerRecordListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLedgerRecordListLogic) GetLedgerRecordList(in *ledger.LedgerRecordListRequest) (*ledger.LedgerRecordListResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.LedgerRecordListResponse{}, nil
}
