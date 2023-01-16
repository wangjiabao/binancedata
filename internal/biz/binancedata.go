package biz

import (
	v1 "binancedata/api/binancedata/v1"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type BinanceDataRepo interface {
}

// BinanceDataUsecase is a BinanceData usecase.
type BinanceDataUsecase struct {
	repo BinanceDataRepo
	log  *log.Helper
}

// NewBinanceDataUsecase new a BinanceData usecase.
func NewBinanceDataUsecase(repo BinanceDataRepo, logger log.Logger) *BinanceDataUsecase {
	return &BinanceDataUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (b *BinanceDataUsecase) DownloadBinancedata(ctx context.Context, req *v1.DownloadBinancedataRequest) (*v1.DownloadBinancedataReply, error) {

	return &v1.DownloadBinancedataReply{}, nil
}


