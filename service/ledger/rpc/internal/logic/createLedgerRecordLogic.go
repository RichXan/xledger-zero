package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLedgerRecordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLedgerRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLedgerRecordLogic {
	return &CreateLedgerRecordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 记账记录操作
func (l *CreateLedgerRecordLogic) CreateLedgerRecord(in *ledger.CreateLedgerRecordRequest) (*ledger.CreateLedgerRecordResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.CreateLedgerRecordResponse{}, nil
}
