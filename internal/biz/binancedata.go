package biz

import (
	v1 "binancedata/api/binancedata/v1"
	"context"
	"fmt"
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

type OperationData2 struct {
	StartTime     int64
	EndTime       int64
	StartPrice    float64
	TopPrice      float64
	LowPrice      float64
	EndPrice      float64
	AvgEndPrice   float64
	Amount        int64
	Type          string
	Status        string
	Action        string
	CloseEndPrice string
	Rate          float64
	CloseSubPrice float64
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

type KLineMOneRepo interface {
	GetKLineMOneOrderByEndTimeLast() (*KLineMOne, error)
	GetKLineMOneByStartTime(start int64, end int64) ([]*KLineMOne, error)
	InsertKLineMOne(ctx context.Context, kLineMOne []*KLineMOne) (bool, error)
	RequestBinanceMinuteKLinesData(symbol string, startTime string, endTime string, interval string, limit string) ([]*KLineMOne, error)
	NewMACDData(list []*KLineMOne) ([]*MACDPoint, error)
}

// BinanceDataUsecase is a BinanceData usecase.
type BinanceDataUsecase struct {
	klineMOneRepo KLineMOneRepo
	repo          BinanceDataRepo
	tx            Transaction
	log           *log.Helper
}

// NewBinanceDataUsecase new a BinanceData usecase.
func NewBinanceDataUsecase(repo BinanceDataRepo, klineMOneRepo KLineMOneRepo, tx Transaction, logger log.Logger) *BinanceDataUsecase {
	return &BinanceDataUsecase{repo: repo, klineMOneRepo: klineMOneRepo, tx: tx, log: log.NewHelper(logger)}
}

// XNIntervalMAvgEndPriceData x个<n个间隔m时间的平均收盘价>数据 .
func (b *BinanceDataUsecase) XNIntervalMAvgEndPriceData(ctx context.Context, req *v1.XNIntervalMAvgEndPriceDataRequest) (*v1.XNIntervalMAvgEndPriceDataReply, error) {
	fmt.Println(req)
	var (
		reqStart         time.Time
		reqEnd           time.Time
		klineMOne        []*KLineMOne
		n1               int
		n2               int
		ma5M15           []*Ma // todo 各种图优化
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

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.SendBody.Start) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.SendBody.End) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	n1 = int(req.SendBody.N1)
	n2 = int(req.SendBody.N2)

	var (
		maxMxN = int64(n2 * 60) // todo 现在是写死的，改成可识别参与计算数据最大值 60分钟的参数是最大的
	)

	// 数据时间限制
	dataLimitTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		reqStart = dataLimitTime
	}

	//x = req.SendBody.X
	//for _, vX := range x {
	//	fmt.Println(vX.M, vX.Method, vX.N)
	//	if vX.N*vX.M > maxMxN { // n条m分钟的均线，n*m分钟
	//		maxMxN = vX.N * vX.M
	//	}
	//}

	// 获取时间范围内的k线分钟数据
	//if 1 <= maxMxN {
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	//}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneByStartTime(
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

	// 遍历k线数据，也是每分钟数据

	var (
		openActionTag    string
		lastActionTag    string
		tmpLastActionTag string
	)
	operationData = make(map[string]*OperationData2, 0)
	for kKlineMOne, vKlineMOne := range klineMOne {
		var tagNum int64
		if maxMxN-1 > int64(kKlineMOne) { // 多的线，注意查到的前置数据个数>=maxMxN，不然有bug
			continue
		}

		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC)
		//fmt.Println(tmpNow, tmpNow0, tmpNow.Sub(tmpNow0).Minutes())
		tmpNowSubNow0 := int(tmpNow.Sub(tmpNow0).Minutes()) + 1 // 这里都是开始时间的差值，而我们需要知道这个范围的差值，例如00：00-01：00我们需要的是2，两个1分钟
		// todo 改成可调节
		// 遍历分钟线
		//for _, vX := range x {
		//
		//}

		// 计算5根5分钟线
		tmpMa5M5 := handleManMnWithKLineMineData(n1, 5, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma5M5 = append(ma5M5, tmpMa5M5)
		// 计算5根10分钟线
		tmpMa10M5 := handleManMnWithKLineMineData(n2, 5, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma10M5 = append(ma10M5, tmpMa10M5)
		//// 计算5根15分钟线
		tmpMa5M15 := handleManMnWithKLineMineData(n1, 15, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma5M15 = append(ma5M15, tmpMa5M15)
		//// 计算10根15分钟线
		tmpMa10M15 := handleManMnWithKLineMineData(n2, 15, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma10M15 = append(ma10M15, tmpMa10M15)
		//// 计算5根60分钟线
		tmpMa5M60 := handleManMnWithKLineMineData(n1, 60, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma5M60 = append(ma5M60, tmpMa5M60)
		//// 计算10根60分钟线
		tmpMa10M60 := handleManMnWithKLineMineData(n2, 60, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		ma10M60 = append(ma10M60, tmpMa10M60)

		// 开 全/半的空多 平 仓
		// 相交 ma5M15和ma5M15
		if maxMxN < int64(kKlineMOne) { // 第一单跳过
			//fmt.Println(kKlineMOne, vKlineMOne)
			//fmt.Println("ma10m15", tmpMa10M15)
			//fmt.Println("ma5m15", tmpMa5M15)
			//fmt.Println("ma5m5", tmpMa5M5)
			//fmt.Println("ma10m5", tmpMa10M5)
			//fmt.Println("ma5m60", tmpMa5M60)
			//fmt.Println("ma10m60", tmpMa10M60)

			// 平仓 立即更新到操作数据
			// 平空仓
			if tmpMa5M15.AvgEndPrice > tmpMa10M15.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// 本次关，下单逻辑判断上一单
					tmpDo := false // 假定为开仓数据，可能是全/半仓的
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

			// 平多仓
			if tmpMa5M15.AvgEndPrice < tmpMa10M15.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// 本次关，下单逻辑判断上一单
					tmpDo := false // 假定为开仓数据，可能是全/半仓的
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

			// 开多
			//fmt.Println(len(ma5M15))
			if tmpMa5M15.AvgEndPrice > tmpMa10M15.AvgEndPrice { // 条件1
				lastMa5M15 := ma5M15[len(ma5M15)-2]    // 上一单
				lastMa10M15 := ma10M15[len(ma10M15)-2] // 上一单
				//last2Ma5M15 := ma5M15[len(ma5M15)-3]    // 上一单
				//last2Ma10M15 := ma10M15[len(ma10M15)-3] // 上一单
				//fmt.Println("last_ma10m15", tmpMa10M15)
				//fmt.Println("last_ma5m15", lastMa5M15)
				//fmt.Println("last2_ma5m15", last2Ma5M15)
				//fmt.Println("last2_ma10m15", last2Ma10M15)
				if lastMa5M15.AvgEndPrice < lastMa10M15.AvgEndPrice { // 条件1
					if 0 < klineMOne[kKlineMOne-1].EndPrice-klineMOne[kKlineMOne-15].StartPrice { //  一定有15根了，条件判断阳线
						//if tmpMa5M60.AvgEndPrice > tmpMa10M60.AvgEndPrice { // 条件2
						// 本次开，下单逻辑判断上一单
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
								Status:      "open", // 全开状态
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

			// 平多半仓
			if tmpMa5M5.AvgEndPrice < tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// 本次关，下单逻辑判断上一单
					tmpDo := false // 假定为开仓数据，可能是全/半仓的
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

			// 加半仓
			if tmpMa5M15.AvgEndPrice > tmpMa10M15.AvgEndPrice && tmpMa5M5.AvgEndPrice > tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// 本次开，下单逻辑判断上一单
					tmpDo := false // 假定为开仓数据，可能是全/半仓的
					if "more" == tmpOpenLastOperationData2.Type && "half" == tmpOpenLastOperationData2.Status {
						tmpDo = true
					}
					if tmpDo {
						// 本次开
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

			// 开空
			if tmpMa5M15.AvgEndPrice < tmpMa10M15.AvgEndPrice { // 条件1
				lastMa5M15 := ma5M15[len(ma5M15)-2]    // 上一单
				lastMa10M15 := ma10M15[len(ma10M15)-2] // 上一单
				//last2Ma5M15 := ma5M15[len(ma5M15)-3]                  // 上一单
				//last2Ma10M15 := ma10M15[len(ma10M15)-3]               // 上一单
				if lastMa5M15.AvgEndPrice > lastMa10M15.AvgEndPrice { // 条件1
					if 0 > klineMOne[kKlineMOne-1].EndPrice-klineMOne[kKlineMOne-15].StartPrice { //  一定有15根了，条件判断阴线
						//if tmpMa5M60.AvgEndPrice < tmpMa10M60.AvgEndPrice { // 条件2
						// 本次开，下单逻辑判断上一单
						tmpDo := false
						if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
							if "close" == tmpOpenLastOperationData2.Status {
								tmpDo = true
							}
						} else {
							tmpDo = true
						}
						// 本次开
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

			// 平空半仓
			if tmpMa5M5.AvgEndPrice > tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// 本次关，下单逻辑判断上一单
					tmpDo := false // 假定为开仓数据，可能是全/半仓的
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

			// 加半仓
			if tmpMa5M15.AvgEndPrice < tmpMa10M15.AvgEndPrice && tmpMa5M5.AvgEndPrice < tmpMa10M5.AvgEndPrice {
				if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					// 本次关，下单逻辑判断上一单
					tmpDo := false // 假定为开仓数据，可能是全/半仓的
					if "empty" == tmpOpenLastOperationData2.Type && "half" == tmpOpenLastOperationData2.Status {
						tmpDo = true
					}

					if tmpDo {
						// 本次开
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

	// 排序
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

	// 得到最后一个关仓
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for k, vOperationData := range resOperationData {
		if k > tmpLastCloseK { // 结束查询到最后一个，默认-1不会被查到
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

// KAnd2NIntervalMAvgEndPriceData k线和x个<2个间隔m时间的平均收盘价>数据 .
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
		maNMFirst           []*Ma // todo 各种图优化
		maNMSecond          []*Ma
		resOperationData    OperationData2Slice
		closeCondition      = 1
		closeCondition2Rate float64
		err                 error
		//x         []*v1.XNIntervalMAvgEndPriceDataRequest_SendBody_List
	)

	// 返回结果
	res := &v1.KAnd2NIntervalMAvgEndPriceDataReply{
		DataListK:          make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListK, 0),
		DataListMaNMFirst:  make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMFirst, 0),
		DataListMaNMSecond: make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMSecond, 0),
		BackGround:         make([]*v1.KAnd2NIntervalMAvgEndPriceDataReply_ListBackGround, 0),
	}

	// 简单的参数限制，解决程序达不到的操作
	reqStart, err = time.Parse("2006-01-02 15:04:05", req.SendBody.Start) // 时间进行格式校验
	if nil != err {
		return res, nil
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.SendBody.End) // 时间进行格式校验
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

	// 条件2目前是止盈止损点
	if 2 == req.SendBody.CloseCondition {
		if maxMxN < 5 {
			return res, nil
		}
		closeCondition = 2
		closeCondition2Rate = req.SendBody.CloseCondition2Rate
	}

	// 获取时间范围内的k线分钟数据
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo 数据时间限制，先应该随着maxMxN改变而改变
	dataLimitTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// 时间查不出数据
	if reqStart.After(reqEnd) {
		return res, nil
	}

	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneByStartTime(
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
	// 遍历k线数据，也是每分钟数据
	for kKlineMOne, vKlineMOne := range klineMOne {
		var tagNum int64
		if maxMxN-1 > int64(kKlineMOne) { // 多的线，注意查到的前置数据个数>=maxMxN，不然有bug
			continue
		}

		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC)
		//fmt.Println(tmpNow, tmpNow0, tmpNow.Sub(tmpNow0).Minutes())
		tmpNowSubNow0 := int(tmpNow.Sub(tmpNow0).Minutes()) + 1 // 这里都是开始时间的差值，而我们需要知道这个范围的差值，例如00：00-01：00我们需要的是2，两个1分钟

		// 计算N根M分钟线，第一条
		tmpMaNMFirst := handleManMnWithKLineMineData(n1, m1, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		maNMFirst = append(maNMFirst, tmpMaNMFirst)
		res.DataListMaNMFirst = append(res.DataListMaNMFirst, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMFirst{X1: tmpMaNMFirst.AvgEndPrice})

		// 计算N根M分钟线，第二条
		tmpMaNMSecond := handleManMnWithKLineMineData(n2, m2, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		maNMSecond = append(maNMSecond, tmpMaNMSecond)
		res.DataListMaNMSecond = append(res.DataListMaNMSecond, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListMaNMSecond{X1: tmpMaNMSecond.AvgEndPrice})

		// 背景颜色
		tmpBackGround := "white"

		// 第一单跳过
		if maxMxN < int64(kKlineMOne) {
			// 平多
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				// 本次关，下单逻辑判断上一单
				tmpDo := false // 假定为开仓数据，可能是全/半仓的
				if "more" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					tmpBackGround = "green"

					if 1 == closeCondition {
						// 更新比较的最高价
						if vKlineMOne.TopPrice > compareTopPrice {
							compareTopPrice = vKlineMOne.TopPrice
						}

						// 最高价-最高价*x% > 最新价 关仓条件
						if compareTopPrice-compareTopPrice*topX > vKlineMOne.EndPrice {
							tmpDo = true
						}
					} else if 2 == closeCondition {
						// 止损
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

			// 平空
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				// 本次关，下单逻辑判断上一单
				tmpDo := false // 假定为开仓数据，可能是全/半仓的
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					tmpBackGround = "red"

					if 1 == closeCondition {
						// 更新比较的最高价
						if vKlineMOne.LowPrice < compareLowPrice {
							compareLowPrice = vKlineMOne.LowPrice
						}

						// 最高价-最高价*x% > 最新价 关仓条件
						if compareLowPrice+compareLowPrice*lowX < vKlineMOne.EndPrice {
							tmpDo = true
						}
					} else if 2 == closeCondition {
						// 止损
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

			// 开多
			//fmt.Println(len(ma5M15))
			if tmpMaNMFirst.AvgEndPrice > tmpMaNMSecond.AvgEndPrice { // 条件1
				lastMaNMFirst := maNMFirst[len(maNMFirst)-2]                // 上一单
				lastMaNMSecond := maNMSecond[len(maNMSecond)-2]             // 上一单
				if lastMaNMFirst.AvgEndPrice < lastMaNMSecond.AvgEndPrice { // 条件1
					// 本次开，下单逻辑判断上一单
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
					// 平空 如果有先平掉上一单
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

					// 找到前5k线数据的最低价
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
						Status:        "open", // 全开状态
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

			// 开空
			//fmt.Println(len(ma5M15))
			if tmpMaNMFirst.AvgEndPrice < tmpMaNMSecond.AvgEndPrice { // 条件1
				lastMaNMFirst := maNMFirst[len(maNMFirst)-2]                // 上一单
				lastMaNMSecond := maNMSecond[len(maNMSecond)-2]             // 上一单
				if lastMaNMFirst.AvgEndPrice > lastMaNMSecond.AvgEndPrice { // 条件1
					// 本次开，下单逻辑判断上一单
					//tmpDo := false
					//if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
					//	if "close" == tmpOpenLastOperationData2.Status {
					//		tmpDo = true
					//	}
					//} else {
					//	tmpDo = true
					//}

					//if tmpDo {
					// 平多 如果有先平掉上一单
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

					// 找到前5k线数据的最低价
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
						Status:        "open", // 全开状态
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

		// 结果
		res.DataListK = append(res.DataListK, &v1.KAnd2NIntervalMAvgEndPriceDataReply_ListK{
			X1: vKlineMOne.StartPrice,
			X2: vKlineMOne.EndPrice,
			X3: vKlineMOne.TopPrice,
			X4: vKlineMOne.LowPrice,
			X5: vKlineMOne.EndTime,
		})
	}

	// 排序
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

	// 得到最后一个关仓
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for k, vOperationData := range resOperationData {
		if k > tmpLastCloseK { // 结束查询到最后一个，默认-1不会被查到
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

// IntervalMAvgEndPriceData k线和间隔m时间的平均收盘价数据 .
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

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // 时间进行格式校验
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

	// 获取时间范围内的k线分钟数据
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo 数据时间限制，先应该随着maxMxN改变而改变
	dataLimitTime := time.Date(2020, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// 时间查不出数据
	if reqStart.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	var (
		lastActionTag string
		operationData map[string]*OperationData2
		maNMFirst     []*Ma // todo 各种图优化
	)
	operationData = make(map[string]*OperationData2, 0)
	// 遍历数据
	for kKlineMOne, vKlineMOne := range klineMOne {
		var tagNum int64
		if maxMxN-1 > kKlineMOne { // 多的线，注意查到的前置数据个数>=maxMxN，不然有bug
			continue
		}

		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC)
		//fmt.Println(tmpNow, tmpNow0, tmpNow.Sub(tmpNow0).Minutes())
		tmpNowSubNow0 := int(tmpNow.Sub(tmpNow0).Minutes()) + 1 // 这里都是开始时间的差值，而我们需要知道这个范围的差值，例如00：00-01：00我们需要的是2，两个1分钟

		// 计算N根M分钟线，第一条
		tmpMaNMFirst := handleManMnWithKLineMineData(n, m, tmpNowSubNow0, kKlineMOne, vKlineMOne, klineMOne)
		maNMFirst = append(maNMFirst, tmpMaNMFirst)
		res.DataListMaNMFirst = append(res.DataListMaNMFirst, &v1.IntervalMAvgEndPriceDataReply_ListMaNMFirst{X1: tmpMaNMFirst.AvgEndPrice})

		// 第一单跳过
		if maxMxN < kKlineMOne {
			// 亏损关单
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "open" == tmpOpenLastOperationData2.Status {

					if "empty" == tmpOpenLastOperationData2.Type {
						tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineMOne.EndPrice)/tmpOpenLastOperationData2.EndPrice - fee
						if tmpRate < -targetCloseRate {
							// 关
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
							// 关
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

			// 开多
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
							Status:     "open", // 全开状态
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
							Status:     "open", // 全开状态
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
						Status:     "open", // 全开状态
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				}

			}

			// 开空
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
							Status:     "open", // 全开状态
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
							Status:     "open", // 全开状态
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
						Status:     "open", // 全开状态
					}
					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineMOne.EndTime, 10)
					operationData[lastActionTag] = currentOperationData
				}
			}
		}

		// 结果
		res.DataListK = append(res.DataListK, &v1.IntervalMAvgEndPriceDataReply_ListK{
			X1: vKlineMOne.StartPrice,
			X2: vKlineMOne.EndPrice,
			X3: vKlineMOne.TopPrice,
			X4: vKlineMOne.LowPrice,
			X5: vKlineMOne.EndTime,
		})
	}

	// 排序
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

	// 得到最后一个关仓
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for k, vOperationData := range resOperationData {
		if k > tmpLastCloseK { // 结束查询到最后一个，默认-1不会被查到
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

// IntervalMMACDData k线和间隔m时间的平均收盘价数据 .
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

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // 时间进行格式校验
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
	maxMxN := 200 * m // macd计算，至少需要数据源头数据条数

	// 获取时间范围内的k线分钟数据
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo 数据时间限制，先应该随着maxMxN改变而改变
	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// 时间查不出数据
	if reqStart.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	// 截取掉不需要的数据，之后计算
	klineM := handleMKData(klineMOne, m)

	//fmt.Println(len(klineM), klineM[199], len(tmpKlineMOne), tmpKlineMOne[199])

	var (
		macdData []*MACDPoint
	)
	operationData := make(map[string]*OperationData2, 0)

	// macd返回数据
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

	// 遍历数据
	var (
		lastActionTag string
		closeEmptyTag int
		closeMoreTag  int
	)
	tmpKlineM := klineM[maxMxN/m-1:]
	for _, vKlineM := range tmpKlineM {
		var tagNum int64

		// 结果
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

		// 关多
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
					// 关
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

		// 关空
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
					// 关
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

		// 开多
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
						Status:     "open", // 全开状态
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
						Status:     "open", // 全开状态
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
					Status:     "open", // 全开状态
				}
				tagNum++
				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				operationData[lastActionTag] = currentOperationData
			}
		}

		// 开空
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
						Status:     "open", // 全开状态
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
						Status:     "open", // 全开状态
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
					Status:     "open", // 全开状态
				}
				tagNum++
				lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				operationData[lastActionTag] = currentOperationData
			}

		}

	}

	// 排序
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

	// 得到最后一个关仓
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for kOperationData, vOperationData := range resOperationData {
		if kOperationData > tmpLastCloseK { // 结束查询到最后一个，默认-1不会被查到
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

// IntervalMKAndMACDData k线和间隔m时间的平均收盘价数据 .
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

	reqStart, err = time.Parse("2006-01-02 15:04:05", req.Start) // 时间进行格式校验
	if nil != err {
		return nil, err
	}
	reqEnd, err = time.Parse("2006-01-02 15:04:05", req.End) // 时间进行格式校验
	if nil != err {
		return nil, err
	}

	res := &v1.IntervalMKAndMACDDataReply{
		DataListK:     make([]*v1.IntervalMKAndMACDDataReply_ListK, 0),
		OperationData: make([]*v1.IntervalMKAndMACDDataReply_List2, 0),
		DataListMacd:  make([]*v1.IntervalMKAndMACDDataReply_ListMacd, 0),
	}

	m = int(req.M)
	k = int(req.K)
	//n = int(req.N)
	maxMxN := 200 * 60 // macd计算，至少需要数据源头数据条数，本次最大查询60分钟

	// 获取时间范围内的k线分钟数据
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo 数据时间限制，先应该随着maxMxN改变而改变
	dataLimitTime := time.Date(2020, 2, 1, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		return res, nil
	}
	// 时间查不出数据
	if reqStart.After(reqEnd) {
		return res, nil
	}
	fmt.Println(maxMxN, reqStart, reqEnd, reqStart.Add(-8*time.Hour).UnixMilli(), reqEnd.Add(-8*time.Hour).UnixMilli())
	klineMOne, err = b.klineMOneRepo.GetKLineMOneByStartTime(
		reqStart.Add(-8*time.Hour).UnixMilli(),
		reqEnd.Add(-8*time.Hour).UnixMilli(),
	)

	// 截取掉不需要的数据，之后计算
	klineM := handleMKData(klineMOne, m)
	var (
		macdData []*MACDPoint
	)
	// macd返回数据
	macdData, err = b.klineMOneRepo.NewMACDData(klineM)
	if nil != err {
		return res, nil
	}
	macdDataMap := make(map[int64]*MACDPoint, 0)
	for _, vMacdData := range macdData[maxMxN/m-(k+1):] {
		macdDataMap[vMacdData.Time] = &MACDPoint{
			Time: vMacdData.Time,
			DIF:  vMacdData.DIF,
			DEA:  vMacdData.DEA,
			MACD: vMacdData.MACD,
		}
		res.DataListMacd = append(res.DataListMacd, &v1.IntervalMKAndMACDDataReply_ListMacd{
			X1: vMacdData.MACD,
			X2: vMacdData.DIF,
			X3: vMacdData.DEA,
			X4: vMacdData.Time,
		})
	}

	klineM3 := handleMKData(klineMOne, 3)
	var (
		macdM3Data []*MACDPoint
	)
	// macd返回数据
	macdM3Data, err = b.klineMOneRepo.NewMACDData(klineM3)
	if nil != err {
		return res, nil
	}
	macdM3DataMap := make(map[int64]*MACDPoint, 0)
	for _, vMacdData := range macdM3Data[maxMxN/3-(k+1):] {
		macdM3DataMap[vMacdData.Time] = &MACDPoint{
			Time: vMacdData.Time,
			DIF:  vMacdData.DIF,
			DEA:  vMacdData.DEA,
			MACD: vMacdData.MACD,
		}

		res.DataListMacd3 = append(res.DataListMacd3, &v1.IntervalMKAndMACDDataReply_ListMacd3{
			X1: vMacdData.MACD,
			X2: vMacdData.DIF,
			X3: vMacdData.DEA,
			X4: vMacdData.Time,
		})
	}

	klineM60 := handleMKData(klineMOne, 60)
	var (
		macdM60Data []*MACDPoint
	)
	// macd返回数据
	macdM60Data, err = b.klineMOneRepo.NewMACDData(klineM60)
	if nil != err {
		return res, nil
	}
	macdM60DataMap := make(map[int64]*MACDPoint, 0)
	for _, vMacdData := range macdM60Data[maxMxN/60-(k+1):] {
		macdM60DataMap[vMacdData.Time] = &MACDPoint{
			Time: vMacdData.Time,
			DIF:  vMacdData.DIF,
			DEA:  vMacdData.DEA,
			MACD: vMacdData.MACD,
		}

		res.DataListMacd60 = append(res.DataListMacd60, &v1.IntervalMKAndMACDDataReply_ListMacd60{
			X1: vMacdData.MACD,
			X2: vMacdData.DIF,
			X3: vMacdData.DEA,
			X4: vMacdData.Time,
		})
	}
	//fmt.Println(len(klineM), klineM[199], len(tmpKlineMOne), tmpKlineMOne[199])

	// 遍历数据
	var (
		lastActionTag    string
		tmpLastActionTag string
		openActionTag    string
	)
	operationData := make(map[string]*OperationData2, 0)
	tmpKlineM := klineM[maxMxN/m-(k+1):]
	for kKlineM, vKlineM := range tmpKlineM {
		if kKlineM < k {
			continue
		}

		var tagNum int64

		// 结果
		res.DataListK = append(res.DataListK, &v1.IntervalMKAndMACDDataReply_ListK{
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

		if _, ok := macdM3DataMap[vKlineM.StartTime]; !ok {
			continue
		}

		var (
			openMoreOne  int
			openMoreTwo  int
			openEmptyOne int
			openEmptyTwo int
		)
		// 找到macd3分钟
		//for i := 1; i <= k; i++ {
		//	// 3分钟
		//	tmpStartTimeKey := vKlineM.StartTime - int64(3*60000*i)
		//	if _, ok := macdM3DataMap[tmpStartTimeKey]; !ok {
		//		break
		//	}
		//}

		// 当前分钟
		for i := 1; i <= k; i++ {
			if _, ok := macdDataMap[tmpKlineM[kKlineM-i].StartTime]; !ok {
				break
			}

			if macdDataMap[tmpKlineM[kKlineM-i].StartTime].DEA > macdDataMap[tmpKlineM[kKlineM-i].StartTime].DIF &&
				macdDataMap[tmpKlineM[kKlineM-i].StartTime].DIF > 0 {
				openMoreOne += 1
			}

			if macdDataMap[tmpKlineM[kKlineM-i].StartTime].DEA < macdDataMap[tmpKlineM[kKlineM-i].StartTime].DIF &&
				macdDataMap[tmpKlineM[kKlineM-i].StartTime].DIF < 0 {
				openEmptyOne += 1
			}
		}

		// 60分钟
		tmpNow := time.UnixMilli(vKlineM.StartTime).UTC()
		tmpNow0 := time.Date(tmpNow.Year(), tmpNow.Month(), tmpNow.Day(), tmpNow.Hour(), 0, 0, 0, time.UTC).UnixMilli()
		if vKlineM.StartTime == tmpNow0 {
			tmpNow0 -= int64(60 * 60000)
		}
		for i := 0; i < k; i++ {
			// 60分钟
			tmpStartTimeKey := tmpNow0 - int64(60*60000*i)
			if _, ok := macdM60DataMap[tmpStartTimeKey]; !ok {
				break
			}

			fmt.Println(macdM60DataMap[tmpStartTimeKey])
			if macdM60DataMap[tmpStartTimeKey].DIF > macdM60DataMap[tmpStartTimeKey].DEA &&
				macdM60DataMap[tmpStartTimeKey].DEA > 0 {
				openMoreTwo += 1
			}

			if macdM60DataMap[tmpStartTimeKey].DIF < macdM60DataMap[tmpStartTimeKey].DEA &&
				macdM60DataMap[tmpStartTimeKey].DEA < 0 {
				openEmptyTwo += 1
			}
		}

		// 平多仓
		if macdM3DataMap[vKlineM.StartTime].DIF < 0 {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "more" == tmpOpenLastOperationData2.Type {

					tmpRate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// 关
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "more",
						Status:     "close",
						Rate:       tmpRate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}
		if macdDataMap[vKlineM.StartTime].DIF < macdM3DataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "more" == tmpOpenLastOperationData2.Type {
					tmpRate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// 关
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "more",
						Status:     "close",
						Rate:       tmpRate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}

		// 平空仓
		if macdM3DataMap[vKlineM.StartTime].DIF > 0 {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "empty" == tmpOpenLastOperationData2.Type {
					tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// 关
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "empty",
						Status:     "close",
						Rate:       tmpRate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}
		if macdDataMap[vKlineM.StartTime].DIF > macdM3DataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if ("open" == tmpOpenLastOperationData2.Status || "half" == tmpOpenLastOperationData2.Status) && "empty" == tmpOpenLastOperationData2.Type {
					tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// 关
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "empty",
						Status:     "close",
						Rate:       tmpRate,
					}

					tagNum++
					lastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[lastActionTag] = tmpCloseLastOperationData
					openActionTag = ""
				}
			}
		}

		// 平多半仓
		if macdM3DataMap[vKlineM.StartTime].DIF < macdM3DataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "open" == tmpOpenLastOperationData2.Status && "more" == tmpOpenLastOperationData2.Type {
					tmpRate := (vKlineM.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// 关
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "more",
						Status:     "half",
						Rate:       tmpRate,
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// 平空半仓
		if macdM3DataMap[vKlineM.StartTime].DIF > macdM3DataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "open" == tmpOpenLastOperationData2.Status && "empty" == tmpOpenLastOperationData2.Type {
					tmpRate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
					// 关
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "empty",
						Status:     "half",
						Rate:       tmpRate,
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// 加多半仓
		if macdM3DataMap[vKlineM.StartTime].DIF > macdM3DataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "half" == tmpOpenLastOperationData2.Status && "more" == tmpOpenLastOperationData2.Type {
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "more",
						Action:     "add",
						Status:     "open",
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// 开多
		fmt.Println(openMoreOne, openMoreTwo, k, macdDataMap[vKlineM.StartTime].DIF, macdDataMap[vKlineM.StartTime].DEA)
		if openMoreOne >= k && openMoreTwo >= k && macdDataMap[vKlineM.StartTime].DIF > macdDataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "empty" == tmpOpenLastOperationData2.Type && "open" == tmpOpenLastOperationData2.Status {
					rate := (tmpOpenLastOperationData2.EndPrice-vKlineM.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
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
						Status:     "open", // 全开状态
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				} else if "empty" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "more",
						Status:     "open", // 全开状态
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
					openActionTag = ""
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:  vKlineM.StartTime,
					EndTime:    vKlineM.EndTime,
					StartPrice: vKlineM.StartPrice,
					EndPrice:   vKlineM.EndPrice,
					Amount:     2,
					Type:       "more",
					Status:     "open", // 全开状态
				}
				tagNum++
				tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
				openActionTag = tmpLastActionTag
				operationData[tmpLastActionTag] = currentOperationData
			}
		}

		// 加空半仓
		if macdM3DataMap[vKlineM.StartTime].DIF < macdM3DataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
				if "half" == tmpOpenLastOperationData2.Status && "empty" == tmpOpenLastOperationData2.Type {
					tmpCloseLastOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     tmpOpenLastOperationData2.Amount - int64(1),
						Type:       "empty",
						Action:     "add",
						Status:     "open",
					}

					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					operationData[tmpLastActionTag] = tmpCloseLastOperationData
				}
			}
		}

		// 开空
		if openEmptyOne >= k && openEmptyTwo >= k && macdDataMap[vKlineM.StartTime].DIF < macdDataMap[vKlineM.StartTime].DEA {
			if tmpOpenLastOperationData2, ok := operationData[openActionTag]; ok && nil != tmpOpenLastOperationData2 {
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
					openActionTag = ""

					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "empty",
						Status:     "open", // 全开状态
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				} else if "more" == tmpOpenLastOperationData2.Type && "close" == tmpOpenLastOperationData2.Status {
					currentOperationData := &OperationData2{
						StartTime:  vKlineM.StartTime,
						EndTime:    vKlineM.EndTime,
						StartPrice: vKlineM.StartPrice,
						EndPrice:   vKlineM.EndPrice,
						Amount:     2,
						Type:       "empty",
						Status:     "open", // 全开状态
					}
					tagNum++
					tmpLastActionTag = strconv.FormatInt(tagNum, 10) + strconv.FormatInt(vKlineM.EndTime, 10)
					openActionTag = tmpLastActionTag
					operationData[tmpLastActionTag] = currentOperationData
				}

			} else {
				currentOperationData := &OperationData2{
					StartTime:  vKlineM.StartTime,
					EndTime:    vKlineM.EndTime,
					StartPrice: vKlineM.StartPrice,
					EndPrice:   vKlineM.EndPrice,
					Amount:     2,
					Type:       "empty",
					Status:     "open", // 全开状态
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

	// 排序
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

	// 得到最后一个关仓
	for i := len(resOperationData) - 1; i >= 0; i-- {
		if "close" == resOperationData[i].Status {
			tmpLastCloseK = i
			break
		}
	}

	for kOperationData, vOperationData := range resOperationData {
		if kOperationData > tmpLastCloseK { // 结束查询到最后一个，默认-1不会被查到
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

		res.OperationData = append(res.OperationData, &v1.IntervalMKAndMACDDataReply_List2{
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

	start, err = time.Parse("2006-01-02 15:04:05", req.Start)       // 时间进行格式校验
	end = time.Now().UTC().Add(8 * time.Hour).Add(-1 * time.Minute) // 上一分钟
	end = time.Date(end.Year(), end.Month(), end.Day(), end.Hour(), end.Minute(), 59, 0, time.UTC)
	if nil != err {
		return nil, err
	}
	fmt.Println(start, end)

	// 获取数据库最后一条数据的时间
	lastKlineMOne, err = b.klineMOneRepo.GetKLineMOneOrderByEndTimeLast()
	if nil != lastKlineMOne {
		lastKlineMOneEndTime = time.UnixMilli(lastKlineMOne.EndTime).UTC().Add(8 * time.Hour).Add(1 * time.Millisecond)
		if start.Before(lastKlineMOneEndTime) { // 置换时间，数据库中已有数据
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

		if tmpEnd.Before(tmpStart) { // 健壮性
			break
		}

		tmpKlineMOne, err = b.klineMOneRepo.RequestBinanceMinuteKLinesData("BTCUSDT",
			strconv.FormatInt(tmpStart.Add(-8*time.Hour).UnixMilli(), 10),
			strconv.FormatInt(tmpEnd.Add(-8*time.Hour).UnixMilli(), 10),
			strconv.FormatInt(m, 10)+"m",
			strconv.FormatInt(limit, 10))
		if nil != err {
			return nil, err
		}

		if err = b.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
			_, err = b.klineMOneRepo.InsertKLineMOne(ctx, tmpKlineMOne)
			if nil != err {
				return err
			}

			return nil
		}); err != nil {
			fmt.Println(err)
			break
		}

		tmpStart = tmpStart.Add(time.Duration(m*limit) * time.Minute)
		fmt.Println(tmpStart, tmpEnd)
		if end.Before(tmpStart) {
			break
		}
	}

	return &v1.PullBinanceDataReply{}, nil
}

func handleManMnWithKLineMineData(n int, interval int, current int, kKlineMOne int, vKlineMOne *KLineMOne, klineMOne []*KLineMOne) *Ma {
	var (
		need                int
		tmpMaNEndPriceTotal float64 // 最新这条永远是最后一条
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
		// 整除
		if 0 == (need-i)%interval {
			tmpMaNEndPriceTotal += klineMOne[kKlineMOne-i].EndPrice // 累加
			//if need <= 25 {
			//	fmt.Println(need, time.UnixMilli(klineMOne[kKlineMOne-i].EndTime))
			//}

		} else if 0 == i {
			tmpMaNEndPriceTotal += klineMOne[kKlineMOne-i].EndPrice // 累加
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
	// 截取掉不需要的数据，之后计算
	subK := 0
	for kKlineMOne, vKlineMOne := range klineMOne {
		// 当前时间
		tmpNow := time.UnixMilli(vKlineMOne.StartTime).UTC().Add(8 * time.Hour)
		if tmpNow.Minute()%m == 0 {
			subK = kKlineMOne
			break
		}
	}
	tmpKlineMOne := klineMOne[subK:]
	lenTmpKlineMOne := len(tmpKlineMOne)
	for tmpKKlineMOne, tmpVKlineMOne := range tmpKlineMOne {
		// 当前时间
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
