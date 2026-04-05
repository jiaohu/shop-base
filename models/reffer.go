package models

import "encoding/json"

type RegisterByRefferCode struct {
	TxHash         string `json:"txHash"`
	RefferedByCode string `json:"refferedByCode"`
}

func (r RegisterByRefferCode) String() string {
	res, err := json.Marshal(r)
	if err != nil {
		return ""
	}
	return string(res)
}
