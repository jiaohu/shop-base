package models

import "encoding/json"

type OrderPayEventData struct {
	OrderId uint64 `json:"orderId"`
	TxHash  string `json:"txHash"`
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
