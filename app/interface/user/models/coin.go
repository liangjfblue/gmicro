package models

type CoinAddRequest struct {
	Uid   string `json:"uid"`
	Value int32  `json:"value"`
}

type CoinAddRespond struct {
	Code int32 `json:"code"`
}

type CoinGetRequest struct {
	Uid string `json:"uid"`
}

type CoinGetRespond struct {
	Value int32 `json:"value"`
}
