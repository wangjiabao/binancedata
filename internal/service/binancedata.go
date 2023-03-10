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

func (b *BinanceDataService) XNIntervalMAvgEndPriceData(ctx context.Context, req *v1.XNIntervalMAvgEndPriceDataRequest) (*v1.XNIntervalMAvgEndPriceDataReply, error) {
	return b.uc.XNIntervalMAvgEndPriceData(ctx, req)
}

func (b *BinanceDataService) KAnd2NIntervalMAvgEndPriceData(ctx context.Context, req *v1.KAnd2NIntervalMAvgEndPriceDataRequest) (*v1.KAnd2NIntervalMAvgEndPriceDataReply, error) {
	return b.uc.KAnd2NIntervalMAvgEndPriceData(ctx, req)
}

func (b *BinanceDataService) IntervalMAvgEndPriceData(ctx context.Context, req *v1.IntervalMAvgEndPriceDataRequest) (*v1.IntervalMAvgEndPriceDataReply, error) {
	return b.uc.IntervalMAvgEndPriceData(ctx, req)
}

func (b *BinanceDataService) IntervalMMACDData(ctx context.Context, req *v1.IntervalMMACDDataRequest) (*v1.IntervalMMACDDataReply, error) {
	return b.uc.IntervalMMACDData(ctx, req)
}

func (b *BinanceDataService) IntervalMKAndMACDData(ctx context.Context, req *v1.IntervalMKAndMACDDataRequest) (*v1.IntervalMKAndMACDDataReply, error) {
	return b.uc.IntervalMKAndMACDData(ctx, req)
}

func (b *BinanceDataService) AreaPointIntervalMAvgEndPriceData(ctx context.Context, req *v1.AreaPointIntervalMAvgEndPriceDataRequest) (*v1.AreaPointIntervalMAvgEndPriceDataReply, error) {
	return b.uc.AreaPointIntervalMAvgEndPriceData(ctx, req)
}

func (b *BinanceDataService) IntervalMAvgEndPriceMacdAndAtrData(ctx context.Context, req *v1.IntervalMAvgEndPriceMacdAndAtrDataRequest) (*v1.IntervalMAvgEndPriceMacdAndAtrDataReply, error) {
	//return b.uc.IntervalMAvgEndPriceMacdAndAtrData(ctx, req)
	return &v1.IntervalMAvgEndPriceMacdAndAtrDataReply{}, nil
}

func (b *BinanceDataService) PullBinanceData(ctx context.Context, req *v1.PullBinanceDataRequest) (*v1.PullBinanceDataReply, error) {
	return b.uc.PullBinanceData(ctx, req)
}

func (b *BinanceDataService) Order(ctx context.Context, req *v1.OrderRequest) (*v1.OrderReply, error) {
	return b.uc.Order(ctx, req)
}
