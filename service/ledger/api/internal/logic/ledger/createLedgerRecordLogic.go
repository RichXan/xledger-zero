// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package ledger

import (
	"context"

	"xledger/service/ledger/api/internal/svc"
	"xledger/service/ledger/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLedgerRecordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLedgerRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLedgerRecordLogic {
	return &CreateLedgerRecordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLedgerRecordLogic) CreateLedgerRecord(req *types.CreateLedgerRecordRequest) (resp *types.LedgerRecordResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
