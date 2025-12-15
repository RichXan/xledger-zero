package logic

import (
	"context"

	"xledger/service/ledger/rpc/internal/svc"
	"xledger/service/ledger/rpc/ledger"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLedgerStatisticsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLedgerStatisticsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLedgerStatisticsLogic {
	return &GetLedgerStatisticsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 统计操作
func (l *GetLedgerStatisticsLogic) GetLedgerStatistics(in *ledger.LedgerStatisticsRequest) (*ledger.LedgerStatisticsResponse, error) {
	// todo: add your logic here and delete this line

	return &ledger.LedgerStatisticsResponse{}, nil
}
