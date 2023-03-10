package biz

import (
	v1 "binancedata/api/binancedata/v1"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"math"
	"sort"
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

type Order struct {
	ID            int64
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

type OperationData2 struct {
	StartTime      int64
	EndTime        int64
	StartPrice     float64
	TopPrice       float64
	LowPrice       float64
	EndPrice       float64
	AvgEndPrice    float64
	Amount         int64
	Type           string
	Status         string
	Action         string
	CloseEndPrice  string
	Rate           float64
	CloseSubPrice  float64
	ListMacdData   []*v1.IntervalMKAndMACDDataReply_List2_ListMacd
	ListMacd3Data  []*v1.IntervalMKAndMACDDataReply_List2_ListMacd3
	ListMacd60Data []*v1.IntervalMKAndMACDDataReply_List2_ListMacd60
}

type OperationData3 struct {
	StartTime      int64
	EndTime        int64
	StartPrice     float64
	TopPrice       float64
	LowPrice       float64
	EndPrice       float64
	Amount         int64
	Type           string
	Status         string
	Rate           float64
	LastMacd       float64
	Macd           float64
	LowKLowPrice   float64
	LowMacd        float64
	MaxKTopPrice   float64
	MaxMacd        float64
	Tag            int64
	CloseStatus    string
	ClosePriceWin  float64
	ClosePriceLost float64
}

type OperationData2Slice []*OperationData2

func (o OperationData2Slice) Len() int { return len(o) }

func (o OperationData2Slice) Less(i, j int) bool {
	if o[i].EndTime < o[j].EndTime {
		return true
	} else if o[i].EndTime == o[j].EndTime && "close" == o[i].Status {
		return true
	}
	return false
}

func (o OperationData2Slice) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

type KLineMOne struct {
	ID                  int64
	StartTime           int64
	EndTime             int64
	StartPrice          float64
	TopPrice            float64
	LowPrice            float64
	EndPrice            float64
	DealTotalAmount     float64
	DealAmount          float64
	DealTotal           int64
	DealSelfTotalAmount float64
	DealSelfAmount      float64
}

type OrderPolicyPointCompare struct {
	ID     int64
	InfoId int64
	Type   string
	Value  int64
}

type OrderPolicyPointCompareInfo struct {
	ID      int64
	OrderId int64
	Type    string
	Num     float64
}

type OrderData struct {
	Symbol          string
	Side            string
	OrderType       string
	PositionSide    string
	Quantity        string
	QuantityFloat64 float64
}

type MACDPoint struct {
	Time int64
	DIF  float64
	DEA  float64
	MACD float64
}

type Ma struct {
	AvgEndPrice float64
}

type BinanceDataRepo interface {
}

type OrderPolicyPointCompareRepo interface {
	RequestBinanceGetOrder(symbol string) (*Order, error)
	RequestBinanceOrder(symbol string, side string, orderType string, positionSide string, quantity string, user int64) (*Order, error)
	GetLastOrderPolicyPointCompareByInfoIdAndType(infoId int64, policyPointType string, user int64) (*OrderPolicyPointCompare, error)
	GetLastOrderPolicyPointCompareInfo(user int64) (*OrderPolicyPointCompareInfo, error)
	InsertOrderPolicyPointCompareInfo(ctx context.Context, orderPolicyPointCompareInfoData *OrderPolicyPointCompareInfo, user int64) (*OrderPolicyPointCompareInfo, error)
	InsertOrderPolicyPointCompare(ctx context.Context, orderPolicyPointCompareData *OrderPolicyPointCompare, user int64) (bool, error)
}

type KLineMOneRepo interface {
	GetKLineMOneOrderByEndTimeLast() (*KLineMOne, error)
	GetEthKLineMOneOrderByEndTimeLast() (*KLineMOne, error)
	GetFilKLineMOneOrderByEndTimeLast() (*KLineMOne, error)
	GetKLineMOneBtcByStartTime(start int64, end int64) ([]*KLineMOne, error)
	GetKLineMOneFilByStartTime(start int64, end int64) ([]*KLineMOne, error)
	GetKLineMOneEthByStartTime(start int64, end int64) ([]*KLineMOne, error)
	InsertKLineMOne(ctx context.Context, kLineMOne []*KLineMOne) (bool, error)
	InsertFilKLineMOne(ctx context.Context, kLineMOne []*KLineMOne) (bool, error)
	InsertEthKLineMOne(ctx context.Context, kLineMOne []*KLineMOne) (bool, error)
	RequestBinanceMinuteKLinesData(symbol string, startTime string, endTime string, interval string, limit string) ([]*KLineMOne, error)
	NewMACDData(list []*KLineMOne) ([]*MACDPoint, error)
}

// BinanceDataUsecase is a BinanceData usecase.
type BinanceDataUsecase struct {
	klineMOneRepo               KLineMOneRepo
	repo                        BinanceDataRepo
	orderPolicyPointCompareRepo OrderPolicyPointCompareRepo
	tx                          Transaction
	log                         *log.Helper
}

// NewBinanceDataUsecase new a BinanceData usecase.
func NewBinanceDataUsecase(repo BinanceDataRepo, klineMOneRepo KLineMOneRepo, orderPolicyPointCompareRepo OrderPolicyPointCompareRepo, tx Transaction, logger log.Logger) *BinanceDataUsecase {
	return &BinanceDataUsecase{repo: repo, klineMOneRepo: klineMOneRepo, orderPolicyPointCompareRepo: orderPolicyPointCompareRepo, tx: tx, log: log.NewHelper(logger)}
}

// XNIntervalMAvgEndPriceData x???<n?????????m????????????????????????>?????? .
func (b *BinanceDataUsecase) XNIntervalMAvgEndPriceData(ctx context.Context, req *v1.XNIntervalMAvgEndPriceDataRequest) (*v1.XNIntervalMAvgEndPriceDataReply, error) {
	fmt.Println(req)
	var (
		reqStart         time.Time
		reqEnd           time.Time
		klineMOne        []*KLineMOne
		n1               int
		n2               int
		ma5M15           []*Ma // todo ???????????????
		ma10M15          []*Ma
		ma5M5            []*Ma
		ma10M5           []*Ma
		ma5M60           []*Ma
		ma10M60          []*Ma
		operationData    map[string]*OperationData2
		resOperationData OperationData2Slice
		err              error
		//x         []*v1.XNIntervalMAvgEndPriceDataRequest_SendBody_List
	)

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.SendBody.Start) // ????????????????????????
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.SendBody.End) // ????????????????????????
	if nil != err {
		return nil, err
	}
	n1 = int(req.SendBody.N1)
	n2 = int(req.SendBody.N2)

	var (
		maxMxN = int64(n2 * 60) // todo ??????????????????????????????????????????????????????????????? 60???????????????????????????
	)

	// ??????????????????
	dataLimitTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		reqStart = dataLimitTime
	}

	//x = req.SendBody.X
	//for _, vX := range x {
	//	fmt.Println(vX.M, vX.Method, vX.N)
	//	if vX.N*vX.M > maxMxN { // n???m??????????????????n*m??????
	//		maxMxN = vX.N * vX.M
	//	}
	//}

	// ????????????????????????k???????????????
	//if 1 <= maxMxN {
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	//}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	res := &v1.XNIntervalMAvgEndPriceDataReply{
		DataListK:           make([]*v1.XNIntervalMAvgEndPriceDataReply_ListK, 0),
		DataListMa5M5:       make([]*v1.XNIntervalMAvgEndPriceDataReply_ListMa5M5, 0),
		DataListMa10M5:      make([]*v1.XNIntervalMAvgEndPriceDataReply_ListMa10M5, 0),
		DataListMa5M15:      make([]*v1.XNIntervalMAvgEndPriceDataReply_ListMa5M15, 0),
		DataListMa10M15:     make([]*v1.XNIntervalMAvgEndPriceDataReply_ListMa10M15, 0),
		DataListMa5M60:      make([]*v1.XNIntervalMAvgEndPriceDataReply_ListMa5M60, 0),
		DataListMa10M60:     make([]*v1.XNIntervalMAvgEndPriceDataReply_ListMa10M60, 0),
		OperationData:       make([]*v1.XNIntervalMAvgEndPriceDataReply_List2, 0),
		OperationOrderTotal: 0,
		OperationWinRate:    "",
		OperationWinAmount:  "",
	}

	// ??????k?????????????????????????????????

	var (
		openActionTag    string
		lastActionTag    string
		tmpLastActionTag string
	)
	operationData = make(map[string]*OperationData2, 0)
	for kKlineMOne, vKlineMOne := range klineMOne {
		var tagNum int64
		if maxMxN-1 > int64(kKlineMOne) { // ?????????????????????????????????????????????>=maxMxN????????????bug
			continue
		}

		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC)
		//fmt.Println(tmpNow, tmpNow0, tmpNow.Sub(tmpNow0).Minutes())
		tmpNowSubNow0 := int(tmpNow.Sub(tmpNow0).Minutes()) + 1 // ???????????????????????????????????????????????????????????????????????????????????????00???00-01???00??????????????????2?????????1??????
		// todo ???????????????
		// ???????????????
		//for _, vX := range x {
		//
		//}

		// ??????5???5?????????
		tmpMa5M5 := handleManMnWithKLineMineData(n1, 5, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma5M5 = append(ma5M5, tmpMa5M5)
		// ??????5???10?????????
		tmpMa10M5 := handleManMnWithKLineMineData(n2, 5, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma10M5 = append(ma10M5, tmpMa10M5)
		//// ??????5???15?????????
		tmpMa5M15 := handleManMnWithKLineMineData(n1, 15, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma5M15 = append(ma5M15, tmpMa5M15)
		//// ??????10???15?????????
		tmpMa10M15 := handleManMnWithKLineMineData(n2, 15, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma10M15 = append(ma10M15, tmpMa10M15)
		//// ??????5???60?????????
		tmpMa5M60 := handleManMnWithKLineMineData(n1, 60, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma5M60 = append(ma5M60, tmpMa5M60)
		//// ??????10???60?????????
		tmpMa10M60 := handleManMnWithKLineMineData(n2, 60, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma10M60 = append(ma10M60, tmpMa10M60)

		// ??? ???/???????????? ??? ???
		// ?????? ma5M15???ma5M15
		if maxMxN < int64(kKlineMOne) { // ???????????????
			//fmt.Println(kKlineMOne, vKlineMOne)
			//fmt.Println("ma10m15", tmpMa10M15)
			//fmt.Println("ma5m15", tmpMa5M15)
			//fmt.Println("ma5m5", tmpMa5M5)
			//fmt.Println("ma10m5", tmpMa10M5)
			//fmt.Println("ma5m60", tmpMa5M60)
			//fmt.Println("ma10m60", tmpMa10M60)

			// ?????? ???????????????????????????
			// ?????????
			if tmpMa5M15.AvgEndPrice > tmpMa10M15.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// ???????????????????????????????????????
					tmpDo := false // ????????????????????????????????????/?????????
					if "empty" == tmpOpenLastOperationData2.Type && ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) {
						tmpDo = true
					}
					if tmpDo {
						rate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
						tmpCloseLastOperationData := &OperationData2{
							StartTime:   vKlineMOne.StartTime,
							EndTime:     vKlineMOne.EndTime,
							StartPrice:  vKlineMOne.StartPrice,
							EndPrice:    vKlineMOne.EndPrice,
							AvgEndPrice: tmpMa5M15.AvgEndPrice,
							Amount:      0,
							Type:        "empty",
							Status:      "close",
							Rate:        rate,
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = tmpCloseLastOperationData
						openActionTag = ""
					}
				}

			}

			// ?????????
			if tmpMa5M15.AvgEndPrice < tmpMa10M15.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// ???????????????????????????????????????
					tmpDo := false // ????????????????????????????????????/?????????
					if "more" == tmpOpenLastOperationData2.Type && ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) {
						tmpDo = true
					}

					if tmpDo {
						rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
						tmpCloseLastOperationData := &OperationData2{
							StartTime:   vKlineMOne.StartTime,
							EndTime:     vKlineMOne.EndTime,
							StartPrice:  vKlineMOne.StartPrice,
							EndPrice:    vKlineMOne.EndPrice,
							AvgEndPrice: tmpMa5M15.AvgEndPrice,
							Amount:      0,
							Type:        "more",
							Status:      "close",
							Rate:        rate,
						}

						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = tmpCloseLastOperationData
						openActionTag = ""
					}
				}
			}

			// ??????
			//fmt.Println(len(ma5M15))
			if tmpMa5M15.AvgEndPrice > tmpMa10M15.AvgEndPrice { // ??????1
				lastMa5M15 := ma5M15[len(ma5M15)-2]    // ?????????
				lastMa10M15 := ma10M15[len(ma10M15)-2] // ?????????
				//last2Ma5M15 := ma5M15[len(ma5M15)-3]    // ?????????
				//last2Ma10M15 := ma10M15[len(ma10M15)-3] // ?????????
				//fmt.Println("last_ma10m15", tmpMa10M15)
				//fmt.Println("last_ma5m15", lastMa5M15)
				//fmt.Println("last2_ma5m15", last2Ma5M15)
				//fmt.Println("last2_ma10m15", last2Ma10M15)
				if lastMa5M15.AvgEndPrice < lastMa10M15.AvgEndPrice { // ??????1
					if 0 < klineMOne[kKlineMOne-1].EndPrice-klineMOne[kKlineMOne-15].StartPrice { //  ?????????15???????????????????????????
						//if tmpMa5M60.AvgEndPrice > tmpMa10M60.AvgEndPrice { // ??????2
						// ???????????????????????????????????????
						tmpDo := false
						if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
							if "close" == tmpOpenLastOperationData2.Status {
								tmpDo = true
							}
						} else {
							tmpDo = true
						}

						if tmpDo {
							//fmt.Println(vKlineMOne, time.UnixMilli(vKlineMOne.StartTime))
							//fmt.Println(vKlineMOne, lastMa5M15, lastMa10M15, tmpMa5M60, tmpMa10M60, last2Ma5M15, last2Ma10M15)
							currentOperationData := &OperationData2{
								StartTime:   vKlineMOne.StartTime,
								EndTime:     vKlineMOne.EndTime,
								StartPrice:  vKlineMOne.StartPrice,
								EndPrice:    vKlineMOne.EndPrice,
								AvgEndPrice: tmpMa5M15.AvgEndPrice,
								Amount:      2,
								Type:        "more",
								Status:      "open", // ????????????
							}
							//tmpOperationData = append(tmpOperationData, currentOperationData)

							tagNum++
							tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
							openActionTag = tmpLastActionTag
							//fmt.Println(openActionTag)
							operationData[tmpLastActionTag] = currentOperationData
						}

						//}
					}
				}
			}

			// ????????????
			if tmpMa5M5.AvgEndPrice < tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// ???????????????????????????????????????
					tmpDo := false // ????????????????????????????????????/?????????
					if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
						tmpDo = true
					}

					if tmpDo {
						rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
						tmpCloseLastOperationData := &OperationData2{
							StartTime:   vKlineMOne.StartTime,
							EndTime:     vKlineMOne.EndTime,
							StartPrice:  vKlineMOne.StartPrice,
							EndPrice:    vKlineMOne.EndPrice,
							AvgEndPrice: tmpMa5M15.AvgEndPrice,
							Amount:      tmpOpenLastOperationData2.Amount - int64(1),
							Type:        "more",
							Status:      "half",
							Rate:        rate,
						}
						//tmpOperationData = append(tmpOperationData, tmpCloseLastOperationData)
						tagNum++
						tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[tmpLastActionTag] = tmpCloseLastOperationData
					}

				}
			}

			// ?????????
			if tmpMa5M15.AvgEndPrice > tmpMa10M15.AvgEndPrice && tmpMa5M5.AvgEndPrice > tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// ???????????????????????????????????????
					tmpDo := false // ????????????????????????????????????/?????????
					if "more" == tmpOpenLastOperationData2.Type && "half" == tmpOpenLastOperationData2.Status {
						tmpDo = true
					}
					if tmpDo {
						// ?????????
						currentOperationData := &OperationData2{
							StartTime:   vKlineMOne.StartTime,
							EndTime:     vKlineMOne.EndTime,
							StartPrice:  vKlineMOne.StartPrice,
							EndPrice:    vKlineMOne.EndPrice,
							AvgEndPrice: tmpMa5M15.AvgEndPrice,
							Amount:      tmpOpenLastOperationData2.Amount + int64(1),
							Type:        "more",
							Action:      "add",
							Status:      "open",
						}
						//tmpOperationData = append(tmpOperationData, currentOperationData)
						tagNum++
						tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[tmpLastActionTag] = currentOperationData
					}
				}
			}

			// ??????
			if tmpMa5M15.AvgEndPrice < tmpMa10M15.AvgEndPrice { // ??????1
				lastMa5M15 := ma5M15[len(ma5M15)-2]    // ?????????
				lastMa10M15 := ma10M15[len(ma10M15)-2] // ?????????
				//last2Ma5M15 := ma5M15[len(ma5M15)-3]                  // ?????????
				//last2Ma10M15 := ma10M15[len(ma10M15)-3]               // ?????????
				if lastMa5M15.AvgEndPrice > lastMa10M15.AvgEndPrice { // ??????1
					if 0 > klineMOne[kKlineMOne-1].EndPrice-klineMOne[kKlineMOne-15].StartPrice { //  ?????????15???????????????????????????
						//if tmpMa5M60.AvgEndPrice < tmpMa10M60.AvgEndPrice { // ??????2
						// ???????????????????????????????????????
						tmpDo := false
						if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
							if "close" == tmpOpenLastOperationData2.Status {
								tmpDo = true
							}
						} else {
							tmpDo = true
						}
						// ?????????
						if tmpDo {
							currentOperationData := &OperationData2{
								StartTime:   vKlineMOne.StartTime,
								EndTime:     vKlineMOne.EndTime,
								StartPrice:  vKlineMOne.StartPrice,
								EndPrice:    vKlineMOne.EndPrice,
								AvgEndPrice: tmpMa5M15.AvgEndPrice,
								Amount:      2,
								Type:        "empty",
								Status:      "open",
							}

							tagNum++
							tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
							openActionTag = tmpLastActionTag
							operationData[tmpLastActionTag] = currentOperationData
						}
						//}
					}
				}
			}

			// ????????????
			if tmpMa5M5.AvgEndPrice > tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// ???????????????????????????????????????
					tmpDo := false // ????????????????????????????????????/?????????
					if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
						tmpDo = true
					}

					if tmpDo {
						rate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
						tmpCloseLastOperationData := &OperationData2{
							StartTime:   vKlineMOne.StartTime,
							EndTime:     vKlineMOne.EndTime,
							StartPrice:  vKlineMOne.StartPrice,
							EndPrice:    vKlineMOne.EndPrice,
							AvgEndPrice: tmpMa5M15.AvgEndPrice,
							Amount:      tmpOpenLastOperationData2.Amount - int64(1),
							Type:        "empty",
							Status:      "half",
							Rate:        rate,
						}
						tagNum++
						tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[tmpLastActionTag] = tmpCloseLastOperationData
					}
				}
			}

			// ?????????
			if tmpMa5M15.AvgEndPrice < tmpMa10M15.AvgEndPrice && tmpMa5M5.AvgEndPrice < tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// ???????????????????????????????????????
					tmpDo := false // ????????????????????????????????????/?????????
					if "empty" == tmpOpenLastOperationData2.Type && "half" == tmpOpenLastOperationData2.Status {
						tmpDo = true
					}

					if tmpDo {
						// ?????????
						currentOperationData := &OperationData2{
							StartTime:   vKlineMOne.StartTime,
							EndTime:     vKlineMOne.EndTime,
							StartPrice:  vKlineMOne.StartPrice,
							EndPrice:    vKlineMOne.EndPrice,
							AvgEndPrice: tmpMa5M15.AvgEndPrice,
							Amount:      tmpOpenLastOperationData2.Amount + int64(1),
							Type:        "empty",
							Action:      "add",
							Status:      "open",
						}
						tagNum++
						tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[tmpLastActionTag] = currentOperationData
					}
				}
			}

			if "" != tmpLastActionTag {
				lastActionTag = tmpLastActionTag
			}
			tmpLastActionTag = ""
		}

		res.DataListK = append(res.DataListK, &v1.XNIntervalMAvgEndPriceDataReply_ListK{
			StartPrice: vKlineMOne.StartPrice,
			EndPrice:   vKlineMOne.EndPrice,
			TopPrice:   vKlineMOne.TopPrice,
			LowPrice:   vKlineMOne.LowPrice,
			Time:       vKlineMOne.EndTime,
		})

	}

	for _, v := range ma5M15 {
		res.DataListMa5M15 = append(res.DataListMa5M15, &v1.XNIntervalMAvgEndPriceDataReply_ListMa5M15{
			AvgEndPrice: v.AvgEndPrice,
		})
	}

	for _, v := range ma10M15 {
		res.DataListMa10M15 = append(res.DataListMa10M15, &v1.XNIntervalMAvgEndPriceDataReply_ListMa10M15{
			AvgEndPrice: v.AvgEndPrice,
		})
	}

	for _, v := range ma5M5 {
		res.DataListMa5M5 = append(res.DataListMa5M5, &v1.XNIntervalMAvgEndPriceDataReply_ListMa5M5{
			AvgEndPrice: v.AvgEndPrice,
		})
	}

	for _, v := range ma10M5 {
		res.DataListMa10M5 = append(res.DataListMa10M5, &v1.XNIntervalMAvgEndPriceDataReply_ListMa10M5{
			AvgEndPrice: v.AvgEndPrice,
		})
	}

	for _, v := range ma5M60 {
		res.DataListMa5M60 = append(res.DataListMa5M60, &v1.XNIntervalMAvgEndPriceDataReply_ListMa5M60{
			AvgEndPrice: v.AvgEndPrice,
		})
	}

	for _, v := range ma10M60 {
		res.DataListMa10M60 = append(res.DataListMa10M60, &v1.XNIntervalMAvgEndPriceDataReply_ListMa10M60{
			AvgEndPrice: v.AvgEndPrice,
		})
	}

	// ??????
	for _, vOperationData := range operationData {
		resOperationData = append(resOperationData, vOperationData)
	}
	sort.Sort(resOperationData)

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
		tmpLastCloseK = -1
	)

	// ????????????????????????
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for k, vOperationData := range resOperationData {
		if k > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

		if "open" == vOperationData.Status {
			res.OperationOrderTotal++
		}

		if "close" == vOperationData.Status || "half" == vOperationData.Status {
			tmpCloseTotal++
			if 0 < vOperationData.Rate {
				tmpWinTotal++
			}
		}

		tmpRate += vOperationData.Rate

		res.OperationData = append(res.OperationData, &v1.XNIntervalMAvgEndPriceDataReply_List2{
			StartPrice: vOperationData.StartPrice,
			EndPrice:   vOperationData.EndPrice,
			StartTime:  vOperationData.StartTime,
			EndTime:    vOperationData.EndTime,
			Type:       vOperationData.Type,
			Action:     vOperationData.Action,
			Status:     vOperationData.Status,
			Rate:       vOperationData.Rate,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	return res, nil
}

// KAnd2NIntervalMAvgEndPriceData k??????x???<2?????????m????????????????????????>?????? .
func (b *BinanceDataUsecase) KAnd2NIntervalMAvgEndPriceData(ctx context.Context, req *v1.KAnd2NIntervalMAvgEndPriceDataRequest) (*v1.KAnd2NIntervalMAvgEndPriceDataReply, error) {
	fmt.Println(req)
	var (
		reqStart            time.Time
		reqEnd              time.Time
		klineMOne           []*KLineMOne
		n1                  int
		n2                  int
		m1                  int
		m2                  int
		topX                float64
		lowX                float64
		fee                 float64
		maNMFirst           []*Ma // todo ???????????????
		maNMSecond          []*Ma
		resOperationData    OperationData2Slice
		closeCondition      = 1
		closeCondition2Rate float64
		err                 error
		//x         []*v1.XNIntervalMAvgEndPriceDataRequest_SendBody_List
	)

	// ????????????
	res := &v1.KAnd2NIntervalMAvgEndPriceDataReply{
		DataListK:          make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListK, 0),
		DataListMaNMFirst:  make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMFirst, 0),
		DataListMaNMSecond: make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMSecond, 0),
		BackGround:         make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListBackGround, 0),
	}

	// ??????????????????????????????????????????????????????
	reqStart, err = time.Parse("2006-01-02 15:04:05", req.SendBody.Start) // ????????????????????????
	if nil != err {
		return res, nil
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.SendBody.End) // ????????????????????????
	if nil != err {
		return res, nil
	}

	n1 = int(req.SendBody.N1)
	m1 = int(req.SendBody.M1)
	n2 = int(req.SendBody.N2)
	m2 = int(req.SendBody.M2)
	topX = req.SendBody.TopX
	lowX = req.SendBody.LowX
	fee = req.SendBody.Fee
	maxMxN := int64(n2 * m2)
	if n1 > n2 || m1 > m2 {
		return res, nil
	}

	// ??????2????????????????????????
	if 2 == req.SendBody.CloseCondition {
		if maxMxN < 5 {
			return res, nil
		}
		closeCondition = 2
		closeCondition2Rate = req.SendBody.CloseCondition2Rate
	}

	// ????????????????????????k???????????????
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo ????????????????????????????????????maxMxN???????????????
	dataLimitTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// ?????????????????????
	if reqStart.After(reqEnd) {
		return res, nil
	}

	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	var (
		openActionTag   string
		lastActionTag   string
		compareTopPrice float64
		compareLowPrice float64
		operationData   map[string]*OperationData2
	)
	operationData = make(map[string]*OperationData2, 0)
	// ??????k?????????????????????????????????
	for kKlineMOne, vKlineMOne := range klineMOne {
		var tagNum int64
		if maxMxN-1 > int64(kKlineMOne) { // ?????????????????????????????????????????????>=maxMxN????????????bug
			continue
		}

		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC)
		//fmt.Println(tmpNow, tmpNow0, tmpNow.Sub(tmpNow0).Minutes())
		tmpNowSubNow0 := int(tmpNow.Sub(tmpNow0).Minutes()) + 1 // ???????????????????????????????????????????????????????????????????????????????????????00???00-01???00??????????????????2?????????1??????

		// ??????N???M?????????????????????
		tmpMaNMFirst := handleManMnWithKLineMineData(n1, m1, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		maNMFirst = append(maNMFirst, tmpMaNMFirst)
		res.DataListMaNMFirst = append(res.DataListMaNMFirst, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMFirst{X1: tmpMaNMFirst.AvgEndPrice})

		// ??????N???M?????????????????????
		tmpMaNMSecond := handleManMnWithKLineMineData(n2, m2, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		maNMSecond = append(maNMSecond, tmpMaNMSecond)
		res.DataListMaNMSecond = append(res.DataListMaNMSecond, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMSecond{X1: tmpMaNMSecond.AvgEndPrice})

		// ????????????
		tmpBackGround := "white"

		// ???????????????
		if maxMxN < int64(kKlineMOne) {
			// ??????
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				// ???????????????????????????????????????
				tmpDo := false // ????????????????????????????????????/?????????
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					tmpBackGround = "green"

					if 1 == closeCondition {
						// ????????????????????????
						if vKlineMOne.TopPrice > compareTopPrice {
							compareTopPrice = vKlineMOne.TopPrice
						}

						// ?????????-?????????*x% > ????????? ????????????
						if compareTopPrice-compareTopPrice*topX > vKlineMOne.EndPrice {
							tmpDo = true
						}
					} else if 2 == closeCondition {
						// ??????
						if vKlineMOne.EndPrice+tmpOpenLastOperationData2.CloseSubPrice*closeCondition2Rate <= tmpOpenLastOperationData2.EndPrice {
							tmpDo = true
						}

					}

				}
				if tmpDo {
					rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineMOne.StartTime,
						EndTime:    vKlineMOne.EndTime,
						StartPrice: vKlineMOne.StartPrice,
						EndPrice:   vKlineMOne.EndPrice,
						Amount:     0,
						Type:       "more",
						Status:     "close",
						Rate:       rate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
					compareTopPrice = 0
				}
			}

			// ??????
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				// ???????????????????????????????????????
				tmpDo := false // ????????????????????????????????????/?????????
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					tmpBackGround = "red"

					if 1 == closeCondition {
						// ????????????????????????
						if vKlineMOne.LowPrice < compareLowPrice {
							compareLowPrice = vKlineMOne.LowPrice
						}

						// ?????????-?????????*x% > ????????? ????????????
						if compareLowPrice+compareLowPrice*lowX < vKlineMOne.EndPrice {
							tmpDo = true
						}
					} else if 2 == closeCondition {
						// ??????
						if vKlineMOne.EndPrice-tmpOpenLastOperationData2.CloseSubPrice*closeCondition2Rate >= tmpOpenLastOperationData2.EndPrice {
							tmpDo = true
						}

					}
				}
				if tmpDo {
					rate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineMOne.StartTime,
						EndTime:    vKlineMOne.EndTime,
						StartPrice: vKlineMOne.StartPrice,
						EndPrice:   vKlineMOne.EndPrice,
						Amount:     0,
						Type:       "empty",
						Status:     "close",
						Rate:       rate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
					compareLowPrice = 0
				}
			}

			// ??????
			//fmt.Println(len(ma5M15))
			if tmpMaNMFirst.AvgEndPrice > tmpMaNMSecond.AvgEndPrice { // ??????1
				lastMaNMFirst := maNMFirst[len(maNMFirst)-2]                // ?????????
				lastMaNMSecond := maNMSecond[len(maNMSecond)-2]             // ?????????
				if lastMaNMFirst.AvgEndPrice < lastMaNMSecond.AvgEndPrice { // ??????1
					// ???????????????????????????????????????
					//tmpDo := false
					//if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					//	if "close" == tmpOpenLastOperationData2.Status {
					//		tmpDo = true
					//	}
					//} else {
					//	tmpDo = true
					//}
					//
					//if tmpDo {
					// ?????? ???????????????????????????
					if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
						if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
							rate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
							tmpCloseLastOperationData := &OperationData2{
								StartTime:  vKlineMOne.StartTime,
								EndTime:    vKlineMOne.EndTime,
								StartPrice: vKlineMOne.StartPrice,
								EndPrice:   vKlineMOne.EndPrice,
								Amount:     0,
								Type:       "empty",
								Status:     "close",
								Rate:       rate,
							}

							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
							operationData[lastActionTag] = tmpCloseLastOperationData
							openActionTag = ""
							compareLowPrice = 0
						}
					}

					// ?????????5k?????????????????????
					tmpFiveLowPrice := klineMOne[kKlineMOne-1].LowPrice
					for tmpI := 2; tmpI <= 5; tmpI++ {
						if tmpFiveLowPrice > klineMOne[kKlineMOne-tmpI].LowPrice {
							tmpFiveLowPrice = klineMOne[kKlineMOne-tmpI].LowPrice
						}
					}

					currentOperationData := &OperationData2{
						StartTime:     vKlineMOne.StartTime,
						EndTime:       vKlineMOne.EndTime,
						StartPrice:    vKlineMOne.StartPrice,
						EndPrice:      vKlineMOne.EndPrice,
						Amount:        2,
						Type:          "more",
						Status:        "open", // ????????????
						CloseSubPrice: math.Abs(vKlineMOne.StartPrice - tmpFiveLowPrice),
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					openActionTag = lastActionTag
					//fmt.Println(openActionTag)
					operationData[lastActionTag] = currentOperationData
					compareTopPrice = vKlineMOne.TopPrice

					tmpBackGround = "green"
				}
				//}
			}

			// ??????
			//fmt.Println(len(ma5M15))
			if tmpMaNMFirst.AvgEndPrice < tmpMaNMSecond.AvgEndPrice { // ??????1
				lastMaNMFirst := maNMFirst[len(maNMFirst)-2]                // ?????????
				lastMaNMSecond := maNMSecond[len(maNMSecond)-2]             // ?????????
				if lastMaNMFirst.AvgEndPrice > lastMaNMSecond.AvgEndPrice { // ??????1
					// ???????????????????????????????????????
					//tmpDo := false
					//if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					//	if "close" == tmpOpenLastOperationData2.Status {
					//		tmpDo = true
					//	}
					//} else {
					//	tmpDo = true
					//}

					//if tmpDo {
					// ?????? ???????????????????????????
					if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
						if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
							rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
							tmpCloseLastOperationData := &OperationData2{
								StartTime:  vKlineMOne.StartTime,
								EndTime:    vKlineMOne.EndTime,
								StartPrice: vKlineMOne.StartPrice,
								EndPrice:   vKlineMOne.EndPrice,
								Amount:     0,
								Type:       "more",
								Status:     "close",
								Rate:       rate,
							}

							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
							operationData[lastActionTag] = tmpCloseLastOperationData
							openActionTag = ""
							compareTopPrice = 0
						}
					}

					// ?????????5k?????????????????????
					tmpFiveTopPrice := klineMOne[kKlineMOne-1].TopPrice
					for tmpI := 2; tmpI <= 5; tmpI++ {
						if tmpFiveTopPrice < klineMOne[kKlineMOne-tmpI].TopPrice {
							tmpFiveTopPrice = klineMOne[kKlineMOne-tmpI].TopPrice
						}
					}

					currentOperationData := &OperationData2{
						StartTime:     vKlineMOne.StartTime,
						EndTime:       vKlineMOne.EndTime,
						StartPrice:    vKlineMOne.StartPrice,
						EndPrice:      vKlineMOne.EndPrice,
						Amount:        2,
						Type:          "empty",
						Status:        "open", // ????????????
						CloseSubPrice: math.Abs(vKlineMOne.StartPrice - tmpFiveTopPrice),
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					openActionTag = lastActionTag
					//fmt.Println(openActionTag)
					operationData[lastActionTag] = currentOperationData
					compareLowPrice = vKlineMOne.LowPrice

					tmpBackGround = "red"
				}

				//}
			}
		}
		res.BackGround = append(res.BackGround, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListBackGround{X1: tmpBackGround})

		// ??????
		res.DataListK = append(res.DataListK, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListK{
			X1: vKlineMOne.StartPrice,
			X2: vKlineMOne.EndPrice,
			X3: vKlineMOne.TopPrice,
			X4: vKlineMOne.LowPrice,
			X5: vKlineMOne.EndTime,
		})
	}

	// ??????
	for _, vOperationData := range operationData {
		resOperationData = append(resOperationData, vOperationData)
	}
	sort.Sort(resOperationData)

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
		tmpLastCloseK = -1
	)

	// ????????????????????????
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for k, vOperationData := range resOperationData {
		if k > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

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

		res.OperationData = append(res.OperationData, &v1.KAnd2NIntervalMAvgEndPriceDataReply_List2{
			StartPrice: vOperationData.StartPrice,
			EndPrice:   vOperationData.EndPrice,
			StartTime:  vOperationData.StartTime,
			EndTime:    vOperationData.EndTime,
			Type:       vOperationData.Type,
			Action:     vOperationData.Action,
			Status:     vOperationData.Status,
			Rate:       vOperationData.Rate,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)

	return res, nil
}

// IntervalMAvgEndPriceData k????????????m?????????????????????????????? .
func (b *BinanceDataUsecase) IntervalMAvgEndPriceData(ctx context.Context, req *v1.IntervalMAvgEndPriceDataRequest) (*v1.IntervalMAvgEndPriceDataReply, error) {
	var (
		resOperationData OperationData2Slice
		klineMOne        []*KLineMOne
		reqStart         time.Time
		reqEnd           time.Time
		m                int
		n                int
		fee              float64
		targetCloseRate  float64
		err              error
	)

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // ????????????????????????
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // ????????????????????????
	if nil != err {
		return nil, err
	}

	res := &v1.IntervalMAvgEndPriceDataReply{
		DataListK:     make([]*v1.IntervalMAvgEndPriceDataReply_ListK, 0),
		OperationData: make([]*v1.IntervalMAvgEndPriceDataReply_List2, 0),
	}

	m = int(req.M)
	n = int(req.N)
	maxMxN := n * m

	fee = req.Fee
	targetCloseRate = req.TargetCloseRate

	// ????????????????????????k???????????????
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo ????????????????????????????????????maxMxN???????????????
	dataLimitTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// ?????????????????????
	if reqStart.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	var (
		lastActionTag string
		operationData map[string]*OperationData2
		maNMFirst     []*Ma // todo ???????????????
	)
	operationData = make(map[string]*OperationData2, 0)
	// ????????????
	for kKlineMOne, vKlineMOne := range klineMOne {
		var tagNum int64
		if maxMxN-1 > kKlineMOne { // ?????????????????????????????????????????????>=maxMxN????????????bug
			continue
		}

		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC)
		//fmt.Println(tmpNow, tmpNow0, tmpNow.Sub(tmpNow0).Minutes())
		tmpNowSubNow0 := int(tmpNow.Sub(tmpNow0).Minutes()) + 1 // ???????????????????????????????????????????????????????????????????????????????????????00???00-01???00??????????????????2?????????1??????

		// ??????N???M?????????????????????
		tmpMaNMFirst := handleManMnWithKLineMineData(n, m, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		maNMFirst = append(maNMFirst, tmpMaNMFirst)
		res.DataListMaNMFirst = append(res.DataListMaNMFirst, &v1.IntervalMAvgEndPriceDataReply_ListMaNMFirst{X1: tmpMaNMFirst.AvgEndPrice})

		// ???????????????
		if maxMxN < kKlineMOne {
			// ????????????
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "open" == tmpOpenLastOperationData2.Status {

					if "empty" == tmpOpenLastOperationData2.Type {
						tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
						if tmpRate < -targetCloseRate {
							// ???
							tmpCloseLastOperationData := &OperationData2{
								StartTime:  vKlineMOne.StartTime,
								EndTime:    vKlineMOne.EndTime,
								StartPrice: vKlineMOne.StartPrice,
								EndPrice:   vKlineMOne.EndPrice,
								Amount:     0,
								Type:       "empty",
								Status:     "close",
								Rate:       tmpRate,
							}

							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
							operationData[lastActionTag] = tmpCloseLastOperationData
						}
					} else if "more" == tmpOpenLastOperationData2.Type {
						tmpRate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
						if tmpRate < -targetCloseRate {
							// ???
							tmpCloseLastOperationData := &OperationData2{
								StartTime:  vKlineMOne.StartTime,
								EndTime:    vKlineMOne.EndTime,
								StartPrice: vKlineMOne.StartPrice,
								EndPrice:   vKlineMOne.EndPrice,
								Amount:     0,
								Type:       "more",
								Status:     "close",
								Rate:       tmpRate,
							}

							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
							operationData[lastActionTag] = tmpCloseLastOperationData
						}
					}

				}
			}

			// ??????
			if vKlineMOne.EndPrice < tmpMaNMFirst.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
						rate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
						tmpCloseLastOperationData := &OperationData2{
							StartTime:  vKlineMOne.StartTime,
							EndTime:    vKlineMOne.EndTime,
							StartPrice: vKlineMOne.StartPrice,
							EndPrice:   vKlineMOne.EndPrice,
							Amount:     0,
							Type:       "empty",
							Status:     "close",
							Rate:       rate,
						}

						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = tmpCloseLastOperationData

						currentOperationData := &OperationData2{
							StartTime:  vKlineMOne.StartTime,
							EndTime:    vKlineMOne.EndTime,
							StartPrice: vKlineMOne.StartPrice,
							EndPrice:   vKlineMOne.EndPrice,
							Amount:     2,
							Type:       "more",
							Status:     "open", // ????????????
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = currentOperationData
					} else if "empty" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
						currentOperationData := &OperationData2{
							StartTime:  vKlineMOne.StartTime,
							EndTime:    vKlineMOne.EndTime,
							StartPrice: vKlineMOne.StartPrice,
							EndPrice:   vKlineMOne.EndPrice,
							Amount:     2,
							Type:       "more",
							Status:     "open", // ????????????
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = currentOperationData
					}

				} else {
					currentOperationData := &OperationData2{
						StartTime:  vKlineMOne.StartTime,
						EndTime:    vKlineMOne.EndTime,
						StartPrice: vKlineMOne.StartPrice,
						EndPrice:   vKlineMOne.EndPrice,
						Amount:     2,
						Type:       "more",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				}

			}

			// ??????
			if vKlineMOne.EndPrice > tmpMaNMFirst.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
						rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
						tmpCloseLastOperationData := &OperationData2{
							StartTime:  vKlineMOne.StartTime,
							EndTime:    vKlineMOne.EndTime,
							StartPrice: vKlineMOne.StartPrice,
							EndPrice:   vKlineMOne.EndPrice,
							Amount:     0,
							Type:       "more",
							Status:     "close",
							Rate:       rate,
						}

						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = tmpCloseLastOperationData

						currentOperationData := &OperationData2{
							StartTime:  vKlineMOne.StartTime,
							EndTime:    vKlineMOne.EndTime,
							StartPrice: vKlineMOne.StartPrice,
							EndPrice:   vKlineMOne.EndPrice,
							Amount:     2,
							Type:       "empty",
							Status:     "open", // ????????????
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = currentOperationData
					} else if "more" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
						currentOperationData := &OperationData2{
							StartTime:  vKlineMOne.StartTime,
							EndTime:    vKlineMOne.EndTime,
							StartPrice: vKlineMOne.StartPrice,
							EndPrice:   vKlineMOne.EndPrice,
							Amount:     2,
							Type:       "empty",
							Status:     "open", // ????????????
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
						operationData[lastActionTag] = currentOperationData
					}

				} else {
					currentOperationData := &OperationData2{
						StartTime:  vKlineMOne.StartTime,
						EndTime:    vKlineMOne.EndTime,
						StartPrice: vKlineMOne.StartPrice,
						EndPrice:   vKlineMOne.EndPrice,
						Amount:     2,
						Type:       "empty",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				}
			}
		}

		// ??????
		res.DataListK = append(res.DataListK, &v1.IntervalMAvgEndPriceDataReply_ListK{
			X1: vKlineMOne.StartPrice,
			X2: vKlineMOne.EndPrice,
			X3: vKlineMOne.TopPrice,
			X4: vKlineMOne.LowPrice,
			X5: vKlineMOne.EndTime,
		})
	}

	// ??????
	for _, vOperationData := range operationData {
		resOperationData = append(resOperationData, vOperationData)
	}
	sort.Sort(resOperationData)

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
		tmpLastCloseK = -1
	)

	// ????????????????????????
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for k, vOperationData := range resOperationData {
		if k > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

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
			StartPrice: vOperationData.StartPrice,
			EndPrice:   vOperationData.EndPrice,
			StartTime:  vOperationData.StartTime,
			EndTime:    vOperationData.EndTime,
			Type:       vOperationData.Type,
			Action:     vOperationData.Action,
			Status:     vOperationData.Status,
			Rate:       vOperationData.Rate,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	return res, nil
}

// IntervalMMACDData k????????????m?????????????????????????????? .
func (b *BinanceDataUsecase) IntervalMMACDData(ctx context.Context, req *v1.IntervalMMACDDataRequest) (*v1.IntervalMMACDDataReply, error) {
	var (
		resOperationData OperationData2Slice
		klineMOne        []*KLineMOne
		reqStart         time.Time
		reqEnd           time.Time
		m                int
		k                int
		//n                int
		err error
	)

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // ????????????????????????
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // ????????????????????????
	if nil != err {
		return nil, err
	}

	res := &v1.IntervalMMACDDataReply{
		DataListK:     make([]*v1.IntervalMMACDDataReply_ListK, 0),
		OperationData: make([]*v1.IntervalMMACDDataReply_List2, 0),
		DataListMacd:  make([]*v1.IntervalMMACDDataReply_ListMacd, 0),
	}

	m = int(req.M)
	k = int(req.K)
	//n = int(req.N)
	maxMxN := 200 * m // macd?????????????????????????????????????????????

	// ????????????????????????k???????????????
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo ????????????????????????????????????maxMxN???????????????
	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// ?????????????????????
	if reqStart.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	// ??????????????????????????????????????????
	klineM := handleMKData(klineMOne, m)

	//fmt.Println(len(klineM), klineM[199], len(tmpKlineMOne), tmpKlineMOne[199])

	var (
		macdData []*MACDPoint
	)
	operationData := make(map[string]*OperationData2, 0)

	// macd????????????
	macdData, err = b.klineMOneRepo.NewMACDData(klineM)
	if nil != err {
		return res, nil
	}
	macdDataMap := make(map[int64]*MACDPoint, 0)
	for _, vMacdData := range macdData[maxMxN/m-1:] {
		macdDataMap[vMacdData.Time] = &MACDPoint{
			Time: vMacdData.Time,
			DIF:  vMacdData.DIF,
			DEA:  vMacdData.DEA,
			MACD: vMacdData.MACD,
		}
		res.DataListMacd = append(res.DataListMacd, &v1.IntervalMMACDDataReply_ListMacd{
			X1: vMacdData.MACD,
			X2: vMacdData.DIF,
			X3: vMacdData.DEA,
			X4: vMacdData.Time,
		})
	}

	// ????????????
	var (
		lastActionTag string
		closeEmptyTag int
		closeMoreTag  int
	)
	tmpKlineM := klineM[maxMxN/m-1:]
	for _, vKlineM := range tmpKlineM {
		var tagNum int64

		// ??????
		res.DataListK = append(res.DataListK, &v1.IntervalMMACDDataReply_ListK{
			X1: vKlineM.StartPrice,
			X2: vKlineM.EndPrice,
			X3: vKlineM.TopPrice,
			X4: vKlineM.LowPrice,
			X5: vKlineM.EndTime,
			X6: vKlineM.StartTime,
		})

		if _, ok := macdDataMap[vKlineM.StartTime]; !ok {
			continue
		}

		// ??????
		if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
			if "open" == tmpOpenLastOperationData2.Status && "more" == tmpOpenLastOperationData2.Type {
				if vKlineM.StartPrice < vKlineM.EndPrice {
					closeMoreTag++
				}
				closeMore := false
				if k <= closeMoreTag {
					closeMore = true
				}
				if closeMore {
					tmpRate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     0,
						Type:       "more",
						Status:     "close",
						Rate:       tmpRate,
					}

					tagNum++
					closeMoreTag = 0
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// ??????
		if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
			if "open" == tmpOpenLastOperationData2.Status && "empty" == tmpOpenLastOperationData2.Type {
				if vKlineM.StartPrice > vKlineM.EndPrice {
					closeEmptyTag++
				}
				closeEmpty := false
				if k <= closeEmptyTag {
					closeEmpty = true
				}
				if closeEmpty {
					tmpRate := (tmpOpenLastOperationData2.EndPrice - vKlineM.EndPrice) / tmpOpenLastOperationData2.EndPrice
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     0,
						Type:       "empty",
						Status:     "close",
						Rate:       tmpRate,
					}

					tagNum++
					closeEmptyTag = 0
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// ??????
		if macdDataMap[vKlineM.StartTime].MACD > 0 {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					rate := (tmpOpenLastOperationData2.EndPrice - vKlineM.EndPrice) / tmpOpenLastOperationData2.EndPrice
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     0,
						Type:       "empty",
						Status:     "close",
						Rate:       rate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData

					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "more",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				} else if "empty" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "more",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:  vKlineM.StartTime,
					EndTime:    vKlineM.EndTime,
					StartPrice: vKlineM.StartPrice,
					EndPrice:   vKlineM.EndPrice,
					Amount:     2,
					Type:       "more",
					Status:     "open", // ????????????
				}
				tagNum++
				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				operationData[lastActionTag] = currentOperationData
			}
		}

		// ??????
		if macdDataMap[vKlineM.StartTime].MACD < 0 {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     0,
						Type:       "more",
						Status:     "close",
						Rate:       rate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData

					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "empty",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				} else if "more" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "empty",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:  vKlineM.StartTime,
					EndTime:    vKlineM.EndTime,
					StartPrice: vKlineM.StartPrice,
					EndPrice:   vKlineM.EndPrice,
					Amount:     2,
					Type:       "empty",
					Status:     "open", // ????????????
				}
				tagNum++
				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				operationData[lastActionTag] = currentOperationData
			}

		}

	}

	// ??????
	for _, vOperationData := range operationData {
		resOperationData = append(resOperationData, vOperationData)
	}
	sort.Sort(resOperationData)

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
		tmpLastCloseK = -1
	)

	// ????????????????????????
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for kOperationData, vOperationData := range resOperationData {
		if kOperationData > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

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

		res.OperationData = append(res.OperationData, &v1.IntervalMMACDDataReply_List2{
			StartPrice: vOperationData.StartPrice,
			EndPrice:   vOperationData.EndPrice,
			StartTime:  vOperationData.StartTime,
			EndTime:    vOperationData.EndTime,
			Type:       vOperationData.Type,
			Action:     vOperationData.Action,
			Status:     vOperationData.Status,
			Rate:       vOperationData.Rate,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	return res, nil
}

// IntervalMKAndMACDData k????????????m?????????????????????????????? .
func (b *BinanceDataUsecase) IntervalMKAndMACDData(ctx context.Context, req *v1.IntervalMKAndMACDDataRequest) (*v1.IntervalMKAndMACDDataReply, error) {
	var (
		resOperationData OperationData2Slice
		klineMOne        []*KLineMOne
		reqStart         time.Time
		reqEnd           time.Time
		m                int
		k                int
		//n                int
		err error
	)

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // ????????????????????????
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // ????????????????????????
	if nil != err {
		return nil, err
	}

	res := &v1.IntervalMKAndMACDDataReply{
		DataListK:     make([]*v1.IntervalMKAndMACDDataReply_ListK, 0),
		OperationData: make([]*v1.IntervalMKAndMACDDataReply_List2, 0),
	}

	m = int(req.M)
	k = int(req.K)
	//n = int(req.N)
	maxMxN := 201 * 60 // macd??????????????????????????????????????????????????????????????????60??????

	// ????????????????????????k???????????????
	startTime := reqStart.Add(-time.Duration(maxMxN) * time.Minute)
	// todo ????????????????????????????????????maxMxN???????????????
	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	if startTime.Before(dataLimitTime) {
		return res, nil
	}
	// ?????????????????????
	if startTime.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, startTime, reqEnd, startTime.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
		startTime.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	// ????????????
	var (
		lastActionTag    string
		tmpLastActionTag string
		openActionTag    string

		kLineDataMLive   []*KLineMOne
		kLineData3MLive  []*KLineMOne
		kLineData60MLive []*KLineMOne
		macdData         []*MACDPoint
		macdM3Data       []*MACDPoint
		macdM60Data      []*MACDPoint
	)
	operationData := make(map[string]*OperationData2, 0)
	macdDataLiveMap := make(map[int64]*MACDPoint, 0)
	macdM60DataLiveMap := make(map[int64]*MACDPoint, 0)
	//tmpKlineM := klineM[maxMxN/m-(k+1):]

	reqStartMilli := reqStart.Add(-8 * time.Hour).UnixMilli()
	for kKlineM, vKlineM := range klineMOne {
		// ????????????
		tmpNow := time.UnixMilli(vKlineM.StartTime).UTC().Add(8 * time.Hour)

		var (
			lastKeyMLive   int
			lastKey3MLive  int
			lastKey60MLive int
		)
		if 0 == tmpNow.Minute()%m {
			//len(kLineData15MLive)-1
			kLineDataMLive = append(kLineDataMLive, &KLineMOne{
				StartPrice: vKlineM.StartPrice,
				StartTime:  vKlineM.StartTime,
				TopPrice:   vKlineM.TopPrice,
				LowPrice:   vKlineM.LowPrice,
				EndPrice:   vKlineM.EndPrice,
				EndTime:    vKlineM.EndTime,
			})
			//if 1675180800000 <= vKlineM.StartTime {
			//	fmt.Println(vKlineM)
			//}
		} else {
			lastKeyMLive = len(kLineDataMLive) - 1
			if 0 <= lastKeyMLive {
				kLineDataMLive[lastKeyMLive].EndPrice = vKlineM.EndPrice
				kLineDataMLive[lastKeyMLive].EndTime = vKlineM.EndTime
				if kLineDataMLive[lastKeyMLive].TopPrice < vKlineM.TopPrice {
					kLineDataMLive[lastKeyMLive].TopPrice = vKlineM.TopPrice
				}
				if kLineDataMLive[lastKeyMLive].LowPrice > vKlineM.LowPrice {
					kLineDataMLive[lastKeyMLive].LowPrice = vKlineM.LowPrice
				}
			}
			//if 1675180800000 <= kLineData15MLive[lastKey].StartTime {
			//	fmt.Println(kLineData15MLive[lastKey], lastKey)
			//}
		}

		if 0 == tmpNow.Minute()%3 {
			//len(kLineData15MLive)-1
			kLineData3MLive = append(kLineData3MLive, &KLineMOne{
				StartPrice: vKlineM.StartPrice,
				StartTime:  vKlineM.StartTime,
				TopPrice:   vKlineM.TopPrice,
				LowPrice:   vKlineM.LowPrice,
				EndPrice:   vKlineM.EndPrice,
				EndTime:    vKlineM.EndTime,
			})
		} else {
			lastKey3MLive = len(kLineData3MLive) - 1
			if 0 <= lastKey3MLive {
				kLineData3MLive[lastKey3MLive].EndPrice = vKlineM.EndPrice
				kLineData3MLive[lastKey3MLive].EndTime = vKlineM.EndTime
				if kLineData3MLive[lastKey3MLive].TopPrice < vKlineM.TopPrice {
					kLineData3MLive[lastKey3MLive].TopPrice = vKlineM.TopPrice
				}
				if kLineData3MLive[lastKey3MLive].LowPrice > vKlineM.LowPrice {
					kLineData3MLive[lastKey3MLive].LowPrice = vKlineM.LowPrice
				}
			}
		}

		if 0 == tmpNow.Minute()%60 {
			//len(kLineData15MLive)-1
			kLineData60MLive = append(kLineData60MLive, &KLineMOne{
				StartPrice: vKlineM.StartPrice,
				StartTime:  vKlineM.StartTime,
				TopPrice:   vKlineM.TopPrice,
				LowPrice:   vKlineM.LowPrice,
				EndPrice:   vKlineM.EndPrice,
				EndTime:    vKlineM.EndTime,
			})
		} else {
			lastKey60MLive = len(kLineData60MLive) - 1
			if 0 <= lastKey60MLive {
				kLineData60MLive[lastKey60MLive].EndPrice = vKlineM.EndPrice
				kLineData60MLive[lastKey60MLive].EndTime = vKlineM.EndTime
				if kLineData60MLive[lastKey60MLive].TopPrice < vKlineM.TopPrice {
					kLineData60MLive[lastKey60MLive].TopPrice = vKlineM.TopPrice
				}
				if kLineData60MLive[lastKey60MLive].LowPrice > vKlineM.LowPrice {
					kLineData60MLive[lastKey60MLive].LowPrice = vKlineM.LowPrice
				}
			}
		}

		lastKeyMLive = len(kLineDataMLive) - 1
		lastKey3MLive = len(kLineData3MLive) - 1
		lastKey60MLive = len(kLineData60MLive) - 1

		// ???????????????????????????????????????
		if reqStartMilli > vKlineM.StartTime {

			// ???????????????????????????macd??????
			if reqStartMilli-int64(k*60000) <= vKlineM.StartTime {
				// macd????????????200???k???????????????
				macdData, err = b.klineMOneRepo.NewMACDData(kLineDataMLive[lastKeyMLive-199:])
				if nil != err {
					continue
				}
				macdDataLiveMap[vKlineM.StartTime] = macdData[199]

				// macd????????????200???60m???k???????????????
				macdM60Data, err = b.klineMOneRepo.NewMACDData(kLineData60MLive[lastKey60MLive-199:])
				if nil != err {
					continue
				}
				macdM60DataLiveMap[vKlineM.StartTime] = macdM60Data[199]
			}
			continue
		}

		//fmt.Println(vKlineM.StartTime, kLineDataMLive[lastKeyMLive], kLineData3MLive[lastKey3MLive], kLineData60MLive[lastKey60MLive])

		// macd????????????200???k???????????????
		macdData, err = b.klineMOneRepo.NewMACDData(kLineDataMLive[lastKeyMLive-199:])
		if nil != err {
			continue
		}
		macdDataLiveMap[vKlineM.StartTime] = macdData[199]

		// macd????????????200???3m???k???????????????
		macdM3Data, err = b.klineMOneRepo.NewMACDData(kLineData3MLive[lastKey3MLive-199:])
		if nil != err {
			continue
		}

		// macd????????????200???60m???k???????????????
		macdM60Data, err = b.klineMOneRepo.NewMACDData(kLineData60MLive[lastKey60MLive-199:])
		if nil != err {
			continue
		}
		macdM60DataLiveMap[vKlineM.StartTime] = macdM60Data[199]

		//fmt.Println(macdData[199], macdM3Data[199], macdM60Data[199])

		var tagNum int64

		// ??????

		tmpResDataListK := &v1.IntervalMKAndMACDDataReply_ListK{
			X1: kLineDataMLive[lastKeyMLive].StartPrice,
			X2: kLineDataMLive[lastKeyMLive].EndPrice,
			X3: kLineDataMLive[lastKeyMLive].TopPrice,
			X4: kLineDataMLive[lastKeyMLive].LowPrice,
			X5: kLineDataMLive[lastKeyMLive].EndTime,
			X6: kLineDataMLive[lastKeyMLive].StartTime,

			X151: macdData[199].DIF,
			X152: macdData[199].DEA,
			X153: macdData[199].MACD,
			X154: macdData[199].Time,
		}

		lastResDataListKKey := len(res.DataListK) - 1
		if 0 == tmpNow.Minute()%m {
			res.DataListK = append(res.DataListK, tmpResDataListK)
		} else {
			if 0 <= lastResDataListKKey {
				res.DataListK[lastResDataListKKey] = tmpResDataListK
			}
		}

		tmpResDataListK3 := &v1.IntervalMKAndMACDDataReply_ListKMacd3{
			X1: kLineData3MLive[lastKey3MLive].StartPrice,
			X2: kLineData3MLive[lastKey3MLive].EndPrice,
			X3: kLineData3MLive[lastKey3MLive].TopPrice,
			X4: kLineData3MLive[lastKey3MLive].LowPrice,
			X5: kLineData3MLive[lastKey3MLive].EndTime,
			X6: kLineData3MLive[lastKey3MLive].StartTime,

			X31: macdM3Data[199].DIF,
			X32: macdM3Data[199].DEA,
			X33: macdM3Data[199].MACD,
			X34: macdM3Data[199].Time,
		}

		lastResDataListKMacd3Key := len(res.DataListKMacd3) - 1
		if 0 == tmpNow.Minute()%3 {
			res.DataListKMacd3 = append(res.DataListKMacd3, tmpResDataListK3)
		} else {
			if 0 <= lastResDataListKMacd3Key {
				res.DataListKMacd3[lastResDataListKMacd3Key] = tmpResDataListK3
			}
		}

		tmpResDataListK60 := &v1.IntervalMKAndMACDDataReply_ListKMacd60{
			X1: kLineData60MLive[lastKey60MLive].StartPrice,
			X2: kLineData60MLive[lastKey60MLive].EndPrice,
			X3: kLineData60MLive[lastKey60MLive].TopPrice,
			X4: kLineData60MLive[lastKey60MLive].LowPrice,
			X5: kLineData60MLive[lastKey60MLive].EndTime,
			X6: kLineData60MLive[lastKey60MLive].StartTime,

			X601: macdM60Data[199].DIF,
			X602: macdM60Data[199].DEA,
			X603: macdM60Data[199].MACD,
			X604: macdM60Data[199].Time,
		}

		lastResDataListKMacd60Key := len(res.DataListKMacd60) - 1
		if 0 == tmpNow.Minute()%60 {
			res.DataListKMacd60 = append(res.DataListKMacd60, tmpResDataListK60)
		} else {
			if 0 <= lastResDataListKMacd60Key {
				res.DataListKMacd60[lastResDataListKMacd60Key] = tmpResDataListK60
			}
		}

		var (
			openMoreOne  int
			openMoreTwo  int
			openEmptyOne int
			openEmptyTwo int

			tmpListMacdData3  []*v1.IntervalMKAndMACDDataReply_List2_ListMacd3
			tmpListMacdData60 []*v1.IntervalMKAndMACDDataReply_List2_ListMacd60
			tmpListMacdData   []*v1.IntervalMKAndMACDDataReply_List2_ListMacd
		)

		// ????????????macd???????????????????????????
		tmpListMacdData3 = append(tmpListMacdData3, &v1.IntervalMKAndMACDDataReply_List2_ListMacd3{
			X31: macdM3Data[199].DIF,
			X32: macdM3Data[199].DEA,
			X33: macdM3Data[199].MACD,
			X34: macdM3Data[199].Time,
		})

		// ????????????macd???????????????????????????
		tmpListMacdData = append(tmpListMacdData, &v1.IntervalMKAndMACDDataReply_List2_ListMacd{
			X31: macdData[199].DIF,
			X32: macdData[199].DEA,
			X33: macdData[199].MACD,
			X34: macdData[199].Time,
		})

		// ????????????
		for i := 1; i <= k; i++ {
			tmpMacdData := macdDataLiveMap[klineMOne[kKlineM-i].StartTime]
			//if 1675683900000 == vKlineM.StartTime {
			//	fmt.Println(tmpMacdData, 3, i)
			//}
			if tmpMacdData.DEA > tmpMacdData.DIF &&
				tmpMacdData.DIF > 0 {
				openMoreOne += 1
			}

			if tmpMacdData.DEA < tmpMacdData.DIF &&
				tmpMacdData.DIF < 0 {
				openEmptyOne += 1
			}

			// ????????????macd???????????????????????????
			tmpListMacdData = append(tmpListMacdData, &v1.IntervalMKAndMACDDataReply_List2_ListMacd{
				X31: tmpMacdData.DIF,
				X32: tmpMacdData.DEA,
				X33: tmpMacdData.MACD,
				X34: tmpMacdData.Time,
			})
		}

		for i := 1; i <= k; i++ {
			tmpMacdData60 := macdM60DataLiveMap[klineMOne[kKlineM-i].StartTime]
			//if 1675683900000 == vKlineM.StartTime {
			//	fmt.Println(tmpMacdData60, 2, i)
			//}
			// 60??????
			if tmpMacdData60.DIF > tmpMacdData60.DEA &&
				tmpMacdData60.DEA > 0 {
				openMoreTwo += 1
			}

			if tmpMacdData60.DIF < tmpMacdData60.DEA &&
				tmpMacdData60.DEA < 0 {
				openEmptyTwo += 1
			}

			// ????????????macd???????????????????????????
			tmpListMacdData60 = append(tmpListMacdData60, &v1.IntervalMKAndMACDDataReply_List2_ListMacd60{
				X31: tmpMacdData60.DIF,
				X32: tmpMacdData60.DEA,
				X33: tmpMacdData60.MACD,
				X34: tmpMacdData60.Time,
			})
		}

		// ?????????
		if macdM3Data[199].DIF < 0 {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "more" == tmpOpenLastOperationData2.Type {

					tmpRate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "more",
						Status:         "close",
						Rate:           tmpRate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}
		if macdData[199].DIF < macdData[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "more" == tmpOpenLastOperationData2.Type {
					tmpRate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "more",
						Status:         "close",
						Rate:           tmpRate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}

		// ?????????
		if macdM3Data[199].DIF > 0 {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "empty" == tmpOpenLastOperationData2.Type {
					tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "empty",
						Status:         "close",
						Rate:           tmpRate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}

		if macdData[199].DIF > macdData[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "empty" == tmpOpenLastOperationData2.Type {
					tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "empty",
						Status:         "close",
						Rate:           tmpRate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}

		// ????????????
		if macdM3Data[199].DIF < macdM3Data[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "open" == tmpOpenLastOperationData2.Status && "more" == tmpOpenLastOperationData2.Type {
					tmpRate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "more",
						Status:         "half",
						Rate:           tmpRate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// ????????????
		if macdM3Data[199].DIF > macdM3Data[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "open" == tmpOpenLastOperationData2.Status && "empty" == tmpOpenLastOperationData2.Type {
					tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// ???
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "empty",
						Status:         "half",
						Rate:           tmpRate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// ????????????
		if macdM3Data[199].DIF > macdM3Data[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "half" == tmpOpenLastOperationData2.Status && "more" == tmpOpenLastOperationData2.Type {
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "more",
						Action:         "add",
						Status:         "open",
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// ??????
		//if 1675683900000 == vKlineM.StartTime {
		//	fmt.Println(macdData[199])
		//}
		if openMoreOne >= k && openMoreTwo >= k && macdData[199].DIF > macdData[199].DEA {
			//fmt.Println(openMoreOne, openMoreTwo, k, macdData[199].DIF, macdData[199].DEA)
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         0,
						Type:           "empty",
						Status:         "close",
						Rate:           rate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData

					currentOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         2,
						Type:           "more",
						Status:         "open", // ????????????
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				} else if "empty" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
					currentOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         2,
						Type:           "more",
						Status:         "open", // ????????????
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:      vKlineM.StartTime,
					EndTime:        vKlineM.EndTime,
					StartPrice:     vKlineM.StartPrice,
					EndPrice:       vKlineM.EndPrice,
					Amount:         2,
					Type:           "more",
					Status:         "open", // ????????????
					ListMacd3Data:  tmpListMacdData3,
					ListMacdData:   tmpListMacdData,
					ListMacd60Data: tmpListMacdData60,
				}
				tagNum++
				tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				openActionTag = tmpLastActionTag
				operationData[tmpLastActionTag] = currentOperationData
			}
		}

		// ????????????
		if macdM3Data[199].DIF < macdM3Data[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "half" == tmpOpenLastOperationData2.Status && "empty" == tmpOpenLastOperationData2.Type {
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         tmpOpenLastOperationData2.Amount - int64(1),
						Type:           "empty",
						Action:         "add",
						Status:         "open",
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// ??????
		if openEmptyOne >= k && openEmptyTwo >= k && macdData[199].DIF < macdData[199].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
					tmpCloseLastOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         0,
						Type:           "more",
						Status:         "close",
						Rate:           rate,
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData

					currentOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         2,
						Type:           "empty",
						Status:         "open", // ????????????
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				} else if "more" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
					currentOperationData := &OperationData2{
						StartTime:      vKlineM.StartTime,
						EndTime:        vKlineM.EndTime,
						StartPrice:     vKlineM.StartPrice,
						EndPrice:       vKlineM.EndPrice,
						Amount:         2,
						Type:           "empty",
						Status:         "open", // ????????????
						ListMacd3Data:  tmpListMacdData3,
						ListMacdData:   tmpListMacdData,
						ListMacd60Data: tmpListMacdData60,
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:      vKlineM.StartTime,
					EndTime:        vKlineM.EndTime,
					StartPrice:     vKlineM.StartPrice,
					EndPrice:       vKlineM.EndPrice,
					Amount:         2,
					Type:           "empty",
					Status:         "open", // ????????????
					ListMacd3Data:  tmpListMacdData3,
					ListMacdData:   tmpListMacdData,
					ListMacd60Data: tmpListMacdData60,
				}
				tagNum++
				tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				openActionTag = tmpLastActionTag
				operationData[tmpLastActionTag] = currentOperationData
			}
		}

		if "" != tmpLastActionTag {
			lastActionTag = tmpLastActionTag
		}
		tmpLastActionTag = ""

	}

	// ??????
	for _, vOperationData := range operationData {
		resOperationData = append(resOperationData, vOperationData)
	}
	sort.Sort(resOperationData)

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
		tmpLastCloseK = -1
	)

	// ????????????????????????
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for kOperationData, vOperationData := range resOperationData {
		if kOperationData > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

		if k > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

		if "open" == vOperationData.Status {
			res.OperationOrderTotal++
		}

		if "close" == vOperationData.Status || "half" == vOperationData.Status {
			tmpCloseTotal++
			if 0 < vOperationData.Rate {
				tmpWinTotal++
			}
		}

		tmpRate += vOperationData.Rate

		res.OperationData = append(res.OperationData, &v1.IntervalMKAndMACDDataReply_List2{
			StartPrice: vOperationData.StartPrice,
			EndPrice:   vOperationData.EndPrice,
			StartTime:  vOperationData.StartTime,
			EndTime:    vOperationData.EndTime,
			Type:       vOperationData.Type,
			Action:     vOperationData.Action,
			Status:     vOperationData.Status,
			Rate:       vOperationData.Rate,
			MacdData:   vOperationData.ListMacdData,
			Macd3Data:  vOperationData.ListMacd3Data,
			Macd60Data: vOperationData.ListMacd60Data,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	return res, nil
}

// AreaPointIntervalMAvgEndPriceData k????????????m?????????????????????????????? .
func (b *BinanceDataUsecase) AreaPointIntervalMAvgEndPriceData(ctx context.Context, req *v1.AreaPointIntervalMAvgEndPriceDataRequest) (*v1.AreaPointIntervalMAvgEndPriceDataReply, error) {
	var (
		resOperationData OperationData2Slice
		klineMOne        []*KLineMOne
		reqStart         time.Time
		reqEnd           time.Time
		m                int
		n                int
		pointFirst       float64
		pointInterval    float64
		err              error
	)

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // ????????????????????????
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // ????????????????????????
	if nil != err {
		return nil, err
	}

	res := &v1.AreaPointIntervalMAvgEndPriceDataReply{
		DataListK:         make([]*v1.AreaPointIntervalMAvgEndPriceDataReply_ListK, 0),
		DataListMaNMFirst: make([]*v1.AreaPointIntervalMAvgEndPriceDataReply_ListMaNMFirst, 0),
		DataListSubPoint:  make([]*v1.AreaPointIntervalMAvgEndPriceDataReply_ListSubPoint, 0),
	}

	m = int(req.M)
	n = int(req.N)
	pointFirst = req.PointFirst
	pointInterval = req.PointInterval

	maxMxN := m*n + 2*n // macd??????????????????????????????????????????????????????????????????60??????

	// ????????????????????????k???????????????
	startTime := reqStart.Add(-time.Duration(maxMxN) * time.Minute)
	// todo ????????????????????????????????????maxMxN???????????????
	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	if startTime.Before(dataLimitTime) {
		return res, nil
	}
	// ?????????????????????
	if startTime.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, startTime, reqEnd, startTime.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	if "BTC" == req.CoinType {
		klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
			startTime.Add(-8*time.Hour).UnixMilli(),
			reqEnd.Add(-8*time.Hour).UnixMilli(),
		)
	} else if "ETH" == req.CoinType {
		klineMOne, err = b.klineMOneRepo.GetKLineMOneEthByStartTime(
			startTime.Add(-8*time.Hour).UnixMilli(),
			reqEnd.Add(-8*time.Hour).UnixMilli(),
		)
	} else if "FIL" == req.CoinType {
		klineMOne, err = b.klineMOneRepo.GetKLineMOneFilByStartTime(
			startTime.Add(-8*time.Hour).UnixMilli(),
			reqEnd.Add(-8*time.Hour).UnixMilli(),
		)
	}

	// ????????????
	var (
		lastActionTag string

		kLineDataMLive []*KLineMOne
	)
	operationData := make(map[string]*OperationData2, 0)
	operationDataToPointSecond := make(map[string]int, 0)
	operationDataToPointSecondKeep := make(map[string]int, 0) // ??????????????????
	operationDataToPointThirdKeep := make(map[string]int, 0)  // ??????????????????

	operationDataToPointSecondConfirm := make(map[string]int, 0) // ??????????????????
	operationDataToPointThirdConfirm := make(map[string]int, 0)  // ??????????????????
	maNDataMLiveMap := make(map[int64]*Ma, 0)

	reqStartMilli := reqStart.Add(-8 * time.Hour).UnixMilli()
	for _, vKlineM := range klineMOne {
		// ????????????
		tmpNow := time.UnixMilli(vKlineM.StartTime).UTC().Add(8 * time.Hour)
		var (
			lastKeyMLive int
		)
		if 0 == tmpNow.Minute()%m {
			kLineDataMLive = append(kLineDataMLive, &KLineMOne{
				StartPrice: vKlineM.StartPrice,
				StartTime:  vKlineM.StartTime,
				TopPrice:   vKlineM.TopPrice,
				LowPrice:   vKlineM.LowPrice,
				EndPrice:   vKlineM.EndPrice,
				EndTime:    vKlineM.EndTime,
			})
			//if 1675180800000 <= vKlineM.StartTime {
			//	fmt.Println(vKlineM)
			//}
			continue
		} else {
			lastKeyMLive = len(kLineDataMLive) - 1
			if 0 <= lastKeyMLive { // ?????????????????????????????????n?????????????????????
				kLineDataMLive[lastKeyMLive].EndPrice = vKlineM.EndPrice
				kLineDataMLive[lastKeyMLive].EndTime = vKlineM.EndTime
				if kLineDataMLive[lastKeyMLive].TopPrice < vKlineM.TopPrice {
					kLineDataMLive[lastKeyMLive].TopPrice = vKlineM.TopPrice
				}
				if kLineDataMLive[lastKeyMLive].LowPrice > vKlineM.LowPrice {
					kLineDataMLive[lastKeyMLive].LowPrice = vKlineM.LowPrice
				}
			}
			//if 1675180800000 <= kLineData15MLive[lastKey].StartTime {
			//	fmt.Println(kLineData15MLive[lastKey], lastKey)
			//}

			if m-1 != tmpNow.Minute()%m {
				continue
			}

		}

		lastKeyMLive = len(kLineDataMLive) - 1 // ????????????

		// ???????????????????????????????????????
		if reqStartMilli > vKlineM.StartTime {
			// ????????????3??????????????????ma???
			if reqStartMilli-int64(2*n*60000) <= vKlineM.StartTime {
				var tmpTotalEndPrice float64
				var tmpAvgEndPrice float64

				for i := 0; i < n; i++ {
					// ???i???m???????????????k???
					tmpTotalEndPrice += kLineDataMLive[lastKeyMLive-i].EndPrice
				}
				tmpAvgEndPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpTotalEndPrice/float64(n)), 64)
				maNDataMLiveMap[kLineDataMLive[lastKeyMLive].StartTime] = &Ma{AvgEndPrice: tmpAvgEndPrice}
			}

			continue
		}

		//fmt.Println(vKlineM.StartTime, kLineDataMLive[lastKeyMLive], kLineData3MLive[lastKey3MLive], kLineData60MLive[lastKey60MLive])

		// ??????????????????ma??????
		var tmpTotalEndPrice float64
		var tmpAvgEndPrice float64
		for i := 0; i < n; i++ {
			// ???i???m???????????????k???
			tmpTotalEndPrice += kLineDataMLive[lastKeyMLive-i].EndPrice
		}
		tmpAvgEndPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpTotalEndPrice/float64(n)), 64)
		maNDataMLiveMap[kLineDataMLive[lastKeyMLive].StartTime] = &Ma{AvgEndPrice: tmpAvgEndPrice}

		var tagNum int64
		// ??????
		tmpResDataListK := &v1.AreaPointIntervalMAvgEndPriceDataReply_ListK{
			X1: kLineDataMLive[lastKeyMLive].StartPrice,
			X2: kLineDataMLive[lastKeyMLive].EndPrice,
			X3: kLineDataMLive[lastKeyMLive].TopPrice,
			X4: kLineDataMLive[lastKeyMLive].LowPrice,
			X5: kLineDataMLive[lastKeyMLive].EndTime,
			X6: kLineDataMLive[lastKeyMLive].StartTime,
		}

		// ????????????
		tmpPointFirstSub := maNDataMLiveMap[kLineDataMLive[lastKeyMLive].StartTime].AvgEndPrice - maNDataMLiveMap[kLineDataMLive[lastKeyMLive-1].StartTime].AvgEndPrice
		tmpPointSecondSub := maNDataMLiveMap[kLineDataMLive[lastKeyMLive-1].StartTime].AvgEndPrice - maNDataMLiveMap[kLineDataMLive[lastKeyMLive-2].StartTime].AvgEndPrice

		tmpResDataListMaNMFirst := &v1.AreaPointIntervalMAvgEndPriceDataReply_ListMaNMFirst{X1: tmpAvgEndPrice}
		tmpResDataListSubPoint := &v1.AreaPointIntervalMAvgEndPriceDataReply_ListSubPoint{X1: tmpPointFirstSub}

		res.DataListK = append(res.DataListK, tmpResDataListK)
		res.DataListMaNMFirst = append(res.DataListMaNMFirst, tmpResDataListMaNMFirst)
		res.DataListSubPoint = append(res.DataListSubPoint, tmpResDataListSubPoint)

		// ????????????????????????????????????
		if tmpPointSecondSub > pointFirst && pointFirst > tmpPointFirstSub {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ????????????????????????????????????
					tmpOpen := false
					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; okTwo {
						if 2 <= operationDataToPointSecondKeep[lastActionTag] {
							tmpOpen = true
						} else {
							operationDataToPointSecondKeep[lastActionTag] = 0
						}
					} else if _, okThird := operationDataToPointThirdKeep[lastActionTag]; okThird {
						tmpOpen = true
					} else if _, okThirdConfirm := operationDataToPointThirdConfirm[lastActionTag]; okThirdConfirm {
						tmpOpen = true
					} else if _, okSecondConfirm := operationDataToPointSecondConfirm[lastActionTag]; okSecondConfirm {
						tmpOpen = true
					}

					if tmpOpen {
						rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
						tmpCloseLastOperationData := &OperationData2{
							StartTime:  vKlineM.StartTime,
							EndTime:    vKlineM.EndTime,
							StartPrice: vKlineM.StartPrice,
							EndPrice:   vKlineM.EndPrice,
							Amount:     0,
							Type:       "more",
							Status:     "close",
							Rate:       rate,
						}

						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
						operationData[lastActionTag] = tmpCloseLastOperationData

						currentOperationData := &OperationData2{
							StartTime:  vKlineM.StartTime,
							EndTime:    vKlineM.EndTime,
							StartPrice: vKlineM.StartPrice,
							EndPrice:   vKlineM.EndPrice,
							Amount:     2,
							Type:       "empty",
							Status:     "open", // ????????????
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
						operationData[lastActionTag] = currentOperationData
					}

				} else if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ????????????????????????????????????????????????????????????
					if _, okTwo := operationDataToPointSecond[lastActionTag]; okTwo {
						operationDataToPointSecond[lastActionTag] = 0
					}
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:  vKlineM.StartTime,
					EndTime:    vKlineM.EndTime,
					StartPrice: vKlineM.StartPrice,
					EndPrice:   vKlineM.EndPrice,
					Amount:     2,
					Type:       "empty",
					Status:     "open", // ????????????
				}
				tagNum++
				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				operationData[lastActionTag] = currentOperationData
			}
		}

		// ????????????
		if (pointFirst + pointInterval) < tmpPointFirstSub {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ??????????????????????????????

					// ???????????????
					if _, okTwo := operationDataToPointSecond[lastActionTag]; !okTwo {
						operationDataToPointSecond[lastActionTag] = 1
					} else if 0 == operationDataToPointSecond[lastActionTag] {
						operationDataToPointSecond[lastActionTag] = 1
					} else {
						// ??????
						operationDataToPointSecond[lastActionTag] += 1
						if 2 <= operationDataToPointSecond[lastActionTag] {
							// ????????????????????????
							rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
							rate = -rate - 0.0003
							tmpCloseLastOperationData := &OperationData2{
								StartTime:  vKlineM.StartTime,
								EndTime:    vKlineM.EndTime,
								StartPrice: vKlineM.StartPrice,
								EndPrice:   vKlineM.EndPrice,
								Amount:     0,
								Type:       "empty",
								Status:     "close",
								Rate:       rate,
							}

							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
							operationData[lastActionTag] = tmpCloseLastOperationData

							currentOperationData := &OperationData2{
								StartTime:  vKlineM.StartTime,
								EndTime:    vKlineM.EndTime,
								StartPrice: vKlineM.StartPrice,
								EndPrice:   vKlineM.EndPrice,
								Amount:     2,
								Type:       "more",
								Status:     "open", // ????????????
							}
							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
							operationData[lastActionTag] = currentOperationData

							// ????????????
							operationDataToPointSecondConfirm[lastActionTag] = 1
						}
					}
				} else if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ??????????????????????????????????????????

					// ???????????????
					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; !okTwo {
						operationDataToPointSecondKeep[lastActionTag] = 1 // ???????????????
					} else {
						operationDataToPointSecondKeep[lastActionTag] += 1 // ???n?????????
					}
				}
			}
		}

		// ???????????????????????????
		if (pointFirst + 2*pointInterval) < tmpPointFirstSub {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ????????????????????????
					rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
					rate = -rate - 0.0003
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     0,
						Type:       "empty",
						Status:     "close",
						Rate:       rate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData

					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "more",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = currentOperationData

					// ????????????
					operationDataToPointThirdConfirm[lastActionTag] = 1
				} else if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					if _, okThird := operationDataToPointThirdKeep[lastActionTag]; !okThird {
						operationDataToPointThirdKeep[lastActionTag] = 1 // ???????????????
					}
				}
			}
		}

		// ????????????????????????????????????
		if tmpPointSecondSub < -pointFirst && -pointFirst < tmpPointFirstSub {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {

					// ????????????????????????????????????
					tmpOpen := false
					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; okTwo {
						if 2 <= operationDataToPointSecondKeep[lastActionTag] {
							tmpOpen = true
						} else {
							operationDataToPointSecondKeep[lastActionTag] = 0
						}

					} else if _, okThird := operationDataToPointThirdKeep[lastActionTag]; okThird {
						tmpOpen = true
					} else if _, okThirdConfirm := operationDataToPointThirdConfirm[lastActionTag]; okThirdConfirm {
						tmpOpen = true
					} else if _, okSecondConfirm := operationDataToPointSecondConfirm[lastActionTag]; okSecondConfirm {
						tmpOpen = true
					}

					if tmpOpen {
						rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
						rate = -rate - 0.0003
						tmpCloseLastOperationData := &OperationData2{
							StartTime:  vKlineM.StartTime,
							EndTime:    vKlineM.EndTime,
							StartPrice: vKlineM.StartPrice,
							EndPrice:   vKlineM.EndPrice,
							Amount:     0,
							Type:       "empty",
							Status:     "close",
							Rate:       rate,
						}

						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
						operationData[lastActionTag] = tmpCloseLastOperationData

						currentOperationData := &OperationData2{
							StartTime:  vKlineM.StartTime,
							EndTime:    vKlineM.EndTime,
							StartPrice: vKlineM.StartPrice,
							EndPrice:   vKlineM.EndPrice,
							Amount:     2,
							Type:       "more",
							Status:     "open", // ????????????
						}
						tagNum++
						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
						operationData[lastActionTag] = currentOperationData
					}

				} else if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ????????????????????????????????????????????????????????????
					if _, okTwo := operationDataToPointSecond[lastActionTag]; okTwo {
						operationDataToPointSecond[lastActionTag] = 0
					}
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:  vKlineM.StartTime,
					EndTime:    vKlineM.EndTime,
					StartPrice: vKlineM.StartPrice,
					EndPrice:   vKlineM.EndPrice,
					Amount:     2,
					Type:       "more",
					Status:     "open", // ????????????
				}
				tagNum++
				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				operationData[lastActionTag] = currentOperationData
			}
		}

		// ????????????
		if (-pointFirst - pointInterval) > tmpPointFirstSub {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				// ????????????????????????????????????????????????
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ??????????????????????????????

					// ???????????????
					if _, okTwo := operationDataToPointSecond[lastActionTag]; !okTwo {
						operationDataToPointSecond[lastActionTag] = 1
					} else if 0 == operationDataToPointSecond[lastActionTag] {
						operationDataToPointSecond[lastActionTag] = 1
					} else {
						operationDataToPointSecond[lastActionTag] += 1
						if 2 <= operationDataToPointSecond[lastActionTag] {
							// ????????????????????????
							rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
							tmpCloseLastOperationData := &OperationData2{
								StartTime:  vKlineM.StartTime,
								EndTime:    vKlineM.EndTime,
								StartPrice: vKlineM.StartPrice,
								EndPrice:   vKlineM.EndPrice,
								Amount:     0,
								Type:       "more",
								Status:     "close",
								Rate:       rate,
							}

							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
							operationData[lastActionTag] = tmpCloseLastOperationData

							currentOperationData := &OperationData2{
								StartTime:  vKlineM.StartTime,
								EndTime:    vKlineM.EndTime,
								StartPrice: vKlineM.StartPrice,
								EndPrice:   vKlineM.EndPrice,
								Amount:     2,
								Type:       "empty",
								Status:     "open", // ????????????
							}
							tagNum++
							lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
							operationData[lastActionTag] = currentOperationData

							// ????????????
							operationDataToPointSecondConfirm[lastActionTag] = 2
						}
					}
				} else if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ??????????????????????????????????????????

					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; !okTwo {
						operationDataToPointSecondKeep[lastActionTag] = 1 // ???????????????
					} else {
						operationDataToPointSecondKeep[lastActionTag] += 1 // ???n?????????
					}
				}
			}
		}

		// ???????????????????????????
		if (-pointFirst - 2*pointInterval) > tmpPointFirstSub {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					// ????????????????????????
					rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     0,
						Type:       "more",
						Status:     "close",
						Rate:       rate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData

					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "empty",
						Status:     "open", // ????????????
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = currentOperationData

					// ????????????
					operationDataToPointThirdConfirm[lastActionTag] = 1
				} else if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					if _, okThird := operationDataToPointThirdKeep[lastActionTag]; !okThird {
						operationDataToPointThirdKeep[lastActionTag] = 1 // ???????????????
					}
				}
			}
		}

	}

	// ??????
	for _, vOperationData := range operationData {
		resOperationData = append(resOperationData, vOperationData)
	}
	sort.Sort(resOperationData)

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
		tmpLastCloseK = -1
	)

	// ????????????????????????
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for kOperationData, vOperationData := range resOperationData {
		if kOperationData > tmpLastCloseK { // ????????????????????????????????????-1???????????????
			break
		}

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

		res.OperationData = append(res.OperationData, &v1.AreaPointIntervalMAvgEndPriceDataReply_List2{
			StartPrice: vOperationData.StartPrice,
			EndPrice:   vOperationData.EndPrice,
			StartTime:  vOperationData.StartTime,
			EndTime:    vOperationData.EndTime,
			Type:       vOperationData.Type,
			Action:     vOperationData.Action,
			Status:     vOperationData.Status,
			Rate:       vOperationData.Rate,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	return res, nil
}

// AreaPointIntervalMAvgEndPriceDataBack k????????????m?????????????????????????????? .
//func (b *BinanceDataUsecase) AreaPointIntervalMAvgEndPriceDataBack(ctx context.Context, req *v1.AreaPointIntervalMAvgEndPriceDataRequest) (*v1.AreaPointIntervalMAvgEndPriceDataReply, error) {
//	var (
//		resOperationData OperationData2Slice
//		klineMOne        []*KLineMOne
//		reqStart         time.Time
//		reqEnd           time.Time
//		m                int
//		n                int
//		pointFirst       float64
//		pointInterval    float64
//		err              error
//	)
//
//	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // ????????????????????????
//	if nil != err {
//		return nil, err
//	}
//	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // ????????????????????????
//	if nil != err {
//		return nil, err
//	}
//
//	res := &v1.AreaPointIntervalMAvgEndPriceDataReply{
//		DataListK:         make([]*v1.AreaPointIntervalMAvgEndPriceDataReply_ListK, 0),
//		DataListMaNMFirst: make([]*v1.AreaPointIntervalMAvgEndPriceDataReply_ListMaNMFirst, 0),
//	}
//
//	m = int(req.M)
//	n = int(req.N)
//	pointFirst = req.PointFirst
//	pointInterval = req.PointInterval
//
//	maxMxN := 4 * 60 // macd??????????????????????????????????????????????????????????????????60??????
//
//	// ????????????????????????k???????????????
//	startTime := reqStart.Add(-time.Duration(maxMxN) * time.Minute)
//	// todo ????????????????????????????????????maxMxN???????????????
//	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
//	if startTime.Before(dataLimitTime) {
//		return res, nil
//	}
//	// ?????????????????????
//	if startTime.After(reqEnd) {
//		return res, nil
//	}
//	fmt.Println(maxMxN, startTime, reqEnd, startTime.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
//	klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
//		startTime.Add(-8*time.Hour).UnixMilli(),
//		reqEnd.Add(-8*time.Hour).UnixMilli(),
//	)
//	// ????????????
//	var (
//		lastActionTag string
//
//		kLineDataMLive []*KLineMOne
//	)
//	operationData := make(map[string]*OperationData2, 0)
//	operationDataToPointSecond := make(map[string]int, 0)     // ??????????????????
//	operationDataToPointSecondKeep := make(map[string]int, 0) // ??????????????????
//	operationDataToPointThirdKeep := make(map[string]int, 0)  // ??????????????????
//	maNDataMLiveMap := make(map[int64]*Ma, 0)
//
//	reqStartMilli := reqStart.Add(-8 * time.Hour).UnixMilli()
//	for kKlineM, vKlineM := range klineMOne {
//		// ????????????
//		tmpNow := time.UnixMilli(vKlineM.StartTime).UTC().Add(8 * time.Hour)
//		var (
//			lastKeyMLive int
//		)
//		if 0 == tmpNow.Minute()%m {
//			kLineDataMLive = append(kLineDataMLive, &KLineMOne{
//				StartPrice: vKlineM.StartPrice,
//				StartTime:  vKlineM.StartTime,
//				TopPrice:   vKlineM.TopPrice,
//				LowPrice:   vKlineM.LowPrice,
//				EndPrice:   vKlineM.EndPrice,
//				EndTime:    vKlineM.EndTime,
//			})
//			//if 1675180800000 <= vKlineM.StartTime {
//			//	fmt.Println(vKlineM)
//			//}
//		} else {
//			lastKeyMLive = len(kLineDataMLive) - 1
//			if 0 <= lastKeyMLive { // ?????????????????????????????????n?????????????????????
//				kLineDataMLive[lastKeyMLive].EndPrice = vKlineM.EndPrice
//				kLineDataMLive[lastKeyMLive].EndTime = vKlineM.EndTime
//				if kLineDataMLive[lastKeyMLive].TopPrice < vKlineM.TopPrice {
//					kLineDataMLive[lastKeyMLive].TopPrice = vKlineM.TopPrice
//				}
//				if kLineDataMLive[lastKeyMLive].LowPrice > vKlineM.LowPrice {
//					kLineDataMLive[lastKeyMLive].LowPrice = vKlineM.LowPrice
//				}
//			}
//			//if 1675180800000 <= kLineData15MLive[lastKey].StartTime {
//			//	fmt.Println(kLineData15MLive[lastKey], lastKey)
//			//}
//		}
//		lastKeyMLive = len(kLineDataMLive) - 1 // ????????????
//
//		// ???????????????????????????????????????
//		if reqStartMilli > vKlineM.StartTime {
//			// ????????????3??????????????????ma???
//			if reqStartMilli-int64(3*60000) <= vKlineM.StartTime {
//				var tmpTotalEndPrice float64
//				var tmpAvgEndPrice float64
//
//				for i := 0; i < n; i++ {
//					// ???i???m???????????????k???
//					tmpTotalEndPrice += kLineDataMLive[lastKeyMLive-i].EndPrice
//				}
//				tmpAvgEndPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpTotalEndPrice/float64(n)), 64)
//				maNDataMLiveMap[vKlineM.StartTime] = &Ma{AvgEndPrice: tmpAvgEndPrice}
//			}
//
//			continue
//		}
//
//		//fmt.Println(vKlineM.StartTime, kLineDataMLive[lastKeyMLive], kLineData3MLive[lastKey3MLive], kLineData60MLive[lastKey60MLive])
//
//		// ??????????????????ma??????
//		var tmpTotalEndPrice float64
//		var tmpAvgEndPrice float64
//		for i := 0; i < n; i++ {
//			// ???i???m???????????????k???
//			tmpTotalEndPrice += kLineDataMLive[lastKeyMLive-i].EndPrice
//		}
//		tmpAvgEndPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpTotalEndPrice/float64(n)), 64)
//		maNDataMLiveMap[vKlineM.StartTime] = &Ma{AvgEndPrice: tmpAvgEndPrice}
//
//		var tagNum int64
//		// ??????
//		tmpResDataListK := &v1.AreaPointIntervalMAvgEndPriceDataReply_ListK{
//			X1: kLineDataMLive[lastKeyMLive].StartPrice,
//			X2: kLineDataMLive[lastKeyMLive].EndPrice,
//			X3: kLineDataMLive[lastKeyMLive].TopPrice,
//			X4: kLineDataMLive[lastKeyMLive].LowPrice,
//			X5: kLineDataMLive[lastKeyMLive].EndTime,
//			X6: kLineDataMLive[lastKeyMLive].StartTime,
//		}
//
//		tmpResDataListMaNMFirst := &v1.AreaPointIntervalMAvgEndPriceDataReply_ListMaNMFirst{X1: tmpAvgEndPrice}
//		lastResDataListKKey := len(res.DataListK) - 1
//		if 0 == tmpNow.Minute()%m {
//			res.DataListK = append(res.DataListK, tmpResDataListK)
//			res.DataListMaNMFirst = append(res.DataListMaNMFirst, tmpResDataListMaNMFirst)
//		} else {
//			if 0 <= lastResDataListKKey {
//				res.DataListK[lastResDataListKKey] = tmpResDataListK
//				res.DataListMaNMFirst[lastResDataListKKey] = tmpResDataListMaNMFirst
//			}
//		}
//
//		// ????????????
//		tmpPointFirstSub := maNDataMLiveMap[vKlineM.StartTime].AvgEndPrice - maNDataMLiveMap[klineMOne[kKlineM-1].StartTime].AvgEndPrice
//		tmpPointSecondSub := maNDataMLiveMap[klineMOne[kKlineM-1].StartTime].AvgEndPrice - maNDataMLiveMap[klineMOne[kKlineM-2].StartTime].AvgEndPrice
//
//		// ????????????????????????????????????
//		if tmpPointSecondSub > pointFirst && pointFirst > tmpPointFirstSub {
//			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
//				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ????????????????????????????????????
//					tmpOpen := false
//					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; okTwo && 2 == operationDataToPointSecondKeep[lastActionTag] {
//						tmpOpen = true
//					} else if _, okThird := operationDataToPointThirdKeep[lastActionTag]; okThird {
//						tmpOpen = true
//					}
//
//					if tmpOpen {
//						rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
//						tmpCloseLastOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     0,
//							Type:       "more",
//							Status:     "close",
//							Rate:       rate,
//						}
//
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = tmpCloseLastOperationData
//
//						currentOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     2,
//							Type:       "empty",
//							Status:     "open", // ????????????
//						}
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = currentOperationData
//					}
//
//				} else if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ????????????????????????????????????????????????????????????
//					if _, okTwo := operationDataToPointSecond[lastActionTag]; okTwo {
//						operationDataToPointSecond[lastActionTag] = 0
//					}
//				}
//
//			} else {
//				currentOperationData := &OperationData2{
//					StartTime:  vKlineM.StartTime,
//					EndTime:    vKlineM.EndTime,
//					StartPrice: vKlineM.StartPrice,
//					EndPrice:   vKlineM.EndPrice,
//					Amount:     2,
//					Type:       "empty",
//					Status:     "open", // ????????????
//				}
//				tagNum++
//				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//				operationData[lastActionTag] = currentOperationData
//			}
//		}
//
//		// ????????????
//		if (pointFirst + pointInterval) < tmpPointFirstSub {
//			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
//				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ??????????????????????????????
//
//					// ???????????????
//					if _, okTwo := operationDataToPointSecond[lastActionTag]; !okTwo {
//						operationDataToPointSecond[lastActionTag] = 1
//					} else {
//						operationDataToPointSecond[lastActionTag] = 2
//						// ????????????????????????
//						rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
//						rate = -rate - 0.0003
//						tmpCloseLastOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     0,
//							Type:       "empty",
//							Status:     "close",
//							Rate:       rate,
//						}
//
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = tmpCloseLastOperationData
//
//						currentOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     2,
//							Type:       "more",
//							Status:     "open", // ????????????
//						}
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = currentOperationData
//					}
//				} else if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ??????????????????????????????????????????
//
//					// ???????????????
//					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; !okTwo {
//						operationDataToPointSecondKeep[lastActionTag] = 1 // ???????????????
//					} else {
//						operationDataToPointSecondKeep[lastActionTag] = 2 // ???n?????????
//					}
//				}
//			}
//		}
//
//		// ???????????????????????????
//		if (pointFirst + 2*pointInterval) < tmpPointFirstSub {
//			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
//				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ????????????????????????
//					rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
//					rate = -rate - 0.0003
//					tmpCloseLastOperationData := &OperationData2{
//						StartTime:  vKlineM.StartTime,
//						EndTime:    vKlineM.EndTime,
//						StartPrice: vKlineM.StartPrice,
//						EndPrice:   vKlineM.EndPrice,
//						Amount:     0,
//						Type:       "empty",
//						Status:     "close",
//						Rate:       rate,
//					}
//
//					tagNum++
//					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//					operationData[lastActionTag] = tmpCloseLastOperationData
//
//					currentOperationData := &OperationData2{
//						StartTime:  vKlineM.StartTime,
//						EndTime:    vKlineM.EndTime,
//						StartPrice: vKlineM.StartPrice,
//						EndPrice:   vKlineM.EndPrice,
//						Amount:     2,
//						Type:       "more",
//						Status:     "open", // ????????????
//					}
//					tagNum++
//					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//					operationData[lastActionTag] = currentOperationData
//				} else if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					if _, okThird := operationDataToPointThirdKeep[lastActionTag]; !okThird {
//						operationDataToPointThirdKeep[lastActionTag] = 1 // ???????????????
//					}
//				}
//			}
//		}
//
//		// ????????????????????????????????????
//		if tmpPointSecondSub < -pointFirst && -pointFirst < tmpPointFirstSub {
//			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
//				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//
//					// ????????????????????????????????????
//					tmpOpen := false
//					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; okTwo && 2 == operationDataToPointSecondKeep[lastActionTag] {
//						tmpOpen = true
//						fmt.Println(111, vKlineM.EndTime)
//					} else if _, okThird := operationDataToPointThirdKeep[lastActionTag]; okThird {
//						tmpOpen = true
//						fmt.Println(222, vKlineM.EndTime)
//					}
//
//					if tmpOpen {
//						rate := (vKlineM.EndPrice - tmpOpenLastOperationData2.EndPrice) / tmpOpenLastOperationData2.EndPrice
//						rate = -rate - 0.0003
//						tmpCloseLastOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     0,
//							Type:       "empty",
//							Status:     "close",
//							Rate:       rate,
//						}
//
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = tmpCloseLastOperationData
//
//						currentOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     2,
//							Type:       "more",
//							Status:     "open", // ????????????
//						}
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = currentOperationData
//					} else if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//						// ????????????????????????????????????????????????????????????
//						if _, okTwo := operationDataToPointSecond[lastActionTag]; okTwo {
//							operationDataToPointSecond[lastActionTag] = 0
//						}
//					}
//				}
//
//			} else {
//				currentOperationData := &OperationData2{
//					StartTime:  vKlineM.StartTime,
//					EndTime:    vKlineM.EndTime,
//					StartPrice: vKlineM.StartPrice,
//					EndPrice:   vKlineM.EndPrice,
//					Amount:     2,
//					Type:       "more",
//					Status:     "open", // ????????????
//				}
//				tagNum++
//				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//				operationData[lastActionTag] = currentOperationData
//			}
//		}
//
//		// ????????????
//		if (-pointFirst - pointInterval) > tmpPointFirstSub {
//			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
//				// ????????????????????????????????????????????????
//				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ??????????????????????????????
//
//					// ???????????????
//					if _, okTwo := operationDataToPointSecond[lastActionTag]; !okTwo {
//						operationDataToPointSecond[lastActionTag] = 1
//					} else {
//						operationDataToPointSecond[lastActionTag] = 2
//						// ????????????????????????
//						rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
//						tmpCloseLastOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     0,
//							Type:       "more",
//							Status:     "close",
//							Rate:       rate,
//						}
//
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = tmpCloseLastOperationData
//
//						currentOperationData := &OperationData2{
//							StartTime:  vKlineM.StartTime,
//							EndTime:    vKlineM.EndTime,
//							StartPrice: vKlineM.StartPrice,
//							EndPrice:   vKlineM.EndPrice,
//							Amount:     2,
//							Type:       "empty",
//							Status:     "open", // ????????????
//						}
//						tagNum++
//						lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//						operationData[lastActionTag] = currentOperationData
//					}
//				} else if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ??????????????????????????????????????????
//
//					if _, okTwo := operationDataToPointSecondKeep[lastActionTag]; !okTwo {
//						operationDataToPointSecondKeep[lastActionTag] = 1 // ???????????????
//					} else {
//						operationDataToPointSecondKeep[lastActionTag] = 2 // ???n?????????
//					}
//				}
//			}
//		}
//
//		// ???????????????????????????
//		if (-pointFirst - 2*pointInterval) > tmpPointFirstSub {
//			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
//				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					// ????????????????????????
//					rate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
//					tmpCloseLastOperationData := &OperationData2{
//						StartTime:  vKlineM.StartTime,
//						EndTime:    vKlineM.EndTime,
//						StartPrice: vKlineM.StartPrice,
//						EndPrice:   vKlineM.EndPrice,
//						Amount:     0,
//						Type:       "more",
//						Status:     "close",
//						Rate:       rate,
//					}
//
//					tagNum++
//					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//					operationData[lastActionTag] = tmpCloseLastOperationData
//
//					currentOperationData := &OperationData2{
//						StartTime:  vKlineM.StartTime,
//						EndTime:    vKlineM.EndTime,
//						StartPrice: vKlineM.StartPrice,
//						EndPrice:   vKlineM.EndPrice,
//						Amount:     2,
//						Type:       "empty",
//						Status:     "open", // ????????????
//					}
//					tagNum++
//					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
//					operationData[lastActionTag] = currentOperationData
//				} else if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
//					if _, okThird := operationDataToPointThirdKeep[lastActionTag]; !okThird {
//						operationDataToPointThirdKeep[lastActionTag] = 1 // ???????????????
//					}
//				}
//			}
//		}
//
//	}
//
//	// ??????
//	for _, vOperationData := range operationData {
//		resOperationData = append(resOperationData, vOperationData)
//	}
//	sort.Sort(resOperationData)
//
//	var (
//		tmpWinTotal   int64
//		tmpCloseTotal int64
//		tmpRate       float64
//		winRate       float64
//		tmpLastCloseK = -1
//	)
//
//	// ????????????????????????
//	for i := len(resOperationData) - 1; i >= 0; i-- {
//		if "close" == resOperationData[i].Status {
//			tmpLastCloseK = i
//			break
//		}
//	}
//
//	for kOperationData, vOperationData := range resOperationData {
//		if kOperationData > tmpLastCloseK { // ????????????????????????????????????-1???????????????
//			break
//		}
//
//		if "open" == vOperationData.Status {
//			res.OperationOrderTotal++
//		}
//
//		if "close" == vOperationData.Status {
//			tmpCloseTotal++
//			if 0 < vOperationData.Rate {
//				tmpWinTotal++
//			}
//		}
//
//		tmpRate += vOperationData.Rate
//
//		res.OperationData = append(res.OperationData, &v1.AreaPointIntervalMAvgEndPriceDataReply_List2{
//			StartPrice: vOperationData.StartPrice,
//			EndPrice:   vOperationData.EndPrice,
//			StartTime:  vOperationData.StartTime,
//			EndTime:    vOperationData.EndTime,
//			Type:       vOperationData.Type,
//			Action:     vOperationData.Action,
//			Status:     vOperationData.Status,
//			Rate:       vOperationData.Rate,
//		})
//	}
//
//	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
//		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
//	}
//	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
//	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
//	return res, nil
//}

// IntervalMAvgEndPriceMacdAndAtrData k????????????m?????????????????????????????? .
func (b *BinanceDataUsecase) IntervalMAvgEndPriceMacdAndAtrData(ctx context.Context, req *v1.IntervalMAvgEndPriceMacdAndAtrDataRequest) (*v1.IntervalMAvgEndPriceMacdAndAtrDataReply, error) {
	var (
		klineMOne []*KLineMOne
		reqStart  time.Time
		reqEnd    time.Time
		m         int
		max1      float64
		max2      float64
		low1      float64
		low2      float64
		atr1N     int
		atr2N     int
		err       error
	)

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // ????????????????????????
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // ????????????????????????
	if nil != err {
		return nil, err
	}

	res := &v1.IntervalMAvgEndPriceMacdAndAtrDataReply{
		DataListK: make([]*v1.IntervalMAvgEndPriceMacdAndAtrDataReply_ListK, 0),
	}

	m = int(req.M)
	max1 = req.Max1
	max2 = req.Max2
	low1 = req.Low1
	low2 = req.Low2
	atr1N = int(req.Atr1N)
	atr2N = int(req.Atr2N)

	maxMxN := m * 201 // macd??????????????????????????????????????????????????????????????????60??????

	// ????????????????????????k???????????????
	startTime := reqStart.Add(-time.Duration(maxMxN) * time.Minute)
	// todo ????????????????????????????????????maxMxN???????????????
	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	if startTime.Before(dataLimitTime) {
		return res, nil
	}
	// ?????????????????????
	if startTime.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, startTime, reqEnd, startTime.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	if "BTC" == req.CoinType {
		klineMOne, err = b.klineMOneRepo.GetKLineMOneBtcByStartTime(
			startTime.Add(-8*time.Hour).UnixMilli(),
			reqEnd.Add(-8*time.Hour).UnixMilli(),
		)
	} else if "ETH" == req.CoinType {
		klineMOne, err = b.klineMOneRepo.GetKLineMOneEthByStartTime(
			startTime.Add(-8*time.Hour).UnixMilli(),
			reqEnd.Add(-8*time.Hour).UnixMilli(),
		)
	} else if "FIL" == req.CoinType {
		klineMOne, err = b.klineMOneRepo.GetKLineMOneFilByStartTime(
			startTime.Add(-8*time.Hour).UnixMilli(),
			reqEnd.Add(-8*time.Hour).UnixMilli(),
		)
	}

	// ????????????
	var (
		//lastActionTag string
		kLineDataMLive []*KLineMOne
		macdData       []*MACDPoint

		maxK    *KLineMOne // ?????????????????????K???
		maxMacd *MACDPoint // ?????????????????????Macd
		lowK    *KLineMOne // ?????????????????????K???
		lowMacd *MACDPoint // ?????????????????????Macd

		lastMacd *MACDPoint // ?????????????????????K???

		operationData []*OperationData3
	)
	//operationData := make(map[string]*OperationData2, 0)
	macdDataLiveMap := make(map[int64]*MACDPoint, 0)

	reqStartMilli := reqStart.Add(-8 * time.Hour).UnixMilli()
	for _, vKlineM := range klineMOne {
		// ????????????
		tmpNow := time.UnixMilli(vKlineM.StartTime).UTC().Add(8 * time.Hour)
		var (
			lastKeyMLive int
		)
		if 0 == tmpNow.Minute()%m {
			kLineDataMLive = append(kLineDataMLive, &KLineMOne{
				StartPrice: vKlineM.StartPrice,
				StartTime:  vKlineM.StartTime,
				TopPrice:   vKlineM.TopPrice,
				LowPrice:   vKlineM.LowPrice,
				EndPrice:   vKlineM.EndPrice,
				EndTime:    vKlineM.EndTime,
			})
			continue
		} else {
			lastKeyMLive = len(kLineDataMLive) - 1
			if 0 <= lastKeyMLive { // ?????????????????????????????????n?????????????????????
				kLineDataMLive[lastKeyMLive].EndPrice = vKlineM.EndPrice
				kLineDataMLive[lastKeyMLive].EndTime = vKlineM.EndTime
				if kLineDataMLive[lastKeyMLive].TopPrice < vKlineM.TopPrice {
					kLineDataMLive[lastKeyMLive].TopPrice = vKlineM.TopPrice
				}
				if kLineDataMLive[lastKeyMLive].LowPrice > vKlineM.LowPrice {
					kLineDataMLive[lastKeyMLive].LowPrice = vKlineM.LowPrice
				}
			}

			if m-1 != tmpNow.Minute()%m {
				continue
			}

		}

		lastKeyMLive = len(kLineDataMLive) - 1 // ????????????

		// ???????????????????????????????????????
		if reqStartMilli > vKlineM.StartTime {
			continue
		}

		// ???????????????????????????
		if reqStartMilli == kLineDataMLive[lastKeyMLive].StartTime {
			// macd????????????200???k???????????????
			lastLastKeyMLive := lastKeyMLive - 1
			var tmpMacdData []*MACDPoint
			tmpMacdData, err = b.klineMOneRepo.NewMACDData(kLineDataMLive[lastLastKeyMLive-199:])
			if nil != err {
				continue
			}
			lastMacd = tmpMacdData[199]
		}

		// macd????????????200???k???????????????
		macdData, err = b.klineMOneRepo.NewMACDData(kLineDataMLive[lastKeyMLive-199:])
		if nil != err {
			continue
		}
		macdDataLiveMap[vKlineM.StartTime] = macdData[199]

		//var tagNum int64
		// ??????
		tmpResDataListK := &v1.IntervalMAvgEndPriceMacdAndAtrDataReply_ListK{
			X1: kLineDataMLive[lastKeyMLive].StartPrice,
			X2: kLineDataMLive[lastKeyMLive].EndPrice,
			X3: kLineDataMLive[lastKeyMLive].TopPrice,
			X4: kLineDataMLive[lastKeyMLive].LowPrice,
			X5: kLineDataMLive[lastKeyMLive].EndTime,
			X6: kLineDataMLive[lastKeyMLive].StartTime,

			Xc1: macdData[199].DIF,
			Xc2: macdData[199].DEA,
			Xc3: macdData[199].MACD,
			Xc4: macdData[199].Time,
		}

		res.DataListK = append(res.DataListK, tmpResDataListK)

		// ?????? ?????????
		if macdData[199].MACD > 0 && max1 < macdData[199].MACD { // ??????&???????????????
			maxMacd = macdData[199]
			maxK = kLineDataMLive[lastKeyMLive]
		}
		if macdData[199].MACD < 0 && low1 > macdData[199].MACD {
			lowMacd = macdData[199]
			lowK = kLineDataMLive[lastKeyMLive]
		}

		var (
			tmpAtr1Total float64
			tmpAtr1Avg   float64
			tmpAtr2Total float64
			tmpAtr2Avg   float64
		)
		// ??????atr1N???k?????????
		for i := 0; i <= atr1N-1; i++ {
			tmpAtr1Total += kLineDataMLive[lastKeyMLive-i].TopPrice - kLineDataMLive[lastKeyMLive-i].LowPrice
		}
		tmpAtr1Avg, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpAtr1Total/float64(atr1N)), 64)

		// ??????atr2N???k?????????
		for i := 0; i <= atr2N-1; i++ {
			tmpAtr2Total += kLineDataMLive[lastKeyMLive-i].TopPrice - kLineDataMLive[lastKeyMLive-i].LowPrice
		}
		tmpAtr2Avg, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpAtr2Total/float64(atr2N)), 64)

		// ??????

		for kOperationData, vOperationData := range operationData {
			if "open" != vOperationData.Status {
				continue
			}

			if "empty" == vOperationData.Type && "" == vOperationData.CloseStatus {
				if kLineDataMLive[lastKeyMLive].EndPrice <= vOperationData.ClosePriceWin || kLineDataMLive[lastKeyMLive].EndPrice >= vOperationData.ClosePriceLost { // ???????????????
					operationData[kOperationData].CloseStatus = "ok" // ???

					rate := (kLineDataMLive[lastKeyMLive].EndPrice - vOperationData.EndPrice) / vOperationData.EndPrice
					rate = -rate - 0.0003

					operationData = append(operationData, &OperationData3{
						StartTime:      kLineDataMLive[lastKeyMLive].StartTime,
						EndTime:        kLineDataMLive[lastKeyMLive].EndTime,
						StartPrice:     kLineDataMLive[lastKeyMLive].StartPrice,
						TopPrice:       kLineDataMLive[lastKeyMLive].TopPrice,
						LowPrice:       kLineDataMLive[lastKeyMLive].LowPrice,
						EndPrice:       kLineDataMLive[lastKeyMLive].EndPrice,
						Type:           "empty",
						Status:         "close",
						Rate:           rate,
						Tag:            vOperationData.StartTime,
						ClosePriceWin:  vOperationData.ClosePriceWin,  // ????????????
						ClosePriceLost: vOperationData.ClosePriceLost, // ????????????
					})
				}
			} else if "more" == vOperationData.Type && "" == vOperationData.CloseStatus {
				if kLineDataMLive[lastKeyMLive].EndPrice >= vOperationData.ClosePriceWin || kLineDataMLive[lastKeyMLive].EndPrice <= vOperationData.ClosePriceLost { // ???????????????
					operationData[kOperationData].CloseStatus = "ok" // ???
					rate := (kLineDataMLive[lastKeyMLive].EndPrice-vOperationData.EndPrice)/vOperationData.EndPrice - 0.0003
					operationData = append(operationData, &OperationData3{
						StartTime:      kLineDataMLive[lastKeyMLive].StartTime,
						EndTime:        kLineDataMLive[lastKeyMLive].EndTime,
						StartPrice:     kLineDataMLive[lastKeyMLive].StartPrice,
						TopPrice:       kLineDataMLive[lastKeyMLive].TopPrice,
						LowPrice:       kLineDataMLive[lastKeyMLive].LowPrice,
						EndPrice:       kLineDataMLive[lastKeyMLive].EndPrice,
						Type:           "more",
						Status:         "close",
						Rate:           rate,
						Tag:            vOperationData.StartTime,
						ClosePriceWin:  vOperationData.ClosePriceWin,  // ????????????
						ClosePriceLost: vOperationData.ClosePriceLost, // ????????????
					})
				}
			}

		}

		// ?????????????????????????????????
		if macdData[199].MACD > 0 && max2 < macdData[199].MACD && // ??????&???????????????
			lastMacd.MACD < macdData[199].MACD && // ????????????????????????????????????????????????????????????
			nil != maxMacd && nil != maxK &&
			macdData[199].MACD < maxMacd.MACD && kLineDataMLive[lastKeyMLive].TopPrice > maxK.TopPrice { // ??????????????????????????? ???????????????????????????????????????
			// ??????
			operationData = append(operationData, &OperationData3{
				StartTime:      kLineDataMLive[lastKeyMLive].StartTime,
				EndTime:        kLineDataMLive[lastKeyMLive].EndTime,
				StartPrice:     kLineDataMLive[lastKeyMLive].StartPrice,
				TopPrice:       kLineDataMLive[lastKeyMLive].TopPrice,
				LowPrice:       kLineDataMLive[lastKeyMLive].LowPrice,
				EndPrice:       kLineDataMLive[lastKeyMLive].EndPrice,
				Type:           "empty",
				Status:         "open",
				LastMacd:       lastMacd.MACD,
				Macd:           macdData[199].MACD,
				MaxKTopPrice:   maxK.TopPrice,
				MaxMacd:        maxMacd.MACD,
				Tag:            kLineDataMLive[lastKeyMLive].StartTime,
				ClosePriceWin:  kLineDataMLive[lastKeyMLive].EndPrice - tmpAtr1Avg, // ????????????
				ClosePriceLost: kLineDataMLive[lastKeyMLive].EndPrice + tmpAtr2Avg, // ????????????
			})
		}

		// ?????????????????????????????????
		if macdData[199].MACD < 0 && low2 > macdData[199].MACD && // ??????&???????????????
			lastMacd.MACD > macdData[199].MACD && // ????????????????????????????????????????????????????????????
			nil != lowMacd && nil != lowK &&
			macdData[199].MACD > lowMacd.MACD && kLineDataMLive[lastKeyMLive].LowPrice < lowK.LowPrice { // ??????????????????????????? ???????????????????????????????????????
			// ??????
			operationData = append(operationData, &OperationData3{
				StartTime:      kLineDataMLive[lastKeyMLive].StartTime,
				EndTime:        kLineDataMLive[lastKeyMLive].EndTime,
				StartPrice:     kLineDataMLive[lastKeyMLive].StartPrice,
				TopPrice:       kLineDataMLive[lastKeyMLive].TopPrice,
				LowPrice:       kLineDataMLive[lastKeyMLive].LowPrice,
				EndPrice:       kLineDataMLive[lastKeyMLive].EndPrice,
				Type:           "more",
				Status:         "open",
				LastMacd:       lastMacd.MACD,
				Macd:           macdData[199].MACD,
				LowKLowPrice:   lowK.LowPrice,
				LowMacd:        lowMacd.MACD,
				Tag:            kLineDataMLive[lastKeyMLive].StartTime,
				ClosePriceWin:  kLineDataMLive[lastKeyMLive].EndPrice + tmpAtr1Avg, // ????????????
				ClosePriceLost: kLineDataMLive[lastKeyMLive].EndPrice - tmpAtr2Avg, // ????????????
			})
		}

		lastMacd = macdData[199]

	}

	var (
		tmpWinTotal   int64
		tmpCloseTotal int64
		tmpRate       float64
		winRate       float64
	)

	for _, vOperationData := range operationData {
		if "" == vOperationData.CloseStatus && "open" == vOperationData.Status {
			continue
		}

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

		res.OperationData = append(res.OperationData, &v1.IntervalMAvgEndPriceMacdAndAtrDataReply_List2{
			StartPrice:     vOperationData.StartPrice,
			EndPrice:       vOperationData.EndPrice,
			StartTime:      vOperationData.StartTime,
			EndTime:        vOperationData.EndTime,
			Type:           vOperationData.Type,
			Status:         vOperationData.Status,
			Rate:           vOperationData.Rate,
			Tag:            vOperationData.Tag,
			LastMacd:       vOperationData.LastMacd,
			TopPrice:       vOperationData.TopPrice,
			LowPrice:       vOperationData.LowPrice,
			Macd:           vOperationData.Macd,
			LowKLowPrice:   vOperationData.LowKLowPrice,
			LowMacd:        vOperationData.LowMacd,
			MaxKTopPrice:   vOperationData.MaxKTopPrice,
			MaxMacd:        vOperationData.Macd,
			CloseStatus:    vOperationData.CloseStatus,
			ClosePriceWin:  vOperationData.ClosePriceWin,
			ClosePriceLost: vOperationData.ClosePriceLost,
		})
	}

	if 0 < tmpWinTotal && 0 < tmpCloseTotal {
		winRate = float64(tmpWinTotal) / float64(tmpCloseTotal)
	}
	res.OperationWinRate = fmt.Sprintf("%.2f", winRate)
	res.OperationWinAmount = strconv.FormatFloat(tmpRate, 'f', -1, 64)
	return res, nil
}

func (b *BinanceDataUsecase) OrderAreaPoint(ctx context.Context, req *v1.OrderAreaPointRequest, test string, endTime time.Time) (*v1.OrderAreaPointReply, error) {
	var (
		start         time.Time
		end           time.Time
		pointFirst    = 0.3
		pointInterval = 0.1
		kLineMOne     []*KLineMOne
		user          int64
		err           error
	)

	// ??????????????? todo ??????????????????????????????????????????
	if 0 < req.User {
		user = req.User
	}

	// ????????????n???15??????
	endNow := time.Now().UTC().Add(8 * time.Hour)
	if "test" == test {
		endNow = endTime
	}

	endM := endNow.Minute() / 15 * 15
	end = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), endNow.Hour(), endM, 59, 0, time.UTC).Add(-1 * time.Minute)
	start = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), endNow.Hour(), endM, 0, 0, time.UTC).Add(-3030 * time.Minute)

	// ??????????????????????????????????????????
	kLineMOne, err = b.klineMOneRepo.RequestBinanceMinuteKLinesData("ETHUSDT",
		strconv.FormatInt(start.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(end.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(15, 10)+"m",
		strconv.FormatInt(1500, 10))
	if nil != err {
		return nil, err
	}

	// ?????????????????????
	var (
		maPoint []*Ma
	)
	for kKLineMOne, _ := range kLineMOne {
		if kKLineMOne < 199 {
			continue
		}

		// ??????ma
		var (
			tmpTotalEndPrice float64
			tmpAvgEndPrice   float64
		)
		for i := 0; i <= 199; i++ {
			tmpTotalEndPrice += kLineMOne[kKLineMOne-i].EndPrice
		}

		tmpAvgEndPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.8f", tmpTotalEndPrice/float64(200)), 64)
		maPoint = append(maPoint, &Ma{AvgEndPrice: tmpAvgEndPrice})
	}

	// ?????????????????????????????????
	var (
		secondKeep    *OrderPolicyPointCompare
		thirdKeep     *OrderPolicyPointCompare
		secondConfirm *OrderPolicyPointCompare
		thirdConfirm  *OrderPolicyPointCompare
		second        *OrderPolicyPointCompare
		lastOrderInfo *OrderPolicyPointCompareInfo
	)
	// ??????????????????
	lastOrderInfo, err = b.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareInfo(user)
	if nil != err {
		return nil, err
	}

	if nil != lastOrderInfo {
		second, err = b.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "second", user)
		if nil != err {
			return nil, err
		}
		secondKeep, err = b.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "second_keep", user)
		if nil != err {
			return nil, err
		}
		thirdKeep, err = b.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "third_keep", user)
		if nil != err {
			return nil, err
		}
		secondConfirm, err = b.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "second_confirm", user)
		if nil != err {
			return nil, err
		}
		thirdConfirm, err = b.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "third_confirm", user)
		if nil != err {
			return nil, err
		}
	}

	// ??????????????????
	tmpPointFirstSub := maPoint[2].AvgEndPrice - maPoint[1].AvgEndPrice
	tmpPointSecondSub := maPoint[1].AvgEndPrice - maPoint[0].AvgEndPrice

	// ????????????
	var (
		order        []*OrderData
		pointCompare []*OrderPolicyPointCompare
	)

	if tmpPointSecondSub > pointFirst && pointFirst > tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "buy_long" == lastOrderInfo.Type { // ??????
				// ????????????????????????????????????
				tmpOpen := false
				if nil != secondKeep {
					if 2 <= secondKeep.Value {
						tmpOpen = true
					} else {
						// ???????????????0?????????
						pointCompare = append(pointCompare, &OrderPolicyPointCompare{
							InfoId: secondKeep.InfoId,
							Type:   secondKeep.Type,
							Value:  0,
						})
					}
				} else if nil != thirdKeep {
					tmpOpen = true
				} else if nil != thirdConfirm {
					tmpOpen = true
				} else if nil != secondConfirm {
					tmpOpen = true
				}

				if tmpOpen {
					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})
					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})
				}

			} else if "sell_short" == lastOrderInfo.Type {
				// ????????????????????????????????????????????????????????????
				if nil != second {
					// ???????????????0?????????
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: second.InfoId,
						Type:   second.Type,
						Value:  0,
					})
				}
			}

		} else {
			// ??????
			tmpAmount := 0.01
			tmpAmountStr := "0.01"
			if 1 == user {
				tmpAmount = 0.1
				tmpAmountStr = "0.1"
			}
			order = append(order, &OrderData{
				Symbol:          "ETHUSDT",
				Side:            "SELL",
				OrderType:       "MARKET",
				PositionSide:    "SHORT",
				Quantity:        tmpAmountStr,
				QuantityFloat64: tmpAmount,
			})
		}
	}

	if tmpPointSecondSub < -pointFirst && -pointFirst < tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "sell_short" == lastOrderInfo.Type { // ??????
				// ????????????????????????????????????
				tmpOpen := false
				if nil != secondKeep {
					if 2 <= secondKeep.Value {
						tmpOpen = true
					} else {
						// ???????????????0?????????
						pointCompare = append(pointCompare, &OrderPolicyPointCompare{
							InfoId: secondKeep.InfoId,
							Type:   secondKeep.Type,
							Value:  0,
						})
					}
				} else if nil != thirdKeep {
					tmpOpen = true
				} else if nil != thirdConfirm {
					tmpOpen = true
				} else if nil != secondConfirm {
					tmpOpen = true
				}

				if tmpOpen {
					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})
					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})
				}

			} else if "buy_long" == lastOrderInfo.Type {
				// ????????????????????????????????????????????????????????????
				if nil != second {
					// ???????????????0?????????
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: second.InfoId,
						Type:   second.Type,
						Value:  0,
					})
				}
			}

		} else {
			// ??????
			tmpAmount := 0.01
			tmpAmountStr := "0.01"
			if 1 == user {
				tmpAmount = 0.1
				tmpAmountStr = "0.1"
			}
			order = append(order, &OrderData{
				Symbol:          "ETHUSDT",
				Side:            "BUY",
				OrderType:       "MARKET",
				PositionSide:    "LONG",
				Quantity:        tmpAmountStr,
				QuantityFloat64: tmpAmount,
			})
		}
	}

	// ????????????
	if (pointFirst + pointInterval) < tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "sell_short" == lastOrderInfo.Type {
				// ??????????????????????????????
				// ???????????????
				if nil == second {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: lastOrderInfo.ID,
						Type:   "second",
						Value:  1,
					})
				} else if 0 == second.Value {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: second.InfoId,
						Type:   second.Type,
						Value:  1,
					})
				} else if 1 <= second.Value {
					// ????????????????????????
					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// ??????lastOrder???????????????
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: 0, // ?????? ?????????????????? ???0??????
						Type:   "second_confirm",
						Value:  1,
					})
				}
			} else if "buy_long" == lastOrderInfo.Type {
				// ??????????????????????????????????????????
				// ???????????????
				if nil == secondKeep {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: lastOrderInfo.ID,
						Type:   "second_keep",
						Value:  1,
					})
				} else {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: secondKeep.InfoId,
						Type:   secondKeep.Type,
						Value:  secondKeep.Value + 1,
					})
				}
			}
		}
	}

	if (-pointFirst - pointInterval) > tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "buy_long" == lastOrderInfo.Type {
				// ??????????????????????????????
				// ???????????????
				if nil == second {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: lastOrderInfo.ID,
						Type:   "second",
						Value:  1,
					})
				} else if 0 == second.Value {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: second.InfoId,
						Type:   second.Type,
						Value:  1,
					})
				} else if 1 <= second.Value {
					// ????????????????????????
					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// ??????
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// ??????lastOrder???????????????
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: 0, // ?????? ?????????????????? ???0??????
						Type:   "second_confirm",
						Value:  1,
					})
				}
			} else if "sell_short" == lastOrderInfo.Type {
				// ??????????????????????????????????????????
				// ???????????????
				if nil == secondKeep {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: lastOrderInfo.ID,
						Type:   "second_keep",
						Value:  1,
					})
				} else {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: secondKeep.InfoId,
						Type:   secondKeep.Type,
						Value:  secondKeep.Value + 1,
					})
				}
			}
		}
	}

	// ???????????????????????????
	if (pointFirst + 2*pointInterval) < tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "sell_short" == lastOrderInfo.Type {
				// ??????
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "BUY",
					OrderType:       "MARKET",
					PositionSide:    "SHORT",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// ??????
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "BUY",
					OrderType:       "MARKET",
					PositionSide:    "LONG",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// ????????????
				// ??????lastOrder???????????????
				pointCompare = append(pointCompare, &OrderPolicyPointCompare{
					InfoId: 0, // ?????? ?????????????????? ???0??????
					Type:   "third_confirm",
					Value:  1,
				})
			} else if "buy_long" == lastOrderInfo.Type {
				if nil == thirdKeep {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: lastOrderInfo.ID,
						Type:   "third_keep",
						Value:  1,
					})
				}
			}
		}
	}

	if (-pointFirst - 2*pointInterval) > tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "buy_long" == lastOrderInfo.Type {
				// ??????
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "SELL",
					OrderType:       "MARKET",
					PositionSide:    "LONG",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// ??????
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "SELL",
					OrderType:       "MARKET",
					PositionSide:    "SHORT",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// ????????????
				// ??????lastOrder???????????????
				pointCompare = append(pointCompare, &OrderPolicyPointCompare{
					InfoId: 0, // ?????? ?????????????????? ???0??????
					Type:   "third_confirm",
					Value:  1,
				})
			} else if "sell_short" == lastOrderInfo.Type {
				if nil == thirdKeep {
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: lastOrderInfo.ID,
						Type:   "third_keep",
						Value:  1,
					})
				}
			}
		}
	}

	// ??????
	var (
		orderData    *OrderPolicyPointCompareInfo
		orderBinance *Order
		newOrderData *OrderPolicyPointCompareInfo
	)
	for _, v := range order {
		fmt.Println(start, end)
		fmt.Println(v, 22)
		fmt.Println(tmpPointSecondSub, tmpPointFirstSub, 23323)
		for _, vMaPoint := range maPoint {
			fmt.Println(vMaPoint, 233235)
		}
		orderBinance, err = b.orderPolicyPointCompareRepo.RequestBinanceOrder(v.Symbol, v.Side, v.OrderType, v.PositionSide, v.Quantity, user)
		if nil != err {
			b.log.Error(err)
			return nil, err
		}

		orderData = &OrderPolicyPointCompareInfo{
			OrderId: orderBinance.OrderId,
			Type:    "",
			Num:     v.QuantityFloat64,
		}

		// ?????????????????????????????????????????????
		if "BUY" == v.Side {
			if "LONG" == v.PositionSide {
				orderData.Type = "buy_long"
			} else {
				orderData.Type = "buy_short"
			}
		} else {
			if "LONG" == v.PositionSide {
				orderData.Type = "sell_long"
			} else {
				orderData.Type = "sell_short"
			}
		}

	}

	if err = b.tx.ExecTx(ctx, func(ctx context.Context) error { // ??????
		// ??????????????????
		if nil != orderData {
			newOrderData, err = b.orderPolicyPointCompareRepo.InsertOrderPolicyPointCompareInfo(ctx, orderData, user)
			if nil != err {
				return err
			}
		}

		if 0 < len(pointCompare) {
			for _, vPointCompare := range pointCompare {
				tmpInfoId := vPointCompare.InfoId
				if "second_confirm" == vPointCompare.Type || "third_confirm" == vPointCompare.Type {
					if nil == newOrderData || 0 >= newOrderData.ID {
						// ?????????????????????
						return errors.New(500, "????????????", "??????????????????????????????")
					}
					tmpInfoId = newOrderData.ID
				}

				tmpOrderPolicyPointCompare := &OrderPolicyPointCompare{
					InfoId: tmpInfoId,
					Type:   vPointCompare.Type,
					Value:  vPointCompare.Value,
				}

				//fmt.Println(tmpOrderPolicyPointCompare, 44)
				_, err = b.orderPolicyPointCompareRepo.InsertOrderPolicyPointCompare(ctx, tmpOrderPolicyPointCompare, user)
				if nil != err {
					return err
				}
			}
		}

		// ??????????????????
		return nil
	}); err != nil {
		b.log.Error(err)
	}

	return &v1.OrderAreaPointReply{}, nil
}

func (b *BinanceDataUsecase) PullBinanceData(ctx context.Context, req *v1.PullBinanceDataRequest) (*v1.PullBinanceDataReply, error) {
	var (
		start                time.Time
		end                  time.Time
		tmpKlineMOne         []*KLineMOne
		lastKlineMOne        *KLineMOne
		lastKlineMOneEndTime time.Time
		m                    = int64(1)
		limit                = int64(1500)
		err                  error
	)

	start, err = time.Parse("2006-01-02 15:04:05", req.Start)       // ????????????????????????
	end = time.Now().UTC().Add(8 * time.Hour).Add(-1 * time.Minute) // ????????????
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), 59, 0, time.UTC)
	if nil != err {
		return nil, err
	}

	// ??????????????????????????????????????????
	if "BTCUSDT" == req.Coin {
		lastKlineMOne, err = b.klineMOneRepo.GetKLineMOneOrderByEndTimeLast()
	} else if "FILUSDT" == req.Coin {
		lastKlineMOne, err = b.klineMOneRepo.GetFilKLineMOneOrderByEndTimeLast()
	} else if "ETHUSDT" == req.Coin {
		lastKlineMOne, err = b.klineMOneRepo.GetEthKLineMOneOrderByEndTimeLast()
	}
	if nil != lastKlineMOne {
		lastKlineMOneEndTime = time.UnixMilli(lastKlineMOne.EndTime).UTC().Add(8 * time.Hour).Add(1 * time.Millisecond)
		if start.Before(lastKlineMOneEndTime) { // ???????????????????????????????????????
			start = lastKlineMOneEndTime
		}
	}

	tmpStart := start
	for {
		var tmpEnd time.Time
		if end.After(tmpStart.Add(time.Duration(m*limit) * time.Minute)) {
			tmpEnd = tmpStart.Add(time.Duration(m*limit) * time.Minute).Add(-1 * time.Millisecond)
		} else {
			tmpEnd = end
		}

		if tmpEnd.Before(tmpStart) { // ?????????
			break
		}

		// BTCUSDT
		tmpKlineMOne, err = b.klineMOneRepo.RequestBinanceMinuteKLinesData(req.Coin,
			strconv.FormatInt(tmpStart.Add(-8*time.Hour).UnixMilli(), 10),
			strconv.FormatInt(tmpEnd.Add(-8*time.Hour).UnixMilli(), 10),
			strconv.FormatInt(m, 10)+"m",
			strconv.FormatInt(limit, 10))
		if nil != err {
			return nil, err
		}

		if err = b.tx.ExecTx(ctx, func(ctx context.Context) error { // ??????

			if "BTCUSDT" == req.Coin {
				_, err = b.klineMOneRepo.InsertKLineMOne(ctx, tmpKlineMOne)
			} else if "FILUSDT" == req.Coin {
				_, err = b.klineMOneRepo.InsertFilKLineMOne(ctx, tmpKlineMOne)
			} else if "ETHUSDT" == req.Coin {
				_, err = b.klineMOneRepo.InsertEthKLineMOne(ctx, tmpKlineMOne)
			}

			if nil != err {
				return err
			}

			return nil
		}); err != nil {
			fmt.Println(err)
			break
		}

		tmpStart = tmpStart.Add(time.Duration(m*limit) * time.Minute)
		if end.Before(tmpStart) {
			break
		}
	}

	return &v1.PullBinanceDataReply{}, nil
}

func handleManMnWithKLineMineData(n int, interval int, current int, kKlineMOne int, vKlineMOne *KLineMOne, klineMOne []*KLineMOne) *Ma {
	var (
		need                int
		tmpMaNEndPriceTotal float64 // ?????????????????????????????????
	)

	tmp := 0
	if current%interval > 0 {
		tmp = current % interval
	} else if current/interval > 0 {
		tmp = interval
	} else {
		tmp = 1
	}

	need = (n-1)*interval + tmp

	//fmt.Println(need, current%interval)

	for i := need - 1; i >= 0; i-- {
		//fmt.Println(klineMOne[kKlineMOne-i], i)
		// ??????
		if 0 == (need-i)%interval {
			tmpMaNEndPriceTotal += klineMOne[kKlineMOne-i].EndPrice // ??????
			//if need <= 25 {
			//	fmt.Println(need, time.UnixMilli(klineMOne[kKlineMOne-i].EndTime))
			//}

		} else if 0 == i {
			tmpMaNEndPriceTotal += klineMOne[kKlineMOne-i].EndPrice // ??????
			//if need <= 25 {
			//	fmt.Println(need, time.UnixMilli(klineMOne[kKlineMOne-i].EndTime))
			//}
		}
	}

	tmpMaNAvgEndPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", tmpMaNEndPriceTotal/float64(n)), 64)
	return &Ma{
		AvgEndPrice: tmpMaNAvgEndPrice,
	}
}

func handleMKData(klineMOne []*KLineMOne, m int) []*KLineMOne {
	var (
		klineM []*KLineMOne
	)
	// ??????????????????????????????????????????
	subK := 0
	for kKlineMOne, vKlineMOne := range klineMOne {
		// ????????????
		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		if tmpNow.Minute()%m == 0 {
			subK = kKlineMOne
			break
		}
	}
	tmpKlineMOne := klineMOne[subK:]
	lenTmpKlineMOne := len(tmpKlineMOne)
	for tmpKKlineMOne, tmpVKlineMOne := range tmpKlineMOne {
		// ????????????
		tmpNow := time.UnixMilli(tmpVKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		if tmpNow.Minute()%m != 0 {
			continue
		}

		tmpKlineM := &KLineMOne{
			StartPrice: tmpVKlineMOne.StartPrice,
			StartTime:  tmpVKlineMOne.StartTime,
			TopPrice:   tmpVKlineMOne.TopPrice,
			LowPrice:   tmpVKlineMOne.LowPrice,
			EndPrice:   tmpVKlineMOne.EndPrice,
			EndTime:    tmpVKlineMOne.EndTime,
		}

		tmpLimitK := tmpKKlineMOne + m - 1
		if lenTmpKlineMOne-1 < tmpLimitK {
			tmpLimitK = lenTmpKlineMOne - 1
		}
		for i := tmpKKlineMOne; i <= tmpLimitK; i++ {
			if nil == tmpKlineMOne[i] {
				break
			}
			if tmpKlineMOne[i].TopPrice > tmpKlineM.TopPrice {
				tmpKlineM.TopPrice = tmpKlineMOne[i].TopPrice
			}
			if tmpKlineMOne[i].LowPrice < tmpKlineM.LowPrice {
				tmpKlineM.LowPrice = tmpKlineMOne[i].LowPrice
			}
			tmpKlineM.EndPrice = tmpKlineMOne[i].EndPrice
			tmpKlineM.EndTime = tmpKlineMOne[i].EndTime
		}
		klineM = append(klineM, tmpKlineM)
	}

	return klineM
}
