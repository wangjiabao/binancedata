package data

//type Kline struct {
//	Amount    float64   `orm:"column(amount)"`    // 成交量
//	Count     int64     `orm:"column(count)"`     // 成交笔数
//	Open      float64   `orm:"column(open)"`      // 开盘价
//	Close     float64   `orm:"column(close)"`     // 收盘价, 当K线为最晚的一根时, 时最新成交价
//	Low       float64   `orm:"column(low)"`       // 最低价
//	High      float64   `orm:"column(high)"`      // 最高价
//	Vol       float64   `orm:"column(vol)"`       // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
//	KlineTime time.Time `orm:"column(klinetime)"` // k线时间
//}

// Kline struct
type Kline struct {
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

type Point struct {
	Time  int64
	Value float64 //数值
}

// EMA struct
type EMA struct {
	Period int //默认计算几天的EMA
	points []*EMAPoint
	kline  []*KLineMOne
}

type EMAPoint struct {
	Point
}

// NewEMA new Func
func NewEMA(list []*KLineMOne, period int) *EMA {
	m := &EMA{kline: list, Period: period}
	return m
}

// Calculation Func
func (e *EMA) Calculation() *EMA {
	for _, v := range e.kline {
		e.Add(v.StartTime, v.EndPrice)
	}
	return e
}

// GetPoints return Point
func (e *EMA) GetPoints() []*EMAPoint {
	return e.points
}

// Add adds a new Value to Ema
// 使用方法，先添加最早日期的数据,最后一条应该是当前日期的数据，结果与 AICoin 对比完全一致
func (e *EMA) Add(timestamp int64, value float64) {
	p := &EMAPoint{}
	p.Time = timestamp

	//平滑指数，一般取作2/(N+1)
	alpha := 2.0 / (float64(e.Period) + 1.0)

	// fmt.Println(alpha)

	emaTminusOne := value
	if len(e.points) > 0 {
		emaTminusOne = e.points[len(e.points)-1].Value
	}

	// 计算 EMA指数
	emaT := alpha*value + (1-alpha)*emaTminusOne
	p.Value = emaT
	e.points = append(e.points, p)
}
