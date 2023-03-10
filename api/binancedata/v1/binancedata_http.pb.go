// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.5.0
// - protoc             v3.21.7
// source: api/binancedata/v1/binancedata.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationBinanceDataAreaPointIntervalMAvgEndPriceData = "/api.binancedata.v1.BinanceData/AreaPointIntervalMAvgEndPriceData"
const OperationBinanceDataIntervalMAvgEndPriceData = "/api.binancedata.v1.BinanceData/IntervalMAvgEndPriceData"
const OperationBinanceDataIntervalMAvgEndPriceMacdAndAtrData = "/api.binancedata.v1.BinanceData/IntervalMAvgEndPriceMacdAndAtrData"
const OperationBinanceDataIntervalMKAndMACDData = "/api.binancedata.v1.BinanceData/IntervalMKAndMACDData"
const OperationBinanceDataIntervalMMACDData = "/api.binancedata.v1.BinanceData/IntervalMMACDData"
const OperationBinanceDataKAnd2NIntervalMAvgEndPriceData = "/api.binancedata.v1.BinanceData/KAnd2NIntervalMAvgEndPriceData"
const OperationBinanceDataOrder = "/api.binancedata.v1.BinanceData/Order"
const OperationBinanceDataPullBinanceData = "/api.binancedata.v1.BinanceData/PullBinanceData"
const OperationBinanceDataXNIntervalMAvgEndPriceData = "/api.binancedata.v1.BinanceData/XNIntervalMAvgEndPriceData"

type BinanceDataHTTPServer interface {
	AreaPointIntervalMAvgEndPriceData(context.Context, *AreaPointIntervalMAvgEndPriceDataRequest) (*AreaPointIntervalMAvgEndPriceDataReply, error)
	IntervalMAvgEndPriceData(context.Context, *IntervalMAvgEndPriceDataRequest) (*IntervalMAvgEndPriceDataReply, error)
	IntervalMAvgEndPriceMacdAndAtrData(context.Context, *IntervalMAvgEndPriceMacdAndAtrDataRequest) (*IntervalMAvgEndPriceMacdAndAtrDataReply, error)
	IntervalMKAndMACDData(context.Context, *IntervalMKAndMACDDataRequest) (*IntervalMKAndMACDDataReply, error)
	IntervalMMACDData(context.Context, *IntervalMMACDDataRequest) (*IntervalMMACDDataReply, error)
	KAnd2NIntervalMAvgEndPriceData(context.Context, *KAnd2NIntervalMAvgEndPriceDataRequest) (*KAnd2NIntervalMAvgEndPriceDataReply, error)
	Order(context.Context, *OrderRequest) (*OrderReply, error)
	PullBinanceData(context.Context, *PullBinanceDataRequest) (*PullBinanceDataReply, error)
	XNIntervalMAvgEndPriceData(context.Context, *XNIntervalMAvgEndPriceDataRequest) (*XNIntervalMAvgEndPriceDataReply, error)
}

func RegisterBinanceDataHTTPServer(s *http.Server, srv BinanceDataHTTPServer) {
	r := s.Route("/")
	r.POST("/api/binancedata/x_n_interval_m_avg_end_price_data", _BinanceData_XNIntervalMAvgEndPriceData0_HTTP_Handler(srv))
	r.POST("/api/binancedata/k_and_2_n_interval_m_avg_end_price_data", _BinanceData_KAnd2NIntervalMAvgEndPriceData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/pull", _BinanceData_PullBinanceData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/interval_m_avg_end_price_data", _BinanceData_IntervalMAvgEndPriceData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/interval_m_macd_data", _BinanceData_IntervalMMACDData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/interval_m_k_and_macd_data", _BinanceData_IntervalMKAndMACDData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/area_point_interval_m_avg_end_price_data", _BinanceData_AreaPointIntervalMAvgEndPriceData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/interval_m_avg_end_price_macd_and_atr_data", _BinanceData_IntervalMAvgEndPriceMacdAndAtrData0_HTTP_Handler(srv))
	r.GET("/api/binancedata/order", _BinanceData_Order0_HTTP_Handler(srv))
}

func _BinanceData_XNIntervalMAvgEndPriceData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in XNIntervalMAvgEndPriceDataRequest
		if err := ctx.Bind(&in.SendBody); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataXNIntervalMAvgEndPriceData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.XNIntervalMAvgEndPriceData(ctx, req.(*XNIntervalMAvgEndPriceDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*XNIntervalMAvgEndPriceDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_KAnd2NIntervalMAvgEndPriceData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in KAnd2NIntervalMAvgEndPriceDataRequest
		if err := ctx.Bind(&in.SendBody); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataKAnd2NIntervalMAvgEndPriceData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.KAnd2NIntervalMAvgEndPriceData(ctx, req.(*KAnd2NIntervalMAvgEndPriceDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*KAnd2NIntervalMAvgEndPriceDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_PullBinanceData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PullBinanceDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataPullBinanceData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.PullBinanceData(ctx, req.(*PullBinanceDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*PullBinanceDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_IntervalMAvgEndPriceData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IntervalMAvgEndPriceDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataIntervalMAvgEndPriceData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.IntervalMAvgEndPriceData(ctx, req.(*IntervalMAvgEndPriceDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IntervalMAvgEndPriceDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_IntervalMMACDData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IntervalMMACDDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataIntervalMMACDData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.IntervalMMACDData(ctx, req.(*IntervalMMACDDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IntervalMMACDDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_IntervalMKAndMACDData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IntervalMKAndMACDDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataIntervalMKAndMACDData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.IntervalMKAndMACDData(ctx, req.(*IntervalMKAndMACDDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IntervalMKAndMACDDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_AreaPointIntervalMAvgEndPriceData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AreaPointIntervalMAvgEndPriceDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataAreaPointIntervalMAvgEndPriceData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AreaPointIntervalMAvgEndPriceData(ctx, req.(*AreaPointIntervalMAvgEndPriceDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AreaPointIntervalMAvgEndPriceDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_IntervalMAvgEndPriceMacdAndAtrData0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in IntervalMAvgEndPriceMacdAndAtrDataRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataIntervalMAvgEndPriceMacdAndAtrData)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.IntervalMAvgEndPriceMacdAndAtrData(ctx, req.(*IntervalMAvgEndPriceMacdAndAtrDataRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*IntervalMAvgEndPriceMacdAndAtrDataReply)
		return ctx.Result(200, reply)
	}
}

func _BinanceData_Order0_HTTP_Handler(srv BinanceDataHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in OrderRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationBinanceDataOrder)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Order(ctx, req.(*OrderRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*OrderReply)
		return ctx.Result(200, reply)
	}
}

type BinanceDataHTTPClient interface {
	AreaPointIntervalMAvgEndPriceData(ctx context.Context, req *AreaPointIntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (rsp *AreaPointIntervalMAvgEndPriceDataReply, err error)
	IntervalMAvgEndPriceData(ctx context.Context, req *IntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (rsp *IntervalMAvgEndPriceDataReply, err error)
	IntervalMAvgEndPriceMacdAndAtrData(ctx context.Context, req *IntervalMAvgEndPriceMacdAndAtrDataRequest, opts ...http.CallOption) (rsp *IntervalMAvgEndPriceMacdAndAtrDataReply, err error)
	IntervalMKAndMACDData(ctx context.Context, req *IntervalMKAndMACDDataRequest, opts ...http.CallOption) (rsp *IntervalMKAndMACDDataReply, err error)
	IntervalMMACDData(ctx context.Context, req *IntervalMMACDDataRequest, opts ...http.CallOption) (rsp *IntervalMMACDDataReply, err error)
	KAnd2NIntervalMAvgEndPriceData(ctx context.Context, req *KAnd2NIntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (rsp *KAnd2NIntervalMAvgEndPriceDataReply, err error)
	Order(ctx context.Context, req *OrderRequest, opts ...http.CallOption) (rsp *OrderReply, err error)
	PullBinanceData(ctx context.Context, req *PullBinanceDataRequest, opts ...http.CallOption) (rsp *PullBinanceDataReply, err error)
	XNIntervalMAvgEndPriceData(ctx context.Context, req *XNIntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (rsp *XNIntervalMAvgEndPriceDataReply, err error)
}

type BinanceDataHTTPClientImpl struct {
	cc *http.Client
}

func NewBinanceDataHTTPClient(client *http.Client) BinanceDataHTTPClient {
	return &BinanceDataHTTPClientImpl{client}
}

func (c *BinanceDataHTTPClientImpl) AreaPointIntervalMAvgEndPriceData(ctx context.Context, in *AreaPointIntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (*AreaPointIntervalMAvgEndPriceDataReply, error) {
	var out AreaPointIntervalMAvgEndPriceDataReply
	pattern := "/api/binancedata/area_point_interval_m_avg_end_price_data"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataAreaPointIntervalMAvgEndPriceData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) IntervalMAvgEndPriceData(ctx context.Context, in *IntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (*IntervalMAvgEndPriceDataReply, error) {
	var out IntervalMAvgEndPriceDataReply
	pattern := "/api/binancedata/interval_m_avg_end_price_data"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataIntervalMAvgEndPriceData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) IntervalMAvgEndPriceMacdAndAtrData(ctx context.Context, in *IntervalMAvgEndPriceMacdAndAtrDataRequest, opts ...http.CallOption) (*IntervalMAvgEndPriceMacdAndAtrDataReply, error) {
	var out IntervalMAvgEndPriceMacdAndAtrDataReply
	pattern := "/api/binancedata/interval_m_avg_end_price_macd_and_atr_data"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataIntervalMAvgEndPriceMacdAndAtrData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) IntervalMKAndMACDData(ctx context.Context, in *IntervalMKAndMACDDataRequest, opts ...http.CallOption) (*IntervalMKAndMACDDataReply, error) {
	var out IntervalMKAndMACDDataReply
	pattern := "/api/binancedata/interval_m_k_and_macd_data"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataIntervalMKAndMACDData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) IntervalMMACDData(ctx context.Context, in *IntervalMMACDDataRequest, opts ...http.CallOption) (*IntervalMMACDDataReply, error) {
	var out IntervalMMACDDataReply
	pattern := "/api/binancedata/interval_m_macd_data"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataIntervalMMACDData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) KAnd2NIntervalMAvgEndPriceData(ctx context.Context, in *KAnd2NIntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (*KAnd2NIntervalMAvgEndPriceDataReply, error) {
	var out KAnd2NIntervalMAvgEndPriceDataReply
	pattern := "/api/binancedata/k_and_2_n_interval_m_avg_end_price_data"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationBinanceDataKAnd2NIntervalMAvgEndPriceData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in.SendBody, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) Order(ctx context.Context, in *OrderRequest, opts ...http.CallOption) (*OrderReply, error) {
	var out OrderReply
	pattern := "/api/binancedata/order"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataOrder))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) PullBinanceData(ctx context.Context, in *PullBinanceDataRequest, opts ...http.CallOption) (*PullBinanceDataReply, error) {
	var out PullBinanceDataReply
	pattern := "/api/binancedata/pull"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationBinanceDataPullBinanceData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *BinanceDataHTTPClientImpl) XNIntervalMAvgEndPriceData(ctx context.Context, in *XNIntervalMAvgEndPriceDataRequest, opts ...http.CallOption) (*XNIntervalMAvgEndPriceDataReply, error) {
	var out XNIntervalMAvgEndPriceDataReply
	pattern := "/api/binancedata/x_n_interval_m_avg_end_price_data"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationBinanceDataXNIntervalMAvgEndPriceData))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in.SendBody, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
