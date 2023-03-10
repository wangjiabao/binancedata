package data

import (
	"binancedata/internal/biz"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

type OrderPolicyPointCompare struct {
	ID        int64     `gorm:"primarykey;type:int"`
	InfoId    int64     `gorm:"type:int;not null"`
	Type      string    `gorm:"type:varchar(100)"`
	Value     int64     `gorm:"type:int;not null"`
	CreatedAt time.Time `gorm:"type:datetime;not null"`
	UpdatedAt time.Time `gorm:"type:datetime;not null"`
}

type OrderPolicyPointCompareInfo struct {
	ID        int64     `gorm:"primarykey;type:int"`
	OrderId   int64     `gorm:"type:bigint;not null"`
	Type      string    `gorm:"type:varchar(100)"`
	Num       float64   `gorm:"type:decimal(65,20);not null"`
	CreatedAt time.Time `gorm:"type:datetime;not null"`
	UpdatedAt time.Time `gorm:"type:datetime;not null"`
}

type Order struct {
	OrderId       int64
	ClientOrderId string
	Symbol        string
	Status        string
	OrigQty       string
	Side          string
	PositionSide  string
	OrderType     string
	OrderOrigType string
}

type BinanceDataRepo struct {
	data *Data
	log  *log.Helper
}

type KLineMOneRepo struct {
	data *Data
	log  *log.Helper
}

type OrderPolicyPointCompareRepo struct {
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

func NewOrderPolicyPointCompareRepoRepo(data *Data, logger log.Logger) biz.OrderPolicyPointCompareRepo {
	return &OrderPolicyPointCompareRepo{
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

	//fmt.Println(u.String())
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

func (o *OrderPolicyPointCompareRepo) RequestBinanceOrder(symbol string, side string, orderType string, positionSide string, quantity string, user int64) (*biz.Order, error) {
	var (
		client *http.Client
		req    *http.Request
		resp   *http.Response
		res    *biz.Order
		data   string
		b      []byte
		err    error
		apiUrl = "https://fapi.binance.com/fapi/v1/order"

		apiKey    = "2eNaMVDIN4kdBVmSdZDkXyeucfwLBteLRwFSmUNHVuGhFs18AeVGDRZvfpTGDToX"
		secretKey = "w2xOINea6jMBJOqq9kWAvB0TWsKRWJdrM70wPbYeCMn2C1W89GxyBigbg1JSVojw"
	)

	if 1 == user {
		apiKey = "MvzfRAnEeU46efaLYeaRms0r92d2g20iXVDQoJ8Ma5UvqH1zkJDMGB1WbSZ30P0W"
		secretKey = "bjGtZYExnHEcNBivXmJ8dLzGfMzr8SkW4ATmxLC1ZCrszbb5YJDulaiJLAgZ7L7h"
	}

	// ??????
	now := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	// ???????????????
	data = "symbol=" + symbol + "&side=" + side + "&type=" + orderType + "&positionSide=" + positionSide + "&quantity=" + quantity + "&timestamp=" + now

	// ??????
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))
	// ????????????

	req, err = http.NewRequest("POST", apiUrl, strings.NewReader(data+"&signature="+signature))
	if err != nil {
		return nil, err
	}
	// ???????????????
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-MBX-APIKEY", apiKey)

	// ????????????
	client = &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	// ??????
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		o.log.Error(err)
		return nil, err
	}

	var i Order
	err = json.Unmarshal(b, &i)
	if err != nil {
		o.log.Error(err)
		return nil, err
	}

	res = &biz.Order{
		ID:            0,
		OrderId:       i.OrderId,
		ClientOrderId: i.ClientOrderId,
		Symbol:        i.Symbol,
		Status:        i.Status,
		OrigQty:       i.OrigQty,
		Side:          i.Side,
		PositionSide:  i.PositionSide,
		OrderType:     i.OrderType,
		OrderOrigType: i.OrderOrigType,
	}

	o.log.Info(res)
	return res, nil
}

func (o *OrderPolicyPointCompareRepo) RequestBinanceGetOrder(symbol string) (*biz.Order, error) {

	return &biz.Order{}, nil
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

// GetFilKLineMOneOrderByEndTimeLast .
func (k *KLineMOneRepo) GetFilKLineMOneOrderByEndTimeLast() (*biz.KLineMOne, error) {
	var kLineMOne KLineMOne
	if err := k.data.db.Order("end_time desc").Table("kline_m_one_fil_usdt").First(&kLineMOne).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("KLINE_M_ONE_FIL_NOT_FOUND", "kline m one fil not found")
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

// GetEthKLineMOneOrderByEndTimeLast .
func (k *KLineMOneRepo) GetEthKLineMOneOrderByEndTimeLast() (*biz.KLineMOne, error) {
	var kLineMOne KLineMOne
	if err := k.data.db.Order("end_time desc").Table("kline_m_one_eth_usdt").First(&kLineMOne).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("KLINE_M_ONE_ETH_NOT_FOUND", "kline m one eth not found")
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
		return false, errors.New(500, "CREATE_KLINE_M_ONE_BTC_USEDT_ERROR", "??????k???????????????")
	}

	return true, nil
}

// InsertFilKLineMOne .
func (k *KLineMOneRepo) InsertFilKLineMOne(ctx context.Context, kLineMOne []*biz.KLineMOne) (bool, error) {
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
	res := k.data.DB(ctx).Table("kline_m_one_fil_usdt").Create(&insertKLineMOne)
	if res.Error != nil {
		return false, errors.New(500, "CREATE_KLINE_M_ONE_FIL_USEDT_ERROR", "??????k???????????????")
	}

	return true, nil
}

// InsertEthKLineMOne .
func (k *KLineMOneRepo) InsertEthKLineMOne(ctx context.Context, kLineMOne []*biz.KLineMOne) (bool, error) {
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
	res := k.data.DB(ctx).Table("kline_m_one_eth_usdt").Create(&insertKLineMOne)
	if res.Error != nil {
		return false, errors.New(500, "CREATE_KLINE_M_ONE_ETH_USEDT_ERROR", "??????k???????????????")
	}

	return true, nil
}

// GetKLineMOneBtcByStartTime .
func (k *KLineMOneRepo) GetKLineMOneBtcByStartTime(start int64, end int64) ([]*biz.KLineMOne, error) {
	var kLineMOnes []*KLineMOne
	// btc
	if err := k.data.db.Where("start_time>=? and start_time<=?", start, end).Table("kline_m_one_btc_usdt").Find(&kLineMOnes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("KLINE_M_ONE_NOT_FOUND", "kline m one not found")
		}

		return nil, errors.New(500, "KLINE M ONE ERROR", err.Error())
	}

	res := make([]*biz.KLineMOne, 0)
	for _, kLineMOne := range kLineMOnes {
		res = append(res, &biz.KLineMOne{
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
		})
	}

	return res, nil
}

// GetKLineMOneEthByStartTime .
func (k *KLineMOneRepo) GetKLineMOneEthByStartTime(start int64, end int64) ([]*biz.KLineMOne, error) {
	var kLineMOnes []*KLineMOne
	// btc
	if err := k.data.db.Where("start_time>=? and start_time<=?", start, end).Table("kline_m_one_eth_usdt").Find(&kLineMOnes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("KLINE_M_ONE_NOT_FOUND", "kline m one not found")
		}

		return nil, errors.New(500, "KLINE M ONE ERROR", err.Error())
	}

	res := make([]*biz.KLineMOne, 0)
	for _, kLineMOne := range kLineMOnes {
		res = append(res, &biz.KLineMOne{
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
		})
	}

	return res, nil
}

// GetKLineMOneFilByStartTime .
func (k *KLineMOneRepo) GetKLineMOneFilByStartTime(start int64, end int64) ([]*biz.KLineMOne, error) {
	var kLineMOnes []*KLineMOne
	// btc
	if err := k.data.db.Where("start_time>=? and start_time<=?", start, end).Table("kline_m_one_fil_usdt").Find(&kLineMOnes).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound("KLINE_M_ONE_NOT_FOUND", "kline m one not found")
		}

		return nil, errors.New(500, "KLINE M ONE ERROR", err.Error())
	}

	res := make([]*biz.KLineMOne, 0)
	for _, kLineMOne := range kLineMOnes {
		res = append(res, &biz.KLineMOne{
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
		})
	}

	return res, nil
}

func (k *KLineMOneRepo) NewMACDData(list []*biz.KLineMOne) ([]*biz.MACDPoint, error) {
	points := make([]*KLineMOne, 0)
	for _, v := range list {
		points = append(points, &KLineMOne{
			StartTime:  v.StartTime,
			EndTime:    v.EndTime,
			StartPrice: v.StartPrice,
			EndPrice:   v.EndPrice,
		})
	}
	return NewMACD(points).Calculation().GetPoints(), nil
}

// GetLastOrderPolicyPointCompareByInfoIdAndType .
func (o *OrderPolicyPointCompareRepo) GetLastOrderPolicyPointCompareByInfoIdAndType(infoId int64, policyPointType string, user int64) (*biz.OrderPolicyPointCompare, error) {
	var orderPolicyPointCompare *OrderPolicyPointCompare

	db := o.data.db.Where("info_id=? and type=?", infoId, policyPointType).Order("created_at desc")
	if 1 == user {
		db = db.Table("order_policy_point_compare_one")
	} else {
		db = db.Table("order_policy_point_compare")
	}

	if err := db.First(&orderPolicyPointCompare).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.New(500, "ORDER POLICY POINT COMPARE ERROR", err.Error())
	}

	return &biz.OrderPolicyPointCompare{
		ID:     orderPolicyPointCompare.ID,
		InfoId: orderPolicyPointCompare.InfoId,
		Type:   orderPolicyPointCompare.Type,
		Value:  orderPolicyPointCompare.Value,
	}, nil
}

// GetLastOrderPolicyPointCompareInfo .
func (o *OrderPolicyPointCompareRepo) GetLastOrderPolicyPointCompareInfo(user int64) (*biz.OrderPolicyPointCompareInfo, error) {
	var orderPolicyPointCompareInfo *OrderPolicyPointCompareInfo
	db := o.data.db.Order("created_at desc")

	if 1 == user {
		db = db.Table("order_policy_point_compare_info_one")
	} else {
		db = db.Table("order_policy_point_compare_info")
	}

	if err := db.First(&orderPolicyPointCompareInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, errors.New(500, "ORDER POLICY POINT COMPARE INFO ERROR", err.Error())
	}

	return &biz.OrderPolicyPointCompareInfo{
		ID:      orderPolicyPointCompareInfo.ID,
		Type:    orderPolicyPointCompareInfo.Type,
		Num:     orderPolicyPointCompareInfo.Num,
		OrderId: orderPolicyPointCompareInfo.OrderId,
	}, nil
}

// InsertOrderPolicyPointCompareInfo .
func (o *OrderPolicyPointCompareRepo) InsertOrderPolicyPointCompareInfo(ctx context.Context, orderPolicyPointCompareInfoData *biz.OrderPolicyPointCompareInfo, user int64) (*biz.OrderPolicyPointCompareInfo, error) {
	orderPolicyPointCompareInfo := &OrderPolicyPointCompareInfo{
		OrderId: orderPolicyPointCompareInfoData.OrderId,
		Type:    orderPolicyPointCompareInfoData.Type,
		Num:     orderPolicyPointCompareInfoData.Num,
	}

	db := o.data.DB(ctx)
	if 1 == user {
		db = db.Table("order_policy_point_compare_info_one")
	} else {
		db = db.Table("order_policy_point_compare_info")
	}

	res := db.Create(&orderPolicyPointCompareInfo)
	if res.Error != nil {
		return nil, errors.New(500, "CREATE_ORDER_POLICY_POINT_COMPARE_INFO_ERROR", "??????????????????")
	}

	return &biz.OrderPolicyPointCompareInfo{
		ID:      orderPolicyPointCompareInfo.ID,
		OrderId: orderPolicyPointCompareInfo.OrderId,
		Type:    orderPolicyPointCompareInfo.Type,
		Num:     orderPolicyPointCompareInfo.Num,
	}, nil
}

// InsertOrderPolicyPointCompare .
func (o *OrderPolicyPointCompareRepo) InsertOrderPolicyPointCompare(ctx context.Context, orderPolicyPointCompareData *biz.OrderPolicyPointCompare, user int64) (bool, error) {
	orderPolicyPointCompare := &OrderPolicyPointCompare{
		InfoId: orderPolicyPointCompareData.InfoId,
		Type:   orderPolicyPointCompareData.Type,
		Value:  orderPolicyPointCompareData.Value,
	}

	db := o.data.DB(ctx)
	if 1 == user {
		db = db.Table("order_policy_point_compare_one")
	} else {
		db = db.Table("order_policy_point_compare")
	}

	res := db.Create(&orderPolicyPointCompare)
	if res.Error != nil {
		return false, errors.New(500, "CREATE_ORDER_POLICY_POINT_COMPARE_ERROR", "??????????????????")
	}

	return true, nil
}
