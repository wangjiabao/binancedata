package biz

import (
	v1 "binancedata/api/binancedata/v1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type BinanceData struct {
	StartTime           string
	StartPrice          string
	EndPrice            string
	TopPrice            string
	LowPrice            string
	EndTime             string
	DealTotalAmount     string
	DealAmount          string
	DealTotal           string
	DealSelfTotalAmount string
	DealSelfAmount      string
}

type OperationData struct {
	StartTime     string
	StartPrice    string
	EndPrice      string
	EndTime       string
	Time          string
	Type          string
	Status        string
	CloseEndPrice string
	Rate          float64
}

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

func (b *BinanceDataUsecase) DownloadBinanceData(ctx context.Context, req *v1.DownloadBinanceDataRequest) (*v1.DownloadBinanceDataReply, error) {

	return &v1.DownloadBinanceDataReply{}, nil
}

func (b *BinanceDataUsecase) IntervalMAvgEndPriceData(ctx context.Context, req *v1.IntervalMAvgEndPriceDataRequest) (*v1.IntervalMAvgEndPriceDataReply, error) {
	var (
		binanceData       []*BinanceData
		beforeBinanceData []*BinanceData
		operationData     []*OperationData
		reqStart          time.Time
		reqEnd            time.Time
		err               error
	)

	m := req.M
	n := req.N // 辅助前置数据数量 & 平均值的计算数量
	startTime := req.Start
	endTime := req.End

	reqStart, err = time.Parse("2006-01-02 15:04:05", startTime) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", endTime) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	//startTime := time.Date(2022, 1, 16, 18, 59, 59, 0, time.Local)
	//endTime := startTime.Add(999 * time.Minute) // 为了每分钟一条的数据，获取最大限制1000条

	limit := int64(1500)
	// 获取数据

	start := reqStart
	end := reqEnd
	fmt.Println(reqStart, reqEnd, start, end)
	for i := 1; i <= 100; i++ { // 每分钟最大请求次数 1200次，最大限度留一次给后边，这里目前够查询15w条
		var tmpBinanceData []*BinanceData
		if reqEnd.After(start.Add(time.Duration(m*limit) * time.Minute)) {
			end = start.Add(time.Duration(m*limit) * time.Minute).Add(-1 * time.Second)
		} else {
			end = reqEnd
		}

		tmpBinanceData, err = requestBinanceMinuteKLinesData(start.Add(-8*time.Hour).UnixMilli(), end.Add(-8*time.Hour).UnixMilli(), m, limit)
		if nil != err {
			return nil, err
		}
		binanceData = append(binanceData, tmpBinanceData...)

		start = start.Add(time.Duration(m*limit) * time.Minute)
		if reqEnd.Before(start) {
			break
		}
	}

	// 前置数据
	beforeEnd := reqStart.Add(-1 * time.Second)
	beforeStart := reqStart.Add(-time.Duration(m*(n-1)) * time.Minute)
	beforeBinanceData, err = requestBinanceMinuteKLinesData(beforeStart.Add(-8*time.Hour).UnixMilli(), beforeEnd.Add(-8*time.Hour).UnixMilli(), m, limit)
	if nil != err {
		return nil, err
	}

	res := &v1.IntervalMAvgEndPriceDataReply{
		Data:          make([]*v1.IntervalMAvgEndPriceDataReply_List, 0),
		OperationData: make([]*v1.IntervalMAvgEndPriceDataReply_List2, 0),
	}

	status := 1

	// 遍历数据
	for k, v := range binanceData {
		var (
			tmpEndPrice      float64
			tmpTotalEndPrice float64
			tmpAvgEndPrice   float64
		)
		if tmpEndPrice, err = strconv.ParseFloat(v.EndPrice, 64); nil != err {
			return nil, errors.New(500, "string to float64 error", "string to float64 error")
		}

		tmpTotalEndPrice = tmpEndPrice

		// 前置数据遍历参与计算
		if (k+1)-int(n) < 0 { // 游标位置判断是否需要前置数据
			for kBefore, vBefore := range beforeBinanceData {
				if k > kBefore { // 游标位置底线
					continue
				}
				var tmpBeforeEndPrice float64
				if tmpBeforeEndPrice, err = strconv.ParseFloat(vBefore.EndPrice, 64); nil != err {
					return nil, errors.New(500, "string to float64 error", "string to float64 error")
				}

				tmpTotalEndPrice += tmpBeforeEndPrice
			}
		}

		// 遍历数据参与计算
		if 0 < k {
			for kBak, vBak := range binanceData {
				if k == kBak { // 游标和到了和本身相等的值
					break
				}

				if kBak < (k+1)-int(n) { // 游标小跳过
					continue
				}

				var tmpBakEndPrice float64
				if tmpBakEndPrice, err = strconv.ParseFloat(vBak.EndPrice, 64); nil != err {
					return nil, errors.New(500, "string to float64 error", "string to float64 error")
				}

				tmpTotalEndPrice += tmpBakEndPrice
			}
		}

		tmpAvgEndPrice = tmpTotalEndPrice / float64(n)

		if 0 == k {
			if tmpEndPrice > tmpAvgEndPrice {
				status = 1
			} else {
				status = -1
			}
		}

		if tmpEndPrice > tmpAvgEndPrice && -1 == status {
			status = 1
			// 开空，有平多
			if 0 < len(operationData) {
				var (
					tmpOpenLastOperationData         *OperationData
					tmpCloseLastOperationData        *OperationData
					tmpOpenLastOperationDataEndPrice float64
				)
				// 开仓的信息
				tmpOpenLastOperationData = operationData[len(operationData)-1]

				if tmpOpenLastOperationDataEndPrice, err = strconv.ParseFloat(tmpOpenLastOperationData.EndPrice, 64); nil != err {
					return nil, errors.New(500, "string to float64 error", "string to float64 error")
				}

				rate := (tmpEndPrice - tmpOpenLastOperationDataEndPrice) / tmpOpenLastOperationDataEndPrice
				// 关
				tmpCloseLastOperationData = &OperationData{
					StartTime:     tmpOpenLastOperationData.StartTime,
					StartPrice:    tmpOpenLastOperationData.StartPrice,
					EndPrice:      tmpOpenLastOperationData.EndPrice,
					EndTime:       tmpOpenLastOperationData.EndTime,
					Time:          tmpOpenLastOperationData.EndTime,
					Type:          tmpOpenLastOperationData.Type,
					CloseEndPrice: v.EndPrice, // 关时候收盘价
					Status:        "close",
					Rate:          rate,
				}
				operationData = append(operationData, tmpCloseLastOperationData)
			}

			// 本次开
			currentOperationData := &OperationData{
				StartTime:  v.StartTime,
				StartPrice: v.StartPrice,
				EndPrice:   v.EndPrice, // 开时候收盘价
				EndTime:    v.EndTime,
				Time:       v.EndTime,
				Type:       "empty",
				Status:     "open",
			}
			operationData = append(operationData, currentOperationData)
		}

		if tmpEndPrice < tmpAvgEndPrice && 1 == status {
			// 开多。有平空
			status = -1

			if 0 < len(operationData) {
				var (
					tmpOpenLastOperationData         *OperationData
					tmpCloseLastOperationData        *OperationData
					tmpOpenLastOperationDataEndPrice float64
				)
				// 开仓的信息
				tmpOpenLastOperationData = operationData[len(operationData)-1]

				if tmpOpenLastOperationDataEndPrice, err = strconv.ParseFloat(tmpOpenLastOperationData.EndPrice, 64); nil != err {
					return nil, errors.New(500, "string to float64 error", "string to float64 error")
				}

				rate := (tmpOpenLastOperationDataEndPrice - tmpEndPrice) / tmpOpenLastOperationDataEndPrice
				// 关
				tmpCloseLastOperationData = &OperationData{
					StartTime:     tmpOpenLastOperationData.StartTime,
					StartPrice:    tmpOpenLastOperationData.StartPrice,
					EndPrice:      tmpOpenLastOperationData.EndPrice,
					EndTime:       tmpOpenLastOperationData.EndTime,
					Time:          tmpOpenLastOperationData.EndTime,
					Type:          tmpOpenLastOperationData.Type,
					CloseEndPrice: v.EndPrice, // 关时候收盘价
					Status:        "close",
					Rate:          rate,
				}
				operationData = append(operationData, tmpCloseLastOperationData)
			}

			// 本次开
			currentOperationData := &OperationData{
				StartTime:  v.StartTime,
				StartPrice: v.StartPrice,
				EndPrice:   v.EndPrice, // 开时候收盘价
				EndTime:    v.EndTime,
				Time:       v.EndTime,
				Type:       "more",
				Status:     "open",
			}
			operationData = append(operationData, currentOperationData)
		}

		res.Data = append(res.Data, &v1.IntervalMAvgEndPriceDataReply_List{
			StartPrice:            v.StartPrice,
			EndPrice:              v.EndPrice,
			TopPrice:              v.TopPrice,
			LowPrice:              v.LowPrice,
			Time:                  v.StartTime,
			WithBeforeAvgEndPrice: strconv.FormatFloat(tmpAvgEndPrice, 'f', -1, 64),
		})
	}

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
	)
	for _, vOperationData := range operationData {
		if "open" == vOperationData.Status {
			res.OperationOrderTotal++
		}

		if "close" == vOperationData.Status {
			tmpCloseTotal++
			if 0 < vOperationData.Rate {
				tmpWinTotal++
			}
		}

		tmpRate += vOperationData.Rate

		res.OperationData = append(res.OperationData, &v1.IntervalMAvgEndPriceDataReply_List2{
			StartPrice:    vOperationData.StartPrice,
			EndPrice:      vOperationData.EndPrice,
			StartTime:     vOperationData.StartTime,
			Time:          vOperationData.Time,
			EndTime:       vOperationData.EndTime,
			Type:          vOperationData.Type,
			Status:        vOperationData.Status,
			Rate:          strconv.FormatFloat(vOperationData.Rate, 'f', -1, 64),
			CloseEndPrice: vOperationData.CloseEndPrice,
		})
	}

	res.OperationWinRate = fmt.Sprintf("%.2f", float64(tmpWinTotal)/float64(tmpCloseTotal))
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	fmt.Println(len(binanceData), len(beforeBinanceData))
	return res, nil
}

func requestBinanceMinuteKLinesData(startTime int64, endTime int64, interval int64, limit int64) ([]*BinanceData, error) {
	fmt.Println(startTime, endTime)
	apiUrl := "https://fapi.binance.com/fapi/v1/klines"
	// URL param
	data := url.Values{}
	data.Set("symbol", "BTCUSDT")
	data.Set("interval", strconv.FormatInt(interval, 10)+"m")
	data.Set("startTime", strconv.FormatInt(startTime, 10))
	data.Set("endTime", strconv.FormatInt(endTime, 10))
	data.Set("limit", strconv.FormatInt(limit, 10))

	u, err := url.ParseRequestURI(apiUrl)
	if err != nil {
		return nil, err
	}
	u.RawQuery = data.Encode() // URL encode
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
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

	res := make([]*BinanceData, 0)
	for _, v := range i {
		tmp := &BinanceData{
			StartTime:           strconv.FormatFloat(v[0].(float64), 'f', -1, 64),
			StartPrice:          v[1].(string),
			EndPrice:            v[4].(string),
			TopPrice:            v[2].(string),
			LowPrice:            v[3].(string),
			EndTime:             strconv.FormatFloat(v[6].(float64), 'f', -1, 64),
			DealTotalAmount:     v[5].(string),
			DealAmount:          v[7].(string),
			DealTotal:           strconv.FormatFloat(v[8].(float64), 'f', -1, 64),
			DealSelfTotalAmount: v[9].(string),
			DealSelfAmount:      v[10].(string),
		}
		res = append(res, tmp)
	}

	return res, err
}
