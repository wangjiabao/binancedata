package data

import (
	"binancedata/internal/biz"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

type Order struct {
	ID            int64
	OrderId       string
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

func (k *KLineMOneRepo) RequestBinanceOrder(symbol string, side string, orderType string, positionSide string, quantity string) (*biz.Order, error) {
	var (
		client    *http.Client
		req       *http.Request
		resp      *http.Response
		b         []byte
		err       error
		apiUrl    = "https://fapi.binance.com//fapi/v1/order/test"
		apiKey    = "2eNaMVDIN4kdBVmSdZDkXyeucfwLBteLRwFSmUNHVuGhFs18AeVGDRZvfpTGDToX"
		secretKey = "w2xOINea6jMBJOqq9kWAvB0TWsKRWJdrM70wPbYeCMn2C1W89GxyBigbg1JSVojw"
	)
	// 时间
	now := strconv.FormatInt(time.Now().UTC().UnixMilli(), 10)
	// 拼请求数据
	data := "symbol=" + symbol + "&side=" + side + "&type=" + orderType + "&positionSide=" + positionSide + "&quantity=" + quantity + "&timestamp=" + now
	// 加密
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	signature := hex.EncodeToString(h.Sum(nil))
	// 构造请求

	req, err = http.NewRequest("POST", apiUrl, strings.NewReader(data+"&signature="+signature))
	if err != nil {
		return nil, err
	}
	// 添加头信息
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-MBX-APIKEY", apiKey)

	// 请求执行
	client = &http.Client{Timeout: 10 * time.Second}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	// 结果
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(b))

	var i Order
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}

	fmt.Println(i)
	return &biz.Order{
		ID:            0,
		OrderId:       "",
		ClientOrderId: "",
		Symbol:        "",
		Status:        "",
		OrigQty:       "",
		Side:          "",
		PositionSide:  "",
		OrderType:     "",
		OrderOrigType: "",
	}, err
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
		return false, errors.New(500, "CREATE_KLINE_M_ONE_BTC_USEDT_ERROR", "创建k线数据失败")
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
		return false, errors.New(500, "CREATE_KLINE_M_ONE_FIL_USEDT_ERROR", "创建k线数据失败")
	}

	return true, nil
}

// InsertEthKLineMOne .
func (k *KLineMOneRepo) InsertEthKLineMOne(ctx context.Context, kLineMOne []*biz.KLineMOne) (bool, error) {
	var insertKLineMOne []*KLineMOne
	for _, v := range kLineMOne {
		println(v)
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
		return false, errors.New(500, "CREATE_KLINE_M_ONE_ETH_USEDT_ERROR", "创建k线数据失败")
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
