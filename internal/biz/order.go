package biz

import (
	v1 "binancedata/api/binancedata/v1"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"time"
)

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

type OrderPolicyMacdCompare struct {
	ID        int64
	Type      string
	MacdType  string
	KTopPrice float64
	KLowPrice float64
	Value     float64
}

type OrderPolicyMacdLock struct {
	ID        int64
	Type      string
	CreatedAt time.Time
}

type OrderPolicyMacdCompareInfo struct {
	ID                  int64
	OrderId             int64
	Type                string
	Status              string
	OpenEndPrice        float64
	Num                 float64
	ClosePriceWin       float64
	ClosePriceLost      float64
	MacdNow             float64
	KPriceNow           float64
	MacdCompare         float64
	KPriceCompare       float64
	MacdYesterday       float64
	MacdBeforeYesterday float64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type OrderData struct {
	Symbol          string
	Side            string
	OrderType       string
	PositionSide    string
	Quantity        string
	QuantityFloat64 float64
}

type Price struct {
	Symbol string
	Price  string
}

type OrderPolicyPointCompareRepo interface {
	RequestBinanceGetOrder(symbol string) (*Order, error)
	RequestBinancePrice(symbol string) (*Price, error)
	RequestBinanceOrder(symbol string, side string, orderType string, positionSide string, quantity string, apiKey string, secretKey string) (*Order, error)

	GetLastOrderPolicyPointCompareByInfoIdAndType(infoId int64, policyPointType string, user int64) (*OrderPolicyPointCompare, error)
	GetLastOrderPolicyPointCompareInfo(user int64) (*OrderPolicyPointCompareInfo, error)
	InsertOrderPolicyPointCompareInfo(ctx context.Context, orderPolicyPointCompareInfoData *OrderPolicyPointCompareInfo, user int64) (*OrderPolicyPointCompareInfo, error)
	InsertOrderPolicyPointCompare(ctx context.Context, orderPolicyPointCompareData *OrderPolicyPointCompare, user int64) (bool, error)

	GetLastOrderPolicyMacdCompareByCreatedAtAndType(policyType string, user int64) (*OrderPolicyMacdCompare, error)
	GetLastOrderPolicyMacdLockByCreatedAt(user int64) (*OrderPolicyMacdLock, error)
	GetOrdersPolicyMacdCompareOpen(user int64) ([]*OrderPolicyMacdCompareInfo, error)
	CloseOrderPolicyMacdCompareInfo(ctx context.Context, id int64, user int64) (bool, error)
	InsertOrderPolicyMacdCompareInfo(ctx context.Context, orderPolicyMacdCompareInfoData []*OrderPolicyMacdCompareInfo, user int64) ([]*OrderPolicyMacdCompareInfo, error)
	InsertOrderPolicyMacdCompare(ctx context.Context, orderPolicyMacdCompareDataSlice []*OrderPolicyMacdCompare, user int64) (bool, error)
	InsertOrderPolicyMacdLock(ctx context.Context, orderPolicyMacdLockInfo *OrderPolicyMacdLock, user int64) (bool, error)
}

// OrderUsecase is an Order usecase.
type OrderUsecase struct {
	klineMOneRepo               KLineMOneRepo
	orderPolicyPointCompareRepo OrderPolicyPointCompareRepo
	tx                          Transaction
	log                         *log.Helper
}

// NewOrderUsecase new a Order usecase.
func NewOrderUsecase(kLineMOneRepo KLineMOneRepo, orderPolicyPointCompareRepo OrderPolicyPointCompareRepo, tx Transaction, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{klineMOneRepo: kLineMOneRepo, orderPolicyPointCompareRepo: orderPolicyPointCompareRepo, tx: tx, log: log.NewHelper(logger)}
}

func (o *OrderUsecase) TestOrder(req *v1.OrderAreaPointRequest) (*v1.OrderAreaPointReply, error) {
	var (
		err error
	)

	var (
		apiKey    string
		secretKey string
	)
	if 1 == req.User {
		apiKey = "MvzfRAnEeU46efaLYeaRms0r92d2g20iXVDQoJ8Ma5UvqH1zkJDMGB1WbSZ30P0W"
		secretKey = "bjGtZYExnHEcNBivXmJ8dLzGfMzr8SkW4ATmxLC1ZCrszbb5YJDulaiJLAgZ7L7h"
	} else {
		apiKey = "2eNaMVDIN4kdBVmSdZDkXyeucfwLBteLRwFSmUNHVuGhFs18AeVGDRZvfpTGDToX"
		secretKey = "w2xOINea6jMBJOqq9kWAvB0TWsKRWJdrM70wPbYeCMn2C1W89GxyBigbg1JSVojw"
	}

	_, err = o.orderPolicyPointCompareRepo.RequestBinanceOrder("ETHUSDT", "BUY", "MARKET", "LONG", "0.057", apiKey, secretKey)
	if nil != err {
		o.log.Error(err)
		return nil, err
	}

	return nil, nil
}

func (o *OrderUsecase) OrderAreaPoint(ctx context.Context, req *v1.OrderAreaPointRequest, test string, endTime time.Time) (*v1.OrderAreaPointReply, error) {
	var (
		start         time.Time
		end           time.Time
		pointFirst    = 0.3
		pointInterval = 0.1
		kLineMOne     []*KLineMOne
		user          int64
		err           error
	)

	// 使用用户， todo 目前非常简单的给不同的表数据
	if 0 < req.User {
		user = req.User
	}

	// 最近的第n个15分钟
	endNow := time.Now().UTC().Add(8 * time.Hour)
	if "test" == test {
		endNow = endTime
	}

	endM := endNow.Minute() / 15 * 15
	end = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), endNow.Hour(), endM, 59, 0, time.UTC).Add(-1 * time.Minute)
	start = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), endNow.Hour(), endM, 0, 0, time.UTC).Add(-3030 * time.Minute)

	// 获取数据库最后一条数据的时间
	kLineMOne, err = o.klineMOneRepo.RequestBinanceMinuteKLinesData("ETHUSDT",
		strconv.FormatInt(start.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(end.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(15, 10)+"m",
		strconv.FormatInt(1500, 10))
	if nil != err {
		return nil, err
	}

	// 获取上一单信息
	var (
		maPoint []*Ma
	)
	for kKLineMOne := range kLineMOne {
		if kKLineMOne < 199 {
			continue
		}

		// 计算ma
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

	// 数据库中上一组点位数据
	var (
		secondKeep    *OrderPolicyPointCompare
		thirdKeep     *OrderPolicyPointCompare
		secondConfirm *OrderPolicyPointCompare
		thirdConfirm  *OrderPolicyPointCompare
		second        *OrderPolicyPointCompare
		lastOrderInfo *OrderPolicyPointCompareInfo
	)
	// 最新交易记录
	lastOrderInfo, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareInfo(user)
	if nil != err {
		return nil, err
	}

	if nil != lastOrderInfo {
		second, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "second", user)
		if nil != err {
			return nil, err
		}
		secondKeep, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "second_keep", user)
		if nil != err {
			return nil, err
		}
		thirdKeep, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "third_keep", user)
		if nil != err {
			return nil, err
		}
		secondConfirm, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "second_confirm", user)
		if nil != err {
			return nil, err
		}
		thirdConfirm, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyPointCompareByInfoIdAndType(lastOrderInfo.ID, "third_confirm", user)
		if nil != err {
			return nil, err
		}
	}

	// 业务逻辑下单
	tmpPointFirstSub := maPoint[2].AvgEndPrice - maPoint[1].AvgEndPrice
	tmpPointSecondSub := maPoint[1].AvgEndPrice - maPoint[0].AvgEndPrice

	// 开单数据
	var (
		order        []*OrderData
		pointCompare []*OrderPolicyPointCompare
	)

	if tmpPointSecondSub > pointFirst && pointFirst > tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "buy_long" == lastOrderInfo.Type { // 开多
				// 有没有开仓的持仓确认点位
				tmpOpen := false
				if nil != secondKeep {
					if 2 <= secondKeep.Value {
						tmpOpen = true
					} else {
						// 新增一条清0的数据
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
					// 关多
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})
					// 开空
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
				// 此时正拿着空单，如果曾经到达过二档，清空
				if nil != second {
					// 新增一条清0的数据
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: second.InfoId,
						Type:   second.Type,
						Value:  0,
					})
				}
			}

		} else {
			// 开空
			// 获取最新价格
			var (
				price    *Price
				EthPrice float64
			)
			price, err = o.orderPolicyPointCompareRepo.RequestBinancePrice("ETHUSDT")
			if nil != err {
				return nil, err
			}
			EthPrice, err = strconv.ParseFloat(price.Price, 64)
			if nil != err {
				return nil, err
			}

			var tmpNum float64
			tmpNum, err = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(100)/EthPrice), 64)
			if nil != err {
				return nil, err
			}

			tmpAmount := 0.01
			tmpAmountStr := "0.01"
			if 1 == user {
				tmpAmount = tmpNum
				tmpAmountStr = fmt.Sprintf("%.10f", float64(100)/EthPrice)
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
			if "sell_short" == lastOrderInfo.Type { // 开空
				// 有没有开仓的持仓确认点位
				tmpOpen := false
				if nil != secondKeep {
					if 2 <= secondKeep.Value {
						tmpOpen = true
					} else {
						// 新增一条清0的数据
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
					// 关空
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})
					// 开多
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
				// 此时正拿着多单，如果曾经到达过二档，清空
				if nil != second {
					// 新增一条清0的数据
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: second.InfoId,
						Type:   second.Type,
						Value:  0,
					})
				}
			}

		} else {
			// 开多
			// 获取最新价格
			var (
				price    *Price
				EthPrice float64
			)
			price, err = o.orderPolicyPointCompareRepo.RequestBinancePrice("ETHUSDT")
			if nil != err {
				return nil, err
			}
			EthPrice, err = strconv.ParseFloat(price.Price, 64)
			if nil != err {
				return nil, err
			}

			var tmpNum float64
			tmpNum, err = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(100)/EthPrice), 64)
			if nil != err {
				return nil, err
			}

			tmpAmount := 0.01
			tmpAmountStr := "0.01"
			if 1 == user {
				tmpAmount = tmpNum
				tmpAmountStr = fmt.Sprintf("%.10f", float64(100)/EthPrice)
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

	// 到达二档
	if (pointFirst + pointInterval) < tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "sell_short" == lastOrderInfo.Type {
				// 拿着空单，第一次记录
				// 第一次到达
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
					// 第二次到达，换单
					// 关空
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// 开多
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "BUY",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// 新增lastOrder，拿确认点
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: 0, // 新增 给新增的数据 用0占位
						Type:   "second_confirm",
						Value:  1,
					})
				}
			} else if "buy_long" == lastOrderInfo.Type {
				// 拿空时，到达下方区域确认点位
				// 第一次到达
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
				// 拿着空单，第一次记录
				// 第一次到达
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
					// 第二次到达，换单
					// 关多
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "LONG",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// 开空
					order = append(order, &OrderData{
						Symbol:          "ETHUSDT",
						Side:            "SELL",
						OrderType:       "MARKET",
						PositionSide:    "SHORT",
						Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
						QuantityFloat64: lastOrderInfo.Num,
					})

					// 新增lastOrder，拿确认点
					pointCompare = append(pointCompare, &OrderPolicyPointCompare{
						InfoId: 0, // 新增 给新增的数据 用0占位
						Type:   "second_confirm",
						Value:  1,
					})
				}
			} else if "sell_short" == lastOrderInfo.Type {
				// 拿空时，到达下方区域确认点位
				// 第一次到达
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

	// 到达三档，直接换单
	if (pointFirst + 2*pointInterval) < tmpPointFirstSub {
		if nil != lastOrderInfo {
			if "sell_short" == lastOrderInfo.Type {
				// 关空
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "BUY",
					OrderType:       "MARKET",
					PositionSide:    "SHORT",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// 开多
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "BUY",
					OrderType:       "MARKET",
					PositionSide:    "LONG",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// 拿确认点
				// 新增lastOrder，拿确认点
				pointCompare = append(pointCompare, &OrderPolicyPointCompare{
					InfoId: 0, // 新增 给新增的数据 用0占位
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
				// 关多
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "SELL",
					OrderType:       "MARKET",
					PositionSide:    "LONG",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// 开空
				order = append(order, &OrderData{
					Symbol:          "ETHUSDT",
					Side:            "SELL",
					OrderType:       "MARKET",
					PositionSide:    "SHORT",
					Quantity:        strconv.FormatFloat(lastOrderInfo.Num, 'f', -1, 64),
					QuantityFloat64: lastOrderInfo.Num,
				})

				// 拿确认点
				// 新增lastOrder，拿确认点
				pointCompare = append(pointCompare, &OrderPolicyPointCompare{
					InfoId: 0, // 新增 给新增的数据 用0占位
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

	// 开单
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

		var (
			apiKey    string
			secretKey string
		)
		if 1 == user {
			apiKey = "MvzfRAnEeU46efaLYeaRms0r92d2g20iXVDQoJ8Ma5UvqH1zkJDMGB1WbSZ30P0W"
			secretKey = "bjGtZYExnHEcNBivXmJ8dLzGfMzr8SkW4ATmxLC1ZCrszbb5YJDulaiJLAgZ7L7h"
		} else {
			apiKey = "2eNaMVDIN4kdBVmSdZDkXyeucfwLBteLRwFSmUNHVuGhFs18AeVGDRZvfpTGDToX"
			secretKey = "w2xOINea6jMBJOqq9kWAvB0TWsKRWJdrM70wPbYeCMn2C1W89GxyBigbg1JSVojw"
		}

		orderBinance, err = o.orderPolicyPointCompareRepo.RequestBinanceOrder(v.Symbol, v.Side, v.OrderType, v.PositionSide, v.Quantity, apiKey, secretKey)
		if nil != err {
			o.log.Error(err)
			return nil, err
		}

		orderData = &OrderPolicyPointCompareInfo{
			OrderId: orderBinance.OrderId,
			Type:    "",
			Num:     v.QuantityFloat64,
		}

		// 此程序是初始开，后续必然是关开
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

	if err = o.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
		// 新增订单信息
		if nil != orderData {
			newOrderData, err = o.orderPolicyPointCompareRepo.InsertOrderPolicyPointCompareInfo(ctx, orderData, user)
			if nil != err {
				return err
			}
		}

		if 0 < len(pointCompare) {
			for _, vPointCompare := range pointCompare {
				tmpInfoId := vPointCompare.InfoId
				if "second_confirm" == vPointCompare.Type || "third_confirm" == vPointCompare.Type {
					if nil == newOrderData || 0 >= newOrderData.ID {
						// 这里必然开单了
						return errors.New(500, "数据混乱", "数据混乱新增数据失败")
					}
					tmpInfoId = newOrderData.ID
				}

				tmpOrderPolicyPointCompare := &OrderPolicyPointCompare{
					InfoId: tmpInfoId,
					Type:   vPointCompare.Type,
					Value:  vPointCompare.Value,
				}

				//fmt.Println(tmpOrderPolicyPointCompare, 44)
				_, err = o.orderPolicyPointCompareRepo.InsertOrderPolicyPointCompare(ctx, tmpOrderPolicyPointCompare, user)
				if nil != err {
					return err
				}
			}
		}

		// 新增点位数据
		return nil
	}); err != nil {
		o.log.Error(err)
	}

	return &v1.OrderAreaPointReply{}, nil
}

func (o *OrderUsecase) OrderMacdAndKPrice(ctx context.Context, req *v1.OrderMacdAndKPriceRequest, test string, endTime time.Time) (*v1.OrderMacdAndKPriceReply, error) {
	var (
		start      time.Time
		end        time.Time
		startDaily time.Time
		endDaily   time.Time
		kLineMOne  []*KLineMOne
		kLineD     []*KLineMOne
		user       int64
		err        error
	)

	// 使用用户， todo 目前非常简单的给不同的表数据
	if 0 < req.User {
		user = req.User
	}

	// 最近的第n个15分钟
	endNow := time.Now().UTC().Add(8 * time.Hour)
	if "test" == test {
		endNow = endTime
	}

	// 获取数据库最后一条数据的时间
	endM := endNow.Minute() / 15 * 15
	end = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), endNow.Hour(), endM, 59, 0, time.UTC).Add(-1 * time.Minute)
	start = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), endNow.Hour(), endM, 0, 0, time.UTC).Add(-3060 * time.Minute)
	kLineMOne, err = o.klineMOneRepo.RequestBinanceMinuteKLinesData("ETHUSDT",
		strconv.FormatInt(start.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(end.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(15, 10)+"m",
		strconv.FormatInt(1500, 10))
	if nil != err {
		return nil, err
	}

	// 日线200条macd参与计算，昨日的
	endDaily = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), 23, 59, 59, 0, time.UTC).Add(-1 * 24 * time.Hour)
	startDaily = time.Date(endNow.Year(), endNow.Month(), endNow.Day(), 0, 0, 0, 0, time.UTC).Add(-201 * 24 * time.Hour)
	kLineD, err = o.klineMOneRepo.RequestBinanceMinuteKLinesData("ETHUSDT",
		strconv.FormatInt(startDaily.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(endDaily.Add(-8*time.Hour).UnixMilli(), 10),
		strconv.FormatInt(1, 10)+"d",
		strconv.FormatInt(1500, 10))
	if nil != err {
		return nil, err
	}

	// 获取最大macd

	// 获取最小macd

	// 计算macd
	var (
		macdData                []*MACDPoint
		lastMacdData            []*MACDPoint
		last2MacdData           []*MACDPoint
		yesterdayMacdData       []*MACDPoint
		beforeYesterdayMacdData []*MACDPoint

		maxMacd *OrderPolicyMacdCompare
		minMacd *OrderPolicyMacdCompare

		lock *OrderPolicyMacdLock

		ordersOpen []*OrderPolicyMacdCompareInfo

		// 新增
		newMaxOrMinMacd []*OrderPolicyMacdCompare
		closeOrder      []*OrderPolicyMacdCompareInfo
		newOrder        []*OrderPolicyMacdCompareInfo
		newOrderLock    *OrderPolicyMacdLock

		price    *Price
		EthPrice float64
	)

	// macd数据，取200条k线数据计算，当前
	lastKeyMLive := len(kLineMOne) - 1
	macdData, err = o.klineMOneRepo.NewMACDData(kLineMOne[lastKeyMLive-199:])
	if nil != err {
		return nil, err
	}
	// macd数据，取200条k线数据计算，上一根
	lastLastKeyMLive := len(kLineMOne) - 2
	lastMacdData, err = o.klineMOneRepo.NewMACDData(kLineMOne[lastLastKeyMLive-199:])
	if nil != err {
		return nil, err
	}
	// macd数据，取200条k线数据计算，倒数第二根
	lastLast2KeyMLive := len(kLineMOne) - 3
	last2MacdData, err = o.klineMOneRepo.NewMACDData(kLineMOne[lastLast2KeyMLive-199:])
	if nil != err {
		return nil, err
	}
	// macd数据，取200条k线数据计算，昨天
	lastKeyD := len(kLineD) - 1
	yesterdayMacdData, err = o.klineMOneRepo.NewMACDData(kLineD[lastKeyD-199:])
	if nil != err {
		return nil, err
	}
	// macd数据，取200条k线数据计算，前天
	lastLastKeyD := len(kLineD) - 2
	beforeYesterdayMacdData, err = o.klineMOneRepo.NewMACDData(kLineD[lastLastKeyD-199:])
	if nil != err {
		return nil, err
	}
	//fmt.Println(lastMacdData[199], macdData[199], yesterdayMacdData[199], beforeYesterdayMacdData[199], user)

	// 获取数据库保存的最大值，最小值
	maxMacd, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyMacdCompareByCreatedAtAndType("max", user)
	if nil != err {
		return nil, err
	}
	minMacd, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyMacdCompareByCreatedAtAndType("min", user)
	if nil != err {
		return nil, err
	}

	// 获取锁定订单日期
	lock, err = o.orderPolicyPointCompareRepo.GetLastOrderPolicyMacdLockByCreatedAt(user)
	if nil != err {
		return nil, err
	}

	// 获取所有订单
	ordersOpen, err = o.orderPolicyPointCompareRepo.GetOrdersPolicyMacdCompareOpen(user)
	if nil != err {
		return nil, err
	}

	// 获取最新价格
	price, err = o.orderPolicyPointCompareRepo.RequestBinancePrice("ETHUSDT")
	if nil != err {
		return nil, err
	}
	EthPrice, err = strconv.ParseFloat(price.Price, 64)
	if nil != err {
		return nil, err
	}

	var (
		max1 = 3.0
		low1 = -3.0
		max2 = 1.0
		low2 = -1.0

		closeRateWin  = 0.01
		closeRateLost = 0.03

		tmpEmptyCloseLostNum int
		tmpMoreCloseLostNum  int

		closeRateWin2  = 0.05
		closeRateLost2 = 0.01
	)

	// 平单，数据正常时，结果一定指向一个方向
	for _, vOrdersOpen := range ordersOpen {
		tmpClose := false
		tmpLock := ""

		if "empty" == vOrdersOpen.Type {
			if vOrdersOpen.ClosePriceWin >= kLineMOne[lastKeyMLive].EndPrice { // 止盈
				tmpRate := (kLineMOne[lastKeyMLive].EndPrice - vOrdersOpen.OpenEndPrice) / vOrdersOpen.OpenEndPrice
				if 0 < -tmpRate { // 盈利
					tmpLock = "more"
				} else {
					// 关同向单
					tmpLock = "empty"
					tmpEmptyCloseLostNum++
				}
				tmpClose = true
			} else if vOrdersOpen.ClosePriceLost <= kLineMOne[lastKeyMLive].EndPrice { // 止损
				tmpClose = true
				// 关同向单
				tmpLock = "empty"
				tmpEmptyCloseLostNum++
			}

		} else if "more" == vOrdersOpen.Type {
			if vOrdersOpen.ClosePriceWin <= kLineMOne[lastKeyMLive].EndPrice { // 止盈
				tmpRate := (kLineMOne[lastKeyMLive].EndPrice - vOrdersOpen.OpenEndPrice) / vOrdersOpen.OpenEndPrice
				if 0 < tmpRate { // 盈利
					tmpLock = "empty"
				} else {
					// 关同向单
					tmpLock = "more"
					tmpMoreCloseLostNum++
				}
				tmpClose = true
			} else if vOrdersOpen.ClosePriceLost >= kLineMOne[lastKeyMLive].EndPrice { // 止损
				tmpClose = true
				// 关同向单
				tmpLock = "more"
				tmpMoreCloseLostNum++
			}
		}
		if "" != tmpLock {
			fmt.Println(tmpLock)
			newOrderLock = &OrderPolicyMacdLock{
				Type: tmpLock,
			}
		}

		if tmpClose {
			fmt.Println(66666, closeOrder)
			closeOrder = append(closeOrder, &OrderPolicyMacdCompareInfo{
				ID:                  vOrdersOpen.ID,
				OrderId:             vOrdersOpen.OrderId,
				Type:                vOrdersOpen.Type,
				Status:              "close",
				OpenEndPrice:        vOrdersOpen.OpenEndPrice,
				Num:                 vOrdersOpen.Num,
				ClosePriceWin:       vOrdersOpen.ClosePriceWin,
				ClosePriceLost:      vOrdersOpen.ClosePriceLost,
				MacdNow:             vOrdersOpen.MacdNow,
				KPriceNow:           vOrdersOpen.KPriceNow,
				MacdCompare:         vOrdersOpen.MacdCompare,
				KPriceCompare:       vOrdersOpen.KPriceCompare,
				MacdYesterday:       vOrdersOpen.MacdYesterday,
				MacdBeforeYesterday: vOrdersOpen.MacdBeforeYesterday,
				CreatedAt:           vOrdersOpen.CreatedAt,
				UpdatedAt:           time.Now().UTC(),
			})
		}

	}

	// 停掉所有反向单
	if nil != newOrderLock {
		fmt.Println(5555)
		for _, vOrdersOpen := range ordersOpen {
			if newOrderLock.Type == vOrdersOpen.Type {
				closeOrder = append(closeOrder, &OrderPolicyMacdCompareInfo{
					ID:                  vOrdersOpen.ID,
					OrderId:             vOrdersOpen.OrderId,
					Type:                vOrdersOpen.Type,
					Status:              "close",
					OpenEndPrice:        vOrdersOpen.OpenEndPrice,
					Num:                 vOrdersOpen.Num,
					ClosePriceWin:       vOrdersOpen.ClosePriceWin,
					ClosePriceLost:      vOrdersOpen.ClosePriceLost,
					MacdNow:             vOrdersOpen.MacdNow,
					KPriceNow:           vOrdersOpen.KPriceNow,
					MacdCompare:         vOrdersOpen.MacdCompare,
					KPriceCompare:       vOrdersOpen.KPriceCompare,
					MacdYesterday:       vOrdersOpen.MacdYesterday,
					MacdBeforeYesterday: vOrdersOpen.MacdBeforeYesterday,
					CreatedAt:           vOrdersOpen.CreatedAt,
					UpdatedAt:           time.Now().UTC(),
				})
			}
		}
	}

	// 超过2个亏损订单，开个反向单
	if tmpEmptyCloseLostNum >= 2 {
		fmt.Println(time.UnixMilli(kLineMOne[lastKeyMLive].StartTime).UTC().Add(8 * time.Hour))
		fmt.Println("lost more open empty more", tmpEmptyCloseLostNum)
		var tmpNum float64
		tmpNum, err = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(15)/EthPrice), 64)
		if nil != err {
			return nil, err
		}

		newOrder = append(newOrder, &OrderPolicyMacdCompareInfo{
			Type:           "more",
			Status:         "open",
			Num:            tmpNum,
			OpenEndPrice:   kLineMOne[lastKeyMLive].EndPrice,
			ClosePriceWin:  kLineMOne[lastKeyMLive].EndPrice + kLineMOne[lastKeyMLive].EndPrice*closeRateWin2,  // 关仓止盈
			ClosePriceLost: kLineMOne[lastKeyMLive].EndPrice - kLineMOne[lastKeyMLive].EndPrice*closeRateLost2, // 关仓止损
		})
	}

	if tmpMoreCloseLostNum >= 2 {
		fmt.Println(time.UnixMilli(kLineMOne[lastKeyMLive].StartTime).UTC().Add(8 * time.Hour))
		fmt.Println("lost more open empty", tmpEmptyCloseLostNum)
		var tmpNum float64
		tmpNum, err = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(15)/EthPrice), 64)
		if nil != err {
			return nil, err
		}

		newOrder = append(newOrder, &OrderPolicyMacdCompareInfo{
			Type:           "empty",
			Status:         "open",
			Num:            tmpNum,
			OpenEndPrice:   kLineMOne[lastKeyMLive].EndPrice,
			ClosePriceWin:  kLineMOne[lastKeyMLive].EndPrice - kLineMOne[lastKeyMLive].EndPrice*closeRateWin2,  // 关仓止盈
			ClosePriceLost: kLineMOne[lastKeyMLive].EndPrice + kLineMOne[lastKeyMLive].EndPrice*closeRateLost2, // 关仓止损
		})
	}

	// 计算 初始化
	if macdData[199].MACD > 0 && max1 < macdData[199].MACD { // 正的&大于设定值
		if nil == maxMacd || (maxMacd.Value < macdData[199].MACD && lastMacdData[199].MACD < macdData[199].MACD) { // 空的或者，值>原值，实心柱，替换
			maxMacd = &OrderPolicyMacdCompare{
				Type:      "max",
				MacdType:  "15m",
				KTopPrice: kLineMOne[lastKeyMLive].TopPrice,
				KLowPrice: kLineMOne[lastKeyMLive].LowPrice,
				Value:     macdData[199].MACD,
			}
			newMaxOrMinMacd = append(newMaxOrMinMacd, maxMacd)
		}
	}
	if macdData[199].MACD < 0 && low1 > macdData[199].MACD {
		if nil == minMacd || (minMacd.Value > macdData[199].MACD && lastMacdData[199].MACD > macdData[199].MACD) { // 空的或者，值>原值，实心柱，替换
			minMacd = &OrderPolicyMacdCompare{
				Type:      "min",
				MacdType:  "15m",
				KTopPrice: kLineMOne[lastKeyMLive].TopPrice,
				KLowPrice: kLineMOne[lastKeyMLive].LowPrice,
				Value:     macdData[199].MACD,
			}

			newMaxOrMinMacd = append(newMaxOrMinMacd, minMacd)
		}
	}

	// 肯定排除了第一点位开单
	if macdData[199].MACD > 0 && max2 < macdData[199].MACD && // 正的&大于设定值
		lastMacdData[199].MACD > macdData[199].MACD && // 上一根大于当前，代表趋势向下第一根空心柱
		last2MacdData[199].MACD < lastMacdData[199].MACD &&
		nil != maxMacd &&
		macdData[199].MACD < maxMacd.Value && kLineMOne[lastKeyMLive].TopPrice > maxMacd.KTopPrice && // 当前小于最大实心柱 且最高价格大于实心柱最高价
		!(yesterdayMacdData[199].MACD > 0 && beforeYesterdayMacdData[199].MACD < yesterdayMacdData[199].MACD) { // 绿色实心不开，所以取反
		fmt.Println(time.UnixMilli(kLineMOne[lastKeyMLive].StartTime).UTC().Add(8 * time.Hour))
		fmt.Println(3333, maxMacd, macdData[199].MACD, lastMacdData[199].MACD, kLineMOne[lastKeyMLive].TopPrice, yesterdayMacdData[199].MACD, beforeYesterdayMacdData[199].MACD)

		if nil == lock || "empty" != lock.Type { // 无锁定
			// 开空
			var tmpNum float64
			tmpNum, err = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(15)/EthPrice), 64)
			if nil != err {
				return nil, err
			}

			newOrder = append(newOrder, &OrderPolicyMacdCompareInfo{
				Type:                "empty",
				Status:              "open",
				Num:                 tmpNum,
				OpenEndPrice:        kLineMOne[lastKeyMLive].EndPrice,
				ClosePriceWin:       kLineMOne[lastKeyMLive].EndPrice - kLineMOne[lastKeyMLive].EndPrice*closeRateWin,  // 关仓止盈
				ClosePriceLost:      kLineMOne[lastKeyMLive].EndPrice + kLineMOne[lastKeyMLive].EndPrice*closeRateLost, // 关仓止损
				MacdNow:             macdData[199].MACD,
				KPriceNow:           kLineMOne[lastKeyMLive].TopPrice,
				MacdCompare:         maxMacd.Value,
				KPriceCompare:       maxMacd.KTopPrice,
				MacdBeforeYesterday: yesterdayMacdData[199].MACD,
				MacdYesterday:       beforeYesterdayMacdData[199].MACD,
			})
		}
	}

	// 肯定排除了第一点位开单
	if macdData[199].MACD < 0 && low2 > macdData[199].MACD && // 负的&小于设定值
		lastMacdData[199].MACD < macdData[199].MACD && // 上一根小于当前，代表趋势向下第一根空心柱
		last2MacdData[199].MACD > lastMacdData[199].MACD &&
		nil != minMacd &&
		macdData[199].MACD > minMacd.Value && kLineMOne[lastKeyMLive].LowPrice < minMacd.KLowPrice && // 当前大于最大实心柱 且最低价格小于实心柱低高价
		!(yesterdayMacdData[199].MACD < 0 && beforeYesterdayMacdData[199].MACD > yesterdayMacdData[199].MACD) { // 绿色实心不开，所以取反
		fmt.Println(time.UnixMilli(kLineMOne[lastKeyMLive].StartTime).UTC().Add(8 * time.Hour))
		fmt.Println(4444, minMacd, macdData[199].MACD, lastMacdData[199].MACD, kLineMOne[lastKeyMLive].TopPrice, yesterdayMacdData[199].MACD, beforeYesterdayMacdData[199].MACD)

		// 开多
		if nil == lock || "more" != lock.Type { // 无锁定
			var tmpNum float64
			tmpNum, err = strconv.ParseFloat(fmt.Sprintf("%.3f", float64(15)/EthPrice), 64)
			if nil != err {
				return nil, err
			}

			newOrder = append(newOrder, &OrderPolicyMacdCompareInfo{
				Type:                "more",
				Status:              "open",
				Num:                 tmpNum,
				OpenEndPrice:        kLineMOne[lastKeyMLive].EndPrice,
				ClosePriceWin:       kLineMOne[lastKeyMLive].EndPrice + kLineMOne[lastKeyMLive].EndPrice*closeRateWin,  // 关仓止盈
				ClosePriceLost:      kLineMOne[lastKeyMLive].EndPrice - kLineMOne[lastKeyMLive].EndPrice*closeRateLost, // 关仓止损
				MacdNow:             macdData[199].MACD,
				KPriceNow:           kLineMOne[lastKeyMLive].LowPrice,
				MacdCompare:         minMacd.Value,
				KPriceCompare:       minMacd.KLowPrice,
				MacdBeforeYesterday: yesterdayMacdData[199].MACD,
				MacdYesterday:       beforeYesterdayMacdData[199].MACD,
			})
		}
	}

	// 新增
	for kCloseOrder, vCloseOrder := range closeOrder {
		var orderBinance *Order

		tmpSide := ""
		tmpPositionSide := ""
		if "empty" == vCloseOrder.Type && "close" == vCloseOrder.Status {
			tmpSide = "BUY"
			tmpPositionSide = "SHORT"
		} else if "more" == vCloseOrder.Type && "close" == vCloseOrder.Status {
			tmpSide = "SELL"
			tmpPositionSide = "LONG"
		}

		var (
			apiKey    string
			secretKey string
		)
		if 1 == user {
			apiKey = "MvzfRAnEeU46efaLYeaRms0r92d2g20iXVDQoJ8Ma5UvqH1zkJDMGB1WbSZ30P0W"
			secretKey = "bjGtZYExnHEcNBivXmJ8dLzGfMzr8SkW4ATmxLC1ZCrszbb5YJDulaiJLAgZ7L7h"
		} else {
			apiKey = "pswGalfy8OvPgL4vdgjzCNbL4XFnlif3OjsA9vymiDoZD4MC2gO4QGTmGLj0mnqP"
			secretKey = "gcT4X2AcWr8dRag3t0CWg8Dfip9sjOSYmNpEx6bxnkNfTc2StICEoqtNGnkQQzwe"
		}

		orderBinance, err = o.orderPolicyPointCompareRepo.RequestBinanceOrder("ETHUSDT", tmpSide, "MARKET", tmpPositionSide, strconv.FormatFloat(vCloseOrder.Num, 'f', -1, 64), apiKey, secretKey)
		if nil != err {
			o.log.Error(err)
			return nil, err
		}
		fmt.Println(orderBinance)

		closeOrder[kCloseOrder].OrderId = orderBinance.OrderId
	}

	for kNewOrder, vNewOrder := range newOrder {
		var orderBinance *Order

		tmpSide := ""
		tmpPositionSide := ""
		if "empty" == vNewOrder.Type && "open" == vNewOrder.Status {
			tmpSide = "SELL"
			tmpPositionSide = "SHORT"
		} else if "more" == vNewOrder.Type && "open" == vNewOrder.Status {
			tmpSide = "BUY"
			tmpPositionSide = "LONG"
		}

		var (
			apiKey    string
			secretKey string
		)
		if 1 == user {
			apiKey = "MvzfRAnEeU46efaLYeaRms0r92d2g20iXVDQoJ8Ma5UvqH1zkJDMGB1WbSZ30P0W"
			secretKey = "bjGtZYExnHEcNBivXmJ8dLzGfMzr8SkW4ATmxLC1ZCrszbb5YJDulaiJLAgZ7L7h"
		} else {
			apiKey = "pswGalfy8OvPgL4vdgjzCNbL4XFnlif3OjsA9vymiDoZD4MC2gO4QGTmGLj0mnqP"
			secretKey = "gcT4X2AcWr8dRag3t0CWg8Dfip9sjOSYmNpEx6bxnkNfTc2StICEoqtNGnkQQzwe"
		}

		orderBinance, err = o.orderPolicyPointCompareRepo.RequestBinanceOrder("ETHUSDT", tmpSide, "MARKET", tmpPositionSide, strconv.FormatFloat(vNewOrder.Num, 'f', -1, 64), apiKey, secretKey)
		if nil != err {
			o.log.Error(err)
			return nil, err
		}
		fmt.Println(orderBinance)

		newOrder[kNewOrder].OrderId = orderBinance.OrderId
	}

	if err = o.tx.ExecTx(ctx, func(ctx context.Context) error { // 事务
		// 修改订单信息
		if len(closeOrder) > 0 {
			for _, vCloseOrder := range closeOrder {
				_, err = o.orderPolicyPointCompareRepo.CloseOrderPolicyMacdCompareInfo(ctx, vCloseOrder.ID, user)
				if nil != err {
					return err
				}
			}
		}

		// 新增订单信息
		if len(newOrder) > 0 {
			_, err = o.orderPolicyPointCompareRepo.InsertOrderPolicyMacdCompareInfo(ctx, newOrder, user)
			if nil != err {
				return err
			}
		}

		// 新增信息
		if len(newMaxOrMinMacd) > 0 {
			_, err = o.orderPolicyPointCompareRepo.InsertOrderPolicyMacdCompare(ctx, newMaxOrMinMacd, user)
			if nil != err {
				return err
			}
		}

		// 新增信息
		if nil != newOrderLock {
			_, err = o.orderPolicyPointCompareRepo.InsertOrderPolicyMacdLock(ctx, newOrderLock, user)
			if nil != err {
				return err
			}
		}

		// 新增点位数据
		return nil
	}); err != nil {
		o.log.Error(err)
	}

	return &v1.OrderMacdAndKPriceReply{}, nil
}
