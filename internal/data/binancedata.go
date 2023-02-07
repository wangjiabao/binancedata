package data

import (
	"binancedata/internal/biz"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type KLineMOne struct {
	ID                  int64     `gorm:"primarykey;type:int"`
	StartTime           int64     `gorm:"type:bigint;not null"`
	EndTime             int64     `gorm:"type:bigint;not null"`
	StartPrice          float64   `gorm:"type:decimal(65,20);not null"`
	TopPrice            float64   `gorm:"type:decimal(65,20);not null"`
	LowPrice            float64   `gorm:"type:decimal(65,20);not null"`
	EndPrice            float64   `gorm:"type:decimal(65,20);not null"`
	DealTotalAmount     float64   `gorm:"type:decimal(65,20);not null"`
	DealAmount          float64   `gorm:"type:decimal(65,20);not null"`
	DealTotal           int64     `gorm:"type:int;not null"`
	DealSelfTotalAmount float64   `gorm:"type:decimal(65,20);not null"`
	DealSelfAmount      float64   `gorm:"type:decimal(65,20);not null"`
	CreatedAt           time.Time `gorm:"type:datetime;not null"`
	UpdatedAt           time.Time `gorm:"type:datetime;not null"`
}

type BinanceDataRepo struct {
	data *Data
	log  *log.Helper
}

type KLineMOneRepo struct {
	data *Data
	log  *log.Helper
}

func NewBinanceDataRepo(data *Data, logger log.Logger) biz.BinanceDataRepo {
	return &BinanceDataRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func NewKLineMOneRepo(data *Data, logger log.Logger) biz.KLineMOneRepo {
	return &KLineMOneRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (k *KLineMOneRepo) RequestBinanceMinuteKLinesData(symbol string, startTime string, endTime string, interval string, limit string) ([]*biz.KLineMOne, error) {
	apiUrl := "https://fapi.binance.com/fapi/v1/klines"
	// URL param
	data := url.Values{}
	data.Set("symbol", symbol)
	data.Set("interval", interval)
	data.Set("startTime", startTime)
	data.Set("endTime", endTime)
	data.Set("limit", limit)

	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = data.Encode() // URL encode
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	fmt.Println(u.String())
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var i [][]interface{}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}
	res := make([]*biz.KLineMOne, 0)
	for _, v := range i {
		startTimeTmp, _ := strconv.ParseInt(strconv.FormatFloat(v[0].(float64), 'f', -1, 64), 10, 64)
		endTimeTmp, _ := strconv.ParseInt(strconv.FormatFloat(v[6].(float64), 'f', -1, 64), 10, 64)
		dealTotalTmp, _ := strconv.ParseInt(strconv.FormatFloat(v[8].(float64), 'f', -1, 64), 10, 64)
		startPriceTmp, _ := strconv.ParseFloat(v[1].(string), 64)
		endPriceTmp, _ := strconv.ParseFloat(v[4].(string), 64)
		topPriceTmp, _ := strconv.ParseFloat(v[2].(string), 64)
		lowPriceTmp, _ := strconv.ParseFloat(v[3].(string), 64)
		dealTotalAmountTmp, _ := strconv.ParseFloat(v[5].(string), 64)
		dealAmountTmp, _ := strconv.ParseFloat(v[7].(string), 64)
		dealSelfTotalAmountTmp, _ := strconv.ParseFloat(v[9].(string), 64)
		dealSelfAmountTmp, _ := strconv.ParseFloat(v[10].(string), 64)

		res = append(res, &biz.KLineMOne{
			StartTime:           startTimeTmp,
			StartPrice:          startPriceTmp,
			EndPrice:            endPriceTmp,
			TopPrice:            topPriceTmp,
			LowPrice:            lowPriceTmp,
			EndTime:             endTimeTmp,
			DealTotalAmount:     dealTotalAmountTmp,
			DealAmount:          dealAmountTmp,
			DealTotal:           dealTotalTmp,
			DealSelfTotalAmount: dealSelfTotalAmountTmp,
			DealSelfAmount:      dealSelfAmountTmp,
		})
	}

	return res, err
}

// GetKLineMOneOrderByEndTimeLast .
func (k *KLineMOneRepo) GetKLineMOneOrderByEndTimeLast() (*biz.KLineMOne, error) {
	var kLineMOne KLineMOne
	if err := k.data.db.Order("end_time desc").Table("kline_m_one_btc_usdt").First(&kLineMOne).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("KLINE_M_ONE_NOT_FOUND", "kline m one not found")
		}

		return nil, errors.New(500, "KLINE M ONE ERROR", err.Error())
	}

	return &biz.KLineMOne{
		ID:                  kLineMOne.ID,
		StartTime:           kLineMOne.StartTime,
		EndTime:             kLineMOne.EndTime,
		StartPrice:          kLineMOne.StartPrice,
		TopPrice:            kLineMOne.TopPrice,
		LowPrice:            kLineMOne.LowPrice,
		EndPrice:            kLineMOne.EndPrice,
		DealTotalAmount:     kLineMOne.DealTotalAmount,
		DealAmount:          kLineMOne.DealAmount,
		DealTotal:           kLineMOne.DealTotal,
		DealSelfTotalAmount: kLineMOne.DealSelfTotalAmount,
		DealSelfAmount:      kLineMOne.DealSelfAmount,
	}, nil
}

// InsertKLineMOne .
func (k *KLineMOneRepo) InsertKLineMOne(ctx context.Context, kLineMOne []*biz.KLineMOne) (bool, error) {
	var insertKLineMOne []*KLineMOne
	for _, v := range kLineMOne {
		insertKLineMOne = append(insertKLineMOne, &KLineMOne{
			StartTime:           v.StartTime,
			EndTime:             v.EndTime,
			StartPrice:          v.StartPrice,
			TopPrice:            v.TopPrice,
			LowPrice:            v.LowPrice,
			EndPrice:            v.EndPrice,
			DealTotalAmount:     v.DealTotalAmount,
			DealAmount:          v.DealAmount,
			DealTotal:           v.DealTotal,
			DealSelfTotalAmount: v.DealSelfTotalAmount,
			DealSelfAmount:      v.DealSelfAmount,
		})
	}
	res := k.data.DB(ctx).Table("kline_m_one_btc_usdt").Create(&insertKLineMOne)
	if res.Error != nil {
		return false, errors.New(500, "CREATE_KLINE_M_ONE_BTC_USEDT_ERROR", "创建k线数据失败")
	}

	return true, nil
}
