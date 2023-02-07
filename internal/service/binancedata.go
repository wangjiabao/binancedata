package service

import (
	v1 "binancedata/api/binancedata/v1"
	"binancedata/internal/biz"
	"context"
)

// BinanceDataService is a BinanceData service .
type BinanceDataService struct {
	v1.UnimplementedBinanceDataServer

	uc *biz.BinanceDataUsecase
}

// NewBinanceDataService new a BinanceData service.
func NewBinanceDataService(uc *biz.BinanceDataUsecase) *BinanceDataService {
	return &BinanceDataService{uc: uc}
}

func (b *BinanceDataService) PullBinanceData(ctx context.Context, req *v1.PullBinanceDataRequest) (*v1.PullBinanceDataReply, error) {
	return b.uc.PullBinanceData(ctx, req)
}

func (b *BinanceDataService) IntervalMAvgEndPriceData(ctx context.Context, req *v1.IntervalMAvgEndPriceDataRequest) (*v1.IntervalMAvgEndPriceDataReply, error) {
	return b.uc.IntervalMAvgEndPriceData(ctx, req)
}
