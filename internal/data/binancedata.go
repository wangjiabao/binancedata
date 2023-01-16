package data

import (
	"binancedata/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type BinanceDataRepo struct {
	data *Data
	log  *log.Helper
}

func NewBinanceDataRepo(data *Data, logger log.Logger) biz.BinanceDataRepo {
	return &BinanceDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
