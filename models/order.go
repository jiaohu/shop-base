package models

import "encoding/json"

type OrderPayEventData struct {
	TxHash string `json:"txHash"`
}

func (o OrderPayEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

type OrderDetailRefundEventData struct {
	OrderDetailId uint64 `json:"orderDetailId"`
	TxHash        string `json:"txHash"`
}

func (o OrderDetailRefundEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

type RefundStockSyncEventData struct {
	AfterSalesId uint64 `json:"afterSalesId"`
	SkuId        uint64 `json:"skuId"`
	Quantity     int    `json:"quantity"`
}

type OrderReceiptEventData struct {
	TxHash  string `json:"txHash"`
	OrderId uint64 `json:"orderId"`
}

func (o OrderReceiptEventData) String() string {
	res, err := json.Marshal(o)
	if err != nil {
		return ""
	}
	return string(res)
}

type SettlementCompletedEventData struct {
	TxHash        string   `json:"txHash"`
	SettlementIds []uint64 `json:"settlementIds"`
	CallerAddr    string   `json:"callerAddr"`
	SettledAt     int64    `json:"settledAt"`
}

func (s SettlementCompletedEventData) String() string {
	res, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(res)
}
