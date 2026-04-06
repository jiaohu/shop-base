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
