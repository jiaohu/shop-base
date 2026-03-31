package models

import "encoding/json"

type RegisterByRefferCode struct {
	TxHash string `json:"txHash"`
}

func (r RegisterByRefferCode) String() string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(res)
}
