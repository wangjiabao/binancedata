package service

import (
	v1 "binancedata/api/binancedata/v1"
	"binancedata/internal/biz"
	"context"
)

// BinanceDataService is a BinanceData service .
type BinanceDataService struct {
	v1.UnimplementedBinancedataServer

	uc *biz.BinanceDataUsecase
}

// NewBinanceDataService new a BinanceData service.
func NewBinanceDataService(uc *biz.BinanceDataUsecase) *BinanceDataService {
	return &BinanceDataService{uc: uc}
}

func (b *BinanceDataService) DownloadBinancedata(ctx context.Context, req *v1.DownloadBinancedataRequest) (*v1.DownloadBinancedataReply, error) {
	return b.uc.DownloadBinancedata(ctx, req)
}
