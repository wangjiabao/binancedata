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
}

type OperationData2Slice []*OperationData2

func (o OperationData2Slice) Len() int           { return len(o) }
func (o OperationData2Slice) Less(i, j int) bool { return o[i].EndTime < o[j].EndTime }
func (o OperationData2Slice) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }

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
	dataLimitTime := time.Date(2022, 01, 02, 0, 0, 0, 0, time.UTC)
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
						rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.StartPrice)/tmpOpenLastOperationData2.StartPrice - 0.0003
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
						rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.StartPrice)/tmpOpenLastOperationData2.StartPrice - 0.0003
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
						rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.EndPrice)/tmpOpenLastOperationData2.EndPrice - 0.0003
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
		reqStart         time.Time
		reqEnd           time.Time
		klineMOne        []*KLineMOne
		n1               int
		n2               int
		m1               int
		m2               int
		topX             float64
		lowX             float64
		fee              float64
		maNMFirst        []*Ma // todo 各种图优化
		maNMSecond       []*Ma
		resOperationData OperationData2Slice
		err              error
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

	// 获取时间范围内的k线分钟数据
	reqStart = reqStart.Add(-time.Duration(maxMxN-1) * time.Minute)
	// todo 数据时间限制，先应该随着maxMxN改变而改变
	dataLimitTime := time.Date(2022, 01, 02, 0, 0, 0, 0, time.UTC)
	if reqStart.Before(dataLimitTime) {
		reqStart = dataLimitTime
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

					// 更新比较的最高价
					if vKlineMOne.TopPrice > compareTopPrice {
						compareTopPrice = vKlineMOne.TopPrice
					}

					// 最高价-最高价*x% > 最新价 关仓条件
					if compareTopPrice-compareTopPrice*topX > vKlineMOne.EndPrice {
						tmpDo = true
					}

				}
				if tmpDo {
					rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.StartPrice)/tmpOpenLastOperationData2.StartPrice - fee
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

					// 更新比较的最高价
					if vKlineMOne.LowPrice < compareLowPrice {
						compareLowPrice = vKlineMOne.LowPrice
					}

					// 最高价-最高价*x% > 最新价 关仓条件
					if compareLowPrice-compareLowPrice*lowX < vKlineMOne.EndPrice {
						tmpDo = true
					}
				}
				if tmpDo {
					rate := (vKlineMOne.EndPrice-tmpOpenLastOperationData2.StartPrice)/tmpOpenLastOperationData2.StartPrice - fee
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
					tmpDo := false
					if tmpOpenLastOperationData2, ok := operationData[lastActionTag]; ok && nil != tmpOpenLastOperationData2 {
						if "close" == tmpOpenLastOperationData2.Status {
							tmpDo = true
						}
					} else {
						tmpDo = true
					}

					if tmpDo {
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
						openActionTag = lastActionTag
						//fmt.Println(openActionTag)
						operationData[lastActionTag] = currentOperationData
						compareTopPrice = vKlineMOne.TopPrice

						tmpBackGround = "green"
					}
				}
			}

			// 开空
			//fmt.Println(len(ma5M15))
			if tmpMaNMFirst.AvgEndPrice < tmpMaNMSecond.AvgEndPrice { // 条件1
				lastMaNMFirst := maNMFirst[len(maNMFirst)-2]                // 上一单
				lastMaNMSecond := maNMSecond[len(maNMSecond)-2]             // 上一单
				if lastMaNMFirst.AvgEndPrice > lastMaNMSecond.AvgEndPrice { // 条件1
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
						openActionTag = lastActionTag
						//fmt.Println(openActionTag)
						operationData[lastActionTag] = currentOperationData
						compareLowPrice = vKlineMOne.LowPrice

						tmpBackGround = "red"
					}

				}
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

		tmpBinanceData, err = requestBinanceMinuteKLinesData("BTCUSDT",
			strconv.FormatInt(start.Add(-8*time.Hour).UnixMilli(), 10),
			strconv.FormatInt(end.Add(-8*time.Hour).UnixMilli(), 10),
			strconv.FormatInt(m, 10)+"m",
			strconv.FormatInt(limit, 10))
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
	beforeBinanceData, err = requestBinanceMinuteKLinesData("BTCUSDT",
		strconv.FormatInt(beforeStart.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(beforeEnd.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(m, 10)+"m",
		strconv.FormatInt(limit, 10))
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

func requestBinanceMinuteKLinesData(symbol string, startTime string, endTime string, interval string, limit string) ([]*BinanceData, error) {
	fmt.Println(symbol, startTime, endTime, interval, limit)
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

	res := make([]*BinanceData, 0)
	for _, v := range i {
		res = append(res, &BinanceData{
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
		})
	}

	return res, err
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
